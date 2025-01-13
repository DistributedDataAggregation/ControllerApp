package main

import (
	"encoding/json"
	"fmt"
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

type HttpResultType string

const (
	IntResult     HttpResultType = "INT"
	FloatResult   HttpResultType = "FLOAT"
	DoubleResult  HttpResultType = "DOUBLE"
	UnknownResult HttpResultType = "UNKNOWN"
)

type HttpQueryRequest struct {
	TableName     string       `json:"table_name"`    // The name of the table on which the query will be executed.
	GroupColumns  []string     `json:"group_columns"` // The names of the columns on which grouping will be performed.
	SelectColumns []HttpSelect `json:"select"`        // A list of objects describing the columns and the aggregate functions to be executed on them.
}

type HttpSelect struct {
	Column   string                `json:"column"`   // The name of the aggregated column.
	Function HttpAggregateFunction `json:"function"` // The aggregate function.
}

type HttpError struct {
	Message      string `json:"message,omitempty"`       // Error message
	InnerMessage string `json:"inner_message,omitempty"` // Inner error message
}

type HttpPartialResult struct {
	IsNull      bool                  `json:"is_null"`                // Indicates if the result is null.
	IntValue    *int64                `json:"int_value,omitempty"`    // Integer result (nullable).
	FloatValue  *float32              `json:"float_value,omitempty"`  // Float result (nullable).
	DoubleValue *float64              `json:"double_value,omitempty"` // Double result (nullable).
	ResultType  HttpResultType        `json:"result_type"`            // Type of result: "INT", "FLOAT", "DOUBLE".
	Aggregation HttpAggregateFunction `json:"aggregation"`            // The aggregate function.
}

type HttpValue struct {
	GroupingValue string              `json:"grouping_value"` // Grouping value, subsequent grouping column values ​​separated by the '|' character
	Results       []HttpPartialResult `json:"results"`        // List of results of given aggregations for the grouping value.
}

type HttpQueryResponse struct {
	Error  *HttpError   `json:"error,omitempty"`  // Information about the query processing error (if any)
	Values []*HttpValue `json:"values,omitempty"` // List of results of performed aggregations for individual combinations of grouping column values
}

type HttpResult struct {
	Response HttpQueryResponse `json:"result"`          // The result of the processed query.
	Time     int64             `json:"processing_time"` // The total time to process the query in milliseconds.
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

	err = validateQueryRequest(queryReq)

	if err != nil {
		result := HttpResult{Response: HttpQueryResponse{
			Error: &HttpError{
				Message:      "Failed to process request",
				InnerMessage: err.Error(),
			}}}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
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

func validateQueryRequest(queryReq HttpQueryRequest) error {
	if queryReq.TableName == "" {
		log.Println("Validation error: table_name cannot be empty")
		return fmt.Errorf("validation error: table_name cannot be empty")
	}

	if len(queryReq.GroupColumns) == 0 {
		log.Println("Validation error: group_columns cannot be empty")
		return fmt.Errorf("validation error: group_columns cannot be empty")
	}

	if len(queryReq.SelectColumns) == 0 {
		log.Println("Validation error: select cannot be empty")
		return fmt.Errorf("validation error: select cannot be empty")

	}

	seen := make(map[string]bool)
	for _, groupCol := range queryReq.GroupColumns {
		if seen[groupCol] {
			log.Printf("Validation error: duplicate column '%s' in group_columns", groupCol)
			return fmt.Errorf("validation error: duplicate column '%s' in group_columns", groupCol)
		}
		seen[groupCol] = true
	}

	for _, sel := range queryReq.SelectColumns {
		if seen[sel.Column] {
			log.Printf("Validation error: column '%s' cannot be in both group_columns and select", sel.Column)
			return fmt.Errorf("validation error: column '%s' cannot be in both group_columns and select", sel.Column)
		}

		if !sel.Function.IsValid() {
			log.Printf("Invalid aggregate function %s. Supported aggregate functions: Minimum, Maximum, Average, Sum, Count", string(sel.Function))
			return fmt.Errorf("invalid aggregate function %s, supported aggregate functions: Minimum, Maximum, Average, Sum, Count", string(sel.Function))
		}
	}

	return nil
}
