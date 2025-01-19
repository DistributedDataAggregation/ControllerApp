package main

import (
	"controller/protomodels"
	"log"
	"strings"

	"google.golang.org/protobuf/proto"
)

func CreateProtoRequest(guid string, files []string, queryReq HttpQueryRequest, mainExecutor string, mainExecutorPort int32,
	isCurrentNodeMain bool, executorsCount int32) *protomodels.QueryRequest {

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
	}
}

func ReadQueryResultProto(data []byte) (QueueResult, string) {
	if len(data) == 0 {
		log.Printf("Error reading response from main executor: Empty imput data")
		return QueueResult{ErrorMessage: "Failed to read data from executor", HttpErrorCode: 500}, ""
	}

	var queryResult protomodels.QueryResult
	err := proto.Unmarshal(data, &queryResult)
	if err != nil {
		log.Printf("Error unmarshalling QueryResult: %v", err)
		return QueueResult{ErrorMessage: "Failed to read data from executor", HttpErrorCode: 500}, ""
	}

	return mapQueryResult(&queryResult), queryResult.Guid
}

func mapQueryResult(src *protomodels.QueryResult) QueueResult {
	if src == nil {
		return QueueResult{}
	}

	if src.Error != nil {
		return QueueResult{ErrorMessage: src.Error.Message, HttpErrorCode: 500} // TODO sometimes 400
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

	return QueueResult{
		QueryResponse: &HttpQueryResponse{Values: httpValues},
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
				httpResult.ResultType = "INT"
				if intValue, ok := result.GetValue().(*protomodels.CombinedResult_IntValue); ok {
					httpResult.IntValue = &intValue.IntValue
				}
			case protomodels.ResultType_FLOAT:
				httpResult.ResultType = "FLOAT"
				if floatValue, ok := result.GetValue().(*protomodels.CombinedResult_FloatValue); ok {
					httpResult.FloatValue = &floatValue.FloatValue
				}
			case protomodels.ResultType_DOUBLE:
				httpResult.ResultType = "DOUBLE"
				if doubleValue, ok := result.GetValue().(*protomodels.CombinedResult_DoubleValue); ok {
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
