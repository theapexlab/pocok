package create_email_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreateEmail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateEmail Suite")
}
