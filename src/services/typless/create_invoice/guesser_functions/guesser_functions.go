package guesser_functions

import (
	"pocok/src/services/typless"
	"pocok/src/utils"
	"pocok/src/utils/currency"
	"pocok/src/utils/models"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func cutPrefix(str string) string {
	// Example: Payment due date: 2021.09.23.  =>  2021.09.23.
	r := regexp.MustCompile("^[A-Za-z ]*:")
	match := r.FindString(str)
	if match != "" {
		strWithoutPrefix := strings.Replace(str, match, "", 1)
		return strings.TrimSpace(strWithoutPrefix)
	}
	return str
}

func GuessInvoiceNumberFromFilename(filename string, textBlocks *[]typless.TextBlock) string {
	// Gets first value which is included in filename and doesnt contain " " and less then 17 chars
	bankAccountRegex := regexp.MustCompile(`^[_\d\w-]{3,17}$`)
	containsNumberRegex := regexp.MustCompile(`[\d+]`)
	for _, block := range *textBlocks {
		v := block.Value
		isMatching := bankAccountRegex.MatchString(v)
		containsNumber := containsNumberRegex.MatchString(v)
		if strings.Contains(filename, v) &&
			containsNumber &&
			isMatching {
			return v
		}
	}
	return ""
}

func GuessInvoiceNumber(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		if len(block.Value) <= 3 {
			continue
		}
		v := strings.TrimSpace(block.Value)
		r, _ := regexp.Compile("^(# )?[a-zA-Z0-9.-]*$")
		match := r.FindString(v)
		if match != "" {
			return v
		}
	}
	return ""
}

func GuessIban(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		valueWithoutPrefix := cutPrefix(block.Value)
		iban, _ := utils.GetValidIban(valueWithoutPrefix)
		if iban != "" {
			return iban
		}
	}
	return ""
}

func GuessHunBankAccountNumber(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		v := strings.TrimSpace(block.Value)
		r := regexp.MustCompile(models.GUESS_HUN_BANK_ACC)
		match := r.FindString(v)

		if match != "" {
			return match
		}
	}
	return ""
}

func GuessGrossPrice(textBlocks *[]typless.TextBlock) string {
	//  Gets highest price mentioned in texblocks
	highestPrice := "0"
	for _, block := range *textBlocks {
		v := block.Value
		if currency.GetCurrencyFromString(v) == "" {
			continue
		}

		price := currency.GetValueFromPrice(v)

		highestPriceInt, _ := currency.ConvertPriceToFloat(highestPrice)
		priceInt, convertPriceError := currency.ConvertPriceToFloat(price)

		if convertPriceError != nil {
			continue
		}

		if highestPriceInt < priceInt {
			highestPrice = price
		}
	}
	return highestPrice
}

func GuessVendorName(textBlocks *[]typless.TextBlock) string {
	// Gets first "Title Cased" value, which must include atleast one " " in it and no invoice indicator
	for _, block := range *textBlocks {
		v := strings.TrimSpace(block.Value)
		containsSpace := strings.Contains(v, " ")
		if !containsSpace || containsInvoiceIndicator(v) {
			continue
		}
		r, _ := regexp.Compile("^[A-ZÉÁÓÚŐÖŰÜ](([ -][A-ZÉÁÓÚŐÖŰÜ(])?[a-zA-Z)éÉáÁóÓúÚőŐöÖűŰüÜ-]*)*$")
		match := r.FindString(v)
		if match != "" {
			return v
		}
	}
	return ""
}

func containsInvoiceIndicator(v string) bool {
	invoiceIndicators := []string{"invoice", "szamla", "számla"}
	contains := false
	for _, indicator := range invoiceIndicators {
		if strings.Contains(strings.ToLower(v), indicator) {
			contains = true
		}
	}
	return contains
}

func GuessCurrency(textBlocks *[]typless.TextBlock) string {
	// Gets first currency found in text blocks
	for _, block := range *textBlocks {
		currency := currency.GetCurrencyFromString(block.Value)
		if currency != "" {
			return currency
		}
	}
	return ""
}

type timeSlice []time.Time

func (s timeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int           { return len(s) }

func GuessDueDate(textBlocks *[]typless.TextBlock) string {
	// Gets latest parsable date, which is after 2020.01.01 for safety reasons
	months := []string{"jan", "feb", "mar", "apr", "apr", "may", "june", "aug", "sept", "oct", "nov", "dec"}
	var foundDates timeSlice = []time.Time{}
	for _, block := range *textBlocks {
		v := cutPrefix(block.Value)
		v = strings.ReplaceAll(v, " ", "")
		for i, month := range months {
			v = strings.Replace(v, month, strconv.Itoa(i+1), 1)
		}
		v = strings.TrimRight(v, ".") // trailing "." makes dateparsing fail

		date, parseError := dateparse.ParseAny(v)
		if parseError != nil {
			continue
		}

		isAfterTwentyTwenty := date.After(time.Date(2020, 1, 1, 0, 0, 0, 0, &time.Location{}))
		if isAfterTwentyTwenty {
			foundDates = append(foundDates, date)
		}
	}

	if foundDates.Len() > 0 {
		sort.Sort(foundDates)
		return foundDates[len(foundDates)-1].Format("2006-01-02")
	}
	return ""
}
