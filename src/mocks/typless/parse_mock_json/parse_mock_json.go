package parse_mock_json

import (
	"encoding/json"
	"os"
	"path"
	"pocok/src/services/typless"
	"pocok/src/utils"
	"runtime"
)

func Parse(mockFilename string) *typless.ExtractDataFromFileOutput {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	filePath := currentPath + "/../" + mockFilename
	mock, readFileError := os.ReadFile(filePath)

	var extractedData *typless.ExtractDataFromFileOutput

	if readFileError != nil {
		utils.LogError("", readFileError)
		panic("Failed to read mock file")
	}

	if unmarshalError := json.Unmarshal(mock, &extractedData); unmarshalError != nil {
		panic("Failed to unmarshal mock file")
	}

	return extractedData
}
