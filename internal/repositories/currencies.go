package repositories

import (
	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-currencies/internal/models"
)

type Currencies struct {
	db *gorm.DB
}

func NewCurrencies(db *gorm.DB) *Currencies {
	return &Currencies{db: db}
}

// Creates new currency with rate for currency
func (s *Currencies) Create(model *models.Currency) error {
	if err := s.db.Create(model).Error; err != nil {
		return err
	}

	return nil
}

// FindById returns currency by passed id
func (s *Currencies) FindById(id uint32) (*models.Currency, error) {
	var currency models.Currency

	if err := s.db.Where("id = ?", id).First(&currency).Error; err != nil {
		return nil, err
	}

	return &currency, nil
}

// FindByCode returns currency by passed code
func (s *Currencies) FindByCode(code string) (*models.Currency, error) {
	var currency models.Currency

	err := s.db.Where("code = ?", code).FirstOrInit(&currency).Error

	return &currency, err
}

// FindByCode returns currencies by passed feed
func (s *Currencies) FindByFeed(feed string) ([]*models.Currency, error) {
	var currencies []*models.Currency
	if err := s.db.Where("feed = ?", feed).Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// FindByCodes returns currencies by passed codes
func (s *Currencies) FindByCodes(codes []string) ([]*models.Currency, error) {
	var currencies []*models.Currency

	if err := s.db.Where("code IN (?)", codes).Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// Returns main currency
func (s *Currencies) Main() (*models.Currency, error) {
	var currency models.Currency
	query := s.db.Joins("INNER JOIN settings ON currencies.id = settings.main_currency_id")
	if err := query.First(&currency).Error; err != nil {
		return nil, err
	}

	return &currency, nil
}

// All returns all currencies
func (s *Currencies) All(params *ListParams) (*[]models.Currency, error) {
	var currencies []models.Currency

	err := s.applyListParams(s.db, params).Find(&currencies).Error

	return &currencies, err
}

// List returns list of currencies by ListParams
func (s *Currencies) List(params *list_params.ListParams) ([]*models.Currency, error) {
	var currencies []*models.Currency

	query := processDbQuery(params, s.db)
	if err := query.Find(&currencies).Error; err != nil {
		return currencies, err
	}
	return currencies, nil
}

// Update updates the currency model
func (s *Currencies) Update(currency *models.Currency) error {
	return s.db.Model(currency).Updates(currency).Error
}

func (s Currencies) WrapContext(db *gorm.DB) *Currencies {
	s.db = db
	return &s
}

// IdsWithoutRate return a slice of currency ID which does not have rates
func (s *Currencies) IdsWithoutRate(fromCurrencyId uint32) []uint32 {
	var ids []uint32
	s.db.
		Table("currencies").
		Where("currencies.id != ?", fromCurrencyId).
		Where("currencies.id NOT IN (?)", s.db.
			Table("currencies").
			Select("rates.currency_to_id").
			Joins("LEFT JOIN rates ON rates.currency_from_id = currencies.id").
			Where("rates.currency_from_id = ?", fromCurrencyId).
			Where("rates.currency_to_id IS NOT NULL").
			QueryExpr()).
		Pluck("currencies.id", &ids)

	return ids
}

func (s *Currencies) applyListParams(query *gorm.DB, params *ListParams) *gorm.DB {
	return query.Order(params.GetSortCriteriasString())
}
