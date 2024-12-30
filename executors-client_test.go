package main

import (
	"reflect"
	"testing"

	"controller/protomodels"

	"google.golang.org/protobuf/proto"
)

func TestCreateProtoRequest(t *testing.T) {
	tests := []struct {
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
			name:  "multiple executors",
			files: []string{"file1", "file2", "file3"},
			queryReq: HttpQueryRequest{
				TableName:    "multi_exec_table",
				GroupColumns: []string{"group_col1"},
				SelectColumns: []HttpSelect{
					{Column: "metric", Function: "Median"},
				},
			},
			mainExecutor:      "172.20.1.1:8080",
			mainExecutorPort:  8080,
			isCurrentNodeMain: true,
			executorsCount:    5,
			expected: &protomodels.QueryRequest{
				FilesNames:   []string{"file1", "file2", "file3"},
				GroupColumns: []string{"group_col1"},
				Select: []*protomodels.Select{
					{Column: "metric", Function: protomodels.Aggregate(protomodels.Aggregate_value["Median"])},
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
			client := &ExecutorsClient{}
			result := client.createProtoRequest(tt.files, tt.queryReq, tt.mainExecutor, tt.mainExecutorPort, tt.isCurrentNodeMain, tt.executorsCount)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReadResponseFromMainExecutor(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expectedResult HttpResult
		expectedError  bool
	}{
		{
			name: "valid response",
			data: func() []byte {
				queryResponse := &protomodels.QueryResponse{
					Values: []*protomodels.Value{
						{
							GroupingValue: "group1",
							Results: []*protomodels.PartialResult{
								{Value: 100, Count: 2},
							},
						},
					},
				}
				data, _ := proto.Marshal(queryResponse)
				return data
			}(),
			expectedResult: HttpResult{
				Response: protomodels.QueryResponse{
					Values: []*protomodels.Value{
						{
							GroupingValue: "group1",
							Results: []*protomodels.PartialResult{
								{Value: 100, Count: 2},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:           "invalid data",
			data:           []byte("invalid protobuf data"),
			expectedResult: HttpResult{},
			expectedError:  true,
		},
		{
			name:           "empty data",
			data:           []byte{},
			expectedResult: HttpResult{},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &ExecutorsClient{}
			result, err := client.readResponseFromMainExecutor(tt.data)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			if tt.expectedError {
				return // Skip further checks if an error is expected
			}

			// Use proto.Equal for protobuf comparisons
			if !proto.Equal(&result.Response, &tt.expectedResult.Response) {
				t.Errorf("expected response: %v, got: %v", tt.expectedResult.Response, result.Response)
			}
		})
	}
}
