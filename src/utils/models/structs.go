package models

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url", "base64"
	Body string `json:"body"`
}

type EmailAttachment struct {
	ContentType string `json:"contentType"`
	Content_b64 string `json:"content_b64"`
}

type EmailFrom struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
type EmailWebhookBody struct {
	Attachments []*EmailAttachment `json:"attachments"`
	Html        string             `json:"html"`
	From        []*EmailFrom       `json:"from"`
}
