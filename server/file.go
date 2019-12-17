package server

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

//Set file name, if there is a same name file here, rename new file ot add a (1)
func setFileName(path string) string {
	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		log.Fatalln(err)
	}
	for isFileExist(path) {
		suffix := filepath.Ext(path)
		filename := strings.TrimSuffix(path, suffix)
		filename += "(1)"
		path = filename + suffix
	}
	return path
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalln(err)
	}
	return true
}
