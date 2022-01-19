package models

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url","base64"
	Body string `json:"body"`
}

type Service struct {
	Name         string `json:"name" dynamodbav:"name,omitempty"`
	Amount       int    `json:"amount" dynamodbav:"amount,omitempty"`
	UnitNetPrice int    `json:"unitNetPrice" dynamodbav:"unitNetPrice,omitempty"`
	NetPrice     int    `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice   int    `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Currency     string `json:"currency" dynamodbav:"currency,omitempty"`
	Tax          int    `json:"tax" dynamodbav:"tax,omitempty"`
}

type Invoice struct {
	Pk     string `json:"pk" dynamodbav:"pk"`
	Sk     string `json:"sk" dynamodbav:"sk"`
	Lsi1sk string `json:"lsi1sk" dynamodbav:"lsi1sk"`
	Lsi2sk string `json:"lsi2sk" dynamodbav:"lsi2sk"`

	InvoiceId  string `json:"invoiceId" dynamodbav:"invoiceId"`
	EntityType string `json:"entityType" dynamodbav:"entityType"`
	Status     string `json:"status" dynamodbav:"status"`
	ReceivedAt string `json:"receivedAt" dynamodbav:"receivedAt"`
	Filename   string `json:"filename" dynamodbav:"filename"`

	CustomerEmail string `json:"customerEmail" dynamodbav:"customerEmail,omitempty"`

	InvoiceNumber string    `json:"invoiceNumber" dynamodbav:"invoiceNumber,omitempty"`
	CustomerName  string    `json:"customerName" dynamodbav:"customerName,omitempty"`
	AccountNumber string    `json:"accountNumber" dynamodbav:"accountNumber,omitempty"`
	Iban          string    `json:"iban" dynamodbav:"iban,omitempty"`
	NetPrice      int       `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice    int       `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Tax           string    `json:"tax" dynamodbav:"tax,omitempty"`
	Currency      string    `json:"currency" dynamodbav:"currency,omitempty"`
	DueDate       string    `json:"dueDate" dynamodbav:"dueDate,omitempty"`
	Services      []Service `json:"services" dynamodbav:"services,omitempty,omitemptyelem"`
	TextractData  string    `json:"textractData" dynamodbav:"textractData,omitempty"`
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

type Email struct {
	To          string   `json:"to"`
	Subject     string   `json:"subject"`
	Html        string   `json:"html"`
	Attachments []string `json:"attachments"`
}
