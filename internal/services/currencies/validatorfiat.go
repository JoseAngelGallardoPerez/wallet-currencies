package currencies

import (
	"regexp"

	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/inconshreveable/log15"
)

type DefaultFiatCurrencyCodeValidator struct {
	pattern *regexp.Regexp
	logger  log15.Logger
}

func NewDefaultFiatCurrencyCodeValidator(logger log15.Logger) *DefaultFiatCurrencyCodeValidator {
	logger = logger.New("object", "DefaultFiatCurrencyCodeValidator")
	return &DefaultFiatCurrencyCodeValidator{logger: logger}
}

func (d *DefaultFiatCurrencyCodeValidator) IsValid(currencyType CurrencyType, code string) error {
	logger := d.logger.New("method", "IsValid")
	if !d.IsTypeSupported(currencyType) {
		logger.Warn(
			"trying to validate unsupported currency type",
			"supportedType", TypeFiat,
			"givenType", currencyType,
		)
		return errorWrongValidation
	}
	if d.pattern == nil {
		d.pattern = regexp.MustCompile("^[A-Z]{3,3}$")
	}
	if d.pattern.MatchString(code) {
		return nil
	}
	err := errcodes.GetError(errcodes.CurrencyCodeInvalid)
	err.Details = "currency code must consist of 3 capital Latin letters"
	err.Meta = map[string]string{"code": code}
	return err
}

func (DefaultFiatCurrencyCodeValidator) IsTypeSupported(currencyType CurrencyType) bool {
	return currencyType == TypeFiat
}
