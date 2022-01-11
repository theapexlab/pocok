package models

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url", "base64"
	Body string `json:"body"`
}
