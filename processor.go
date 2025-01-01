package main

import (
	"log"
	"path/filepath"
	"sync"
)

type Processor struct {
	Planner         Planner
	ExecutorsClient ExecutorsClient
}

func NewProcessor(planner *Planner, executorsClient *ExecutorsClient) *Processor {
	return &Processor{Planner: *planner, ExecutorsClient: *executorsClient}
}

func (p *Processor) processRequest(queryReq HttpQueryRequest) (HttpResult, error) {

	files, err := p.findDataFiles(queryReq.TableName)
	if err != nil || len(files) == 0 {
		log.Printf("Could not find files %v", err)
		return HttpResult{Response: HttpQueryResponse{Error: &HttpError{Message: "could not find files"}}}, nil
	}

	filesPerExecutorIdx, executorsIdxs := p.Planner.distributeFiles(files, p.ExecutorsClient.Addresses)

	return p.sendToExecutors(filesPerExecutorIdx, executorsIdxs, queryReq)
}

func (p *Processor) findDataFiles(tableName string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(config.DataPath, "*"+tableName+"*"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (p *Processor) sendToExecutors(filesPerExecutorIdx map[int][]string, executorsIdxs []int, queryReq HttpQueryRequest) (HttpResult, error) {

	mainExecutorIdx := p.ExecutorsClient.MainIdx

	err := p.ExecutorsClient.sendTaskToExecutor(filesPerExecutorIdx[mainExecutorIdx], mainExecutorIdx, int32(len(executorsIdxs)), queryReq, nil)
	if err != nil {
		return HttpResult{}, err
	}

	var wg sync.WaitGroup
	for _, executorIdx := range executorsIdxs {
		if executorIdx != mainExecutorIdx {
			wg.Add(1)
			go p.ExecutorsClient.sendTaskToExecutor(filesPerExecutorIdx[executorIdx], executorIdx, int32(len(executorsIdxs)), queryReq, &wg)
		}
	}
	wg.Wait()

	response, err := p.ExecutorsClient.receiveResponseFromMainExecutor()
	if err != nil {
		return HttpResult{}, err
	}

	return response, nil

}
