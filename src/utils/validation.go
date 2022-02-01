package utils

import (
	"errors"
	"pocok/src/utils/currency"
	"pocok/src/utils/models"
	"regexp"
	"strconv"
	"time"

	"github.com/araddon/dateparse"

	"github.com/almerlucke/go-iban/iban"
)

func GetValidAccountNumber(accNr string) (string, error) {
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

func GetValidIban(accNr string) (string, error) {
	validIban, err := iban.NewIBAN(accNr)
	if validIban != nil {
		return validIban.Code, nil
	}
	return "", err
}

func GetValidCurrency(currencyInput string) (string, error) {
	validCurrency := currency.GetCurrencyFromString(currencyInput)
	if validCurrency != "" {
		return validCurrency, nil
	}
	return "", errors.New("invalid currency")
}

func GetValidPrice(price string) (string, error) {
	priceInt, err := strconv.Atoi(price)
	if err == nil && priceInt > 0 {
		return price, nil
	}
	return "", err
}

func GetValidDueDate(dueDate string) (string, error) {
	currentTime := time.Now()
	date, err := dateparse.ParseAny(dueDate)

	if err == nil && date.After(currentTime) {
		return dueDate, nil
	}
	return "", errors.New("invalid date")
}
