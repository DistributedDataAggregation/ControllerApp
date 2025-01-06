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
