package guesser_functions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGuesserFunctions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GuesserFunctions Suite")
}
