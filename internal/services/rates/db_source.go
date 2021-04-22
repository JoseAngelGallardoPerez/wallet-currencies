package rates

import (
	"fmt"

	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/repositories"
)

type DbSource struct {
	ratesRepository      *repositories.Rates
	currenciesRepository *repositories.Currencies
}

func NewDbSource(ratesRepository *repositories.Rates, currenciesRepository *repositories.Currencies) RateAndMarginSource {
	return &DbSource{ratesRepository: ratesRepository, currenciesRepository: currenciesRepository}
}

// FindRateAndMargin returns rate and margin for target currencies
func (s *DbSource) FindRateAndMargin(currencyCodeFrom, currencyCodeTo string) (*RateAndMargin, error) {
	rate, err := s.findByCurrencyCodes(currencyCodeFrom, currencyCodeTo)
	if err != nil {
		return nil, err
	}

	res := &RateAndMargin{
		Base:           currencyCodeFrom,
		Reference:      currencyCodeTo,
		Rate:           rate.Value,
		ExchangeMargin: rate.ExchangeMargin,
	}

	return res, nil
}

// findByCurrencyCodes returns rates model from DB
func (s *DbSource) findByCurrencyCodes(currencyCodeFrom, currencyCodeTo string) (*models.Rate, error) {
	rate, err := s.ratesRepository.FindByCurrencies(currencyCodeFrom, currencyCodeTo)

	if err != nil {
		return nil, err
	}

	if rate == nil || !rate.IsExists() {
		return nil, fmt.Errorf("rate not found: %s -> %s", currencyCodeFrom, currencyCodeTo)
	}

	return rate, nil
}
