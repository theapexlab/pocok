package form_data

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 10 // 10MB

func convertRequest(request events.APIGatewayProxyRequest) (http.Request, error) {
	r := http.Request{}

	r.Header = make(map[string][]string)

	for k, v := range request.Headers {
		if strings.ToLower(k) == "content-type" {
			r.Header.Set("Content-Type", v)
		}
	}

	body, err := base64.StdEncoding.DecodeString(request.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	return r, nil
}

func getFileFromRequest(r *http.Request) ([]byte, error) {
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		panic(err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseFormData(request events.APIGatewayProxyRequest) ([]byte, error) {
	httpRequest, err := convertRequest(request)
	if err != nil {
		return nil, err
	}

	file, _ := getFileFromRequest(&httpRequest)
	if err != nil {
		return nil, err
	}

	return file, nil
}
