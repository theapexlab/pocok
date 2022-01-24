package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"

	"net/http"
	"os"
	"pocok/src/utils"
	"pocok/src/utils/aws_clients"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dependencies struct {
	dbClient  *dynamodb.Client
	tableName string
}

func main() {
	d := &dependencies{
		tableName: os.Getenv("tableName"),
		dbClient:  aws_clients.GetDbClient(),
	}

	lambda.Start(d.handler)
}

func (d *dependencies) handler(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	n, err := base64.StdEncoding.DecodeString(r.Body)
	if err != nil {
		log.Fatal("Error decoding base64 message content", err)
	}

	fmt.Println("------------------------")
	// fmt.Println(string(n))
	reader := bytes.NewReader(n)
	contentType := r.Headers["content-type"]

	mediaType, params, err := mime.ParseMediaType(contentType)
	fmt.Println("-----------mediaType-------------")
	fmt.Println(mediaType)
	fmt.Println("-----------params-------------")
	fmt.Println(params)
	fmt.Println("------------------------")

	if err != nil {
		log.Fatal(err)
	}

	mr := multipart.NewReader(reader, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err != nil {
			log.Fatal(err)
		}
		slurp, err := io.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(p.Header)
		fmt.Printf("Part %q: %q\n", p.Header.Get("Content-Disposition"), slurp)
	}

	// fmt.Println("------------------------")

	// token := r.QueryStringParameters["token"]
	// claims, err := auth.ParseToken(token)
	// if err != nil {
	// 	return utils.MailApiResponse(http.StatusUnauthorized, ""), err
	// }

	// id := ""
	// status := ""
	// updateErr := db.UpdateInvoiceStatus(client, tableName, claims.OrgId, id, status)
	// if updateErr != nil {
	// 	return utils.MailApiResponse(http.StatusInternalServerError, ""), nil
	// }

	return utils.MailApiResponse(http.StatusOK, ""), nil
}
