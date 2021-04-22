package ratesupdater

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-currencies/internal/services/currencies"
	"github.com/Confialink/wallet-currencies/internal/services/exchange"
)

type Service struct {
	currenciesService *currencies.Service
	providerFactory   *ProviderFactory
	rateReceiver      exchange.RateReceiver
}

func NewService(
	currenciesService *currencies.Service,
	providerFactory *ProviderFactory,
	rateReceiver exchange.RateReceiver,
) *Service {
	return &Service{currenciesService, providerFactory, rateReceiver}
}

func (s Service) WrapContext(db *gorm.DB) *Service {
	s.currenciesService = s.currenciesService.WrapContext(db)
	return &s
}

// Call updates rates in DB.
// Source of rates depends on main currency feed.
// It updates rates for all currencies with the same feed.
func (s *Service) Call() error {
	mainCurrency, err := s.currenciesService.Main()
	if err != nil {
		return err
	}

	if !mainCurrency.Feed.Valid || len(mainCurrency.Feed.String) == 0 {
		return errors.New("The main currency has empty feed")
	}

	provider, err := s.providerFactory.Source(mainCurrency.Feed.String)
	if err != nil {
		return err
	}

	// we update currencies which have the same feed
	targetCurrencies, err := s.currenciesService.FindByFeed(mainCurrency.Feed.String)
	if err != nil {
		return errors.Wrap(err, "cannot obtain rates by feed")
	}
	for _, currency := range targetCurrencies {
		if mainCurrency.Code == currency.Code {
			continue
		}
		rate, err := provider.FindRate(mainCurrency.Code, currency.Code)
		if err != nil {
			if err == exchange.ErrRateNotFound {
				continue
			}
			return err
		}

		if err = s.rateReceiver.Set(rate); err != nil {
			return err
		}
	}

	return nil
}
