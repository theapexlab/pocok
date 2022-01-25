package models

import "github.com/golang-jwt/jwt"

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url","base64"
	Body string `json:"body"`
}

type Service struct {
	Name         string `json:"name" dynamodbav:"name,omitempty"`
	Amount       string `json:"amount" dynamodbav:"amount,omitempty"`
	Unit         string `json:"" dynamodbav:"unit,omitempty"`
	UnitNetPrice string `json:"unitNetPrice" dynamodbav:"unitNetPrice,omitempty"`
	NetPrice     string `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice   string `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Currency     string `json:"currency" dynamodbav:"currency,omitempty"`
	VatAmount    string `json:"vatAmount" dynamodbav:"vatAmount,omitempty"`
	VatRate      string `json:"vatRate" dynamodbav:"vatRate,omitempty"`
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

	VendorEmail string `json:"vendorEmail" dynamodbav:"vendorEmail,omitempty"`

	InvoiceNumber string    `json:"invoiceNumber" dynamodbav:"invoiceNumber,omitempty"`
	VendorName    string    `json:"vendorName" dynamodbav:"vendorName,omitempty"`
	AccountNumber string    `json:"accountNumber" dynamodbav:"accountNumber,omitempty"`
	Iban          string    `json:"iban" dynamodbav:"iban,omitempty"`
	NetPrice      string    `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice    string    `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Currency      string    `json:"currency" dynamodbav:"currency,omitempty"`
	DueDate       string    `json:"dueDate" dynamodbav:"dueDate,omitempty"`
	VatAmount     string    `json:"vatAmount" dynamodbav:"vatAmount,omitempty"`
	VatRate       string    `json:"vatRate" dynamodbav:"vatRate,omitempty"`
	Services      []Service `json:"services" dynamodbav:"services,omitempty,omitemptyelem"`
}

type InvoiceResponse struct {
	Items []Invoice `json:"items"`
	Total int       `json:"total"`
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

type EmailResponseData struct {
	Amp         string
	Attachments map[string][]byte
}

type JWTCustomClaims struct {
	OrgId string `json:"orgId"`
}
type JWTClaims struct {
	jwt.StandardClaims
	JWTCustomClaims
}
