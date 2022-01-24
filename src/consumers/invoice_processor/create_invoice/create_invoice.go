package create_invoice

import (
	"pocok/src/services/typless"
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
		case typless.INVOICE_NUMBER:
			invoice.InvoiceNumber = getFieldValue(field)
		case typless.VENDOR_NAME:
			invoice.VendorName = getFieldValue(field)
		case typless.ACCOUNT_NUMBER:
			invoice.AccountNumber = getFieldValue(field)
		case typless.IBAN:
			invoice.Iban = getFieldValue(field)
		case typless.NET_PRICE:
			invoice.NetPrice = getFieldValue(field)
		case typless.GROSS_PRICE:
			invoice.GrossPrice = getFieldValue(field)
		case typless.CURRENCY:
			invoice.Currency = getFieldValue(field)
		case typless.DUE_DATE:
			invoice.DueDate = getFieldValue(field)
		case typless.VAT_RATE:
			invoice.DueDate = getFieldValue(field)
		case typless.VAT_AMOUNT:
			invoice.DueDate = getFieldValue(field)
		}
	}
LineItemsLoop:
	for _, lineItemFields := range extractedData.LineItems {
		service := models.Service{}
		for _, field := range lineItemFields {
			switch field.Name {
			case typless.SERVICE_NAME:
				service.Name = getFieldValue(field)
			case typless.SERVICE_AMOUNT:
				service.Amount = getFieldValue(field)
			case typless.SERVICE_UNIT:
				service.Unit = getFieldValue(field)
			case typless.SERVICE_NET_PRICE:
				service.NetPrice = getFieldValue(field)
			case typless.SERVICE_GROSS_PRICE:
				service.GrossPrice = getFieldValue(field)
			case typless.SERVICE_CURRENCY:
				service.Currency = getFieldValue(field)
			case typless.SERVICE_VAT_RATE:
				service.VatRate = getFieldValue(field)
			case typless.SERVICE_VAT_AMOUNT:
				service.VatAmount = getFieldValue(field)
			}
		}

		if service.Name == "" || (service.GrossPrice == "" && service.Currency == "" && service.NetPrice == "") {
			continue LineItemsLoop
		}

		invoice.Services = append(invoice.Services, service)
	}

	return &invoice
}
