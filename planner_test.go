package main

import (
	"reflect"
	"testing"
)

func TestDistributeFiles(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		executors  []string
		expMapping map[int][]string
		expUsed    []int
	}{
		{
			name:      "even distribution",
			files:     []string{"file1", "file2", "file3", "file4"},
			executors: []string{"exec1", "exec2"},
			expMapping: map[int][]string{
				0: {"file1", "file3"},
				1: {"file2", "file4"},
			},
			expUsed: []int{0, 1},
		},
		{
			name:      "more executors than files",
			files:     []string{"file1"},
			executors: []string{"exec1", "exec2", "exec3"},
			expMapping: map[int][]string{
				0: {"file1"},
			},
			expUsed: []int{0},
		},
		{
			name:      "uneven distribution",
			files:     []string{"file1", "file2", "file3"},
			executors: []string{"exec1", "exec2"},
			expMapping: map[int][]string{
				0: {"file1", "file3"},
				1: {"file2"},
			},
			expUsed: []int{0, 1},
		},
		{
			name:       "no files",
			files:      []string{},
			executors:  []string{"exec1", "exec2"},
			expMapping: map[int][]string{},
			expUsed:    []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := NewPlanner()
			mapping, used := planner.distributeFiles(tt.files, tt.executors)

			if !reflect.DeepEqual(mapping, tt.expMapping) {
				t.Errorf("expected mapping %v, got %v", tt.expMapping, mapping)
			}

			if !reflect.DeepEqual(used, tt.expUsed) {
				t.Errorf("expected used executors %v, got %v", tt.expUsed, used)
			}
		})
	}
}
