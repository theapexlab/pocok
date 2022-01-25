package currency_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCurrency(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Currency Suite")
}
