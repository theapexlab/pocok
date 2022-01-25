package currency

import (
	"regexp"
	"strconv"
	"strings"
)

func GetValueFromPrice(price string) string {
	r := regexp.MustCompile(`[^0-9,\.]*([0-9,\.]*)[^0-9,\.\n]*`)
	firstMatch := r.FindStringSubmatch(price)[1]
	if num, err := strconv.ParseFloat(strings.ReplaceAll(firstMatch, ",", ""), 32); err == nil && num > 0 {
		return firstMatch
	}
	return ""
}

func GetCurrencyFromPrice(price string) string {
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
