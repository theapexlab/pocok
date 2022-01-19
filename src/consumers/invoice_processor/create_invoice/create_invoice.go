package create_invoice

import (
	"pocok/src/services/typless"
	"pocok/src/utils/currency"
	"pocok/src/utils/models"
	"strings"
)

func getFieldValue(field typless.ExtractedField) string {
	firstValueField := field.Values[0]
	if firstValueField.ConfidenceScore > 0 {
		return strings.TrimSpace(firstValueField.Value)
	}
	return ""
}

func CreateInvoice(extractedData *typless.ExtractDataFromFileOutput) *models.Invoice {
	invoice := models.Invoice{}
	for _, field := range extractedData.ExtractedFields {
		switch field.Name {
		case "invoice_number":
			invoice.InvoiceNumber = getFieldValue(field)
		case "customer_name":
			invoice.CustomerName = getFieldValue(field)
		case "account_number":
			invoice.AccountNumber = getFieldValue(field)
		case "iban":
			invoice.Iban = getFieldValue(field)
		case "net_price":
			invoice.NetPrice = currency.TrimCurrencyFromPrice(getFieldValue(field))
		case "gross_price":
			invoice.GrossPrice = currency.TrimCurrencyFromPrice(getFieldValue(field))
			invoice.Currency = currency.GetCurrencyFromPrice(getFieldValue(field))
		case "due_date":
			invoice.DueDate = getFieldValue(field)
		}
	}

	for _, lineItemFields := range extractedData.LineItems {
		service := models.Service{}
		for _, field := range lineItemFields {
			switch field.Name {
			case "service_name":
				service.Name = getFieldValue(field)
			case "service_amount":
				service.Amount = getFieldValue(field)
			case "service_net_price":
				service.NetPrice = currency.GetCurrencyFromPrice(getFieldValue(field))
			case "service_gross_price":
				service.GrossPrice = currency.TrimCurrencyFromPrice(getFieldValue(field))
				service.Currency = currency.GetCurrencyFromPrice(getFieldValue(field))
			case "service_vat":
				service.Tax = getFieldValue(field)
			}
		}
		if service.Name != "" {
			invoice.Services = append(invoice.Services, service)
		}
	}

	return &invoice
}
