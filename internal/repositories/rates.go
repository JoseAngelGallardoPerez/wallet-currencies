package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-pkg-errors"

	"github.com/Confialink/wallet-currencies/internal/models"
)

type Rates struct {
	db *gorm.DB
}

func NewRates(db *gorm.DB) *Rates {
	return &Rates{db: db}
}

func (r *Rates) CreateByCurrencyIds(fromCurrencyId uint32, toCurrencyId uint32) (rate *models.Rate, err error) {
	rate = &models.Rate{CurrencyFromId: fromCurrencyId, CurrencyToId: toCurrencyId}

	if err = r.db.Create(&rate).Error; err != nil {
		return
	}
	return
}

func (r *Rates) IndexActiveFromCurrency(fromCurrencyId uint32) *[]models.Rate {
	var rates []models.Rate

	r.joinsCurrencies().
		Where("currencies_from.id = ?", fromCurrencyId).
		Where("currencies_to.active = ?", true).
		Order("currencies_to.code ASC").
		Preload("CurrencyFrom").
		Preload("CurrencyTo").
		Find(&rates)

	return &rates
}

func (r *Rates) GetByModelFields(rateCondition *models.Rate) (*models.Rate, error) {
	var rate models.Rate

	if err := r.db.Where(rateCondition).Find(&rate).Error; err != nil {
		return nil, err
	}

	return &rate, nil
}

func (r *Rates) UpdateRate(rate *models.Rate) (*models.Rate, error) {
	if err := r.db.Model(rate).Updates(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *Rates) FindByCurrencies(codeFrom, codeTo string) (*models.Rate, error) {
	rate := models.Rate{}
	err := r.joinsCurrencies().
		Where("currencies_from.code = ?", codeFrom).
		Where("currencies_to.code = ?", codeTo).First(&rate).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &rate, nil
}

func (r *Rates) FindByCurrencyIDs(currencyIDFrom, currencyIDTo uint32) (*models.Rate, errors.TypedError) {
	rate := models.Rate{}
	err := r.db.Where("currency_from_id = ?", currencyIDFrom).
		Where("currency_to_id = ?", currencyIDTo).First(&rate).Error

	if err != nil {
		return nil, &errors.PrivateError{OriginalError: err}
	}
	return &rate, nil
}

func (r Rates) WrapContext(db *gorm.DB) *Rates {
	r.db = db
	return &r
}

func (r *Rates) joinsCurrencies() *gorm.DB {
	return r.db.
		Joins("LEFT JOIN currencies AS currencies_from ON rates.currency_from_id = currencies_from.id").
		Joins("LEFT JOIN currencies AS currencies_to ON rates.currency_to_id = currencies_to.id")
}
