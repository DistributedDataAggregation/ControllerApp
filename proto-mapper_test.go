package main

import (
	"reflect"
	"testing"

	"controller/protomodels"

	"google.golang.org/protobuf/proto"
)

func TestCreateProtoRequest(t *testing.T) {
	tests := []struct {
		guid              string
		name              string
		files             []string
		queryReq          HttpQueryRequest
		mainExecutor      string
		mainExecutorPort  int32
		isCurrentNodeMain bool
		executorsCount    int32
		expected          *protomodels.QueryRequest
	}{
		{
			guid:  "guid",
			name:  "basic test",
			files: []string{"file1", "file2"},
			queryReq: HttpQueryRequest{
				TableName:    "test_table",
				GroupColumns: []string{"col1", "col2"},
				SelectColumns: []HttpSelect{
					{Column: "col1", Function: "Minimum"},
					{Column: "col2", Function: "Maximum"},
				},
			},
			mainExecutor:      "172.20.0.2:8080",
			mainExecutorPort:  8080,
			isCurrentNodeMain: true,
			executorsCount:    3,
			expected: &protomodels.QueryRequest{
				Guid:         "guid",
				FilesNames:   []string{"file1", "file2"},
				GroupColumns: []string{"col1", "col2"},
				Select: []*protomodels.Select{
					{Column: "col1", Function: protomodels.Aggregate(protomodels.Aggregate_value["Minimum"])},
					{Column: "col2", Function: protomodels.Aggregate(protomodels.Aggregate_value["Maximum"])},
				},
				Executor: &protomodels.ExecutorInformation{
					IsCurrentNodeMain: true,
					MainIpAddress:     "172.20.0.2",
					MainPort:          8080,
					ExecutorsCount:    3,
				},
			},
		},
		{
			guid:  "guid",
			name:  "no files",
			files: []string{},
			queryReq: HttpQueryRequest{
				TableName:    "test_table",
				GroupColumns: []string{},
				SelectColumns: []HttpSelect{
					{Column: "col1", Function: "Average"},
				},
			},
			mainExecutor:      "172.20.0.3:9090",
			mainExecutorPort:  9090,
			isCurrentNodeMain: false,
			executorsCount:    1,
			expected: &protomodels.QueryRequest{
				Guid:         "guid",
				FilesNames:   []string{},
				GroupColumns: []string{},
				Select: []*protomodels.Select{
					{Column: "col1", Function: protomodels.Aggregate(protomodels.Aggregate_value["Average"])},
				},
				Executor: &protomodels.ExecutorInformation{
					IsCurrentNodeMain: false,
					MainIpAddress:     "172.20.0.3",
					MainPort:          9090,
					ExecutorsCount:    1,
				},
			},
		},
		{
			guid:  "guid",
			name:  "multiple executors",
			files: []string{"file1", "file2", "file3"},
			queryReq: HttpQueryRequest{
				TableName:    "multi_exec_table",
				GroupColumns: []string{"group_col1"},
				SelectColumns: []HttpSelect{
					{Column: "metric", Function: "Minimum"},
				},
			},
			mainExecutor:      "172.20.1.1:8080",
			mainExecutorPort:  8080,
			isCurrentNodeMain: true,
			executorsCount:    5,
			expected: &protomodels.QueryRequest{
				Guid:         "guid",
				FilesNames:   []string{"file1", "file2", "file3"},
				GroupColumns: []string{"group_col1"},
				Select: []*protomodels.Select{
					{Column: "metric", Function: protomodels.Aggregate(protomodels.Aggregate_value["Minimum"])},
				},
				Executor: &protomodels.ExecutorInformation{
					IsCurrentNodeMain: true,
					MainIpAddress:     "172.20.1.1",
					MainPort:          8080,
					ExecutorsCount:    5,
				},
			},
		},
		{
			guid:  "guid",
			name:  "no select columns",
			files: []string{"file1"},
			queryReq: HttpQueryRequest{
				TableName:     "empty_select_table",
				GroupColumns:  []string{"col1"},
				SelectColumns: []HttpSelect{},
			},
			mainExecutor:      "172.20.2.2:7070",
			mainExecutorPort:  7070,
			isCurrentNodeMain: false,
			executorsCount:    2,
			expected: &protomodels.QueryRequest{
				Guid:         "guid",
				FilesNames:   []string{"file1"},
				GroupColumns: []string{"col1"},
				Select:       []*protomodels.Select{},
				Executor: &protomodels.ExecutorInformation{
					IsCurrentNodeMain: false,
					MainIpAddress:     "172.20.2.2",
					MainPort:          7070,
					ExecutorsCount:    2,
				},
			},
		},
		{
			guid:  "guid",
			name:  "empty group and select",
			files: []string{"fileA", "fileB"},
			queryReq: HttpQueryRequest{
				TableName:     "empty_group_select_table",
				GroupColumns:  []string{},
				SelectColumns: []HttpSelect{},
			},
			mainExecutor:      "172.20.3.3:6060",
			mainExecutorPort:  6060,
			isCurrentNodeMain: true,
			executorsCount:    4,
			expected: &protomodels.QueryRequest{
				Guid:         "guid",
				FilesNames:   []string{"fileA", "fileB"},
				GroupColumns: []string{},
				Select:       []*protomodels.Select{},
				Executor: &protomodels.ExecutorInformation{
					IsCurrentNodeMain: true,
					MainIpAddress:     "172.20.3.3",
					MainPort:          6060,
					ExecutorsCount:    4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateProtoRequest(tt.guid, tt.files, tt.queryReq, tt.mainExecutor, tt.mainExecutorPort, tt.isCurrentNodeMain, tt.executorsCount)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReadResponseFromMainExecutor(t *testing.T) {
	tests := []struct {
		name           string
		guid           string
		data           []byte
		expectedResult QueueResult
		expectedError  bool
	}{
		{
			name: "valid response",
			guid: "guid",
			data: func() []byte {
				queryResponse := &protomodels.QueryResult{
					Guid:  "guid",
					Error: nil,
					Values: []*protomodels.ResultValue{
						{
							GroupingValue: "group1",
							Results: []*protomodels.CombinedResult{
								{
									IsNull:   false,
									Type:     protomodels.ResultType_INT,
									Function: protomodels.Aggregate_Maximum,
									Value:    &protomodels.CombinedResult_IntValue{IntValue: 100},
								},
							},
						},
					},
				}
				data, _ := proto.Marshal(queryResponse)
				return data
			}(),
			expectedResult: QueueResult{
				QueryResponse: &HttpQueryResponse{
					Values: []*HttpValue{
						{
							GroupingValue: "group1",
							Results: []HttpPartialResult{
								{
									IsNull:      false,
									IntValue:    ptrInt64(100),
									ResultType:  "INT",
									Aggregation: HttpAggregateFunction(protomodels.Aggregate_name[int32(protomodels.Aggregate_Maximum)]),
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:           "invalid data",
			guid:           "guid",
			data:           []byte("invalid protobuf data"),
			expectedResult: QueueResult{ErrorMessage: "Failed to read data from executor", HttpErrorCode: 500},
			expectedError:  true,
		},
		{
			name:           "empty data",
			guid:           "guid",
			data:           []byte{},
			expectedResult: QueueResult{ErrorMessage: "Failed to read data from executor", HttpErrorCode: 500},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, receivedGuid := ReadQueryResultProto(tt.data)

			if (result.HttpErrorCode != 0) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, result.ErrorMessage)
			}

			if tt.expectedError {
				return // Skip further checks if an error is expected
			}

			if receivedGuid != tt.guid {
				t.Errorf("expected guid: %s, got: %s", tt.guid, receivedGuid)

			}

			if !reflect.DeepEqual(&result, &tt.expectedResult) {
				t.Errorf("expected response: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func ptrInt64(v int64) *int64 {
	return &v
}
