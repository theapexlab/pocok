package open_ai_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOpenAi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OpenAi Suite")
}
