package api

import (
	"fmt"
	"net/http"
	"pocok/src/services/wise/api/models"
	"pocok/src/utils"
)

func (wise *WiseClient) GetProfiles() (*[]models.Profile, error) {
	url := fmt.Sprintf("%s/%s/profiles", wise.baseUrl, V1)

	req, newReqErr := http.NewRequest("GET", url, nil)
	if newReqErr != nil {
		utils.LogError("Error creating GetProfiles request", newReqErr)
		return nil, newReqErr
	}

	res := []models.Profile{}
	if reqErr := wise.sendRequest(req, &res); reqErr != nil {
		utils.LogError("Error sending GetProfiles request", reqErr)
		return nil, reqErr
	}

	return &res, nil
}
