package main

import (
	"controller/protomodels"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/beevik/guid"
)

// type Aggregate int

// const (
// 	AggregateMinimum Aggregate = iota
// 	AggregateMaximum
// 	AggregateAverage
// 	AggregateMedian
// )

// var aggregateName = map[Aggregate]string{
// 	AggregateMinimum: "minimum",
// 	AggregateMaximum: "maximum",
// 	AggregateAverage: "average",
// 	AggregateMedian:  "median",
// }

// func (a Aggregate) String() string {
// 	return aggregateName[a]
// }

type HttpQueryRequest struct {
	TableName     string       `json:"table_name"`
	GroupColumns  []string     `json:"group_columns"`
	SelectColumns []HttpSelect `json:"select"`
}

type HttpSelect struct {
	Column   string `json:"column"`
	Function string `json:"function"`
}

type HttpResult struct {
	Response protomodels.QueryResponse `json:"result"`
	Time     int64                     `json:"processing_time"`
}

type QueryHandler struct {
	Scheduler QueriesScheduler
}

func NewQueryHandler(scheduler *QueriesScheduler) *QueryHandler {
	return &QueryHandler{Scheduler: *scheduler}
}

// @Summary Query data from table
// @Description Queries data with specified grouping and selection
// @Tags query
// @Accept  json
// @Produce  json
// @Param query body HttpQueryRequest true "Query Request"
// @Success 200 {object} HttpResult "Query has been processed"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /query [post]
func (h *QueryHandler) handleQuery(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var queryReq HttpQueryRequest
	err := json.NewDecoder(r.Body).Decode(&queryReq)
	if err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	queueReq := QueueRequest{
		Guid:       guid.New(),
		Request:    queryReq,
		ResultChan: make(chan HttpResult),
		ErrorChan:  make(chan error),
	}

	h.Scheduler.AddQuery(queueReq)

	result, err := <-queueReq.ResultChan, <-queueReq.ErrorChan // TODO ErrorChan

	if err != nil {
		log.Printf("Error processing request %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	result.Time = time.Since(start).Milliseconds()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
