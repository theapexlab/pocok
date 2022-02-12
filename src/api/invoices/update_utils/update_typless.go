package update_utils

import (
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_training_data"
	"pocok/src/utils"
	"pocok/src/utils/models"
)

func UpdateTypless(typlessToken string, typlessDocType string, invoice models.Invoice) error {
	typlessError := typless.AddDocumentFeedback(
		&typless.Config{
			Token:   typlessToken,
			DocType: typlessDocType,
		},
		*create_training_data.CreateTrainingData(&invoice),
	)
	if typlessError != nil {
		utils.LogError("Error adding document feedback to typless", typlessError)
		return typlessError
	}

	return nil
}
