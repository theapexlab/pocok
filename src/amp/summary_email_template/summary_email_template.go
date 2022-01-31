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

	file, fileErr := os.ReadFile(filePath)
	if fileErr != nil {
		utils.LogError("Error while reading in the html file.", fileErr)
		return "", fileErr
	}
	return string(file), nil
}
