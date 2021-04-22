package controllers

import (
	"net/http"

	errors "github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-currencies/internal/abilities"
	"github.com/Confialink/wallet-currencies/internal/abilities/actions"
	"github.com/Confialink/wallet-currencies/internal/abilities/resources"
	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/Confialink/wallet-currencies/internal/serializers"
	"github.com/Confialink/wallet-currencies/internal/services/rates"
)

type Rates struct {
	ratesService    *rates.Service
	ratesCalculator *rates.Rates
}

func NewRatesController(ratesService *rates.Service, ratesCalculator *rates.Rates) *Rates {
	return &Rates{ratesService, ratesCalculator}
}

// Returns currency by id or main currency
func (r *Rates) Index(c *gin.Context) {
	if abilities.CanNot(c, actions.Index, resources.Rates) {
		errcodes.AddError(c, errcodes.Forbidden)
		return
	}

	result, err := r.ratesService.IndexActiveForMainCurrency()
	if err != nil {
		privateErr := errors.PrivateError{OriginalError: err}
		errors.AddErrors(c, &privateErr)
		return
	}
	serialized := serializers.Rates().SerializeList(result)
	formOkResponse(c, serialized)
}

func (r *Rates) GetForCurrencies(c *gin.Context) {
	if abilities.CanNot(c, actions.Read, resources.Rates) {
		errcodes.AddError(c, errcodes.Forbidden)
		return
	}

	baseCurrency, provided := c.GetQuery("baseCurrencyCode")
	if !provided {
		errcodes.AddError(c, errcodes.BadCollectionParams)
		return
	}

	referenceCurrency, provided := c.GetQuery("referenceCurrencyCode")
	if !provided {
		errcodes.AddError(c, errcodes.BadCollectionParams)
		return
	}

	result, err := r.ratesCalculator.CalculateRate(baseCurrency, referenceCurrency)
	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{OriginalError: err})
		return
	}

	formOkResponse(c, result)
}

func (r *Rates) UpdateList(c *gin.Context) {
	if abilities.CanNot(c, actions.Update, resources.Rates) {
		errcodes.AddError(c, errcodes.Forbidden)
		return
	}

	if r.ratesService.UpdateRates(c) {
		c.Status(http.StatusNoContent)
	}
}
