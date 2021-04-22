package currencies

import (
	"github.com/Confialink/wallet-pkg-utils/pointer"
	"github.com/Confialink/wallet-pkg-utils/value"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/repositories"
)

type DefaultManager struct {
	currencyService    *Service
	currencyRepository *repositories.Currencies
	codeValidator      CodeValidator
	logger             log15.Logger
}

func NewDefaultManager(currencyService *Service, currencyRepository *repositories.Currencies, codeValidator CodeValidator, logger log15.Logger) Manager {
	logger = logger.New("Service", "DefaultManager")
	return &DefaultManager{currencyService: currencyService, currencyRepository: currencyRepository, codeValidator: codeValidator, logger: logger}
}

//AddCurrency adds new currency or enables existing one
func (m DefaultManager) AddCurrency(currency Currency) (*models.Currency, error) {
	logger := m.logger.New("method", "AddCurrency")
	existingCurrency, err := m.currencyService.FindByCode(currency.Code)
	if err != nil {
		logger.Crit("failed to find currency by code", "error", err, "currencyCode", currency.Code)
		return nil, err
	}
	if existingCurrency.IsExist() {
		meta := map[string]interface{}{
			"currencyCode": currency.Code,
		}
		isInConflict := false
		if currency.DecimalPlaces != existingCurrency.DecimalPlaces {
			isInConflict = true
			meta["decimalPlaces"] = existingCurrency.DecimalPlaces
			meta["decimalPlacesClaimed"] = currency.DecimalPlaces
		}

		if currency.Type.String() != existingCurrency.Type.String {
			isInConflict = true
			meta["type"] = existingCurrency.Type.String
			meta["typeClaimed"] = currency.Type.String()
		}

		if isInConflict {
			err := errcodes.GetError(errcodes.CurrencyConflict)
			err.Meta = meta
			logger.Warn("adding new currency conflict", "details", meta)
			return nil, err
		}

		if !value.FromBool(existingCurrency.Active) {
			existingCurrency.Active = pointer.ToBool(true)
			err := m.currencyRepository.Update(existingCurrency)
			if err != nil {
				logger.Crit("failed to update currency activity status", "error", err)
				return nil, err
			}
		}

		logger.Info("trying to add currency that is already exist, skipping", "currencyCode", currency.Code)
		return existingCurrency, nil
	}

	if !IsKnownType(currency.Type) {
		err := errcodes.GetError(errcodes.UnknownCurrencyType)
		err.Meta = map[string]string{"type": currency.Type.String()}
		logger.Warn("trying to add currency with unknown type", "details", err.Meta)
		return nil, err
	}

	if err := m.codeValidator.IsValid(currency.Type, currency.Code); err != nil {
		logger.Info("invalid currency code", "error", err)
		return nil, err
	}

	currencyModel := models.Currency{
		Code:          currency.Code,
		DecimalPlaces: currency.DecimalPlaces,
		Active:        &currency.IsActive,
	}
	if currency.Name != nil {
		currencyModel.Name = *currency.Name
	}

	currencyModel.Type.SetValid(currency.Type.String())
	if err := m.currencyService.Create(&currencyModel); err != nil {
		logger.Crit("failed to create currency", "error", err)
		return nil, err
	}

	return &currencyModel, nil
}

//DisableCurrency inactivate existing currency
func (m DefaultManager) DisableCurrency(currencyCode string) (*models.Currency, error) {
	logger := m.logger.New("method", "RemoveCurrency")
	existingCurrency, err := m.currencyService.FindByCode(currencyCode)
	if err != nil {
		logger.Crit("failed to find currency by code", "error", err, "currencyCode", currencyCode)
		return nil, err
	}
	if !existingCurrency.IsExist() {
		err := errcodes.GetError(errcodes.UnknownCurrency)
		err.Meta = map[string]string{"code": currencyCode}
		return nil, err
	}

	if value.FromBool(existingCurrency.Active) {
		existingCurrency.Active = pointer.ToBool(false)
		err := m.currencyRepository.Update(existingCurrency)
		if err != nil {
			logger.Crit("failed to update currency activity status", "error", err)
			return nil, err
		}
	}

	return existingCurrency, nil
}
