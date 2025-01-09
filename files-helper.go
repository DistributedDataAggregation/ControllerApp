package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func findDataDirs() ([]string, error) {
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

	return dirs, err
}

func findDataFiles(tableName string) ([]string, error) {
	targetPath := filepath.Join(config.DataPath, tableName)
	files := []string{}

	err := filepath.Walk(targetPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".parquet" {
			//relPath, _ := filepath.Rel(config.DataPath, path)
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
