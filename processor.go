package main

import (
	"fmt"
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

func (p *Processor) ProcessRequest(guid string, queryReq HttpQueryRequest) HttpResult {

	availableIdxs, err := p.ExecutorsClient.GetAvailableExecutorIdxs()

	if err != nil {
		log.Printf("Could not proccess request: %s", err.Error())

		return HttpResult{Response: HttpQueryResponse{
			Error: &HttpError{
				Message:      "Could not proccess request",
				InnerMessage: err.Error(),
			}}}
	}

	files, err := findDataFiles(queryReq.TableName)
	if err != nil || len(files) == 0 {
		innerMessage := ""
		if err != nil {
			innerMessage = err.Error()
		}
		log.Printf("Could not find files: %s", innerMessage)

		return HttpResult{Response: HttpQueryResponse{
			Error: &HttpError{
				Message:      "Could not find files",
				InnerMessage: innerMessage,
			}}}
	}

	err = p.validateFilesSchema(files)
	if err != nil {
		log.Printf("Failed to process request [%s]: %v", guid, err)
		return HttpResult{Response: HttpQueryResponse{
			Error: &HttpError{
				Message:      "Failed to process request",
				InnerMessage: err.Error(),
			}}}
	}

	filesPerExecutorIdx, executorsIdxs := p.Planner.distributeFiles(files, availableIdxs)

	result, err := p.sendToExecutors(guid, filesPerExecutorIdx, executorsIdxs, queryReq)
	if err != nil {
		log.Printf("Failed to process request [%s]: %v", guid, err)
		return HttpResult{Response: HttpQueryResponse{
			Error: &HttpError{
				Message:      "Failed to process request",
				InnerMessage: err.Error(),
			}}}
	}

	return result
}

func (p *Processor) validateFilesSchema(files []string) error {
	schema, err := GetParquetSchemaByPath(files[0])
	if err != nil {
		return err
	}
	for i := 1; i < len(files); i++ {
		temp, err := GetParquetSchemaByPath(files[i])
		if err != nil {
			return err
		}
		if !EqualsParquetSchema(schema, temp) {
			return fmt.Errorf("files %s and %s have different schema", files[0], files[i])
		}
	}
	return nil
}

func (p *Processor) sendToExecutors(guid string, filesPerExecutorIdx map[int][]string, executorsIdxs []int, queryReq HttpQueryRequest) (HttpResult, error) {

	mainExecutorIdx := *p.ExecutorsClient.MainIdx

	err := p.ExecutorsClient.SendTaskToExecutor(guid, filesPerExecutorIdx[mainExecutorIdx], mainExecutorIdx, int32(len(executorsIdxs)), queryReq)
	if err != nil {
		return HttpResult{}, err
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
			return HttpResult{}, err
		}
	}

	response, err := p.ExecutorsClient.ReceiveResponseFromMainExecutor(guid)
	if err != nil {
		return HttpResult{}, err
	}

	return response, nil

}
