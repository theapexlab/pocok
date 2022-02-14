package api

import (
	"errors"
	"pocok/src/services/wise/api/models"
)

var ErrNoProfileFound = errors.New("no profile found")

func FindRecipient(recipients []models.RecipientAccount, compare func(a models.RecipientAccount) bool) *models.RecipientAccount {
	for _, recipient := range recipients {
		if compare(recipient) {
			return &recipient
		}
	}

	return nil
}

func FindProfile(profiles []models.Profile, compare func(a models.Profile) bool) (*models.Profile, error) {
	for _, profile := range profiles {
		if compare(profile) {
			return &profile, nil
		}
	}

	return nil, ErrNoProfileFound
}
