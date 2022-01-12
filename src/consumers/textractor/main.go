package main

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

var textractSession *textract.Textract

func init() {
	textractSession = textract.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"), // Frankfurt
	})))
}

func main() {
	file, err := ioutil.ReadFile("example-invoice.pdf")
	if err != nil {
		panic(err)
	}

	// fmt.Print(file)

	// featureTypes := [...]string{"FORMS"}

	var strs []string = []string{"FORMS"}

	resp, err := textractSession.AnalyzeDocument(&textract.AnalyzeDocumentInput{
		Document: &textract.Document{
			Bytes: file,
		},
		FeatureTypes: aws.StringSlice(strs),
	})

	if err != nil {
		panic(err)
	}g

	fmt.Println(resp)
}
