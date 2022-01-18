package typless

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func ExtractData(file []byte, documentTypeName string, token string) (*ExtractDataFromFileOutput, error) {
	url := "https://developers.typless.com/api/extract-data"

	payload := ExtractDataFromFileInput{
		DocumentTypeName: documentTypeName,
		FileName:         "invoice.pdf",
		File:             base64.StdEncoding.EncodeToString(file),
	}
	payloadStr, _ := json.Marshal(payload)
	payloadStrReader := strings.NewReader(string(payloadStr))

	req, _ := http.NewRequest("POST", url, payloadStrReader)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	output := ExtractDataFromFileOutput{}

	unmarshalErr := json.Unmarshal(body, &output)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &output, nil
}
