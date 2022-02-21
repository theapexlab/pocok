package read_raw_email

import (
	"os"
	"path"
	"runtime"
)

func Read(filename string) string {
	_, currentFilename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(currentFilename)
	filePath := currentPath + "/../" + filename
	file, readFileError := os.ReadFile(filePath)

	if readFileError != nil {
		panic(readFileError)
	}

	return string(file)
}
