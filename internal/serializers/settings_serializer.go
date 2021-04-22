package serializers

import "github.com/Confialink/wallet-currencies/internal/models"

// Serialize settings
type SettingsSerializer struct {
}

type serializedSettigns struct {
	MainCurrencyId    uint32 `json:"mainCurrencyId"`
	AutoUpdatingRates bool   `json:"autoUpdatingRates"`
}

func (ss *SettingsSerializer) Serialize(model *models.Settings) *serializedSettigns {
	serialized := new(serializedSettigns)
	serialized.MainCurrencyId = model.MainCurrencyId
	serialized.AutoUpdatingRates = *model.AutoUpdatingRates
	return serialized
}
