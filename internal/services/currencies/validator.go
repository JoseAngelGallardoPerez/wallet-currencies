package currencies

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

var errorWrongValidation = errors.New("given validator does not support this type of currency")

type CurrencyCodeValidatorAggregate struct {
	validators []CodeValidator
	logger     log15.Logger
}

func NewCurrencyCodeValidatorAggregate(validators []CodeValidator, logger log15.Logger) *CurrencyCodeValidatorAggregate {
	logger = logger.New("object", "CurrencyCodeValidatorAggregate")
	return &CurrencyCodeValidatorAggregate{validators: validators, logger: logger}
}

func (c CurrencyCodeValidatorAggregate) IsValid(currencyType CurrencyType, code string) error {
	logger := c.logger.New("method", "IsValid")
	for _, validator := range c.validators {
		if validator.IsTypeSupported(currencyType) {
			return validator.IsValid(currencyType, code)
		}
	}
	logger.Warn(
		"trying to validate unsupported currency type",
		"givenType", currencyType,
	)
	return errorWrongValidation
}

func (c CurrencyCodeValidatorAggregate) IsTypeSupported(currencyType CurrencyType) bool {
	for _, validator := range c.validators {
		if validator.IsTypeSupported(currencyType) {
			return true
		}
	}
	return false
}
