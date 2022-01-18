package create_invoice

import (
	"fmt"
	"pocok/src/services/typless"
	"pocok/src/utils/models"
	"regexp"
	"strconv"
)

func getFieldValue(field typless.ExtractedField) string {
	bestValue := field.Values[0]
	if bestValue.ConfidenceScore > 0 {
		return bestValue.Value
	}
	return ""
}

func parsePriceValue(priceValue string) (float32, error) {
	r, err := regexp.Compile(`[0-9]+,?[0-9]*`)
	if err != nil {
		return 0, err
	}
	matches := r.FindStringSubmatch(priceValue)

	numVal, err := strconv.(matches[0])
	if err != nil {
		return 0, err
	}

	return numVal, nil
}

func getCurrencyFromPriceValue(priceValue string) string {
	
}

func CreateInvoice(extractedData *typless.ExtractDataFromFileOutput) (*models.Invoice, error) {
	invoice := models.Invoice{}

	for _, field := range extractedData.ExtractedFields {
		service := models.Service{}
		switch field.Name {
		case "invoice_number":
			invoice.InvoiceNumber = getFieldValue(field)
		case "customer_name":
			invoice.CustomerName = getFieldValue(field)
		case "account_number":
			invoice.AccountNumber = getFieldValue(field)
		case "iban":
			invoice.Iban = getFieldValue(field)
		case "net_price":
			invoice.NetPrice, _ = getFieldValue(field)
		case "gross_price":
			invoice.GrossPrice, _ = getFieldValue(field)
		case "currency":
			// todo: parse currency from gross price
			invoice.Currency = getCurrency(field)
		case "due_date":
			invoice.DueDate = getFieldValue(field)
		case "service_name":
			service.Name = getFieldValue(field)
		case "service_amount":
			fmt.Println("service_amount")
		case "service_net_price":
			fmt.Println("service_net_price")
		case "service_gross_price":
			fmt.Println("service_gross_price")
		case "service_currency":
			fmt.Println("service_currency")
		case "service_vat":
			fmt.Println("service_vat")
		}

		if service.Name != "" {
			invoice.Services = append(invoice.Services, service)
		}

	}

	return &invoice, nil
}
