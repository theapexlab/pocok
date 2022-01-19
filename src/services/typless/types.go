package typless

const (
	INVOICE_NUMBER      string = "invoice_number"
	GROSS_PRICE         string = "gross_price"
	NET_PRICE           string = "net_price"
	VENDOR_NAME         string = "vendor_name"
	ACCOUNT_NUMBER      string = "account_number"
	IBAN                string = "iban"
	DUE_DATE            string = "due_date"
	SERVICE_NAME        string = "service_name"
	SERVICE_AMOUNT      string = "service_amount"
	SERVICE_NET_PRICE   string = "service_net_price"
	SERVICE_GROSS_PRICE string = "service_gross_price"
	SERVICE_VAT         string = "service_vat"
)

type ExtractDataFromFileInput struct {
	DocumentTypeName string `json:"document_type_name"`
	FileName         string `json:"file_name"`
	File             string `json:"file"`
}

type ExtractDataFromFileOutput struct {
	Customer        string             `json:"customer"`
	ExtractedFields []ExtractedField   `json:"extracted_fields"`
	LineItems       [][]ExtractedField `json:"line_items"`
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
