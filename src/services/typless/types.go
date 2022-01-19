package typless

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
