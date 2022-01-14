package parse_expense_analysis

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/textract"
)

func ParseExpenseAnalysis(analysisOutput *textract.GetExpenseAnalysisOutput) {
	for _, document := range analysisOutput.ExpenseDocuments {
		for _, record := range document.SummaryFields {
			switch *record.Type.Text {
			case "INVOICE_RECEIPT_ID":
				fmt.Println("invoiceNr ", *record.ValueDetection.Text)
			case "DUE_DATE":
				fmt.Println("DueDate", *record.ValueDetection.Text)
			case "VENDOR_NAME":
				fmt.Println("CustomerName", *record.ValueDetection.Text)
			case "TAX":
				fmt.Println("Tax", *record.ValueDetection.Text)
			case "SUBTOTAL":
				fmt.Println("NetPrice", *record.ValueDetection.Text)
			case "TOTAL":
				fmt.Println("GrossPrice", *record.ValueDetection.Text)
			}
		}
	}
}

/*
(DB FIELD) -  ( ANALYZE EXPENSE TYPE)
invoiceNr - INVOICE_RECEIPT_ID
DueDate - DUE_DATE
CustomerName - VENDOR_NAME
Tax -TAX
NetPrice - SUBTOTAL
GrossPrice - TOTAL

Missing Invoice struct props =>
 	InvoiceNumber string
	CustomerName  string
	AccountNumber string
	Iban          string
	Currency      string
	DueDate       stringd
	Name
	Amount
	NetPrice
	GrossPrice
	Currency

The following is a list of the standard fields that AnalyzeExpense currently supports:
	Vendor Name: VENDOR_NAME
	Total: TOTAL
	Receiver Address: RECEIVER_ADDRESS
	Invoice/Receipt Date: INVOICE_RECEIPT_DATE
	Invoice/Receipt ID: INVOICE_RECEIPT_ID
	Payment Terms: PAYMENT_TERMS
	Subtotal: SUBTOTAL
	Due Date: DUE_DATE
	Tax: TAX
	Invoice Tax Payer ID (SSN/ITIN or EIN): TAX_PAYER_ID
	Item Name: ITEM_NAME
	Item Price: PRICE
	Item Quantity: QUANTITY
*/
