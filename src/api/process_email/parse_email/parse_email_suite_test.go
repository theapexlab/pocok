package parse_email_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestParseEmail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ParseEmail Suite")
}
