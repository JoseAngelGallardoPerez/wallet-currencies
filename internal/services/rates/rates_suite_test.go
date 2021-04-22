package rates

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rates Suite")
}
