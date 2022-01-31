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
	for _, block := range *textBlocks {
		v := block.Value
		isMatching, _ := regexp.MatchString(`^[_\d\w-]{3,17}$`, v)
		containsNumber, _ := regexp.MatchString(`[\d+]`, v)
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
		iban, _ := iban.NewIBAN(formattedBlock)
		if iban != nil {
			return iban.Code

		}
	}
	return ""
}

func GuessHunBankAccountNumber(textBlocks *[]typless.TextBlock) string {
	for _, block := range *textBlocks {
		v := strings.TrimSpace(block.Value)
		r, _ := regexp.Compile("[0-9]{8}-[0-9]{8}-[0-9]{8}")
		match := r.FindString(v)

		if match != "" {
			return match
		}
		r, _ = regexp.Compile("^[0-9]{8}-[0-9]{8}$")
		match = r.FindString(v)

		if match != "" {
			return match
		}

	}
	return ""
}

func GuessGrossPrice(textBlocks *[]typless.TextBlock) string {
	//  Gest highest price mentioned in texblocks
	highestPrice := ""
	for _, block := range *textBlocks {
		v := block.Value
		if currency.GetCurrencyFromString(v) != "" {
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
	// Gets first "Title Cased" value, which must include atleast one " " in it
	for _, block := range *textBlocks {
		v := strings.TrimSpace(block.Value)
		if !strings.Contains(v, " ") {
			continue
		}
		r, _ := regexp.Compile("^[A-ZÉÁÓÚŐÖŰÜ](([ -][A-ZÉÁÓÚŐÖŰÜ(])?[a-zA-Z)éÉáÁóÓúÚőŐöÖűŰ-]*)*$")
		match := r.FindString(v)
		if match != "" {
			return v
		}
	}
	return ""
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
			if strings.Contains(v, month) {
				v = strings.Replace(v, month, strconv.Itoa(i+1), 1)
			}
		}
		v = strings.TrimRight(v, ".") // trailing "." makes dateparsing fail

		date, err := dateparse.ParseAny(v)
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
		return foundDates[len(foundDates)-1].Format("2006-01-02")
	}

	return ""
}
