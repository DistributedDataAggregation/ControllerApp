package main

import (
	"reflect"
	"testing"
)

func TestGetAvailableExecutorIdxs(t *testing.T) {
	tests := []struct {
		name        string
		client      *ExecutorsClient
		expectedIdx []int
		expectedErr string
	}{
		{
			name: "main executor available",
			client: &ExecutorsClient{
				MainIdx:        ptrInt(0),
				SocketStatuses: []bool{true, false, true},
			},
			expectedIdx: []int{0, 2},
			expectedErr: "",
		},
		{
			name: "main executor unavailable",
			client: &ExecutorsClient{
				MainIdx:        ptrInt(0),
				SocketStatuses: []bool{false, true, true},
			},
			expectedIdx: nil,
			expectedErr: "main executor is unavailable",
		},
		{
			name: "no executors available",
			client: &ExecutorsClient{
				MainIdx:        ptrInt(0),
				SocketStatuses: []bool{false, false, false},
			},
			expectedIdx: nil,
			expectedErr: "main executor is unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idxs, err := tt.client.GetAvailableExecutorIdxs()

			if (err != nil) != (tt.expectedErr != "") {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("expected error message: %v, got: %v", tt.expectedErr, err.Error())
			}

			if !reflect.DeepEqual(idxs, tt.expectedIdx) {
				t.Errorf("expected indices: %v, got: %v", tt.expectedIdx, idxs)
			}
		})
	}
}

func TestGetFirstAvailableExecutor(t *testing.T) {
	tests := []struct {
		name        string
		client      *ExecutorsClient
		expectedIdx int
	}{
		{
			name: "first executor available",
			client: &ExecutorsClient{
				SocketStatuses: []bool{true, false, true},
			},
			expectedIdx: 0,
		},
		{
			name: "second executor available",
			client: &ExecutorsClient{
				SocketStatuses: []bool{false, true, false},
			},
			expectedIdx: 1,
		},
		{
			name: "no executors available",
			client: &ExecutorsClient{
				SocketStatuses: []bool{false, false, false},
			},
			expectedIdx: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.client.getFirstAvailableExecutor()

			if idx != tt.expectedIdx {
				t.Errorf("expected index: %d, got: %d", tt.expectedIdx, idx)
			}
		})
	}
}
