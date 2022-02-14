package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pocok/src/services/wise/api/models"
	"pocok/src/utils"
)

func (wise *WiseClient) GetRecipientAccounts() (*[]models.RecipientAccount, error) {
	url := fmt.Sprintf("%s/%s/accounts", wise.baseUrl, V2)

	req, newReqErr := http.NewRequest("GET", url, nil)
	if newReqErr != nil {
		utils.LogError("Error creating GetRecipientAccounts request", newReqErr)
		return nil, newReqErr
	}

	res := models.RecipientAccountList{}
	if reqErr := wise.sendRequest(req, &res); reqErr != nil {
		utils.LogError("Error sending GetRecipientAccounts request", reqErr)
		return nil, reqErr
	}

	return &res.Content, nil
}

func (wise *WiseClient) GetRecipientAccountById(id int) (*models.RecipientAccount, error) {
	url := fmt.Sprintf("%s/%s/accounts/%d", wise.baseUrl, V2, id)

	req, newReqErr := http.NewRequest("GET", url, nil)
	if newReqErr != nil {
		utils.LogError("Error creating GetRecipientAccountById request", newReqErr)
		return nil, newReqErr
	}

	res := models.RecipientAccount{}
	if reqErr := wise.sendRequest(req, &res); reqErr != nil {
		utils.LogError("Error sending GetRecipientAccountById request", reqErr)
		return nil, reqErr
	}

	return &res, nil
}

func (wise *WiseClient) CreateRecipientAccount(body models.RecipientAccountV1) (*models.RecipientAccountV1, error) {
	url := fmt.Sprintf("%s/%s/accounts", wise.baseUrl, V1)

	jsonBody, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		utils.LogError("Error marshalling CreateRecipientAccount body", jsonErr)
		return nil, jsonErr
	}

	req, newReqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if newReqErr != nil {
		utils.LogError("Error creating CreateRecipientAccount request", newReqErr)
		return nil, newReqErr
	}

	res := models.RecipientAccountV1{}
	if reqErr := wise.sendRequest(req, &res); reqErr != nil {
		utils.LogError("Error sending CreateRecipientAccount request", reqErr)
		return nil, reqErr
	}

	return &res, nil
}
