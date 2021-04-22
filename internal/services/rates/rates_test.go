package rates

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"

	"github.com/Confialink/wallet-currencies/internal/repositories"
)

var _ = Describe("Rates", func() {

	Context("CalculateRate", func() {
		defaultRate, _ := decimal.NewFromString(defaultRateValue)
		dbConnect := func() (sqlmock.Sqlmock, *gorm.DB) {
			db, dbMock, _ := sqlmock.New()
			gormDb, _ := gorm.Open("mysql", db)

			return dbMock, gormDb
		}

		selectCurrencyQuery := "SELECT (.+) FROM `currencies`(.+)"

		When("base currency is equal reference currency", func() {
			It("should return correct rate", func() {
				source := &RateAndMarginSourceMock{}
				_, gormDb := dbConnect()
				currenciesRepo := repositories.NewCurrencies(gormDb)
				service := NewRates(source, currenciesRepo)
				currency := "EUR"

				rate, err := service.CalculateRate(currency, currency)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rate.Base).Should(Equal(currency))
				Expect(rate.Reference).Should(Equal(currency))
				Expect(rate.ExchangeMargin).Should(Equal(decimal.Zero))
				Expect(rate.Rate).Should(Equal(defaultRate))
			})
		})

		When("the system cannot find main currency", func() {
			It("returns an error", func() {
				source := &RateAndMarginSourceMock{}
				dbMock, gormDb := dbConnect()
				dbErr := errors.New("random_db_error")
				dbMock.ExpectQuery(selectCurrencyQuery).WillReturnError(dbErr)
				currenciesRepo := repositories.NewCurrencies(gormDb)
				service := NewRates(source, currenciesRepo)
				_, err := service.CalculateRate("EUR", "USD")
				Expect(err.Error()).To(ContainSubstring(dbErr.Error()))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		Context("find direct rate", func() {
			baseCurrency := "EUR"
			referenceCurrency := "USD"
			mainCurrency := baseCurrency
			When("rate source returns error", func() {
				It("should return the same error", func() {
					source := &RateAndMarginSourceMock{}
					mockErr := errors.New("random_source_errors")
					source.On("FindRateAndMargin", baseCurrency, referenceCurrency).Return(nil, mockErr)
					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					_, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).To(Equal(mockErr))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})

			When("rate source returns correct rate", func() {
				It("should return the same correct rate", func() {
					rateMock := &RateAndMargin{
						Base:           baseCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.NewFromInt(1),
						ExchangeMargin: decimal.NewFromInt(2),
					}
					source := &RateAndMarginSourceMock{}
					source.On("FindRateAndMargin", baseCurrency, referenceCurrency).Return(rateMock, nil)
					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					rate, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).Should(BeNil())
					Expect(rate.ExchangeMargin).To(Equal(rateMock.ExchangeMargin))
					Expect(rate.Rate).To(Equal(rateMock.Rate))
					Expect(rate.Base).To(Equal(rateMock.Base))
					Expect(rate.Reference).To(Equal(rateMock.Reference))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})
		})

		Context("find reverse rate", func() {
			baseCurrency := "USD"
			referenceCurrency := "EUR"
			mainCurrency := referenceCurrency
			When("rate source returns error", func() {
				It("should return the same error", func() {
					source := &RateAndMarginSourceMock{}
					mockErr := errors.New("random_source_errors")
					source.On("FindRateAndMargin", referenceCurrency, baseCurrency).Return(nil, mockErr)
					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					_, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).To(Equal(mockErr))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})

			When("rate source returns correct rate", func() {
				It("should return the same correct rate", func() {
					rateMock := &RateAndMargin{
						Base:           baseCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.NewFromInt(2),
						ExchangeMargin: decimal.NewFromInt(4),
					}
					source := &RateAndMarginSourceMock{}
					source.On("FindRateAndMargin", referenceCurrency, baseCurrency).Return(rateMock, nil)
					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					rate, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).Should(BeNil())
					Expect(rate.Base).To(Equal(baseCurrency))
					Expect(rate.Reference).To(Equal(referenceCurrency))
					Expect(rate.Rate.String()).To(Equal(decimal.NewFromFloat(0.5).String()))
					Expect(rate.ExchangeMargin.String()).To(Equal(rateMock.ExchangeMargin.String()))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})
			When("rate source returns correct rate with zero values", func() {
				It("should return the same correct rate", func() {
					rateMock := &RateAndMargin{
						Base:           baseCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.Zero,
						ExchangeMargin: decimal.NewFromInt(5),
					}
					source := &RateAndMarginSourceMock{}
					source.On("FindRateAndMargin", referenceCurrency, baseCurrency).Return(rateMock, nil)
					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					rate, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).Should(BeNil())
					Expect(rate.Base).To(Equal(baseCurrency))
					Expect(rate.Reference).To(Equal(referenceCurrency))
					Expect(rate.Rate.String()).To(Equal(decimal.Zero.String()))
					Expect(rate.ExchangeMargin.String()).To(Equal(decimal.Zero.String()))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})
		})

		Context("find pivot rate", func() {
			baseCurrency := "USD"
			referenceCurrency := "AUD"
			mainCurrency := "EUR"
			When("rate source returns error during obtain first rate", func() {
				It("should return the same error", func() {
					source := &RateAndMarginSourceMock{}
					mockErr := errors.New("random_source_errors")
					source.On("FindRateAndMargin", mainCurrency, baseCurrency).Return(nil, mockErr)
					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					_, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).To(Equal(mockErr))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})

			When("rate source returns error during obtain second rate", func() {
				It("should return the same error", func() {
					source := &RateAndMarginSourceMock{}
					rateMockFirst := &RateAndMargin{
						Base:           baseCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.NewFromInt(2),
						ExchangeMargin: decimal.NewFromInt(4),
					}
					source.On("FindRateAndMargin", mainCurrency, baseCurrency).Return(rateMockFirst, nil)
					mockErr := errors.New("random_source_errors")
					source.On("FindRateAndMargin", mainCurrency, referenceCurrency).Return(nil, mockErr)

					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					_, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).To(Equal(mockErr))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})

			When("rate source returns correct rate", func() {
				It("should return the same correct rate", func() {
					source := &RateAndMarginSourceMock{}
					rateMockFirst := &RateAndMargin{
						Base:           mainCurrency,
						Reference:      baseCurrency,
						Rate:           decimal.NewFromInt(3),
						ExchangeMargin: decimal.NewFromInt(2),
					}
					rateMockSecond := &RateAndMargin{
						Base:           mainCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.NewFromInt(6),
						ExchangeMargin: decimal.NewFromInt(8),
					}
					source.On("FindRateAndMargin", mainCurrency, baseCurrency).Return(rateMockFirst, nil)
					source.On("FindRateAndMargin", mainCurrency, referenceCurrency).Return(rateMockSecond, nil)

					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					rate, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).Should(BeNil())
					Expect(rate.Base).To(Equal(baseCurrency))
					Expect(rate.Reference).To(Equal(referenceCurrency))
					Expect(rate.Rate.String()).To(Equal(decimal.NewFromFloat(2).String()))
					Expect(rate.ExchangeMargin.String()).To(Equal(rateMockFirst.ExchangeMargin.String()))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})

			When("rate source returns first rate with zero values", func() {
				It("should return the same correct rate with zero values", func() {
					source := &RateAndMarginSourceMock{}
					rateMockFirst := &RateAndMargin{
						Base:           mainCurrency,
						Reference:      baseCurrency,
						Rate:           decimal.NewFromInt(0),
						ExchangeMargin: decimal.NewFromInt(7),
					}
					rateMockSecond := &RateAndMargin{
						Base:           mainCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.NewFromInt(6),
						ExchangeMargin: decimal.NewFromInt(8),
					}
					source.On("FindRateAndMargin", mainCurrency, baseCurrency).Return(rateMockFirst, nil)
					source.On("FindRateAndMargin", mainCurrency, referenceCurrency).Return(rateMockSecond, nil)

					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					rate, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).Should(BeNil())
					Expect(rate.Base).To(Equal(baseCurrency))
					Expect(rate.Reference).To(Equal(referenceCurrency))
					Expect(rate.Rate.String()).To(Equal(decimal.Zero.String()))
					Expect(rate.ExchangeMargin.String()).To(Equal(decimal.Zero.String()))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})

			When("rate source returns first rate with zero values", func() {
				It("should return the same correct rate with zero values", func() {
					source := &RateAndMarginSourceMock{}
					rateMockFirst := &RateAndMargin{
						Base:           mainCurrency,
						Reference:      baseCurrency,
						Rate:           decimal.NewFromInt(3),
						ExchangeMargin: decimal.NewFromInt(2),
					}
					rateMockSecond := &RateAndMargin{
						Base:           mainCurrency,
						Reference:      referenceCurrency,
						Rate:           decimal.NewFromInt(0),
						ExchangeMargin: decimal.NewFromInt(5),
					}
					source.On("FindRateAndMargin", mainCurrency, baseCurrency).Return(rateMockFirst, nil)
					source.On("FindRateAndMargin", mainCurrency, referenceCurrency).Return(rateMockSecond, nil)

					dbMock, gormDb := dbConnect()
					mainCurrencySqlRows := sqlmock.NewRows([]string{"code"}).AddRow(mainCurrency)
					dbMock.ExpectQuery(selectCurrencyQuery).WillReturnRows(mainCurrencySqlRows)
					currenciesRepo := repositories.NewCurrencies(gormDb)
					service := NewRates(source, currenciesRepo)
					rate, err := service.CalculateRate(baseCurrency, referenceCurrency)

					Expect(err).Should(BeNil())
					Expect(rate.Base).To(Equal(baseCurrency))
					Expect(rate.Reference).To(Equal(referenceCurrency))
					Expect(rate.Rate.String()).To(Equal(decimal.Zero.String()))
					Expect(rate.ExchangeMargin.String()).To(Equal(decimal.Zero.String()))
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})
		})
	})
})
