package create_invoice

import (
	"pocok/src/services/typless"
	"pocok/src/utils/models"
	"reflect"
	"strings"
)

func CreateInvoice(extractedData *typless.ExtractDataFromFileOutput) *models.Invoice {
	invoice := models.Invoice{}

	for _, field := range extractedData.ExtractedFields {
		setFieldValue(reflect.ValueOf(&invoice), typless.ExtractDataToInvoiceMap, field)
	}

LineItemsLoop:
	for _, lineItemFields := range extractedData.LineItems {
		service := models.Service{}

		for _, field := range lineItemFields {
			setFieldValue(reflect.ValueOf(&service), typless.LineItemsToServiceMap, field)
		}

		if service.Name == "" || (service.GrossPrice == "" && service.Currency == "" && service.NetPrice == "") {
			continue LineItemsLoop
		}

		invoice.Services = append(invoice.Services, service)
	}

	invoice.TyplessObjectId = extractedData.ObjectId

	return &invoice
}

func setFieldValue(reflectValue reflect.Value, fieldMap map[string]string, field typless.ExtractedField) {
	fieldName := fieldMap[field.Name]
	if fieldName == "" {
		return
	}

	fieldValue := getFieldValue(field)

	reflectValue.Elem().FieldByName(fieldName).SetString(fieldValue)
}

func getFieldValue(field typless.ExtractedField) string {
	firstValueField := field.Values[0]
	if firstValueField.ConfidenceScore > 0 {
		return strings.TrimSpace(firstValueField.Value)
	}
	return ""
}
