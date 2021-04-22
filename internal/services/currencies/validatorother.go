package currencies

import (
	"regexp"

	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/inconshreveable/log15"
)

type DefaultOtherCurrencyValidator struct {
	pattern *regexp.Regexp
	logger  log15.Logger
}

func NewDefaultOtherCurrencyValidator(logger log15.Logger) *DefaultOtherCurrencyValidator {
	logger = logger.New("object", "DefaultOtherCurrencyValidator")
	return &DefaultOtherCurrencyValidator{logger: logger}
}

func (d *DefaultOtherCurrencyValidator) IsValid(currencyType CurrencyType, code string) error {
	logger := d.logger.New("method", "IsValid")
	if !d.IsTypeSupported(currencyType) {
		logger.Warn(
			"trying to validate unsupported currency type",
			"supportedType", TypeOther,
			"givenType", currencyType,
		)
		return errorWrongValidation
	}
	if d.pattern == nil {
		d.pattern = regexp.MustCompile("^[0-9A-Za-z]{2,8}$")
	}
	if d.pattern.MatchString(code) {
		return nil
	}

	err := errcodes.GetError(errcodes.CurrencyCodeInvalid)
	err.Details = "currency code must consist of Latin letters and digits and must contain from 2 to 8 characters"
	err.Meta = map[string]string{"code": code}
	return err
}

func (d *DefaultOtherCurrencyValidator) IsTypeSupported(currencyType CurrencyType) bool {
	return currencyType == TypeOther
}
