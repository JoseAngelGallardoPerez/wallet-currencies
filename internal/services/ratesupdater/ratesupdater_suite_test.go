package ratesupdater_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rates Updater Suite")
}
