package util

import (
	"os"
	"path/filepath"
)

func MkdirForFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	return nil
}
