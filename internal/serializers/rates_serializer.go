package serializers

import "github.com/Confialink/wallet-currencies/internal/models"

// Serialize only fields in SerializedRate
type RatesSerializer struct {
}

type SerializedRate struct {
	Id             uint32             `json:"id"`
	Value          string             `json:"value"`
	ExchangeMargin string             `json:"exchangeMargin"`
	CurrencyFrom   serializedCurrency `json:"currencyFrom"`
	CurrencyTo     serializedCurrency `json:"currencyTo"`
}

func (s *RatesSerializer) Serialize(rate *models.Rate) SerializedRate {
	currencySerializer := Currency()
	var currencyFrom = *currencySerializer.Serialize(&rate.CurrencyFrom).(*serializedCurrency)
	var currencyTo = *currencySerializer.Serialize(&rate.CurrencyTo).(*serializedCurrency)
	var serialized = SerializedRate{Id: rate.ID,
		Value:          rate.Value.String(),
		ExchangeMargin: rate.ExchangeMargin.String(),
		CurrencyFrom:   currencyFrom,
		CurrencyTo:     currencyTo,
	}

	return serialized
}

func (s *RatesSerializer) SerializeList(rates *[]models.Rate) []SerializedRate {
	var serialized = make([]SerializedRate, len(*rates))
	for i, rate := range *rates {
		serialized[i] = s.Serialize(&rate)
	}
	return serialized
}
