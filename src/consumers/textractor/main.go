package main

import (
	"os"
	"pocok/src/utils/aws_clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

type dependencies struct {
	bucketName     string
	s3Client       *s3.Client
	textractClient *textract.Client
	dbClient       *dynamodb.Client
}

func ParsePdf(pdfS3Object s3.GetObjectOutput) {
	// if textractSession == nil {
	// 	textractSession = textract.New(session.Must(session.NewSession(&aws.Config{
	// 		Region: aws.String("eu-central-1"), // Frankfurt
	// 	})))
	// }

	// file, err := ioutil.ReadFile("example-invoice.pdf")
	// if err != nil {
	// 	panic(err)
	// }

	// strs := []string{"FORMS"}

	// resp, err := textractSession.AnalyzeDocument(&textract.AnalyzeDocumentInput{
	// 	Document: &textract.Document{
	// 		// Bytes: file,
	// 		S3Object: pdfS3Object,
	// 	},
	// 	FeatureTypes: aws.StringSlice(strs),
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(resp)
}

func (d *dependencies) handler(event events.SQSEvent) error {
	for _, record := range event.Records {
		filename := record.Body
	}

	return nil
}

func main() {
	d := &dependencies{
		bucketName:     os.Getenv("bucketName"),
		s3Client:       aws_clients.GetS3Client(),
		textractClient: aws_clients.GetTextractClient(),
		dbClient:       aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}
