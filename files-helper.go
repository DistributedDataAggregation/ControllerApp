package main

import (
	"io/fs"
	"path/filepath"
)

func findDataFiles(tableName string) ([]string, error) {
	targetPath := filepath.Join(config.DataPath, tableName)
	files := []string{}

	err := filepath.Walk(targetPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".parquet" {
			relPath, _ := filepath.Rel(config.DataPath, path)
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
