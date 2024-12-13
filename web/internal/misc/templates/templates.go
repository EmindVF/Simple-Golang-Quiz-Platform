package templates

import (
	"os"
	"path/filepath"
)

func LoadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				LoadTemplates(path)
			}
		} else {
			files = append(files, path)
		}
		return err
	})

	return files, err
}
