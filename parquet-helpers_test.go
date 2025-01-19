package main

import (
	"reflect"
	"testing"
)

func TestGetParquetColumnType(t *testing.T) {
	tests := []struct {
		name         string
		colType      string
		expectedType string
	}{
		{"valid INT type", "INT32", "INT"},
		{"valid DOUBLE type", "DOUBLE", "DOUBLE"},
		{"valid FLOAT type", "FLOAT", "FLOAT"},
		{"valid STRING type", "STRING", "STRING"},
		{"valid BOOL type", "BOOL", "BOOL"},
		{"unsupported type", "DATE", "UNSUPPORTED"},
		{"complex type with valid INT", "INT(64)", "INT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getParquetColumnType(tt.colType)

			if result != tt.expectedType {
				t.Errorf("expected type: %v, got: %v", tt.expectedType, result)
			}
		})
	}
}

func TestNormalizeType(t *testing.T) {
	tests := []struct {
		name         string
		colType      string
		expectedType string
	}{
		{"simple INT type", "int32", "INT"},
		{"simple DOUBLE type", "double", "DOUBLE"},
		{"complex type with params", "float(32)", "FLOAT"},
		{"unsupported type", "date", "UNSUPPORTED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeType(tt.colType)

			if result != tt.expectedType {
				t.Errorf("expected type: %v, got: %v", tt.expectedType, result)
			}
		})
	}
}

func TestEqualsParquetSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema1  []ParquetColumnInfo
		schema2  []ParquetColumnInfo
		expected bool
	}{
		{
			name: "schemas are equal",
			schema1: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: STRING},
			},
			schema2: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: STRING},
			},
			expected: true,
		},
		{
			name: "schemas are not equal (different lengths)",
			schema1: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
			},
			schema2: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: STRING},
			},
			expected: false,
		},
		{
			name: "schemas are not equal (different types)",
			schema1: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
			},
			schema2: []ParquetColumnInfo{
				{Name: "col1", Type: DOUBLE},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EqualsParquetSchema(tt.schema1, tt.schema2)

			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestFilterOutUnsupportedParquetColumns(t *testing.T) {
	tests := []struct {
		name           string
		columns        []ParquetColumnInfo
		expectedResult []ParquetColumnInfo
	}{
		{
			name: "filter unsupported columns",
			columns: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: UNSUPPORTED},
				{Name: "col3", Type: STRING},
			},
			expectedResult: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col3", Type: STRING},
			},
		},
		{
			name: "all columns supported",
			columns: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: FLOAT},
			},
			expectedResult: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: FLOAT},
			},
		},
		{
			name: "all columns unsupported",
			columns: []ParquetColumnInfo{
				{Name: "col1", Type: UNSUPPORTED},
				{Name: "col2", Type: UNSUPPORTED},
			},
			expectedResult: []ParquetColumnInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterOutUnsupportedParquetColumns(tt.columns)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("expected: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestFilterSelectParquetColumns(t *testing.T) {
	tests := []struct {
		name           string
		columns        []ParquetColumnInfo
		expectedResult []ParquetColumnInfo
	}{
		{
			name: "filter select columns",
			columns: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: STRING},
				{Name: "col3", Type: BOOL},
				{Name: "col4", Type: FLOAT},
			},
			expectedResult: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col4", Type: FLOAT},
			},
		},
		{
			name: "all columns selectable",
			columns: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: FLOAT},
			},
			expectedResult: []ParquetColumnInfo{
				{Name: "col1", Type: INT},
				{Name: "col2", Type: FLOAT},
			},
		},
		{
			name: "no selectable columns",
			columns: []ParquetColumnInfo{
				{Name: "col1", Type: STRING},
				{Name: "col2", Type: BOOL},
			},
			expectedResult: []ParquetColumnInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterSelectParquetColumns(tt.columns)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("expected: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}
