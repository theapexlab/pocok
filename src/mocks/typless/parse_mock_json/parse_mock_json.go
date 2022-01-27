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
	mock, readFileErr := os.ReadFile(filePath)
	var extractedData *typless.ExtractDataFromFileOutput
	// currentPath := path.Dir(filePath)
	// fmt.Println(currentPath)
	// mock, readFileErr := ioutil.ReadFile(filePath)
	if readFileErr != nil {
		utils.LogError("", readFileErr)
		panic("Failed to read mock file")
	}

	if err := json.Unmarshal(mock, &extractedData); err != nil {
		panic("Failed to unmarshal mock file")
	}

	return extractedData
}
