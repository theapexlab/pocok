package utils

import (
	"encoding/json"
	"pocok/src/utils/models"
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
