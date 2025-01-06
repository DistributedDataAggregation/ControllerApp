package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// @Summary List available tables
// @Description Returns the names of non-empty folders in the data path
// @Tags tables
// @Produce json
// @Success 200 {array} string "List of table names"
// @Failure 500 {strinf} string "Internal server error"
// @Router /tables [get]
func handleTablesQuery(w http.ResponseWriter, r *http.Request) {
	dirs := []string{}

	err := filepath.Walk(config.DataPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			if len(entries) > 0 && path != config.DataPath {
				dirs = append(dirs, filepath.Base(path))
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dirs)
}

// @Summary Get table columns
// @Description Returns column names and their types for a given table
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

	columns, err := GetParquetSchema(filepath.Join(config.DataPath, files[0]))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting parquet schema: %v", err), http.StatusInternalServerError)
		return
	}

	columns = FilterUnsupportedParquetColumns(columns)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}
