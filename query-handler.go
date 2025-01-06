package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/beevik/guid"
)

type HttpAggregateFunction string

const (
	Minimum HttpAggregateFunction = "Minimum"
	Maximum HttpAggregateFunction = "Maximum"
	Average HttpAggregateFunction = "Average"
	Sum     HttpAggregateFunction = "Sum"
	Count   HttpAggregateFunction = "Count"
)

func (h HttpAggregateFunction) IsValid() bool {
	switch h {
	case Minimum, Maximum, Average, Sum, Count:
		return true
	}
	return false
}

type HttpQueryRequest struct {
	TableName     string       `json:"table_name"`
	GroupColumns  []string     `json:"group_columns"`
	SelectColumns []HttpSelect `json:"select"`
}

type HttpSelect struct {
	Column   string                `json:"column"`
	Function HttpAggregateFunction `json:"function"`
}

type HttpError struct {
	Message      string `json:"message"`
	InnerMessage string `json:"inner_message"`
}

type HttpPartialResult struct {
	IsNull      bool     `json:"is_null"`                // Indicates if the result is null.
	Value       *int64   `json:"value,omitempty"`        // Integer value (nullable).
	FloatValue  *float32 `json:"float_value,omitempty"`  // Float value (nullable).
	DoubleValue *float64 `json:"double_value,omitempty"` // Double value (nullable).
	Count       int64    `json:"count"`                  // Count associated with the result.
	ResultType  string   `json:"result_type"`            // Type of result: "INT", "FLOAT", "DOUBLE".
	Aggregation string   `json:"aggregation"`
}

type HttpValue struct {
	GroupingValue string              `json:"grouping_value"`
	Results       []HttpPartialResult `json:"results"`
}

type HttpQueryResponse struct {
	Error  *HttpError   `json:"error"`
	Values []*HttpValue `json:"values"`
}

type HttpResult struct {
	Response HttpQueryResponse `json:"result"`
	Time     int64             `json:"processing_time"`
}

type QueryHandler struct {
	Scheduler QueriesScheduler
}

func NewQueryHandler(scheduler *QueriesScheduler) *QueryHandler {
	return &QueryHandler{Scheduler: *scheduler}
}

// @Summary Query data from a table
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
	}

	h.Scheduler.AddQuery(queueReq)

	result := <-queueReq.ResultChan

	result.Time = time.Since(start).Milliseconds()

	if result.Response.Error == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(result)
}
