package tasks

import (
	"github.com/Confialink/wallet-currencies/internal/services/ratesupdater"
)

// update currency rates
func updateRates(ratesUpdaterService *ratesupdater.Service) error {
	return ratesUpdaterService.Call()
}
