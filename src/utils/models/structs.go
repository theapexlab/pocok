package models

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url","base64"
	Body string `json:"body"`
}

type Service struct {
	Name       string `json:"name" dynamodbav:"name,omitempty"`
	Amount     string `json:"amount" dynamodbav:"amount,omitempty"`
	NetPrice   string `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice string `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Currency   string `json:"currency" dynamodbav:"currency,omitempty"`
	Tax        string `json:"tax" dynamodbav:"tax,omitempty"`
}

type Invoice struct {
	Id           string `json:"id" dynamodbav:"id"`
	Filename     string `json:"filename" dynamodbav:"filename"`
	TextractData string `json:"textractData" dynamodbav:"textractData,omitempty"`

	InvoiceNumber string    `json:"invoiceNumber" dynamodbav:"invoiceNumber,omitempty"`
	CustomerName  string    `json:"customerName" dynamodbav:"customerName,omitempty"`
	AccountNumber string    `json:"accountNumber" dynamodbav:"accountNumber,omitempty"`
	Iban          string    `json:"iban" dynamodbav:"iban,omitempty"`
	NetPrice      string    `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice    string    `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Currency      string    `json:"currency" dynamodbav:"currency,omitempty"`
	DueDate       string    `json:"dueDate" dynamodbav:"dueDate,omitempty"`
	Services      []Service `json:"services" dynamodbav:"services,omitempty,omitemptyelem"`

	CustomerEmail string `json:"custmerEmail" dynamodbav:"custmerEmail,omitempty"`
	Status        string `json:"status" dynamodbav:"status,omitempty"` // InvoiceStatus
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
