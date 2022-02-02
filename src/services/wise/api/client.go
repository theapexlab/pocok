package api

import (
	"net/http"
	"time"
)

type WiseClient struct {
	baseUrl    string
	apiToken   string
	httpClient *http.Client
}

func CreateWiseClient(apiToken string) *WiseClient {
	return &WiseClient{
		baseUrl:  base_url,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
