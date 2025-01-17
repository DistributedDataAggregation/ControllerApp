package main

import (
	"log"
	"sync"
)

type Processor struct {
	Planner         Planner
	ExecutorsClient ExecutorsClient
}

func NewProcessor(planner *Planner, executorsClient *ExecutorsClient) *Processor {
	return &Processor{Planner: *planner, ExecutorsClient: *executorsClient}
}

func (p *Processor) ProcessRequest(guid string, queryReq HttpQueryRequest) QueueResult {

	availableIdxs, err := p.ExecutorsClient.GetAvailableExecutorIdxs()

	if err != nil {
		log.Printf("[%s] Could not proccess request: %s", guid, err.Error())
		return QueueResult{ErrorMessage: err.Error(), HttpErrorCode: 500}
	}

	files, err := findDataFiles(queryReq.TableName)
	if err != nil || len(files) == 0 {
		log.Printf("[%s] Could not retreive data files: %s", guid, err.Error())
		return QueueResult{ErrorMessage: "Could not retreive data files", HttpErrorCode: 500}
	}
	if len(files) == 0 {
		log.Printf("[%s] Could not find data files", guid)
		return QueueResult{ErrorMessage: "Could not find data files", HttpErrorCode: 400}
	}

	valid, err := p.validateFilesSchema(files)
	if err != nil {
		log.Printf("[%s] Failed to process request: %v", guid, err)
		return QueueResult{ErrorMessage: "Could not validate table schemat", HttpErrorCode: 500}
	}
	if !valid {
		log.Printf("[%s] Failed to process request", guid)
		return QueueResult{ErrorMessage: "Invalid table schema", HttpErrorCode: 500}
	}

	filesPerExecutorIdx, executorsIdxs := p.Planner.distributeFiles(files, availableIdxs)

	result := p.sendToExecutors(guid, filesPerExecutorIdx, executorsIdxs, queryReq)
	if result.HttpErrorCode != 0 {
		log.Printf("[%s] Failed to process request: %s", guid, result.ErrorMessage)
	}

	return result
}

func (p *Processor) validateFilesSchema(files []string) (bool, error) {
	schema, err := GetParquetSchemaByPath(files[0])
	if err != nil {
		return false, err
	}
	for i := 1; i < len(files); i++ {
		temp, err := GetParquetSchemaByPath(files[i])
		if err != nil {
			return false, err
		}
		if !EqualsParquetSchema(schema, temp) {
			log.Printf("files %s and %s have different schema", files[0], files[i])
			return false, nil
		}
	}
	return true, nil
}

func (p *Processor) sendToExecutors(guid string, filesPerExecutorIdx map[int][]string, executorsIdxs []int, queryReq HttpQueryRequest) QueueResult {

	mainExecutorIdx := *p.ExecutorsClient.MainIdx

	err := p.ExecutorsClient.SendTaskToExecutor(guid, filesPerExecutorIdx[mainExecutorIdx], mainExecutorIdx, int32(len(executorsIdxs)), queryReq)
	if err != nil {
		return QueueResult{ErrorMessage: err.Error(), HttpErrorCode: 500}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(executorsIdxs))

	for _, executorIdx := range executorsIdxs {
		if executorIdx != mainExecutorIdx {
			wg.Add(1)
			go func(executorIdx int) {
				defer wg.Done()
				err := p.ExecutorsClient.SendTaskToExecutor(
					guid,
					filesPerExecutorIdx[executorIdx],
					executorIdx,
					int32(len(executorsIdxs)),
					queryReq,
				)
				errChan <- err
			}(executorIdx)
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return QueueResult{ErrorMessage: err.Error(), HttpErrorCode: 500}
		}
	}

	return p.ExecutorsClient.ReceiveResponseFromMainExecutor(guid)
}
