package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pocok/src/services/wise/api/models"
	"pocok/src/utils"
)

func (wise *WiseClient) CreateTransfer(body models.Transfer) (*models.Transfer, error) {
	url := fmt.Sprintf("%s/%s/transfers", wise.baseUrl, V1)

	jsonBody, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		utils.LogError("Error marshalling CreateTransfer body", jsonErr)
		return nil, jsonErr
	}

	req, newReqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if newReqErr != nil {
		utils.LogError("Error creating CreateTransfer request", newReqErr)
		return nil, newReqErr
	}

	res := models.Transfer{}
	if reqErr := wise.sendRequest(req, &res); reqErr != nil {
		utils.LogError("Error sending CreateTransfer request", reqErr)
		return nil, reqErr
	}

	return &res, nil
}
