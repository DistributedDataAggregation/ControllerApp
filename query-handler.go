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
)

func (h HttpAggregateFunction) IsValid() bool {
	switch h {
	case Minimum, Maximum, Average:
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
	Value  int64 `json:"value"`
	Count  int64 `json:"count"`
	IsNull bool  `json:"is_null"`
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

func (h *QueryHandler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3006")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
