package ratesupdater

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/Confialink/wallet-currencies/internal/mocks/log15"
	"github.com/Confialink/wallet-currencies/internal/repositories"
	"github.com/Confialink/wallet-currencies/internal/services/accounts"
	"github.com/Confialink/wallet-currencies/internal/services/currencies"
	"github.com/Confialink/wallet-currencies/internal/services/exchange"
	"github.com/Confialink/wallet-currencies/internal/services/exchange/mocks"
	ratesMock "github.com/Confialink/wallet-currencies/internal/services/ratesupdater/mocks"
)

var _ = Describe("RateUpdater", func() {
	Context("Call", func() {
		dbConnect := func() (sqlmock.Sqlmock, *gorm.DB) {
			db, dbMock, _ := sqlmock.New()
			gormDb, _ := gorm.Open("mysql", db)

			return dbMock, gormDb
		}
		selectCurrencyQuery := "SELECT (.+) FROM `currencies`(.+)"
		logger := &log15.Logger{}
		logger.On("New", mock.Anything, mock.Anything).Return(logger)

		When("main currency not found", func() {
			It("returns an error", func() {
				dbMock, gormDb := dbConnect()
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnError(errors.New("random text error"))
				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				service := NewService(currenciesService, NewProviderFactory(), ratesReceiver)

				Expect(service.Call()).Should(HaveOccurred())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("main currency does not have a feed", func() {
			It("returns an error", func() {
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", "")
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				service := NewService(currenciesService, NewProviderFactory(), ratesReceiver)

				err := service.Call()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("The main currency has empty feed"))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("rates provider returns an error", func() {
			It("returns an error", func() {
				randomFeedName := "random_feed_name"
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				service := NewService(currenciesService, NewProviderFactory(), ratesReceiver)

				Expect(service.Call()).Should(HaveOccurred())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("currencies service cannot find currencies with the same feed", func() {
			It("returns an error", func() {
				randomFeedName := "random_feed_name"
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				dbError := errors.New("random db error")
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnError(dbError)

				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				providerFactory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				mockSource := &mocks.RateSource{}
				mockSourceFactory.On("Init").Return(mockSource, nil)

				_ = providerFactory.Add(mockSourceFactory, randomFeedName)
				service := NewService(currenciesService, providerFactory, ratesReceiver)

				err := service.Call()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(dbError.Error()))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("rates provider cannot find a rate", func() {
			It("returns an error", func() {
				randomFeedName := "random_feed_name"
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				currenciesSqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("AUD", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(currenciesSqlRows)

				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				providerFactory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				mockSource := &mocks.RateSource{}
				mockErr := errors.New("random provider err")
				mockSource.On("FindRate", mock.Anything, mock.Anything).Return(exchange.Rate{}, mockErr)
				mockSourceFactory.On("Init").Return(mockSource, nil)

				_ = providerFactory.Add(mockSourceFactory, randomFeedName)
				service := NewService(currenciesService, providerFactory, ratesReceiver)

				err := service.Call()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(mockErr))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())

			})
		})

		When("rates receiver cannot set a rate", func() {
			It("returns an error", func() {
				randomFeedName := "random_feed_name"
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				currenciesSqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("AUD", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(currenciesSqlRows)

				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				mockErr := errors.New("random source err")
				ratesReceiver.On("Set", mock.Anything).Return(mockErr)
				providerFactory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				mockSource := &mocks.RateSource{}
				mockSource.On("FindRate", mock.Anything, mock.Anything).Return(exchange.Rate{}, nil)
				mockSourceFactory.On("Init").Return(mockSource, nil)

				_ = providerFactory.Add(mockSourceFactory, randomFeedName)
				service := NewService(currenciesService, providerFactory, ratesReceiver)

				err := service.Call()
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(mockErr))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("rates provider does not contain a rate", func() {
			It("does not return an error.", func() {
				randomFeedName := "random_feed_name"
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				currenciesSqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("AUD", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(currenciesSqlRows)

				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}

				providerFactory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				mockSource := &mocks.RateSource{}
				mockSource.On("FindRate", mock.Anything, mock.Anything).Return(exchange.Rate{}, exchange.ErrRateNotFound)
				mockSourceFactory.On("Init").Return(mockSource, nil)

				_ = providerFactory.Add(mockSourceFactory, randomFeedName)
				service := NewService(currenciesService, providerFactory, ratesReceiver)

				Expect(service.Call()).ShouldNot(HaveOccurred())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("everything is OK", func() {
			It("does not return an error", func() {
				randomFeedName := "random_feed_name"
				dbMock, gormDb := dbConnect()
				mainCurrencySqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("EUR", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
				currenciesSqlRows := sqlmock.NewRows([]string{"code", "feed"}).AddRow("AUD", randomFeedName)
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(currenciesSqlRows)

				currenciesRepo := repositories.NewCurrencies(gormDb)
				ratesRepo := repositories.NewRates(gormDb)
				currenciesService := currencies.NewService(currenciesRepo, accounts.NewService(logger), ratesRepo)
				ratesReceiver := &mocks.RateReceiver{}
				ratesReceiver.On("Set", mock.Anything).Return(nil)
				providerFactory := NewProviderFactory()
				mockSourceFactory := &ratesMock.CurrencySourceFactory{}
				mockSource := &mocks.RateSource{}
				mockSource.On("FindRate", mock.Anything, mock.Anything).Return(exchange.Rate{}, nil)
				mockSourceFactory.On("Init").Return(mockSource, nil)

				_ = providerFactory.Add(mockSourceFactory, randomFeedName)
				service := NewService(currenciesService, providerFactory, ratesReceiver)

				Expect(service.Call()).ShouldNot(HaveOccurred())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})
	})
})
