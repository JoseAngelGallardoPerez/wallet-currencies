package serializers

import "github.com/Confialink/wallet-currencies/internal/models"

// Serialize all currency fields
type adminCurrencySerializer struct {
	ICurrencySerializer
}

type adminSerializedCurrency struct {
	serializedCurrency
	Active bool `json:"active"`
}

func (acs *adminCurrencySerializer) Serialize(model *models.Currency) interface{} {
	return &adminSerializedCurrency{
		serializedCurrency: serializedCurrency{
			Id:            model.ID,
			Code:          model.Code,
			DecimalPlaces: model.DecimalPlaces,
			Type:          model.Type.String,
			Name:          model.Name,
			Feed:          model.Feed,
		},
		Active: *model.Active,
	}
}

func (acs *adminCurrencySerializer) SerializeList(models []*models.Currency) interface{} {
	serialized := make([]adminSerializedCurrency, len(models))
	for i, model := range models {
		if currency, ok := acs.Serialize(model).(*adminSerializedCurrency); ok {
			(serialized)[i] = *(currency)
		}
	}

	return serialized
}
