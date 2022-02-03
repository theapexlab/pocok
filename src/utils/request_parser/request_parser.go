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
	data, decodeError := base64.StdEncoding.DecodeString(r.Body)
	result := map[string]string{}
	if decodeError != nil {
		utils.LogError("Error decoding base64 message content", decodeError)
		return result, decodeError
	}

	parsedFormData, parseQueryError := url.ParseQuery(string(data))
	if parseQueryError != nil {
		utils.LogError("Error parsing url encoded form data", parseQueryError)
		return result, parseQueryError
	}

	for key, value := range parsedFormData {
		result[key] = value[0]
	}

	return result, nil
}

// Parsing multipart/form-data body from AWS proxy request
func ParseMultipartFormData(r events.APIGatewayProxyRequest) (map[string]string, error) {
	data, decodeError := base64.StdEncoding.DecodeString(r.Body)
	result := map[string]string{}
	if decodeError != nil {
		utils.LogError("Error decoding base64 message content", decodeError)
		return result, decodeError
	}
	reader := bytes.NewReader(data)
	contentType := r.Headers["content-type"]
	_, params, parseMediaTypeError := mime.ParseMediaType(contentType)
	if parseMediaTypeError != nil {
		utils.LogError("Error during parse form media type", parseMediaTypeError)
		return result, parseMediaTypeError
	}
	mr := multipart.NewReader(reader, params["boundary"])
	for {
		part, readError := mr.NextPart()
		if readError == io.EOF {
			break
		}
		if readError != nil {
			return result, readError
		}
		slurp, readAllError := io.ReadAll(part)
		if readAllError != nil {
			return result, readAllError
		}
		key := part.FormName()
		value := string(slurp)
		result[key] = value
	}
	return result, nil
}
