package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/apache/arrow/go/parquet/file"
	"github.com/apache/arrow/go/parquet/schema"
)

type ParquetColumnInfo struct {
	Name string            `json:"name"` // The name of the column.
	Type ParquetColumnType `json:"type"` // The type of the column
}

type ParquetColumnType string

const (
	INT         ParquetColumnType = "INT"
	DOUBLE      ParquetColumnType = "DOUBLE"
	FLOAT       ParquetColumnType = "FLOAT"
	STRING      ParquetColumnType = "STRING"
	BOOL        ParquetColumnType = "BOOL"
	UNSUPPORTED ParquetColumnType = "UNSUPPORTED"
)

func getParquetColumnType(colType string) string {

	allowedTypes := []string{string(INT), string(DOUBLE), string(FLOAT), string(STRING), string(BOOL)}

	for _, allowedType := range allowedTypes {
		if strings.Contains(colType, allowedType) {
			return allowedType
		}
	}

	return string(UNSUPPORTED)
}

func GetParquetSchemaByPath(filePath string) ([]ParquetColumnInfo, error) {

	parquetFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer parquetFile.Close()

	reader, err := file.NewParquetReader(parquetFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create ParquetReader: %v", err)
	}

	return getParquetSchema(reader), nil
}

func GetParquetSchemaByMultipartFile(parquetFile multipart.File) ([]ParquetColumnInfo, error) {

	currentPos, err := parquetFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get current file position")
	}
	defer parquetFile.Seek(currentPos, io.SeekStart)

	reader, err := file.NewParquetReader(parquetFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create ParquetReader: %v", err)
	}

	return getParquetSchema(reader), nil
}

func getParquetSchema(reader *file.Reader) []ParquetColumnInfo {

	schema := reader.MetaData().Schema
	var columns []ParquetColumnInfo

	for i := 0; i < schema.NumColumns(); i++ {
		col := schema.Column(i)
		columnName := col.Name()
		columnType := getColumnType(col)
		columns = append(columns, ParquetColumnInfo{Name: columnName, Type: ParquetColumnType(columnType)})
	}

	return columns
}

func getColumnType(col *schema.Column) string {

	columnLogicalType := col.LogicalType()
	if columnLogicalType != nil && !columnLogicalType.Equals(schema.NoLogicalType{}) && !columnLogicalType.Equals(schema.IntLogicalType{}) {
		return normalizeType(columnLogicalType.String())
	}

	return normalizeType(col.PhysicalType().String())
}

func normalizeType(colType string) string {

	colType = strings.Split(colType, "(")[0]
	colType = strings.ToUpper(colType)
	return getParquetColumnType(colType)

}

func EqualsParquetSchema(this []ParquetColumnInfo, that []ParquetColumnInfo) bool {

	if len(this) != len(that) {
		return false
	}

	for i := 0; i < len(this); i++ {
		if !EqualsParquetColumnInfo(this[i], that[i]) {
			return false
		}
	}

	return true
}

func EqualsParquetColumnInfo(this ParquetColumnInfo, that ParquetColumnInfo) bool {
	if this.Name == that.Name && this.Type == that.Type {
		return true
	}
	return false
}

func FilterOutUnsupportedParquetColumns(columns []ParquetColumnInfo) []ParquetColumnInfo {
	var filteredColumns []ParquetColumnInfo

	for _, col := range columns {
		if col.Type != UNSUPPORTED {
			filteredColumns = append(filteredColumns, col)
		}
	}

	if filteredColumns == nil {
		return []ParquetColumnInfo{}
	}
	return filteredColumns
}

func FilterSelectParquetColumns(columns []ParquetColumnInfo) []ParquetColumnInfo {
	var filteredColumns []ParquetColumnInfo

	for _, col := range columns {
		if col.Type != STRING && col.Type != BOOL && col.Type != UNSUPPORTED {
			filteredColumns = append(filteredColumns, col)
		}
	}

	if filteredColumns == nil {
		return []ParquetColumnInfo{}
	}
	return filteredColumns
}
