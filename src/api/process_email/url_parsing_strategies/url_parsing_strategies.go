package url_parsing_strategies

import (
	"errors"
	"pocok/src/utils/models"
)

type urlParsingStrategy interface {
	Parse(jsonBody *models.EmailWebhookBody) (string, error)
}

var ErrNoUrlParsingStrategyFound = errors.New("no url parsing strategy found")

func GetPdfUrlFromEmail(jsonBody *models.EmailWebhookBody) (string, error) {
	strategy, getStrategyError := getStrategy(jsonBody)
	if getStrategyError != nil {
		return "", getStrategyError
	}

	return strategy.Parse(jsonBody)
}

func getStrategy(jsonBody *models.EmailWebhookBody) (urlParsingStrategy, error) {
	if len(jsonBody.From) == 0 {
		return nil, ErrNoUrlParsingStrategyFound
	}

	if jsonBody.From[0].Address == BillingoAddress {
		return &Billingo{}, nil
	}

	return nil, ErrNoUrlParsingStrategyFound
}
