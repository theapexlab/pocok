package typless

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"pocok/src/utils"
	"strings"
	"time"
)

func ExtractData(config *Config, file []byte, timeout int) (*ExtractDataFromFileOutput, error) {
	url := "https://developers.typless.com/api/extract-data?training=true"

	payload := ExtractDataFromFileInput{
		DocumentTypeName: config.DocType,
		FileName:         "invoice.pdf",
		File:             base64.StdEncoding.EncodeToString(file),
	}
	payloadStr, _ := json.Marshal(payload)
	payloadStrReader := strings.NewReader(string(payloadStr))

	req, requestError := http.NewRequest("POST", url, payloadStrReader)
	if requestError != nil {
		utils.LogError("Failed to http.NewRequest()", requestError)
		return nil, requestError
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+config.Token)

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	res, doRequestError := client.Do(req)
	if doRequestError != nil {
		utils.LogError("Failed to http.Client.Do() request", doRequestError)
		return nil, doRequestError
	}

	defer res.Body.Close()
	body, readAllError := ioutil.ReadAll(res.Body)
	if readAllError != nil {
		utils.LogError("Failed to ioutil.ReadAll() request", readAllError)
		return nil, readAllError
	}

	if res.StatusCode != 200 {
		utils.Logf("Status:  %s  \n", res.Status)
		utils.Logf("Body:  %s  \n", string(body))
		return nil, errors.New("❌ request to typless failed")
	}

	output := ExtractDataFromFileOutput{}

	unmarshalError := json.Unmarshal(body, &output)
	if unmarshalError != nil {
		utils.LogError("Failed to unmarshal", unmarshalError)
		return nil, unmarshalError
	}

	return &output, nil
}

func AddDocumentFeedback(config *Config, trainingData TrainingData) error {
	url := "https://developers.typless.com/api/add-document-feedback"

	payload := AddDocumentFeedbackInput{
		DocumentTypeName: config.DocType,
		DocumentObjectId: trainingData.DocumentObjectId,
		LearningFields:   trainingData.LearningFields,
		LineItems:        trainingData.LineItems,
	}
	payloadStr, _ := json.Marshal(payload)
	payloadStrReader := strings.NewReader(string(payloadStr))

	req, requestError := http.NewRequest("POST", url, payloadStrReader)
	if requestError != nil {
		return requestError
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+config.Token)

	res, doRequestError := http.DefaultClient.Do(req)
	if doRequestError != nil {
		return doRequestError
	}

	if res.StatusCode >= 300 {
		body, readAllError := ioutil.ReadAll(res.Body)
		if readAllError != nil {
			return readAllError
		}
		utils.Logf("Status:  %s  \n", res.Status)
		utils.Logf("Body:  %s  \n", string(body))

		return errors.New("❌ request to typless failed")
	}

	defer res.Body.Close()

	return nil
}
