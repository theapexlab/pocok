package url_parsing_strategies

import (
	"errors"
	"strings"

	"github.com/DusanKasan/parsemail"
)

type urlParsingStrategy interface {
	Parse(email *parsemail.Email) (string, error)
}

var ErrNoUrlParsingStrategyFound = errors.New("no url parsing strategy found")

func GetPdfUrlFromEmail(email *parsemail.Email) (string, error) {
	strategy, getStrategyError := getStrategy(email)
	if getStrategyError != nil {
		return "", getStrategyError
	}

	return strategy.Parse(email)
}

func getStrategy(email *parsemail.Email) (urlParsingStrategy, error) {
	if len(email.From) == 0 {
		return nil, ErrNoUrlParsingStrategyFound
	}

	if email.From[0].Address == BillingoAddress {
		return &Billingo{}, nil
	}

	if strings.Contains(email.From[0].Address, SzamlazzAddress) {
		return &Szamlazz{}, nil
	}

	return nil, ErrNoUrlParsingStrategyFound
}
