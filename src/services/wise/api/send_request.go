package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (wise *WiseClient) sendRequest(req *http.Request, data interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", wise.apiToken))

	res, err := wise.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(res.Status, "2") {
		return errors.New("http wise request error: " + res.Status + " " + string(body))
	}

	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}
