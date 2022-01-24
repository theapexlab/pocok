package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"pocok/src/consumers/invoice_processor/create_invoice"
	"pocok/src/db"
	"pocok/src/services/typless"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"

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

		// extract the text from the invoice
		// todo: in some cases this func makes lambe function disconnect
		//  todo: Failed to send response because the Lambda function is disconnected
		extractedData, err := typless.ExtractData(invoicePdf, typlessConfig)
		fmt.Println(extractedData) // todo: remove this line
		if err != nil {
			utils.LogError("Failed to extract data", err)
			return err
		}

		invoice := create_invoice.CreateInvoice(extractedData)

		// todo: add VendorEmail
		invoice.Filename = filename

		fmt.Println(invoice) // todo: remove this line
		putOutput, dbErr := db.PutInvoice(d.dbClient, d.tableName, invoice)
		fmt.Println(putOutput.Attributes) // todo: remove this line

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
