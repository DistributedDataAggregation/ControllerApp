package main

import (
	"fmt"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/reader"
)

type ParquetColumnInfo struct {
	Name string
	Type string
}

func GetParquetSchema(filePath string) ([]ParquetColumnInfo, error) {

	file, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	parquetReader, err := reader.NewParquetReader(file, nil, 1)
	if err != nil {
		return nil, fmt.Errorf("cannot create ParquetReader: %v", err)
	}
	defer parquetReader.ReadStop()

	schemaHandler := parquetReader.SchemaHandler
	var columns []ParquetColumnInfo

	for _, col := range schemaHandler.SchemaElements {

		if col.GetNumChildren() != 0 {
			continue
		}

		columnName := col.GetName()
		columnType := GetColumnType(col)
		columns = append(columns, ParquetColumnInfo{Name: columnName, Type: columnType})
	}

	return columns, nil
}

func GetColumnType(col *parquet.SchemaElement) string {

	columnLogicalType := col.GetLogicalType()

	if columnLogicalType != nil {

		if columnLogicalType.IsSetINTEGER() {
			return "INTEGER"
		}

		if columnLogicalType.IsSetSTRING() {
			return "STRING"
		}

		if columnLogicalType.IsSetDECIMAL() {
			return "DECIMAL"
		}

		return "UNSUPPORTED"

	}

	if col.Type != nil {
		return col.Type.String()
	}

	return "UNSUPPORTED"
}

func FilterUnsupportedParquetColumns(columns []ParquetColumnInfo) []ParquetColumnInfo {
	var filteredColumns []ParquetColumnInfo
	for _, col := range columns {
		if col.Type != "UNSUPPORTED" {
			filteredColumns = append(filteredColumns, col)
		}
	}
	return filteredColumns
}
