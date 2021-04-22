package rates

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/Confialink/wallet-currencies/internal/repositories"
)

const defaultRateValue string = "1"

type RateAndMargin struct {
	Base           string          `json:"baseCurrencyCode"`
	Reference      string          `json:"referenceCurrencyCode"`
	Rate           decimal.Decimal `json:"rateValue"`
	ExchangeMargin decimal.Decimal `json:"exchangeMargin"`
}

// RateAndMarginSource
type RateAndMarginSource interface {
	FindRateAndMargin(base, reference string) (*RateAndMargin, error)
}

type Rates struct {
	source               RateAndMarginSource
	currenciesRepository *repositories.Currencies
}

func NewRates(source RateAndMarginSource, currenciesRepository *repositories.Currencies) *Rates {
	return &Rates{source, currenciesRepository}
}

// CalculateRate returns information about currency rate and exchange margin for provided currencies
func (s *Rates) CalculateRate(currencyCodeFrom, currencyCodeTo string) (*RateAndMargin, error) {
	res := &RateAndMargin{
		Base:           currencyCodeFrom,
		Reference:      currencyCodeTo,
		ExchangeMargin: decimal.Zero,
		Rate:           decimal.Zero,
	}

	if currencyCodeFrom == currencyCodeTo {
		defaultRate, _ := decimal.NewFromString(defaultRateValue)
		res.Rate = defaultRate

		return res, nil
	}

	mainCurrency, err := s.currenciesRepository.Main()
	if err != nil {
		return nil, errors.Wrap(err, "cannot obtain main currency")
	}
	mainCurrencyCode := mainCurrency.Code

	// Find direct rate
	if currencyCodeFrom == mainCurrencyCode {
		return s.source.FindRateAndMargin(currencyCodeFrom, currencyCodeTo)
	}

	// Find reverse rate
	// It tries to find rate using the given source.
	// In case if no rate found it tries to find reverse rate and calculate required rate with it.
	// a/b = x -> b/a = 1/x
	if currencyCodeTo == mainCurrencyCode {
		rate, err := s.source.FindRateAndMargin(currencyCodeTo, currencyCodeFrom)
		if err != nil {
			return nil, err
		}
		if !rate.Rate.IsZero() {
			res.Rate = decimal.NewFromFloat(1).Div(rate.Rate)
			res.ExchangeMargin = rate.ExchangeMargin
		}

		return res, nil
	}

	// It is used in order to calculate rate with "pivot" currency.
	// For example("a" is pivot): a/b = x; a/c = y; b/c = x/y
	rateFrom, err := s.source.FindRateAndMargin(mainCurrencyCode, currencyCodeFrom)
	if err != nil {
		return nil, err
	}

	rateTo, err := s.source.FindRateAndMargin(mainCurrencyCode, currencyCodeTo)
	if err != nil {
		return nil, err
	}

	if !rateTo.Rate.IsZero() && !rateFrom.Rate.IsZero() {
		res.Rate = rateTo.Rate.Div(rateFrom.Rate)
		res.ExchangeMargin = rateFrom.ExchangeMargin
	}

	return res, nil
}
