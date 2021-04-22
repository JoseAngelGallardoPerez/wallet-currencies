package ratesupdater

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Confialink/wallet-currencies/internal/services/exchange/mocks"
	ratesMock "github.com/Confialink/wallet-currencies/internal/services/ratesupdater/mocks"
)

var _ = Describe("ProviderFactory", func() {
	randomFeedName := "feed_name"
	Context("Add", func() {
		When("provider is successfully added", func() {
			It("does not return an error", func() {
				factory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				err := factory.Add(mockSourceFactory, randomFeedName)
				Expect(err).ShouldNot(HaveOccurred())
			})

		})

		When("provider is already added", func() {
			It("returns an error", func() {
				factory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				_ = factory.Add(mockSourceFactory, randomFeedName)
				err := factory.Add(mockSourceFactory, randomFeedName)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already registered"))
			})
		})
	})

	Context("Source", func() {
		When("there is unknown provider", func() {
			It("returns an error", func() {
				factory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				_ = factory.Add(mockSourceFactory, randomFeedName)
				_, err := factory.Source("another_feed_name")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown provider"))
			})
		})

		When("rateSource is exists", func() {
			It("does not return an error", func() {
				factory := NewProviderFactory()
				mockSource := &mocks.RateSource{}
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				mockSourceFactory.On("Init").Return(mockSource, nil)
				_ = factory.Add(mockSourceFactory, randomFeedName)
				source, err := factory.Source(randomFeedName)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(source).To(Equal(mockSource))
			})
		})
	})
})
