package utils

import (
	"errors"
	"pocok/src/utils/models"
	"regexp"

	"github.com/almerlucke/go-iban/iban"
)

func ValidateAccountNumber(accNr string) (string, error) {
	r, _ := regexp.Compile(models.HUN_BANK_ACC_THREE_PART)
	match := r.FindString(accNr)

	if match != "" {
		return match, nil
	}
	r, _ = regexp.Compile(models.HUN_BANK_ACC_TWO_PART)
	match = r.FindString(accNr)

	if match != "" {
		return match, nil
	}

	return "", errors.New("invalid account number")
}

func ValidateIban(accNr string) (string, error) {
	validIban, err := iban.NewIBAN(accNr)
	if validIban != nil {
		return validIban.Code, nil
	}
	return "", err
}
