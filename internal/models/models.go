// Package for models
package models

import (
	"time"

	"github.com/guregu/null"
	"github.com/shopspring/decimal"
)

const TypeCrypto = "crypto"

type Currency struct {
	ID     uint32      `json:"id"`
	Code   string      `gorm:"not null" json:"code"`
	Active *bool       `gorm:"not null" default:"false" json:"active"`
	Type   null.String `json:"type"`
	//Named source from where currency must update rate
	Feed          null.String `json:"feed"`
	DecimalPlaces uint8       `json:"decimalPlaces" binding:"required"`
	Name          string      `json:"name" binding:"required"`
	LogoFileID    uint64      `json:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c Currency) IsExist() bool {
	return c.ID > 0
}

// Settings model
type Settings struct {
	ID                uint32   `json:"id"`
	MainCurrencyId    uint32   `gorm:"not null" json:"mainCurrencyId"`
	MainCurrency      Currency `gorm:"foreignkey:MainCurrencyId;association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	AutoUpdatingRates *bool    `json:"autoUpdatingRates" binding:"omitempty"`
}

// Rate model
type Rate struct {
	ID             uint32          `json:"id"`
	CurrencyFromId uint32          `json:"currencyFromId"`
	CurrencyFrom   Currency        `gorm:"foreignkey:CurrencyFromId;association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	CurrencyToId   uint32          `json:"currencyToId"`
	CurrencyTo     Currency        `gorm:"foreignkey:CurrencyToId;association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Value          decimal.Decimal `json:"value"`
	ExchangeMargin decimal.Decimal `json:"exchangeMargin"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *Rate) IsExists() bool {
	return s.ID != 0
}
