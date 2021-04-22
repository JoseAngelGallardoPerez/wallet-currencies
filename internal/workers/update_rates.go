package workers

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-currencies/internal/repositories"
	ratesUpdater "github.com/Confialink/wallet-currencies/internal/services/ratesupdater"
)

type UpdateRates struct {
	ratesUpdater       *ratesUpdater.Service
	settingsRepository *repositories.Settings
}

func NewUpdateRates(ratesUpdater *ratesUpdater.Service, settingsRepository *repositories.Settings) *UpdateRates {
	return &UpdateRates{ratesUpdater, settingsRepository}
}

// Update updates currency rates
func (s *UpdateRates) Update() error {
	settings, err := s.settingsRepository.MainSettings()
	if err != nil {
		return errors.Wrap(err, "cannot obtain settings")
	}

	// skip if auto-updating rates is disabled
	if !*settings.AutoUpdatingRates {
		return nil
	}
	if err := s.ratesUpdater.Call(); err != nil {
		return err
	}

	return nil
}

func (s UpdateRates) WrapContext(db *gorm.DB) *UpdateRates {
	s.settingsRepository = s.settingsRepository.WrapContext(db)
	s.ratesUpdater = s.ratesUpdater.WrapContext(db)
	return &s
}
