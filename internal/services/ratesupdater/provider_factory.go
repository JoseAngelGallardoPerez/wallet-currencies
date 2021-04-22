package ratesupdater

import (
	"fmt"

	"github.com/Confialink/wallet-currencies/internal/services/exchange"
)

const (
	EcbProvider = "ECB"
)

type CurrencySourceFactory interface {
	// returns a prepared rate source
	Init() (exchange.RateSource, error)
}

type ProviderFactory struct {
	factories map[string]CurrencySourceFactory
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{factories: make(map[string]CurrencySourceFactory)}
}

// Add adds a new factory to the pool
func (s *ProviderFactory) Add(factory CurrencySourceFactory, code string) error {
	_, ok := s.factories[code]
	if ok {
		return fmt.Errorf("provider %s already registered", code)
	}

	s.factories[code] = factory

	return nil
}

// Source returns a new rate source
func (s *ProviderFactory) Source(code string) (exchange.RateSource, error) {
	factory, ok := s.factories[code]
	if !ok {
		return nil, fmt.Errorf("unknown provider %s", code)
	}

	return factory.Init()
}
