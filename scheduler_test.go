package main

import (
	"testing"

	"github.com/beevik/guid"
)

func TestAddQuery(t *testing.T) {

	t.Run("no available executors", func(t *testing.T) {
		mockPlanner := Planner{}

		mockExecutorsClient := ExecutorsClient{
			MainIdx:        ptrInt(-1),
			SocketStatuses: []bool{},
		}

		processor := NewProcessor(&mockPlanner, &mockExecutorsClient)

		scheduler := NewQueriesScheduler(processor)

		guid := guid.New()
		request := QueueRequest{
			Guid: guid,
			Request: HttpQueryRequest{
				TableName:     "test_table",
				GroupColumns:  []string{"col1"},
				SelectColumns: []HttpSelect{{Column: "col2", Function: Minimum}},
			},
			ResultChan: make(chan QueueResult, 1),
		}

		scheduler.AddQuery(request)

		result := <-request.ResultChan

		if result.HttpErrorCode != 500 || result.ErrorMessage != "main executor is unavailable" {
			t.Errorf("Expected error 'main executor is unavailable' with HTTP 500, got: %+v", result)
		}
	})
}
