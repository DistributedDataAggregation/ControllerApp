package main

import (
	"reflect"
	"testing"
)

func TestDistributeFiles(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		executors  []int
		expMapping map[int][]string
		expUsed    []int
	}{
		{
			name:      "even distribution",
			files:     []string{"file1", "file2", "file3", "file4"},
			executors: []int{3, 5},
			expMapping: map[int][]string{
				3: {"file1", "file3"},
				5: {"file2", "file4"},
			},
			expUsed: []int{3, 5},
		},
		{
			name:      "more executors than files",
			files:     []string{"file1"},
			executors: []int{1, 2, 4},
			expMapping: map[int][]string{
				1: {"file1"},
			},
			expUsed: []int{1},
		},
		{
			name:      "uneven distribution",
			files:     []string{"file1", "file2", "file3"},
			executors: []int{14, 6},
			expMapping: map[int][]string{
				14: {"file1", "file3"},
				6:  {"file2"},
			},
			expUsed: []int{14, 6},
		},
		{
			name:       "no files",
			files:      []string{},
			executors:  []int{1, 3},
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
