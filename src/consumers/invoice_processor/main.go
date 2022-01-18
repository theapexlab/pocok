package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"pocok/src/consumers/invoice_processor/create_invoice"
	"pocok/src/services/typless"
	"pocok/src/utils/aws_clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type dependencies struct {
	bucketName   string
	typlessToken string
	s3Client     *s3.Client
	dbClient     *dynamodb.Client
}

func main() {
	d := &dependencies{
		bucketName:   os.Getenv("bucketName"),
		typlessToken: os.Getenv("typlessToken"),
		s3Client:     aws_clients.GetS3Client(),
		dbClient:     aws_clients.GetDbClient(),
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

		// extract the text from the invoice
		extractedData, err := typless.ExtractData(invoicePdf, "en-invoice", d.typlessToken)
		if err != nil {
			return err
		}

		invoice, err := create_invoice.CreateInvoice(extractedData)
		if err != nil {
			return err
		}

		fmt.Println(invoice)
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
