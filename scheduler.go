package main

import (
	"log"

	"github.com/beevik/guid"
)

type QueueRequest struct {
	Guid       *guid.Guid
	Request    HttpQueryRequest
	ResultChan chan HttpResult
	ErrorChan  chan error
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

		result, err := qs.Processor.processRequest(req.Request)
		if err != nil {
			req.ErrorChan <- err
		} else {
			req.ResultChan <- result
		}

		close(req.ResultChan)
		close(req.ErrorChan)
		log.Printf("Processed request [%v]", req.Guid)
	}
}
