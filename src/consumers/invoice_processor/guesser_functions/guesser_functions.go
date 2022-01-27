package guesser_functions

import (
	"pocok/src/services/typless"
	"pocok/src/utils/currency"
	"regexp"
	"strconv"
	"strings"

	"github.com/almerlucke/go-iban/iban"
)

func GuessInvoiceNumberFromFilename(filename string, textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		v := block.Value
		if strings.Contains(filename, v) && !strings.Contains(" ", v) {
			return v
		}
	}
	return ""
}

func GuessIbanFromTextBlocks(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		iban, _ := iban.NewIBAN(block.Value)
		if iban != nil {
			return iban.Code
		}
	}
	return ""
}

func GuessHunBankAccountNumberFromTextBlocks(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		r, _ := regexp.Compile("[0-9]{8}-[0-9]{8}-[0-9]{8}")
		match := r.FindString(block.Value)
		if match != "" {
			return match
		}
	}
	return ""
}

func GuessGrossPriceFromTextBlocks(textBlocks *[]typless.TextBlock) string {
	highestPrice := ""
	for _, block := range *textBlocks {
		v := block.Value
		if currency.GetCurrencyFromPrice(v) != "" {
			continue
		}

		price := currency.GetValueFromPrice(v)

		highestPriceInt, err := strconv.Atoi(highestPrice)
		priceInt, err := strconv.Atoi(price)
		if err != nil {
			continue
		}

		if highestPriceInt < priceInt {
			highestPrice = price
		}
	}
	return highestPrice
}

func GuessVendorName(textBlocks *[]typless.TextBlock) string {
	// first "Title Case" value with atleast one " " or "-" deliminator, which doesnt contain special characters, except "-"
	for _, block := range *textBlocks {
		v := block.Value
		match, _ := regexp.MatchString("^[A-Z]+(([ -][A-Z ])?[a-z]*)*$", strings.TrimSpace(v))
		if match == true {
			return v
		}
	}
	return ""
}
