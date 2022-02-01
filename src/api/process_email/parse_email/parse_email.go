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

func ParseEmail(body string) (*models.UploadInvoiceMessage, error) {
	var jsonBody models.EmailWebhookBody

	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return nil, models.ErrInvalidJson
	}

	if isSentFromRecipientAddress(&jsonBody) {
		return nil, ErrEmailFromSenderAddress
	}

	if hasPdfAttachment(jsonBody.Attachments) {
		return &models.UploadInvoiceMessage{
			Type:     "base64",
			Body:     jsonBody.Attachments[0].Content_b64,
			Filename: jsonBody.Attachments[0].Filename,
		}, nil
	}

	if ok, url := hasPdfUrl(&jsonBody); ok {
		return &models.UploadInvoiceMessage{
			Type: "url",
			Body: url,
		}, nil
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
	url, err := url_parsing_strategies.GetPdfUrlFromEmail(jsonBody)
	if !errors.Is(err, url_parsing_strategies.ErrNoUrlParsingStrategyFound) {
		utils.LogError("error while parsing url from email", err)
	}
	return err == nil && url != "", url
}

func isSentFromRecipientAddress(jsonBody *models.EmailWebhookBody) bool {
	return jsonBody.From[0].Address == os.Getenv("mailgunSender")
}
