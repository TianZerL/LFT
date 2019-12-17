package client

import (
	"os"
	"path/filepath"
)

func checkDir(src string) (bool, error) {
	f, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

//Get all files and dirs path in to []string
func getDirInfo(src string) (fileInfo []string, dirInfo []string) {
	filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				dirInfo = append(dirInfo, path)
			} else {
				fileInfo = append(fileInfo, path)
			}
			return nil
		})
	return
}
