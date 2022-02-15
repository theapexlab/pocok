package slack_service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"pocok/src/utils"
	"strings"
)

type SlackClient struct {
	Url      string
	Username string
	Channel  string
}

func (slackClient *SlackClient) SendMessage(message string) (http.Response, error) {
	payload := map[string]string{
		"channel":    slackClient.Channel,
		"username":   slackClient.Username,
		"icon_emoji": ":pocok-logo:",
		"text":       message,
	}
	payloadBytes, _ := json.Marshal(payload)
	data := url.Values{"payload": []string{string(payloadBytes)}}

	request, _ := http.NewRequest(http.MethodPost, slackClient.Url, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		utils.LogError("error while sending notification to slack", err)
	}
	return *resp, err
}
