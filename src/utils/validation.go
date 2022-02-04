package utils

import (
	"errors"
	"pocok/src/utils/currency"
	"pocok/src/utils/models"
	"regexp"
	"strings"

	"github.com/araddon/dateparse"

	"github.com/almerlucke/go-iban/iban"
)

func GetValidAccountNumber(accNr string) (string, error) {
	v := strings.TrimSpace(accNr)
	r := regexp.MustCompile(models.VALIDATE_HUN_BANK_ACC)
	match := r.FindString(v)

	if match != "" {
		reg := regexp.MustCompile("[^0-9]+")
		accNrWithOnlyNumbers := reg.ReplaceAllString(match, "")
		return accNrWithOnlyNumbers, nil
	}

	return "", errors.New("invalid account number")
}

func GetValidIban(accNr string) (string, error) {
	validIban, ibanError := iban.NewIBAN(accNr)
	if validIban != nil {
		return validIban.Code, nil
	}
	return "", ibanError
}

func GetValidCurrency(currencyInput string) (string, error) {
	validCurrency := currency.GetCurrencyFromString(currencyInput)
	if validCurrency != "" {
		return validCurrency, nil
	}
	return "", errors.New("invalid currency")
}

func GetValidPrice(price string) (string, error) {
	_, priceConvertError := currency.ConvertPriceToFloat(price)
	if priceConvertError != nil {
		return "", priceConvertError
	}
	return price, nil
}

func GetValidDueDate(dueDate string) (string, error) {
	_, dateParseError := dateparse.ParseAny(dueDate)
	if dateParseError != nil {
		return "", errors.New("invalid date")
	}
	return dueDate, nil
}
