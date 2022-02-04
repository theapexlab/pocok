package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	fmt.Printf("----- Response status : %s ----- \n", res.Status)
	fmt.Printf("%s \n ", body)
	fmt.Println("------------------------------------")

	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}
