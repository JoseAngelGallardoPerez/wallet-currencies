package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-currencies/internal/models"
)

type Settings struct {
	db *gorm.DB
}

func NewSettings(db *gorm.DB) *Settings {
	return &Settings{db: db}
}

// Create creates new settings
func (r *Settings) Create(model *models.Settings) error {
	return r.db.Create(model).Error
}

// MainSettings returns main settings
func (r *Settings) MainSettings() (*models.Settings, error) {
	var settings models.Settings

	if err := r.db.First(&settings).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// Update updates settings
func (r *Settings) Update(settings *models.Settings) error {
	return r.db.Save(&settings).Error
}

// WrapContext makes a new copy of the repository with new GORM instance
func (r Settings) WrapContext(db *gorm.DB) *Settings {
	r.db = db
	return &r
}
