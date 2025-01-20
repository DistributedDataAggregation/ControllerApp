package main

import (
	"testing"
)

func TestValidateQueryRequest(t *testing.T) {
	tests := []struct {
		name        string
		queryReq    HttpQueryRequest
		expectedErr string
	}{
		{
			name:        "valid query request",
			queryReq:    HttpQueryRequest{TableName: "table1", GroupColumns: []string{"col1"}, SelectColumns: []HttpSelect{{Column: "col2", Function: Minimum}}},
			expectedErr: "",
		},
		{
			name:        "missing table name",
			queryReq:    HttpQueryRequest{TableName: "", GroupColumns: []string{"col1"}, SelectColumns: []HttpSelect{{Column: "col2", Function: Minimum}}},
			expectedErr: "validation error: table_name cannot be empty",
		},
		{
			name:        "missing group columns",
			queryReq:    HttpQueryRequest{TableName: "table1", GroupColumns: []string{}, SelectColumns: []HttpSelect{{Column: "col2", Function: Minimum}}},
			expectedErr: "validation error: group_columns cannot be empty",
		},
		{
			name:        "missing select columns",
			queryReq:    HttpQueryRequest{TableName: "table1", GroupColumns: []string{"col1"}, SelectColumns: []HttpSelect{}},
			expectedErr: "validation error: select cannot be empty",
		},
		{
			name:        "duplicate group column",
			queryReq:    HttpQueryRequest{TableName: "table1", GroupColumns: []string{"col1", "col1"}, SelectColumns: []HttpSelect{{Column: "col2", Function: Minimum}}},
			expectedErr: "validation error: duplicate column 'col1' in group_columns",
		},
		{
			name:        "column in both group and select",
			queryReq:    HttpQueryRequest{TableName: "table1", GroupColumns: []string{"col1"}, SelectColumns: []HttpSelect{{Column: "col1", Function: Minimum}}},
			expectedErr: "validation error: column 'col1' cannot be in both group_columns and select",
		},
		{
			name:        "invalid aggregate function",
			queryReq:    HttpQueryRequest{TableName: "table1", GroupColumns: []string{"col1"}, SelectColumns: []HttpSelect{{Column: "col2", Function: HttpAggregateFunction("InvalidFunc")}}},
			expectedErr: "invalid aggregate function InvalidFunc, supported aggregate functions: Minimum, Maximum, Average, Sum, Count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateQueryRequest(tt.queryReq)

			if (err != nil) != (tt.expectedErr != "") {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("expected error message: %v, got: %v", tt.expectedErr, err.Error())
			}
		})
	}
}
