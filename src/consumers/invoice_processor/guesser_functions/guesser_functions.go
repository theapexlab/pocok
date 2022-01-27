package guesser_functions

import (
	"pocok/src/services/typless"
	"pocok/src/utils/currency"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/almerlucke/go-iban/iban"
	"github.com/araddon/dateparse"
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
	for _, block := range *textBlocks {
		v := strings.TrimSpace(block.Value)
		if !strings.Contains(v, " ") {
			continue
		}
		r, _ := regexp.Compile("^[A-ZéÉáÁóÓúÚőŐűŰ](([ -][A-ZéÉáÁóÓúÚőŐűŰ()])?[a-zA-Z()éÉáÁóÓúÚőŐűŰ]*)*$")
		match := r.FindString(v)
		if match != "" {
			return v
		}
	}
	return ""
}

func FindInArray(searchArray []string, textBlocks *[]typless.TextBlock) string {
	sort.Strings(searchArray)
	for _, block := range *textBlocks {
		matchIndex := sort.SearchStrings(searchArray, block.Value)
		if matchIndex < len(searchArray) {
			return searchArray[matchIndex]
		}
	}
	return ""
}

type timeSlice []time.Time

func (s timeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int           { return len(s) }

func GuessDueDate(textBlocks *[]typless.TextBlock) string {
	var foundDates timeSlice = []time.Time{}
	for _, block := range *textBlocks {
		date, err := dateparse.ParseAny(block.Value)
		if err != nil {
			continue
		}
		isAfterTwentyTwenty := date.After(time.Date(2020, 1, 1, 0, 0, 0, 0, &time.Location{}))
		if isAfterTwentyTwenty {
			foundDates = append(foundDates, date)
		}
	}

	if foundDates.Len() > 0 {
		sort.Sort(foundDates)
		return foundDates[0].Format("2006-01-02")
	}

	return ""
}
