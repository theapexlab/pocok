package typless

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pocok/src/utils"
	"strings"
)

func ExtractData(file []byte, config *Config) (*ExtractDataFromFileOutput, error) {
	url := "https://developers.typless.com/api/extract-data"

	payload := ExtractDataFromFileInput{
		DocumentTypeName: config.DocType,
		FileName:         "invoice.pdf",
		File:             base64.StdEncoding.EncodeToString(file),
	}
	payloadStr, _ := json.Marshal(payload)
	payloadStrReader := strings.NewReader(string(payloadStr))

	req, err := http.NewRequest("POST", url, payloadStrReader)

	if err != nil {
		utils.LogError("Failed to http.NewRequest()", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+config.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogError("Failed to http.DefaultClient.Do() request", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		utils.LogError("Failed to ioutil.ReadAll() request", err)
		return nil, err
	}

	output := ExtractDataFromFileOutput{}

	unmarshalErr := json.Unmarshal(body, &output)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &output, nil
}
