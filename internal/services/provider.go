package services

import (
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-currencies/internal/services/accounts"
	"github.com/Confialink/wallet-currencies/internal/services/currencies"
	"github.com/Confialink/wallet-currencies/internal/services/files"
	"github.com/Confialink/wallet-currencies/internal/services/rates"
	"github.com/Confialink/wallet-currencies/internal/services/ratesupdater"
	"github.com/Confialink/wallet-currencies/internal/services/settings"
)

func Providers() []interface{} {
	return []interface{}{
		accounts.NewService,
		currencies.NewService,
		files.NewService,
		rates.NewService,
		rates.NewRates,
		rates.NewDbSource,
		settings.NewService,
		ratesupdater.NewService,
		func(logger log15.Logger) currencies.CodeValidator {
			validators := []currencies.CodeValidator{
				currencies.NewDefaultFiatCurrencyCodeValidator(logger),
				currencies.NewDefaultCryptoCurrencyCodeValidator(logger),
				currencies.NewDefaultOtherCurrencyValidator(logger),
			}
			return currencies.NewCurrencyCodeValidatorAggregate(validators, logger)
		},
		currencies.NewDefaultManager,
		func() *ratesupdater.ProviderFactory {
			factory := ratesupdater.NewProviderFactory()
			if err := factory.Add(ratesupdater.NewEcb(), ratesupdater.EcbProvider); err != nil {
				panic("cannot register provider: " + err.Error())
			}

			return factory
		},
		ratesupdater.NewDbRatesReceiver,
	}
}
