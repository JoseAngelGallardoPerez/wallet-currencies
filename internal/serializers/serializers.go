// Contains serializers to serialize models
package serializers

import "github.com/Confialink/wallet-currencies/internal/models"

// Interface to serialize currencies by different ways
type ICurrencySerializer interface {
	Serialize(*models.Currency) interface{}
	SerializeList([]*models.Currency) interface{}
}

func Currency() ICurrencySerializer {
	return new(currencySerializer)
}

func AdminCurrency() ICurrencySerializer {
	return new(adminCurrencySerializer)
}

func Settings() *SettingsSerializer {
	return new(SettingsSerializer)
}

func Rates() *RatesSerializer {
	return new(RatesSerializer)
}
