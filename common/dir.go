package common

import (
	"os"
	"path/filepath"
)

// ClearDir 清空指定目录中的所有文件
func ClearDir(dirname string) error {
	dir, err := os.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if file.IsDir() {
			if err := os.RemoveAll(filepath.Join(dirname, file.Name())); err != nil {
				return err
			}
		} else {
			if err := os.Remove(filepath.Join(dirname, file.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
