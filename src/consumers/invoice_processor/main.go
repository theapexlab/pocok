package main

import (
	"context"
	"io/ioutil"
	"os"
	"pocok/src/db"
	"pocok/src/services/typless"
	"pocok/src/services/typless/create_invoice"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	bucketName     string
	typlessToken   string
	typlessDocType string
	tableName      string
	lambdaTimeout  string
	s3Client       *s3.Client
	dbClient       *dynamodb.Client
}

func main() {
	d := &dependencies{
		bucketName:     os.Getenv("bucketName"),
		typlessToken:   os.Getenv("typlessToken"),
		typlessDocType: os.Getenv("typlessDocType"),
		tableName:      os.Getenv("tableName"),
		s3Client:       aws_clients.GetS3Client(),
		dbClient:       aws_clients.GetDbClient(),
		lambdaTimeout:  os.Getenv("lambdaTimeout"),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		filename := record.Body

		// get the invoice from s3
		s3Invoice, s3GetError := d.getS3Invoice(filename)
		if s3GetError != nil {
			utils.LogError("Failed to get s3 invoice", s3GetError)
			return s3GetError
		}

		originalFilename := s3Invoice.Metadata["OriginalFilename"]

		invoicePdf, readError := ioutil.ReadAll(s3Invoice.Body)
		if readError != nil {
			utils.LogError("Failed to read invoice body", readError)
			return readError
		}

		typlessConfig := &typless.Config{
			Token:   d.typlessToken,
			DocType: d.typlessDocType,
		}

		lambdaTimeout, atoiError := strconv.Atoi(d.lambdaTimeout)
		if atoiError != nil {
			utils.LogError("Failed convert lambda timeout", atoiError)
			return atoiError
		}

		// to make sure we close http connection before lambda times out
		safetyTimeout := lambdaTimeout - 5

		// extract the text from the invoice
		extractedData, extractError := typless.ExtractData(typlessConfig, invoicePdf, safetyTimeout)

		if extractError != nil {
			utils.LogError("Failed to extract data", extractError)
			return extractError
		}

		if !isInvoiceWithWireTransfer(extractedData) {
			utils.Log("Not an invoice with wire transfer, skipping")
			_, deleteErr := d.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: &d.bucketName,
				Key:    &filename,
			})
			if deleteErr != nil {
				utils.LogError("Failed to delete s3 object", deleteErr)
				return deleteErr
			}
			return nil
		}

		createInvoiceService := create_invoice.CreateInvoiceService{
			OriginalFilename: originalFilename,
			ExtractedData:    extractedData,
		}
		invoice := createInvoiceService.CreateInvoice()
		invoice.Filename = filename

		_, saveInvoiceError := db.PutInvoice(d.dbClient, d.tableName, invoice)

		if saveInvoiceError != nil {
			utils.LogError("Failed to save invoice to db", saveInvoiceError)
			return saveInvoiceError
		}
	}

	return nil
}

func (d *dependencies) getS3Invoice(filename string) (*s3.GetObjectOutput, error) {
	invoice, s3GetError := d.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &d.bucketName,
		Key:    &filename,
	})
	if s3GetError != nil {
		utils.LogError("Failed to get invoice from s3", s3GetError)
		return nil, s3GetError
	}

	return invoice, s3GetError
}

func isInvoiceWithWireTransfer(extractedData *typless.ExtractDataFromFileOutput) bool {
	isWireTransfer := false
	indicators := []string{"átutalás", "wire transfer"}

textBlockLoop:
	for _, textBlock := range extractedData.TextBlocks {
		for _, indicator := range indicators {
			if strings.Contains(strings.ToLower(textBlock.Value), indicator) {
				isWireTransfer = true
				break textBlockLoop
			}
		}
	}

	return isWireTransfer
}
