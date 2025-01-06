package main

import (
	"controller/protomodels"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/protobuf/proto"
)

func CreateProtoRequest(guid string, files []string, queryReq HttpQueryRequest, mainExecutor string, mainExecutorPort int32,
	isCurrentNodeMain bool, executorsCount int32) (*protomodels.QueryRequest, error) {

	if err := validateQueryRequest(queryReq); err != nil {
		return nil, err
	}

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
		for _, groupCol := range queryReq.GroupColumns {
			if sel.Column == groupCol {
				log.Printf("Validation error: column '%s' cannot be in both group_columns and select", sel.Column)
				return fmt.Errorf("validation error: column '%s' cannot be in both group_columns and select", sel.Column)
			}
		}

		if !sel.Function.IsValid() {
			log.Printf("Invalid aggregate function %s. Supported aggregate functions: Minimum, Maximum, Average, Sum, Count", string(sel.Function))
			return fmt.Errorf("invalid aggregate function %s, supported aggregate functions: Minimum, Maximum, Average, Sum, Count", string(sel.Function))
		}
	}

	return nil
}

func ReadProtoResponse(data []byte) (HttpResult, string, error) {

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

			httpResult := HttpPartialResult{
				IsNull:      result.IsNull,
				Count:       result.Count,
				Aggregation: protomodels.Aggregate_name[int32(result.Function)],
			}

			switch result.Type {
			case protomodels.ResultType_INT:
				if intValue, ok := result.GetValue().(*protomodels.PartialResult_IntValue); ok {
					httpResult.ResultType = "INT"
					httpResult.Value = &intValue.IntValue
				}
			case protomodels.ResultType_FLOAT:
				if floatValue, ok := result.GetValue().(*protomodels.PartialResult_FloatValue); ok {
					httpResult.ResultType = "FLOAT"
					httpResult.FloatValue = &floatValue.FloatValue
				}
			case protomodels.ResultType_DOUBLE:
				if doubleValue, ok := result.GetValue().(*protomodels.PartialResult_DoubleValue); ok {
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
