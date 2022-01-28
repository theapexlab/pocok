package typless

const (
	INVOICE_NUMBER string = "invoice_number"
	GROSS_PRICE    string = "total_amount"
	NET_PRICE      string = "net_amount"
	SUPPLIER_NAME  string = "supplier_name"
	VENDOR_NAME    string = "vendor_name"
	ACCOUNT_NUMBER string = "account_number"
	IBAN           string = "iban"
	DUE_DATE       string = "pay_due_date"

	CURRENCY   string = "currency"
	VAT_RATE   string = "vat_rate"
	VAT_AMOUNT string = "vat_amount"

	SERVICE_NAME        string = "product_description"
	SERVICE_AMOUNT      string = "quantity"
	SERVICE_UNIT        string = "unit"
	SERVICE_NET_PRICE   string = "net_price"
	SERVICE_GROSS_PRICE string = "price"
	SERVICE_CURRENCY    string = "currency"
	SERVICE_VAT_RATE    string = "vat_rate"
	SERVICE_VAT_AMOUNT  string = "vat_amount"
)

var ExtractDataToInvoiceMap map[string]string = map[string]string{
	INVOICE_NUMBER: "InvoiceNumber",
	VENDOR_NAME:    "VendorName",
	ACCOUNT_NUMBER: "AccountNumber",
	IBAN:           "Iban",
	NET_PRICE:      "NetPrice",
	GROSS_PRICE:    "GrossPrice",
	CURRENCY:       "Currency",
	DUE_DATE:       "DueDate",
	VAT_RATE:       "VatRate",
	VAT_AMOUNT:     "VatAmount",
}

var LineItemsToServiceMap map[string]string = map[string]string{
	SERVICE_NAME:        "Name",
	SERVICE_UNIT:        "Unit",
	SERVICE_AMOUNT:      "Amount",
	SERVICE_NET_PRICE:   "NetPrice",
	SERVICE_GROSS_PRICE: "GrossPrice",
	SERVICE_VAT_RATE:    "VatRate",
	SERVICE_VAT_AMOUNT:  "VatAmount",
}

type Config struct {
	Token   string
	DocType string
}

type ExtractDataFromFileInput struct {
	DocumentTypeName string `json:"document_type_name"`
	FileName         string `json:"file_name"`
	File             string `json:"file"`
}

type ExtractDataFromFileOutput struct {
	Customer        string             `json:"customer"`
	ExtractedFields []ExtractedField   `json:"extracted_fields"`
	LineItems       [][]ExtractedField `json:"line_items"`
	ObjectId        string             `json:"object_id"`
}

type ExtractedField struct {
	DataType string             `json:"data_type"`
	Name     string             `json:"name"`
	Values   []ExtractionResult `json:"values"`
}

type ExtractionResult struct {
	ConfidenceScore float64 `json:"confidence_score"`
	Value           string  `json:"value"`
}

type TrainingData struct {
	DocumentObjectId string            `json:"document_object_id"`
	LearningFields   []LearningField   `json:"learning_fields"`
	LineItems        [][]LearningField `json:"line_items"`
}

type LearningField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AddDocumentFeedbackInput struct {
	DocumentTypeName string            `json:"document_type_name"`
	LearningFields   []LearningField   `json:"learning_fields"`
	LineItems        [][]LearningField `json:"line_items,omitempty"`
	DocumentObjectId string            `json:"document_object_id"`
}
