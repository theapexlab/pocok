package open_ai

import (
	"context"
	"fmt"
	"pocok/src/utils"

	"github.com/PullRequestInc/go-gpt3"
)

func GetPredictedInvoiceData(textBlocks []string) (string, error) {
	apiKey := "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"
	client := gpt3.NewClient(apiKey)

	resp, err := client.Completion(context.TODO(), gpt3.CompletionRequest{
		Prompt: []string{"Raw data from an invoice.\n", fmt.Sprintln(textBlocks), "\nInvoiceNumber | VendorName | AccountNumber | Iban | Currency | GrossPrice | DueDate ."},
	})
	if err != nil {
		utils.LogError("error getting prediction from openai", err)
		return "", err
	}

	return resp.Choices[0].Text, nil
}
