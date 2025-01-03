package main

import (
	"controller/protomodels"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/protobuf/proto"
)

func createProtoRequest(guid string, files []string, queryReq HttpQueryRequest, mainExecutor string, mainExecutorPort int32,
	isCurrentNodeMain bool, executorsCount int32) (*protomodels.QueryRequest, error) {
	selects := make([]*protomodels.Select, len(queryReq.SelectColumns))
	for i, sel := range queryReq.SelectColumns {

		if !sel.Function.IsValid() {
			log.Printf("Invalid aggreage function %s. Supported aggregate functions: Minimum, Maximum, Average", string(sel.Function))
			return nil, fmt.Errorf("invalid aggreage function %s, supported aggregate functions: Minimum, Maximum, Average", string(sel.Function))
		}

		selects[i] = &protomodels.Select{
			Column:   sel.Column,
			Function: protomodels.Aggregate(protomodels.Aggregate_value[string(sel.Function)]),
		}
	}

	return &protomodels.QueryRequest{
		Guid:         guid,
		FilesNames:   files,
		GroupColumns: queryReq.GroupColumns,
		Select:       selects,
		Executor: &protomodels.ExecutorInformation{
			IsCurrentNodeMain: isCurrentNodeMain,
			MainIpAddress:     strings.Split(mainExecutor, ":")[0],
			MainPort:          mainExecutorPort,
			ExecutorsCount:    executorsCount,
		},
	}, nil
}

func readProtoResponse(data []byte) (HttpResult, string, error) {

	if len(data) == 0 {
		return HttpResult{}, "", fmt.Errorf("empty input data")
	}

	var queryResponse protomodels.QueryResponse
	err := proto.Unmarshal(data, &queryResponse)
	if err != nil {
		log.Printf("Error unmarshalling QueryResponse: %v", err)
		return HttpResult{}, "", err
	}

	httpResult := HttpResult{
		Response: mapQueryResponse(&queryResponse),
	}

	return httpResult, queryResponse.Guid, nil
}

func mapQueryResponse(src *protomodels.QueryResponse) HttpQueryResponse {
	if src == nil {
		return HttpQueryResponse{}
	}

	var httpError *HttpError
	if src.Error != nil {
		httpError = &HttpError{
			Message:      src.Error.Message,
			InnerMessage: src.Error.InnerMessage,
		}
	}

	httpValues := make([]*HttpValue, len(src.Values))
	for i, value := range src.Values {
		if value != nil {
			httpValues[i] = &HttpValue{
				GroupingValue: value.GroupingValue,
				Results:       mapPartialResults(value.Results),
			}
		}
	}

	return HttpQueryResponse{
		Error:  httpError,
		Values: httpValues,
	}
}

func mapPartialResults(results []*protomodels.PartialResult) []HttpPartialResult {
	httpResults := make([]HttpPartialResult, len(results))
	for i, result := range results {
		if result != nil {
			httpResults[i] = HttpPartialResult{
				Value: result.Value,
				Count: result.Count,
			}
		}
	}
	return httpResults
}

func printProtoRequest(queryReq *protomodels.QueryRequest, adress net.Addr) {
	log.Printf("Sent request to %s\n", adress)
	log.Printf("Files:\n")
	for _, file := range queryReq.FilesNames {
		log.Printf("	%s, \n", file)
	}

	log.Printf("Grouping columns:\n")
	for _, col := range queryReq.Select {
		log.Printf("	%s, \n", col)
	}

	log.Printf("Select:\n")
	for _, sel := range queryReq.Select {
		log.Printf("	Column %s, Function %s,\n", sel.Column, sel.Function)
	}

	log.Printf("Main executor:\n")
	log.Printf("	Is main: %t", queryReq.Executor.IsCurrentNodeMain)
	log.Printf("	Main ip address: %s", queryReq.Executor.MainIpAddress)
	log.Printf("	Main port: %d", queryReq.Executor.MainPort)
	log.Printf("	Executors count: %d", queryReq.Executor.ExecutorsCount)

}