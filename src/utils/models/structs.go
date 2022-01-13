package models

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url","base64"
	Body string `json:"body"`
}

// TODO

// type Service struct {
// 	Name       string
// 	Amount     int
// 	NetPrice   int
// 	GrossPrice int
// 	Currency   string
// 	Tax        float64
// }

type Invoice struct {
	Id           string `json:"id" dynamodbav:"id"`
	Filename     string `json:"filename" dynamodbav:"filename"`
	TextractData string `json:"textractData" dynamodbav:"textractData,omitemty"`

	InvoiceNumber string `json:"invoiceNumber" dynamoadbav:"invoiceNumber,omitempty"`
	CustomerName  string `json:"customerName" dynamoadbav:"customerName,omitempty"`
	AccountNumber string `json:"accountNumber" dynamoadbav:"accountNumber,omitempty"`
	Iban          string `json:"iban" dynamoadbav:"iban,omitempty"`
	NetPrice      int    `json:"netPrice" dynamoadbav:"netPrice,omitempty"`
	GrossPrice    int    `json:"grossPrice" dynamoadbav:"grossPrice,omitempty"`
	Currency      string `json:"currency" dynamoadbav:"currency,omitempty"`
	DueDate       string `json:"dueDate" dynamoadbav:"dueDate,omitempty"`

	// TODO Refactor later into Service struct array
	ServiceName       string `json:"serviceName" dynamoadbav:"serviceName,omitempty"`
	ServiceAmount     int    `json:"serviceAmount" dynamoadbav:"serviceAmount,omitempty"`
	ServiceNetPrice   int    `json:"serviceNetPrice" dynamoadbav:"serviceNetPrice,omitempty"`
	ServiceGrossPrice int    `json:"serviceGrossPrice" dynamoadbav:"serviceGrossPrice,omitempty"`
	ServiceCurrency   string `json:"serviceCurrency" dynamoadbav:"serviceCurrency,omitempty"`
	ServiceTax        int    `json:"serviceTax" dynamoadbav:"serviceTax,omitempty"`

	CustomerEmail string `json:"custmerEmail" dynamoadbav:"custmerEmail,omitempty"`
	Status        string `json:"status" dynamoadbav:"status,omitempty"` // InvoiceStatus
}
