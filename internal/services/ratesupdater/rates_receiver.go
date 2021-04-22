package ratesupdater

import (
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-currencies/internal/repositories"
	"github.com/Confialink/wallet-currencies/internal/services/exchange"
)

type DbRatesReceiver struct {
	ratesRepository *repositories.Rates
}

func NewDbRatesReceiver(ratesRepository *repositories.Rates) exchange.RateReceiver {
	return &DbRatesReceiver{ratesRepository: ratesRepository}
}

// Set updates rate in the DB
func (s *DbRatesReceiver) Set(rate exchange.Rate) (err error) {
	dbRate, err := s.ratesRepository.FindByCurrencies(rate.BaseCurrencyCode(), rate.ReferenceCurrencyCode())
	if err != nil {
		return errors.Wrap(err, "cannot obtain rates from DB")
	}
	if dbRate.IsExists() {
		dbRate.Value = rate.Rate()
		_, err = s.ratesRepository.UpdateRate(dbRate)
		if err != nil {
			return errors.Wrap(err, "cannot update rates in DB")
		}
	}

	return nil
}
