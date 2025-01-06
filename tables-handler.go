package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// handleTablesQuery lists all non-empty folders in config.DataPath
// @Summary List available tables
// @Description Returns the names of non-empty folders in the data path
// @Tags tables
// @Produce json
// @Success 200 {array} string "List of table names"
// @Failure 500 {object} map[string]string "Internal server error"
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
