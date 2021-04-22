package tasks

import (
	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/services/currencies"
	settingsService "github.com/Confialink/wallet-currencies/internal/services/settings"
)

const DefaultCurrencyCode = "EUR"

// Creates default currencies and settings
func populateDb(currenciesService *currencies.Service, settingsService *settingsService.Service) error {
	return createSettings(currenciesService, settingsService)
}

func createSettings(currenciesService *currencies.Service, settingsService *settingsService.Service) error {
	currency, _ := currenciesService.FindByCode(DefaultCurrencyCode)

	autoUpdatingRates := false
	settings := &models.Settings{
		MainCurrencyId:    currency.ID,
		AutoUpdatingRates: &autoUpdatingRates,
	}

	return settingsService.Create(settings)
}
