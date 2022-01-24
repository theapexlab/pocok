package db

import (
	"context"
	"pocok/src/utils/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/segmentio/ksuid"
)

func PutInvoice(client *dynamodb.Client, tableName string, filename string) (*dynamodb.PutItemOutput, error) {

	id := ksuid.New().String()
	status := models.PENDING
	reveicedAt := time.Now().Format(time.RFC3339)

	invoice := models.Invoice{
		Pk:         models.ORG + "#" + models.APEX_ID,
		Sk:         models.INVOICE + "#" + id,
		Lsi1sk:     models.STATUS + "#" + status,
		Lsi2sk:     models.VENDOR + "#unknown#" + models.DATE + "#" + reveicedAt,
		ReceivedAt: reveicedAt,
		InvoiceId:  id,
		EntityType: models.INVOICE,
		Status:     status,
		Filename:   filename,
	}

	item, itemErr := attributevalue.MarshalMap(invoice)
	if itemErr != nil {
		return nil, itemErr
	}

	dbResp, dbErr := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return dbResp, dbErr
}
