package main

import (
	"log"

	"github.com/beevik/guid"
)

type QueueResult struct {
	QueryResponse *HttpQueryResponse
	HttpErrorCode int
	ErrorMessage  string
}

type QueueRequest struct {
	Guid       *guid.Guid
	Request    HttpQueryRequest
	ResultChan chan QueueResult
}

type QueriesScheduler struct {
	Queue     chan QueueRequest
	Processor Processor
}

func NewQueriesScheduler(processor *Processor) *QueriesScheduler {
	qs := &QueriesScheduler{
		Queue:     make(chan QueueRequest, 100),
		Processor: *processor,
	}
	go qs.processQueue()
	return qs
}

func (qs *QueriesScheduler) AddQuery(request QueueRequest) {
	qs.Queue <- request
	log.Printf("Added request [%v] to the queue", request.Guid)
}

func (qs *QueriesScheduler) processQueue() {
	for req := range qs.Queue {
		log.Printf("Processing request [%v]", req.Guid)

		result := qs.Processor.ProcessRequest(req.Guid.String(), req.Request)
		req.ResultChan <- result

		close(req.ResultChan)
		log.Printf("Finished processing request [%v]", req.Guid)
	}
}
