package typless

import (
	"bytes"
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

func ExtractData(file []byte, config *Config, timeout int) (*ExtractDataFromFileOutput, error) {
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

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		utils.LogError("Failed to http.DefaultClient.Do() request", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Printf("Status:  %s  \n", res.Status)
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		bodyString := buf.String()
		fmt.Printf("Body:  %s  \n", bodyString)
		return nil, errors.New("Request to typless api errored")
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
		utils.LogError("Failed to unmarshal", err)
		return nil, unmarshalErr
	}

	return &output, nil
}
