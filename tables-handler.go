package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/beevik/guid"
)

// @Summary List available tables
// @Description Returns the names of non-empty folders in the data path
// @Tags tables
// @Produce json
// @Success 200 {array} string "List of table names"
// @Failure 500 {string} string "Internal server error"
// @Router /tables [get]
func handleTablesQuery(w http.ResponseWriter, r *http.Request) {

	dirs, err := findDataDirs()

	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dirs)
}

// @Summary Get table columns
// @Description Returns column names and their types for a given table. Filters out columns of unsupported types.
// @Tags tables
// @Produce json
// @Param name query string true "Table name"
// @Success 200 {array} []ParquetColumnInfo "List of columns with their types"
// @Failure 400 {string} string "Bad request with error message"
// @Failure 500 {string} string "Internal server error"
// @Router /tables/columns [get]
func handleTablesColumnsQuery(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("name")
	if tableName == "" {
		http.Error(w, "Missing table name", http.StatusBadRequest)
		return
	}

	files, err := findDataFiles(tableName)
	if err != nil || len(files) == 0 {
		http.Error(w, fmt.Sprintf("No parquet files found for table: %s", tableName), http.StatusBadRequest)
		return
	}

	columns, err := GetParquetSchemaByPath(files[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting parquet schema: %v", err), http.StatusInternalServerError)
		return
	}

	columns = FilterOutUnsupportedParquetColumns(columns)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}

// @Summary Get table columns that can be aggregated
// @Description Returns column names and their types for a given table. Filters out columns of types unsupported for aggregations.
// @Tags tables
// @Produce json
// @Param name query string true "Table name"
// @Success 200 {array} []ParquetColumnInfo "List of columns with their types"
// @Failure 400 {string} string "Bad request with error message"
// @Failure 500 {string} string "Internal server error"
// @Router /tables/select-columns [get]
func handleTablesSelectColumnsQuery(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("name")
	if tableName == "" {
		http.Error(w, "Missing table name", http.StatusBadRequest)
		return
	}

	files, err := findDataFiles(tableName)
	if err != nil || len(files) == 0 {
		http.Error(w, fmt.Sprintf("No parquet files found for table: %s", tableName), http.StatusBadRequest)
		return
	}

	columns, err := GetParquetSchemaByPath(files[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting parquet schema: %v", err), http.StatusInternalServerError)
		return
	}

	columns = FilterSelectParquetColumns(columns)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}

// @Summary Upload file to a table with given name
// @Description Uploads a Parquet file (max 10 MB) to a table with a given name. If the table does not exist, it is created.
// @Tags tables
// @Accept mpfd
// @Produce json
// @Param name query string true "Table name"
// @Param file formData file true "File to upload (must have .parquet)"
// @Success 200 {string} string "File uploaded successfully"
// @Failure 400 {string} string "Invalid table name or file"
// @Failure 500 {string} string "Internal server error"
// @Router /tables/upload [post]
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("name")
	if tableName == "" {
		http.Error(w, "Missing table name", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file or invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(fileHeader.Filename) != ".parquet" {
		http.Error(w, "File must have .parquet extension", http.StatusBadRequest)
		return
	}

	tablePath := filepath.Join(config.DataPath, tableName)
	if err := os.MkdirAll(tablePath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	if code, err := validateFileSchema(tableName, file); err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	tempPath := filepath.Join(tablePath, "_"+guid.NewString()+"_"+fileHeader.Filename)
	destPath := filepath.Join(tablePath, guid.NewString()+"_"+fileHeader.Filename)
	destFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	os.Rename(tempPath, destPath)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully")
}

func validateFileSchema(tableName string, file multipart.File) (int, error) {

	files, err := findDataFiles(tableName)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to retreive existing table files")
	}
	if len(files) > 0 {
		schema, err := GetParquetSchemaByPath(files[0])
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to get existing table schema")
		}
		newFileSchema, err := GetParquetSchemaByMultipartFile(file)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to file schema")
		}
		if !EqualsParquetSchema(schema, newFileSchema) {
			return http.StatusBadRequest, fmt.Errorf("file schema differs from existing table schema")
		}
	}

	return 0, nil
}
