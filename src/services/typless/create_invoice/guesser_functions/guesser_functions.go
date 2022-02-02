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
	r, _ := regexp.Compile("^[A-Za-z ]*:")
	match := r.FindString(str)
	if match != "" {
		strWithoutPrefix := strings.Replace(str, match, "", 1)
		return strings.TrimSpace(strWithoutPrefix)
	}
	return str
}

func GuessInvoiceNumberFromFilename(filename string, textBlocks *[]typless.TextBlock) string {
	// Gets first value which is included in filename and doesnt contain " " and less then 17 chars
	bankAccountRegex, _ := regexp.Compile(`^[_\d\w-]{3,17}$`)
	containsNumberRegex, _ := regexp.Compile(`[\d+]`)
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

func GuessIban(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		valueWithoutPrefix := cutPrefix(block.Value)
		formattedBlock := strings.ReplaceAll(valueWithoutPrefix, "-", "")
		iban, _ := utils.GetValidIban(formattedBlock)
		if iban != "" {
			return iban
		}
	}
	return ""
}

func GuessHunBankAccountNumber(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		// todo: currently fails to guess if bank account parts are deliminated with " " instead of "-"
		valueParts := strings.Split(block.Value, " ")
		for _, v := range valueParts {
			r, _ := regexp.Compile(models.HUN_BANK_ACC_THREE_PART)
			match := r.FindString(v)

			if match != "" {
				return match
			}
			r, _ = regexp.Compile(models.HUN_BANK_ACC_TWO_PART)
			match = r.FindString(v)

			if match != "" {
				return match
			}
		}
	}
	return ""
}

//  todo: unit test this
func GuessGrossPrice(textBlocks *[]typless.TextBlock) string {
	//  Gets highest price mentioned in texblocks
	highestPrice := ""
	for _, block := range *textBlocks {
		v := block.Value
		if currency.GetCurrencyFromString(v) == "" {
			continue
		}

		price := currency.GetValueFromPrice(v)

		highestPriceInt, convertHighestPriceError := currency.ConvertPriceToFloat(highestPrice)
		priceInt, convertPriceError := currency.ConvertPriceToFloat(price)

		if convertHighestPriceError != nil || convertPriceError != nil {
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
		// todo: match for names with known pre- and suffixes eg.: dr. Test Zoltán ev.
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
