package currencies

import (
	errorsPkg "github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/repositories"
	"github.com/Confialink/wallet-currencies/internal/services/accounts"
)

type UpdateCurrency struct {
	Id     *uint32 `json:"id" binding:"required"`
	Active bool    `json:"active"`
}

type UpdateCurrenciesRequest struct {
	Data []UpdateCurrency `json:"data" binding:"required,dive,required"`
}

type Service struct {
	repository      *repositories.Currencies
	accountsService *accounts.Service
	ratesRepository *repositories.Rates
}

func NewService(repository *repositories.Currencies, accountsService *accounts.Service, ratesRepository *repositories.Rates) *Service {
	return &Service{repository: repository, accountsService: accountsService, ratesRepository: ratesRepository}
}

func (s *Service) Create(model *models.Currency) error {
	if err := s.repository.Create(model); err != nil {
		return err
	}

	main, err := s.Main()
	if err != nil {
		return err
	}
	if main != nil {
		if _, err := s.ratesRepository.CreateByCurrencyIds(main.ID, model.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) FindById(id uint32) (*models.Currency, error) {
	return s.repository.FindById(id)
}

func (s *Service) FindByCode(code string) (*models.Currency, error) {
	return s.repository.FindByCode(code)
}

func (s *Service) FindByFeed(feed string) ([]*models.Currency, error) {
	return s.repository.FindByFeed(feed)
}

func (s *Service) List(params *list_params.ListParams) ([]*models.Currency, error) {
	return s.repository.List(params)
}

func (s *Service) Main() (*models.Currency, error) {
	return s.repository.Main()
}

func (s *Service) GetAll(params *repositories.ListParams) *[]models.Currency {
	currencies, _ := s.repository.All(params)
	return currencies
}

// UpdateAdminCurrencies updates list of currencies with needed checks
func (s *Service) UpdateAdminCurrencies(req *UpdateCurrenciesRequest) (errors []errorsPkg.TypedError) {
	repository := s.repository
	mainCurrency, err := s.Main()
	if err != nil {
		errors = append(errors, &errorsPkg.PrivateError{OriginalError: err})
		return
	}

	for _, currencyReq := range req.Data {
		currencyModel := s.getCurrencyModel(&currencyReq, repository)
		if !*currencyModel.Active {
			if currencyModel.ID == mainCurrency.ID {
				errors = append(errors, errcodes.GetError(errcodes.DeactivatingMainCurrency))
				continue
			}

			if s.hasBelonging(currencyModel.Code) {
				errors = append(errors, errcodes.GetUsedCurrencyError(currencyModel.Code))
				continue
			}
		}
		if err := repository.Update(currencyModel); err != nil {
			errors = append(errors, &errorsPkg.PrivateError{OriginalError: err})
			return
		}
	}
	return
}

func (s Service) WrapContext(db *gorm.DB) *Service {
	s.repository = s.repository.WrapContext(db)
	return &s
}

func (s *Service) getCurrencyModel(currencyRequest *UpdateCurrency,
	repository *repositories.Currencies,
) *models.Currency {
	currency, _ := repository.FindById(*currencyRequest.Id)
	currency.Active = &currencyRequest.Active
	return currency
}

func (s *Service) hasBelonging(code string) bool {
	return !s.accountsService.CanDisableCurrency(code)
}
