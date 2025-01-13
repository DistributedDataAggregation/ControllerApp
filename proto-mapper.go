package main

import (
	"controller/protomodels"
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/proto"
)

func CreateProtoRequest(guid string, files []string, queryReq HttpQueryRequest, mainExecutor string, mainExecutorPort int32,
	isCurrentNodeMain bool, executorsCount int32) (*protomodels.QueryRequest, error) {

	selects := make([]*protomodels.Select, len(queryReq.SelectColumns))
	for i, sel := range queryReq.SelectColumns {
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

func ReadQueryResultProto(data []byte) (HttpResult, string, error) {
	if len(data) == 0 {
		return HttpResult{}, "", fmt.Errorf("empty input data")
	}

	var queryResult protomodels.QueryResult
	err := proto.Unmarshal(data, &queryResult)
	if err != nil {
		log.Printf("Error unmarshalling QueryResponse: %v", err)
		return HttpResult{}, "", err
	}

	httpResult := HttpResult{
		Response: mapQueryResult(&queryResult),
	}

	return httpResult, queryResult.Guid, nil
}

func mapQueryResult(src *protomodels.QueryResult) HttpQueryResponse {
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
				Results:       mapCombinedResults(value.Results),
			}
		}
	}

	return HttpQueryResponse{
		Error:  httpError,
		Values: httpValues,
	}
}

func mapCombinedResults(results []*protomodels.CombinedResult) []HttpPartialResult {
	httpResults := make([]HttpPartialResult, len(results))

	for i, result := range results {
		if result != nil {

			httpResult := HttpPartialResult{
				IsNull:      result.IsNull,
				Aggregation: HttpAggregateFunction(protomodels.Aggregate_name[int32(result.Function)]),
			}

			switch result.Type {
			case protomodels.ResultType_INT:
				if intValue, ok := result.GetValue().(*protomodels.CombinedResult_IntValue); ok {
					httpResult.ResultType = "INT"
					httpResult.IntValue = &intValue.IntValue
				}
			case protomodels.ResultType_FLOAT:
				if floatValue, ok := result.GetValue().(*protomodels.CombinedResult_FloatValue); ok {
					httpResult.ResultType = "FLOAT"
					httpResult.FloatValue = &floatValue.FloatValue
				}
			case protomodels.ResultType_DOUBLE:
				if doubleValue, ok := result.GetValue().(*protomodels.CombinedResult_DoubleValue); ok {
					httpResult.ResultType = "DOUBLE"
					httpResult.DoubleValue = &doubleValue.DoubleValue
				}
			default:
				httpResult.ResultType = "UNKNOWN"
			}

			httpResults[i] = httpResult
		}
	}
	return httpResults
}

func printProtoRequest(queryReq *protomodels.QueryRequest, adress string) {
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
