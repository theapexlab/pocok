package typless

const (
	INVOICE_NUMBER      string = "invoice_number"
	GROSS_PRICE         string = "total_amount"
	NET_PRICE           string = "net_amount"
	VENDOR_NAME         string = "supplier_name"
	ACCOUNT_NUMBER      string = "account_number"
	IBAN                string = "iban"
	DUE_DATE            string = "pay_due_date"

	CURRENCY            string = "currency"
	VAT_RATE            string = "vat_rate"
	VAT_AMOUNT          string = "vat_amount"

	SERVICE_NAME        string = "product_description"
	SERVICE_AMOUNT      string = "quantity"
	SERVICE_AMOUNT_UNIT string = "unit"
	SERVICE_NET_PRICE   string = "price"
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
