package wise

import (
	"pocok/src/services/wise/api"
	"pocok/src/services/wise/api/models"
	"pocok/src/utils"
	base_models "pocok/src/utils/models"
)

const (
	WiseStep1 = "step1:get_profile_id"
	WiseStep2 = "step2:upsert_recipient_account"
	WiseStep3 = "step3:create_quote"
	WiseStep4 = "step4:create_transfer"
)

type WiseMessageData struct {
	RequestType        string              `json:"requestType"`
	ProfileId          int                 `json:"profileId"`
	RecipientAccountId int                 `json:"recipientAccountId"`
	QuoteId            string              `json:"quoteId"`
	TransactionId      string              `json:"transactionId"`
	Invoice            base_models.Invoice `json:"invoice"`
}

type WiseService struct {
	WiseApi *api.WiseClient
}

func CreateWiseService(apiToken string) *WiseService {
	wise := api.CreateWiseClient(apiToken)
	return &WiseService{
		wise,
	}
}

func (s *WiseService) GetBusinessProfile() (*models.Profile, error) {
	profiles, getProfilesErr := s.WiseApi.GetProfiles()
	if getProfilesErr != nil {
		utils.LogError("Error getting profiles", getProfilesErr)
		return nil, getProfilesErr
	}

	businessProfile, findProfileErr := api.FindProfile(*profiles, func(a models.Profile) bool { return a.Type == "business" })
	if businessProfile == nil {
		utils.LogError("", findProfileErr)
		return nil, findProfileErr
	}

	return businessProfile, nil
}

func (s *WiseService) UpsertRecipient(invoice *base_models.Invoice) (*models.RecipientAccount, error) {
	recipients, getRecipientErr := s.WiseApi.GetRecipientAccounts()
	if getRecipientErr != nil {
		return nil, getRecipientErr
	}

	recipient := api.FindRecipient(*recipients, func(a models.RecipientAccount) bool { return a.Name.FullName == invoice.VendorName })
	if recipient == nil {
		recipientInput := mapInvoiceToRecipient(invoice)

		newRecipient, createRecipientErr := s.WiseApi.CreateRecipientAccount(*recipientInput)
		if createRecipientErr != nil {
			utils.LogError("Error creating recipient", createRecipientErr)
			return nil, createRecipientErr
		}

		recipientById, getRecipientByIdErr := s.WiseApi.GetRecipientAccountById(newRecipient.ID)
		if getRecipientByIdErr != nil {
			utils.LogError("Error getting recipient by id", getRecipientByIdErr)
			return nil, getRecipientByIdErr
		}
		recipient = recipientById
	}

	return recipient, nil
}

func mapInvoiceToRecipient(invoice *base_models.Invoice) *models.RecipientAccountV1 {
	return &models.RecipientAccountV1{
		AccountHolderName: invoice.VendorName,
		Currency:          invoice.Currency,
		Type:              "hungarian",
		Details: models.RecipientAccountV1Details{
			LegalType:     "PRIVATE",
			AccountNumber: invoice.AccountNumber,
			Email:         invoice.VendorEmail,
		},
	}
}
