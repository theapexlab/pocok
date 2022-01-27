package create_training_data_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreateTrainingData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateTrainingData Suite")
}
