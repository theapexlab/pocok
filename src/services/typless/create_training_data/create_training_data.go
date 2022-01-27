package create_training_data

import (
	"pocok/src/services/typless"
	"pocok/src/utils/models"
	"reflect"
)

func CreateTrainingData(invoice *models.Invoice) *typless.TrainingData {
	var trainingData typless.TrainingData

	trainingData.DocumentObjectId = invoice.TyplessObjectId

	for typlessField, invoiceField := range typless.ExtractDataToInvoiceMap {
		trainingData.LearningFields = append(trainingData.LearningFields, typless.LearningField{
			Name:  typlessField,
			Value: reflect.ValueOf(invoice).Elem().FieldByName(invoiceField).String(),
		})
	}

	for typlessField, serviceField := range typless.LineItemsToServiceMap {
		for _, service := range invoice.Services {
			trainingData.LineItems = append(trainingData.LineItems, []typless.LearningField{
				{
					Name:  typlessField,
					Value: reflect.ValueOf(&service).Elem().FieldByName(serviceField).String(),
				},
			})
		}
	}

	return &trainingData
}
