package get_pdf_url

import (
	"encoding/json"
	"errors"
)

type Attachment struct {
	Encoding string `json:"encoding"`
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
	Url      string `json:"url"`
}
type webhookBody struct {
	Attachment *Attachment `json:"attachment-1"`
}

var ErrInvalidJson = errors.New("invalid body")
var ErrNoPdfAttachmentFound = errors.New("no pdf attachment found")

func GetPdfUrl(body string) (string, error) {
	var jsonBody webhookBody

	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return "", ErrInvalidJson
	}

	if jsonBody.Attachment == nil || jsonBody.Attachment.Mimetype != "application/pdf" {
		return "", ErrNoPdfAttachmentFound
	}

	return jsonBody.Attachment.Url, nil
}
