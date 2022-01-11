package parse_email

import (
	"encoding/json"
	"errors"
	"pocok/src/utils/models"
)

type attachment struct {
	ContentType string `json:"contentType"`
	Content_b64 string `json:"content_b64"`
}

type from struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type webhookBody struct {
	Attachments []*attachment `json:"attachments"`
	Text        string        `json:"text"`
	From        []*from       `json:"from"`
}

var ErrNoPdfAttachmentFound = errors.New("no pdf attachment found")

func hasPdfAttachment(attachments []*attachment) bool {
	for _, attachment := range attachments {
		if attachment.ContentType == "application/pdf" && attachment.Content_b64 != "" {
			return true
		}
	}

	return false
}

func ParseEmail(body string) (*models.UploadInvoiceMessage, error) {
	var jsonBody webhookBody

	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return nil, models.ErrInvalidJson
	}

	if !hasPdfAttachment(jsonBody.Attachments) {
		return nil, ErrNoPdfAttachmentFound
	}

	return &models.UploadInvoiceMessage{
		Type: "base64",
		Body: jsonBody.Attachments[0].Content_b64,
	}, nil
}
