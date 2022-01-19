package create_email_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInvoiceSummary(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Invoice Summary Suite")
}
