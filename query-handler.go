package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
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

func handleQuery(w http.ResponseWriter, r *http.Request) {
	var queryReq HttpQueryRequest
	err := json.NewDecoder(r.Body).Decode(&queryReq)
	if err != nil {
		log.Fatalf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	files, err := findDataFiles(queryReq.TableName)
	if err != nil {
		log.Fatalf("Could not find files %v", err)
		http.Error(w, "Could not find files", http.StatusInternalServerError)
		return
	}

	err = sendToExecutors(files, queryReq)
	if err != nil {
		log.Fatalf("Error processing request: %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Query has been processed"))
}

func findDataFiles(tableName string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(config.DataPath, "*"+tableName+"*"))
	if err != nil {
		return nil, err
	}
	return files, nil
}
