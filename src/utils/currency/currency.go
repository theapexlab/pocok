package currency

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func GetValueFromPrice(price string) string {
	r := regexp.MustCompile(`[^0-9,\.]*([0-9,\.]*)[^0-9,\.\n]*`)
	priceWithoutSpaces := strings.ReplaceAll(price, " ", "")
	firstMatch := r.FindStringSubmatch(priceWithoutSpaces)[1]
	if _, err := ConvertPriceToFloat(firstMatch); err == nil {
		return firstMatch
	}
	return ""
}

func GetCurrencyFromString(price string) string {
	currencyMap := map[string][]string{
		"EUR": {"€", "EUR"},
		"USD": {"$", "USD"},
		"HUF": {"Ft", "HUF"},
	}

	for currency, patterns := range currencyMap {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(price), strings.ToLower(pattern)) {
				return currency
			}
		}
	}

	return ""
}

func ConvertPriceToFloat(price string) (float64, error) {
	priceWithoutSpaces := strings.ReplaceAll(price, " ", "")
	priceWithoutCommas := strings.ReplaceAll(priceWithoutSpaces, ",", "")
	num, err := strconv.ParseFloat(priceWithoutCommas, 32)
	if err == nil && num > 0 {
		return num, nil
	}
	return 0, errors.New("invalid price")
}
