package wise

import (
	"pocok/src/services/wise/api"
	"pocok/src/services/wise/api/models"
	"pocok/src/utils"
	base_models "pocok/src/utils/models"

	"github.com/google/uuid"
)

type WiseService struct {
	wise *api.WiseClient
}

func (s *WiseService) GetBusinessProfile() (*models.Profile, error) {
	profiles, getProfilesErr := s.wise.GetProfiles()
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
	recipients, getRecipientErr := s.wise.GetRecipientAccounts()
	if getRecipientErr != nil {
		return nil, getRecipientErr
	}

	recipient := api.FindRecipient(*recipients, func(a models.RecipientAccount) bool { return a.Name.FullName == invoice.VendorName })
	if recipient == nil {
		recipientInput := mapInvoiceToRecipient(invoice)

		newRecipient, createRecipientErr := s.wise.CreateRecipientAccount(*recipientInput)
		if createRecipientErr != nil {
			utils.LogError("Error creating recipient", createRecipientErr)
			return nil, createRecipientErr
		}

		recipientById, getRecipientByIdErr := s.wise.GetRecipientAccountById(newRecipient.ID)
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

func (s *WiseService) CreateQuote(profileID int) (*models.Quote, error) {
	var createQuoteInput = models.Quote{
		Profile:        profileID,
		SourceCurrency: "EUR",
		TargetCurrency: "HUF",
		TargetAmount:   200,
	}

	quote, createQuoteErr := s.wise.CreateQuote(createQuoteInput)
	if createQuoteErr != nil {
		utils.LogError("Error creating quote", createQuoteErr)
		return nil, createQuoteErr
	}

	return quote, nil
}

func (s *WiseService) CreateTransfer(quoteID string, recipientID int) (*models.Transfer, error) {
	transferData := models.Transfer{
		TargetAccount: recipientID,
		QuoteUUID:     quoteID,
		Details: struct {
			Reference string `json:"reference"`
		}{Reference: "my ref"},
		CustomerTransactionID: string(uuid.NewString()),
	}

	transfer, createTransferErr := s.wise.CreateTransfer(transferData)
	if createTransferErr != nil {
		utils.LogError("Error creating transfer", createTransferErr)
		return nil, createTransferErr
	}
	return transfer, nil
}
