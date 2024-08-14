package fs

import (
	"fmt"
	"os"
)

func FilesInCurrentDir() ([]string, error) {
	var files []string
	wd, err := os.Getwd()
	if err != nil {
		return files, fmt.Errorf("error to get now directory: %w", err)
	}
	return FilesInDir(wd)
}

func FilesInDir(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files, fmt.Errorf("error to get entries in directory: %w", err)
	}
	for _, e := range entries {
		if info, err := os.Stat(e.Name()); err == nil {
			if !info.IsDir() {
				files = append(files, e.Name())
			}
		}
	}
	return files, nil
}
