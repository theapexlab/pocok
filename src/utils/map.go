package utils

import (
	"encoding/json"
	"pocok/src/utils/models"
	"strconv"
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

func MapUpdateDataToInvoice(data map[string]string) (models.Invoice, error) {
	var invoice models.Invoice
	mapToStructError := MapToStruct(data, &invoice)

	index := 0
	for {
		service := models.Service{}
		serviceMap := map[string]string{}
		found := false
		for key, val := range data {
			parts := strings.Split(key, "_")
			if strings.HasPrefix(parts[0], "service") && parts[2] == strconv.Itoa(index) {
				found = true
				fieldName := parts[1]
				serviceMap[fieldName] = val
			}
		}
		if !found {
			break
		}
		mapToStructError := MapToStruct(serviceMap, &service)
		if mapToStructError != nil {
			LogError("error while parsing service", mapToStructError)
		}
		invoice.Services = append(invoice.Services, service)
		index++
	}

	return invoice, mapToStructError

}
