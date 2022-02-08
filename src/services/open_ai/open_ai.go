package open_ai

import (
	"context"
	"fmt"
	"pocok/src/services/typless"
	"pocok/src/utils"

	"github.com/PullRequestInc/go-gpt3"
)

func GetPredictedInvoiceData(textBlocks []typless.TextBlock) (string, error) {
	apiKey := ""
	client := gpt3.NewClient(apiKey)

	Prompt := []string{"Raw data from an invoice.\n" + fmt.Sprintln(textBlocks) + ".\n InvoiceNumber | VendorName | AccountNumber | Iban | Currency | GrossPrice | DueDate."}
	Temperature := float32(0.0)
	MaxTokens := 64

	resp, err := client.Completion(context.TODO(), gpt3.CompletionRequest{
		Prompt:      Prompt,
		Stop:        []string{"."},
		Temperature: &Temperature,
		MaxTokens:   &MaxTokens,
	})
	if err != nil {
		utils.LogError("error getting prediction from openai", err)
		return "", err
	}

	return resp.Choices[0].Text, nil
}
