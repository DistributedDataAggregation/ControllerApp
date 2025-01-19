package main

import "testing"

func TestProcessRequest(t *testing.T) {

	t.Run("no available executors", func(t *testing.T) {
		mockPlanner := Planner{}

		mockExecutorsClient := ExecutorsClient{
			MainIdx:        ptrInt(-1),
			SocketStatuses: []bool{},
		}

		processor := NewProcessor(&mockPlanner, &mockExecutorsClient)

		guid := "test-guid"
		queryRequest := HttpQueryRequest{
			TableName:     "test_table",
			GroupColumns:  []string{"col1"},
			SelectColumns: []HttpSelect{{Column: "col2", Function: Minimum}},
		}

		result := processor.ProcessRequest(guid, queryRequest)

		if result.HttpErrorCode != 500 || result.ErrorMessage != "main executor is unavailable" {
			t.Errorf("Expected error 'main executor is unavailable' with HTTP 500, got: %+v", result)
		}
	})
}
