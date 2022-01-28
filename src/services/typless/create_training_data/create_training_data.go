package create_training_data

import (
	"pocok/src/services/typless"
	"pocok/src/utils/models"
	"reflect"
)

func CreateTrainingData(invoice *models.Invoice) *typless.TrainingData {
	var trainingData typless.TrainingData

	trainingData.DocumentObjectId = invoice.TyplessObjectId

	trainingData.LearningFields = append(trainingData.LearningFields, typless.LearningField{
		Name:  typless.SUPPLIER_NAME,
		Value: invoice.VendorName,
	})

	for typlessField, invoiceField := range typless.ExtractDataToInvoiceMap {
		trainingData.LearningFields = append(trainingData.LearningFields, typless.LearningField{
			Name:  typlessField,
			Value: reflect.ValueOf(invoice).Elem().FieldByName(invoiceField).String(),
		})
	}

	for _, service := range invoice.Services {
		var lineItems []typless.LearningField

		for typlessField, serviceField := range typless.LineItemsToServiceMap {
			lineItems = append(lineItems, typless.LearningField{
				Name:  typlessField,
				Value: reflect.ValueOf(&service).Elem().FieldByName(serviceField).String(),
			})
		}

		trainingData.LineItems = append(trainingData.LineItems, lineItems)
	}

	return &trainingData
}
