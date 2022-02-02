package create_invoice

import (
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_invoice/guesser_functions"
	"pocok/src/utils/models"
	"reflect"
	"strings"
)

type CreateInvoiceService struct {
	OriginalFilename string
	ExtractedData    *typless.ExtractDataFromFileOutput
}

func (c *CreateInvoiceService) CreateInvoice() *models.Invoice {
	invoice := models.Invoice{}

	for _, field := range c.ExtractedData.ExtractedFields {
		c.setFieldValue(reflect.ValueOf(&invoice), typless.ExtractDataToInvoiceMap, field)
	}

LineItemsLoop:
	for _, lineItemFields := range c.ExtractedData.LineItems {
		service := models.Service{}

		for _, field := range lineItemFields {
			c.setFieldValue(reflect.ValueOf(&service), typless.LineItemsToServiceMap, field)
		}

		if service.Name == "" || (service.GrossPrice == "" && service.NetPrice == "") {
			continue LineItemsLoop
		}

		invoice.Services = append(invoice.Services, service)
	}

	invoice.TyplessObjectId = c.ExtractedData.ObjectId

	return &invoice
}

func (c *CreateInvoiceService) setFieldValue(reflectValue reflect.Value, fieldMap map[string]string, field typless.ExtractedField) {
	// todo: maybe add validation to typless currency?
	fieldName := fieldMap[field.Name]
	if fieldName == "" {
		return
	}

	fieldValue := c.getFieldValue(field)

	reflectValue.Elem().FieldByName(fieldName).SetString(fieldValue)
}

func (c *CreateInvoiceService) getFieldValue(field typless.ExtractedField) string {
	firstValueField := field.Values[0]
	if firstValueField.Value != "" && firstValueField.ConfidenceScore > 0 {
		return strings.TrimSpace(firstValueField.Value)
	}
	return c.getFieldFallbackValue(field.Name)
}

func (c *CreateInvoiceService) getFieldFallbackValue(fieldName string) string {
	switch fieldName {
	case typless.INVOICE_NUMBER:
		return guesser_functions.GuessInvoiceNumberFromFilename(c.OriginalFilename, &c.ExtractedData.TextBlocks)
	case typless.VENDOR_NAME:
		return guesser_functions.GuessVendorName(&c.ExtractedData.TextBlocks)
	case typless.ACCOUNT_NUMBER:
		return guesser_functions.GuessHunBankAccountNumber(&c.ExtractedData.TextBlocks)
	case typless.IBAN:
		return guesser_functions.GuessIban(&c.ExtractedData.TextBlocks)
	case typless.CURRENCY:
		return guesser_functions.GuessCurrency(&c.ExtractedData.TextBlocks)
	case typless.GROSS_PRICE:
		return guesser_functions.GuessGrossPrice(&c.ExtractedData.TextBlocks)
	case typless.DUE_DATE:
		return guesser_functions.GuessDueDate(&c.ExtractedData.TextBlocks)
	default:
		return ""
	}
}
