package rates

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/repositories"
)

// Service service for operations with rates
type Service struct {
	ratesRepository      *repositories.Rates
	currenciesRepository *repositories.Currencies
}

func NewService(ratesRepository *repositories.Rates, currenciesRepository *repositories.Currencies) *Service {
	return &Service{ratesRepository: ratesRepository, currenciesRepository: currenciesRepository}
}

type updateRate struct {
	Id             uint32 `json:"id" binding:"required"`
	Value          string `json:"value" binding:"omitempty,decimal,decimalGT=0"`
	ExchangeMargin string `json:"exchangeMargin" binding:"omitempty,decimal"`
}

type updateRatesRequest struct {
	Data []updateRate `json:"data" binding:"required,min=1,dive,required"`
}

func (s *Service) IndexActiveForMainCurrency() (*[]models.Rate, error) {
	mainCurrency, err := s.currenciesRepository.Main()
	if err != nil {
		return nil, err
	}
	return s.ratesRepository.IndexActiveFromCurrency(mainCurrency.ID), nil
}

func (s *Service) UpdateRates(context *gin.Context) bool {
	var colRates = updateRatesRequest{}

	if err := context.ShouldBindJSON(&colRates); err == nil {
		for _, rate := range colRates.Data {
			rateModel := s.getRateModel(&rate)
			if _, err := s.ratesRepository.UpdateRate(&rateModel); err != nil {
				typedErr := &errors.PrivateError{OriginalError: err}
				errors.AddErrors(context, typedErr)
				return false
			}
		}
	} else {
		errors.AddShouldBindError(context, err)
		return false
	}

	return true
}

func (s *Service) getRateModel(rateRequest *updateRate) models.Rate {
	valueDecimal, _ := decimal.NewFromString(rateRequest.Value)
	exchangeMarginDecimal, _ := decimal.NewFromString(rateRequest.ExchangeMargin)
	return models.Rate{
		ID:             rateRequest.Id,
		Value:          valueDecimal,
		ExchangeMargin: exchangeMarginDecimal,
	}
}
