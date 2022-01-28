package create_invoice

import (
	"pocok/src/consumers/invoice_processor/guesser_functions"
	"pocok/src/services/typless"
	"pocok/src/utils/models"
	"strings"
)

type CreateInvoiceService struct {
	OriginalFilename string
}

func (c *CreateInvoiceService) getFieldFallbackValue(field *typless.ExtractedField, textBlocks *[]typless.TextBlock) string {
	switch field.Name {
	case typless.INVOICE_NUMBER:
		return guesser_functions.GuessInvoiceNumberFromFilename(c.OriginalFilename, textBlocks)
	case typless.VENDOR_NAME:
		return guesser_functions.GuessVendorName(textBlocks)
	case typless.ACCOUNT_NUMBER:
		return guesser_functions.GuessHunBankAccountNumber(textBlocks)
	case typless.IBAN:
		return guesser_functions.GuessIban(textBlocks)
	case typless.CURRENCY:
		return guesser_functions.GuessCurrency(textBlocks)
	case typless.GROSS_PRICE:
		return guesser_functions.GuessGrossPrice(textBlocks)
	case typless.DUE_DATE:
		return guesser_functions.GuessDueDate(textBlocks)
	default:
		return ""
	}
}

func (c *CreateInvoiceService) getExtractedFieldValue(extractedData *typless.ExtractDataFromFileOutput, fieldIndex int) string {
	firstValueField := extractedData.ExtractedFields[fieldIndex].Values[0]
	if firstValueField.ConfidenceScore > 0 {
		return strings.TrimSpace(firstValueField.Value)
	}
	return c.getFieldFallbackValue(&extractedData.ExtractedFields[fieldIndex], &extractedData.TextBlocks)
}

func (c *CreateInvoiceService) getLineItemFieldValue(field typless.ExtractedField) string {
	firstValueField := field.Values[0]
	if firstValueField.ConfidenceScore > 0 {
		return strings.TrimSpace(firstValueField.Value)
	}
	return ""
}

func (c *CreateInvoiceService) CreateInvoice(extractedData *typless.ExtractDataFromFileOutput) *models.Invoice {
	invoice := models.Invoice{}
	for i, field := range extractedData.ExtractedFields {
		switch field.Name {
		case typless.INVOICE_NUMBER:
			invoice.InvoiceNumber = c.getExtractedFieldValue(extractedData, i)
		case typless.VENDOR_NAME:
			invoice.VendorName = c.getExtractedFieldValue(extractedData, i)
		case typless.ACCOUNT_NUMBER:
			invoice.AccountNumber = c.getExtractedFieldValue(extractedData, i)
		case typless.IBAN:
			invoice.Iban = c.getExtractedFieldValue(extractedData, i)
		case typless.NET_PRICE:
			invoice.NetPrice = c.getExtractedFieldValue(extractedData, i)
		case typless.GROSS_PRICE:
			invoice.GrossPrice = c.getExtractedFieldValue(extractedData, i)
		case typless.CURRENCY:
			invoice.Currency = c.getExtractedFieldValue(extractedData, i)
		case typless.DUE_DATE:
			invoice.DueDate = c.getExtractedFieldValue(extractedData, i)
		case typless.VAT_RATE:
			invoice.VatRate = c.getExtractedFieldValue(extractedData, i)
		case typless.VAT_AMOUNT:
			invoice.VatAmount = c.getExtractedFieldValue(extractedData, i)
		}
	}
LineItemsLoop:
	for _, lineItemFields := range extractedData.LineItems {
		service := models.Service{}
		for _, field := range lineItemFields {
			switch field.Name {
			case typless.SERVICE_NAME:
				service.Name = c.getLineItemFieldValue(field)
			case typless.SERVICE_AMOUNT:
				service.Amount = c.getLineItemFieldValue(field)
			case typless.SERVICE_UNIT:
				service.Unit = c.getLineItemFieldValue(field)
			case typless.SERVICE_NET_PRICE:
				service.NetPrice = c.getLineItemFieldValue(field)
			case typless.SERVICE_GROSS_PRICE:
				service.GrossPrice = c.getLineItemFieldValue(field)
			case typless.SERVICE_CURRENCY:
				service.Currency = c.getLineItemFieldValue(field)
			case typless.SERVICE_VAT_RATE:
				service.VatRate = c.getLineItemFieldValue(field)
			case typless.SERVICE_VAT_AMOUNT:
				service.VatAmount = c.getLineItemFieldValue(field)
			}
		}

		if service.Name == "" || (service.GrossPrice == "" && service.Currency == "" && service.NetPrice == "") {
			continue LineItemsLoop
		}

		invoice.Services = append(invoice.Services, service)
	}

	return &invoice
}
