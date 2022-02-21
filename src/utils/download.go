package utils

import (
	"io/ioutil"
	"net/http"
)

func DownloadFile(url string) ([]byte, error) {
	resp, httpGetError := http.Get(url)
	if httpGetError != nil {
		LogError("", httpGetError)
		return nil, httpGetError
	}

	defer resp.Body.Close()

	body, readAllErr := ioutil.ReadAll(resp.Body)
	if readAllErr != nil {
		LogError("", readAllErr)
		return nil, readAllErr
	}

	return body, nil
}
