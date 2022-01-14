package url_parsing_strategies_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUrlParsingStrategies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UrlParsingStrategies Suite")
}
