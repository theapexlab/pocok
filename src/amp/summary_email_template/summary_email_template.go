package summary_email_template

import (
	"os"
	"path"
	"pocok/src/utils"
	"runtime"
)

func Get() (string, error) {
	filePath := "src/amp/templates/summary_email.html"
	if os.Getenv("stage") != "production" {
		_, filename, _, _ := runtime.Caller(0)
		currentPath := path.Dir(filename)
		filePath = currentPath + "/../../../" + filePath
	}

	file, fileError := os.ReadFile(filePath)
	if fileError != nil {
		utils.LogError("Error while reading in the html file.", fileError)
		return "", fileError
	}
	return string(file), nil
}
