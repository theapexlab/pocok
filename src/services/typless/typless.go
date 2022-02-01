package typless

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	req, err := http.NewRequest("POST", url, payloadStrReader)

	if err != nil {
		utils.LogError("Failed to http.NewRequest()", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+config.Token)

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		utils.LogError("Failed to http.Client.Do() request", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		utils.LogError("Failed to ioutil.ReadAll() request", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Printf("Status:  %s  \n", res.Status)
		fmt.Printf("Body:  %s  \n", string(body))
		err := errors.New("❌ request to typless failed")
		return nil, err
	}

	output := ExtractDataFromFileOutput{}

	err = json.Unmarshal(body, &output)
	if err != nil {
		utils.LogError("Failed to unmarshal", err)
		return nil, err
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

	req, err := http.NewRequest("POST", url, payloadStrReader)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+config.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		fmt.Printf("Status:  %s  \n", res.Status)
		fmt.Printf("Body:  %s  \n", string(body))

		return errors.New("❌ request to typless failed")
	}

	defer res.Body.Close()

	return nil
}
