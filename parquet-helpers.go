package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/apache/arrow/go/parquet/file"
	"github.com/apache/arrow/go/parquet/schema"
)

type ParquetColumnInfo struct {
	Name string
	Type string
}

func GetParquetSchema(filePath string) ([]ParquetColumnInfo, error) {

	parquetFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer parquetFile.Close()

	reader, err := file.NewParquetReader(parquetFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create ParquetReader: %v", err)
	}

	schema := reader.MetaData().Schema
	var columns []ParquetColumnInfo

	for i := 0; i < schema.NumColumns(); i++ {
		col := schema.Column(i)
		columnName := col.Name()
		columnType := GetColumnType(col)
		columns = append(columns, ParquetColumnInfo{Name: columnName, Type: columnType})
	}

	return columns, nil
}

func GetColumnType(col *schema.Column) string {

	columnLogicalType := col.LogicalType()
	if columnLogicalType != nil && !columnLogicalType.Equals(schema.NoLogicalType{}) && !columnLogicalType.Equals(schema.IntLogicalType{}) {
		return columnLogicalType.String()
	}

	return col.PhysicalType().String()
}

func FilterUnsupportedParquetColumns(columns []ParquetColumnInfo) []ParquetColumnInfo {
	var filteredColumns []ParquetColumnInfo
	allowedTypes := []string{"INT", "DOUBLE", "FLOAT", "DECIMAL", "BOOL", "STRING"}

	for _, col := range columns {
		normalizedType := normalizeType(col.Type)
		for _, allowedType := range allowedTypes {
			if strings.Contains(normalizedType, allowedType) {
				filteredColumns = append(filteredColumns, ParquetColumnInfo{
					Name: col.Name,
					Type: normalizedType,
				})
				break
			}
		}
	}

	return filteredColumns
}

func normalizeType(colType string) string {
	trimmedType := strings.Split(colType, "(")[0]
	return strings.ToUpper(trimmedType)
}
