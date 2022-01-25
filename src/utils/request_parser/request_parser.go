package request_parser

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"pocok/src/utils"

	"github.com/aws/aws-lambda-go/events"
)

// Parsing x-www-form-urlencoded body from AWS proxy request
func ParseUrlEncodedFormData(r events.APIGatewayProxyRequest) (map[string]string, error) {
	data, err := base64.StdEncoding.DecodeString(r.Body)
	result := map[string]string{}
	if err != nil {
		utils.LogError("Error decoding base64 message content", err)
		return result, err
	}

	parsedFormData, parseErr := url.ParseQuery(string(data))
	if parseErr != nil {
		utils.LogError("Error parsing url encoded form data", parseErr)
		return result, parseErr
	}

	for key, value := range parsedFormData {
		result[key] = value[0]
	}

	return result, nil
}

// Parsing multipart/form-data body from AWS proxy request
func ParseMultipartFormData(r events.APIGatewayProxyRequest) (map[string]string, error) {
	data, err := base64.StdEncoding.DecodeString(r.Body)
	result := map[string]string{}
	if err != nil {
		utils.LogError("Error decoding base64 message content", err)
		return result, err
	}
	reader := bytes.NewReader(data)
	contentType := r.Headers["content-type"]
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		utils.LogError("Error during parse form media type", err)
		return result, err
	}
	mr := multipart.NewReader(reader, params["boundary"])
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
		slurp, err := io.ReadAll(part)
		if err != nil {
			return result, err
		}
		key := part.FormName()
		value := string(slurp)
		result[key] = value
	}
	return result, nil
}
