package settings

import (
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/repositories"
	"github.com/Confialink/wallet-currencies/internal/services/ratesupdater"
)

type Service struct {
	settingsRepo        *repositories.Settings
	currenciesRepo      *repositories.Currencies
	ratesRepo           *repositories.Rates
	ratesUpdaterService *ratesupdater.Service
	logger              log15.Logger
}

func NewService(settingsRepo *repositories.Settings, currenciesRepo *repositories.Currencies, ratesRepo *repositories.Rates, ratesUpdaterService *ratesupdater.Service, logger log15.Logger) *Service {
	logger = logger.New("Service", "Settings")
	return &Service{settingsRepo: settingsRepo, currenciesRepo: currenciesRepo, ratesRepo: ratesRepo, ratesUpdaterService: ratesUpdaterService, logger: logger}
}

func (s *Service) Create(model *models.Settings) error {
	return s.settingsRepo.Create(model)
}

func (s *Service) Update(mainCurrencyId uint32, autoUpdatingRates *bool) (current *models.Settings, err error) {
	current, err = s.MainSettings()
	if err != nil {
		return
	}
	mainCurrency, err := s.currenciesRepo.FindById(mainCurrencyId)
	if err != nil {
		s.logger.Error("cannot find currency", "err", err)
		return nil, errcodes.GetError(errcodes.CurrencyNotFound)
	}

	if !*mainCurrency.Active {
		return nil, errcodes.GetError(errcodes.CurrencyIsNotActive)
	}

	s.assignUpdateValues(mainCurrencyId, autoUpdatingRates, current)

	// allow autoupdate rates only if the currency has a feed
	if *current.AutoUpdatingRates && !mainCurrency.Feed.Valid || len(mainCurrency.Feed.String) == 0 {
		return nil, errcodes.GetError(errcodes.EmptyFeed)
	}

	if err = s.settingsRepo.Update(current); err != nil {
		return nil, errors.Wrap(err, "cannot update settings")
	}

	if err = s.applyNewSettings(current); err != nil {
		return nil, errors.Wrap(err, "cannot apply settings")
	}

	return
}

func (s *Service) MainSettings() (*models.Settings, error) {
	return s.settingsRepo.MainSettings()
}

func (s Service) WrapContext(db *gorm.DB) *Service {
	s.settingsRepo = s.settingsRepo.WrapContext(db)
	s.currenciesRepo = s.currenciesRepo.WrapContext(db)
	s.ratesRepo = s.ratesRepo.WrapContext(db)
	s.ratesUpdaterService = s.ratesUpdaterService.WrapContext(db)
	return &s
}

func (s *Service) assignUpdateValues(mainCurrencyId uint32, autoUpdatingRates *bool, to *models.Settings) {
	if autoUpdatingRates != nil {
		to.AutoUpdatingRates = autoUpdatingRates
	}

	to.MainCurrencyId = mainCurrencyId
}

// applyNewSettings creates missing rates and updates rates if need
func (s *Service) applyNewSettings(settings *models.Settings) (err error) {
	if err = s.createMissingRates(settings.MainCurrencyId); err != nil {
		return
	}

	if settings.AutoUpdatingRates != nil && *settings.AutoUpdatingRates {
		return s.ratesUpdaterService.Call()
	}

	return nil
}

// createMissingRates creates rates for provided currency id
func (s *Service) createMissingRates(fromCurrencyId uint32) error {
	ids := s.currenciesRepo.IdsWithoutRate(fromCurrencyId)
	ratesRepo := s.ratesRepo

	for _, id := range ids {
		if _, err := ratesRepo.CreateByCurrencyIds(fromCurrencyId, id); err != nil {
			return err
		}
	}

	return nil
}
