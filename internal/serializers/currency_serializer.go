package serializers

import (
	"github.com/guregu/null"

	"github.com/Confialink/wallet-currencies/internal/models"
)

// Serialize id and name
type currencySerializer struct {
	ICurrencySerializer
}

type serializedCurrency struct {
	Id            uint32      `json:"id"`
	Code          string      `json:"code"`
	DecimalPlaces uint8       `json:"decimalPlaces"`
	Type          string      `json:"type"`
	Name          string      `json:"name"`
	Feed          null.String `json:"feed"`
}

func (cs *currencySerializer) Serialize(model *models.Currency) interface{} {
	return &serializedCurrency{
		Id:            model.ID,
		Code:          model.Code,
		DecimalPlaces: model.DecimalPlaces,
		Type:          model.Type.String,
		Name:          model.Name,
		Feed:          model.Feed,
	}
}

func (cs *currencySerializer) SerializeList(models []*models.Currency) interface{} {
	serialized := make([]serializedCurrency, len(models))
	for i, model := range models {
		if currency, ok := cs.Serialize(model).(*serializedCurrency); ok {
			(serialized)[i] = *(currency)
		}
	}

	return serialized
}
