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
		invoicePdf, err := d.getInvoicePdf(filename)
		if err != nil {
			return err
		}

		typlessConfig := &typless.Config{
			Token:   d.typlessToken,
			DocType: d.typlessDocType,
		}

		lambdaTimeout, atioErr := strconv.Atoi(d.lambdaTimeout)
		if atioErr != nil {
			return atioErr
		}

		// to make sure we close http connection before lambda times out
		safetyTimeout := lambdaTimeout - 5

		// extract the text from the invoice
		extractedData, err := typless.ExtractData(typlessConfig, invoicePdf, safetyTimeout)

		if err != nil {
			utils.LogError("Failed to extract data", err)
			return err
		}

		invoice := create_invoice.CreateInvoice(extractedData)
		invoice.Filename = filename

		_, dbErr := db.PutInvoice(d.dbClient, d.tableName, invoice)

		if dbErr != nil {
			utils.LogError("Failed to save invoice to db", dbErr)
			return dbErr
		}

	}

	return nil
}

func (d *dependencies) getInvoicePdf(filename string) ([]byte, error) {
	invoice, err := d.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &d.bucketName,
		Key:    &filename,
	})
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(invoice.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
