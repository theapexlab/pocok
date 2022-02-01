package utils

import (
	"encoding/json"
	"pocok/src/utils/models"
	"strings"
)

func MapToStruct(data interface{}, v interface{}) error {
	jsonData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		LogError("error while marshaling json", marshalErr)
		return marshalErr
	}
	unmarshalErr := json.Unmarshal(jsonData, v)
	if unmarshalErr != nil {
		LogError("error while unmarshaling json", unmarshalErr)
		return unmarshalErr
	}
	return nil
}

func MapInvoiceToInvoiceServiceIndexes(invoices []models.Invoice) []models.InvoiceWithServiceIndex {
	indexedInvoices := []models.InvoiceWithServiceIndex{}

	for _, invoice := range invoices {
		newInvoice := models.InvoiceWithServiceIndex{Invoice: invoice}
		for i, service := range invoice.Services {
			newInvoice.Services = append(newInvoice.Services, models.ServiceWithIndex{
				Service: service,
				Index:   i,
			})
		}
	}

	return indexedInvoices
}

func MapUpdateDataToInvoice(data map[string]string) (models.Invoice, error) {
	var invoice models.Invoice
	err := MapToStruct(data, &invoice)

	index := 0
	for {
		service := models.Service{}
		serviceMap := map[string]string{}
		found := false
		for key, val := range data {
			parts := strings.Split(key, "_")
			if strings.HasPrefix(parts[0], "service") && parts[2] == string(index) {
				found = true
				fieldName := parts[1]
				serviceMap[fieldName] = val
			}
		}
		if !found {
			break
		}
		err := MapToStruct(service, &service)
		if err != nil {
			LogError("error while parsing service", err)
		}
		invoice.Services = append(invoice.Services, service)
		index++
	}

	return invoice, err

}
