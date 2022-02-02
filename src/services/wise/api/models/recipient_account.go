package models

import "fmt"

type RecipientAccount struct {
	ID        int `json:"id"`
	CreatorID int `json:"creatorId"`
	ProfileID int `json:"profileId"`
	Name      struct {
		FullName                 string      `json:"fullName"`
		GivenName                interface{} `json:"givenName"`
		FamilyName               interface{} `json:"familyName"`
		MiddleName               interface{} `json:"middleName"`
		PatronymicName           interface{} `json:"patronymicName"`
		CannotHavePatronymicName interface{} `json:"cannotHavePatronymicName"`
	} `json:"name"`
	Email           string `json:"email"`
	Currency        string `json:"currency"`
	Country         string `json:"country"`
	Type            string `json:"type"`
	LegalEntityType string `json:"legalEntityType"`
	Active          bool   `json:"active"`
	Details         struct {
		Iban                       string      `json:"iban"`
		AccountNumber              string      `json:"accountNumber"`
		Bic                        interface{} `json:"bic"`
		HashedByLooseHashAlgorithm string      `json:"hashedByLooseHashAlgorithm"`
	} `json:"details"`
	CommonFieldMap struct {
		AccountNumberField string `json:"accountNumberField"`
	} `json:"commonFieldMap"`
	IsDefaultAccount   bool   `json:"isDefaultAccount"`
	Hash               string `json:"hash"`
	AccountSummary     string `json:"accountSummary"`
	LongAccountSummary string `json:"longAccountSummary"`
	DisplayFields      []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"displayFields"`
	OwnedByCustomer bool `json:"ownedByCustomer"`
}

func (r RecipientAccount) String() string {
	return fmt.Sprintf("%v %v (email: %v)", r.ID, r.Name.FullName, r.Email)
}

type RecipientAccountList struct {
	Content []RecipientAccount `json:"content"`
	Sort    struct {
		Empty    bool `json:"empty"`
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
	} `json:"sort"`
	Size int `json:"size"`
}

type RecipientAccountV1 struct {
	ID                int                       `json:"id"`
	Profile           int                       `json:"profile"`
	AccountHolderName string                    `json:"accountHolderName"`
	Currency          string                    `json:"currency"`
	Type              string                    `json:"type"` // iban, hungarian; examples: https://api-docs.wise.com/#recipient-accounts-create
	Details           RecipientAccountV1Details `json:"details"`
}

type RecipientAccountV1Details struct {
	LegalType     string `json:"legalType"`     // PRIVATE of BUSINESS
	AccountNumber string `json:"accountNumber"` // if type is "hungarian"
	IBAN          string `json:"IBAN"`          // if type is "iban"
	Email         string `json:"email"`
}
