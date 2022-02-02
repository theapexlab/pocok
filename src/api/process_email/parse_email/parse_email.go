package parse_email

import (
	"encoding/json"
	"errors"
	"os"
	"pocok/src/api/process_email/url_parsing_strategies"
	"pocok/src/utils"
	"pocok/src/utils/models"
)

var ErrNoPdfAttachmentFound = errors.New("no pdf attachment found")
var ErrEmailFromSenderAddress = errors.New("email sent from mailgun sender address")

func ParseEmail(body string) ([]models.UploadInvoiceMessage, error) {
	var jsonBody models.EmailWebhookBody

	if unmarshalError := json.Unmarshal([]byte(body), &jsonBody); unmarshalError != nil {
		return nil, models.ErrInvalidJson
	}

	if isSentFromRecipientAddress(&jsonBody) {
		return nil, ErrEmailFromSenderAddress
	}

	if hasPdfAttachment(jsonBody.Attachments) {
		// Map, so only unique values go throuh
		messageMap := map[string]models.UploadInvoiceMessage{}
		for _, attachment := range jsonBody.Attachments {
			messageMap[attachment.Filename+attachment.Content_b64] = models.UploadInvoiceMessage{
				Type:     "base64",
				Body:     attachment.Content_b64,
				Filename: attachment.Filename,
			}
		}

		messages := []models.UploadInvoiceMessage{}
		for _, message := range messageMap {
			messages = append(messages, message)
		}

		return messages, nil
	}

	if ok, url := hasPdfUrl(&jsonBody); ok {
		messages := []models.UploadInvoiceMessage{{
			Type: "url",
			Body: url,
		}}
		return messages, nil
	}

	return nil, ErrNoPdfAttachmentFound
}

func hasPdfAttachment(attachments []*models.EmailAttachment) bool {
	for _, attachment := range attachments {
		if attachment.ContentType == "application/pdf" && attachment.Content_b64 != "" {
			return true
		}
	}

	return false
}

func hasPdfUrl(jsonBody *models.EmailWebhookBody) (bool, string) {
	url, parseError := url_parsing_strategies.GetPdfUrlFromEmail(jsonBody)
	if !errors.Is(parseError, url_parsing_strategies.ErrNoUrlParsingStrategyFound) {
		utils.LogError("error while parsing url from email", parseError)
	}
	return parseError == nil && url != "", url
}

func isSentFromRecipientAddress(jsonBody *models.EmailWebhookBody) bool {
	return len(jsonBody.From) > 0 && jsonBody.From[0].Address == os.Getenv("mailgunSender")
}
