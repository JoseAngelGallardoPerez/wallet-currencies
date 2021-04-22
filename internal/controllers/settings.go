package controllers

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-currencies/internal/abilities"
	"github.com/Confialink/wallet-currencies/internal/abilities/actions"
	"github.com/Confialink/wallet-currencies/internal/abilities/resources"
	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/Confialink/wallet-currencies/internal/serializers"
	settingsService "github.com/Confialink/wallet-currencies/internal/services/settings"
)

type Settings struct {
	settingsService    *settingsService.Service
	settingsSerializer *serializers.SettingsSerializer
	db                 *gorm.DB
}

func NewSettingsController(settingsService *settingsService.Service, settingsSerializer *serializers.SettingsSerializer, db *gorm.DB) *Settings {
	return &Settings{settingsService: settingsService, settingsSerializer: settingsSerializer, db: db}
}

type updateSettingsRequest struct {
	Data *UpdateSettingsParams `json:"data" binding:"required"`
}

type UpdateSettingsParams struct {
	MainCurrencyId    uint32 `json:"mainCurrencyId"`
	AutoUpdatingRates *bool  `json:"autoUpdatingRates" binding:"omitempty"`
}

// ShowMain returns settings
func (c *Settings) ShowMain(context *gin.Context) {
	if abilities.CanNot(context, actions.Read, resources.Settings) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	settings, err := c.settingsService.MainSettings()
	if err != nil {
		_ = context.Error(err)
		return
	}
	serialized := c.settingsSerializer.Serialize(settings)
	formOkResponse(context, serialized)
}

// UpdateMain updates settings
func (c *Settings) UpdateMain(context *gin.Context) {
	if abilities.CanNot(context, actions.Update, resources.Settings) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	requestObject := updateSettingsRequest{}
	if err := context.ShouldBindJSON(&requestObject); err != nil {
		errors.AddShouldBindError(context, err)
		return
	}

	updateData := requestObject.Data
	tx := c.db.Begin()
	settingsSrv := c.settingsService.WrapContext(tx)

	updated, err := settingsSrv.Update(updateData.MainCurrencyId, updateData.AutoUpdatingRates)
	if err != nil {
		tx.Rollback()
		_ = context.Error(err)
		return
	}
	tx.Commit()

	serialized := c.settingsSerializer.Serialize(updated)
	formOkResponse(context, serialized)

}
