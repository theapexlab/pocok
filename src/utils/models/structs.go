package models

type UploadInvoiceMessage struct {
	Type string `json:"type"` // "url","base64"
	Body string `json:"body"`
}

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

	InvoiceNumber string `json:"invoiceNumber" dynamodbav:"invoiceNumber,omitempty"`
	CustomerName  string `json:"customerName" dynamodbav:"customerName,omitempty"`
	AccountNumber string `json:"accountNumber" dynamodbav:"accountNumber,omitempty"`
	Iban          string `json:"iban" dynamodbav:"iban,omitempty"`
	NetPrice      int    `json:"netPrice" dynamodbav:"netPrice,omitempty"`
	GrossPrice    int    `json:"grossPrice" dynamodbav:"grossPrice,omitempty"`
	Currency      string `json:"currency" dynamodbav:"currency,omitempty"`
	DueDate       string `json:"dueDate" dynamodbav:"dueDate,omitempty"`

	// TODO Refactor later into Service struct array
	ServiceName       string `json:"serviceName" dynamodbav:"serviceName,omitempty"`
	ServiceAmount     int    `json:"serviceAmount" dynamodbav:"serviceAmount,omitempty"`
	ServiceNetPrice   int    `json:"serviceNetPrice" dynamodbav:"serviceNetPrice,omitempty"`
	ServiceGrossPrice int    `json:"serviceGrossPrice" dynamodbav:"serviceGrossPrice,omitempty"`
	ServiceCurrency   string `json:"serviceCurrency" dynamodbav:"serviceCurrency,omitempty"`
	ServiceTax        int    `json:"serviceTax" dynamodbav:"serviceTax,omitempty"`

	CustomerEmail string `json:"custmerEmail" dynamodbav:"custmerEmail,omitempty"`
	Status        string `json:"status" dynamodbav:"status,omitempty"` // InvoiceStatus
}
