package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pocok/src/services/wise/api/models"
	"pocok/src/utils"
)

func (wise *WiseClient) CreateQuote(body models.Quote) (*models.Quote, error) {
	url := fmt.Sprintf("%s/%s/quotes", wise.baseUrl, v2)

	jsonBody, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		utils.LogError("Error marshalling CreateQuote body", jsonErr)
		return nil, jsonErr
	}

	req, newReqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if newReqErr != nil {
		utils.LogError("Error creating CreateQuote request", newReqErr)
		return nil, newReqErr
	}

	res := models.Quote{}
	if reqErr := wise.sendRequest(req, &res); reqErr != nil {
		utils.LogError("Error sending CreateQuote request", reqErr)
		return nil, reqErr
	}

	return &res, nil
}
