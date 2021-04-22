package errcodes

import (
	"fmt"
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
)

const (
	Forbidden                = "FORBIDDEN"
	CurrencyNotFound         = "CURRENCY_NOT_FOUND"
	CurrencyIsNotActive      = "CURRENCY_IS_NOT_ACTIVE"
	BadCollectionParams      = "BAD_COLLECTION_PARAMS"
	DeactivatingMainCurrency = "DEACTIVATING_MAIN_CURRENCY"
	UsedCurrency             = "USED_CURRENCY"
	UnknownFeed              = "UNKNOWN_FEED"
	EmptyFeed                = "EMPTY_FEED"
	UnknownCurrencyType      = "UNKNOWN_CURRENCY_TYPE"
	UnknownCurrency          = "UNKNOWN_CURRENCY"
	CurrencyCodeInvalid      = "CURRENCY_CODE_INVALID"
	CurrencyConflict         = "CURRENCY_CONFLICT"
)

var TitlesMap = map[string]string{
	Forbidden:                "Forbidden",
	CurrencyNotFound:         "Currency not found",
	BadCollectionParams:      "Params for collections are not valid",
	DeactivatingMainCurrency: "Deactivating main currency is impossible.",
	UnknownFeed:              "Provided currency feed is unknown.",
	UnknownCurrencyType:      "Provided currency type is unknown.",
	UnknownCurrency:          "Provided currency code is unknown",
	CurrencyCodeInvalid:      "The currency code is invalid.",
	CurrencyConflict:         "The currency is in conflict with the previously added currency.",
}

var StatusCodes = map[string]int{
	Forbidden:                http.StatusForbidden,
	CurrencyNotFound:         http.StatusNotFound,
	BadCollectionParams:      http.StatusBadRequest,
	DeactivatingMainCurrency: http.StatusBadRequest,
	UnknownFeed:              http.StatusBadRequest,
	UnknownCurrencyType:      http.StatusForbidden,
	UnknownCurrency:          http.StatusBadRequest,
	CurrencyCodeInvalid:      http.StatusBadRequest,
	CurrencyConflict:         http.StatusConflict,
}

func AddError(c *gin.Context, code string) {
	publicErr := &errors.PublicError{
		Title:      TitlesMap[code],
		Code:       code,
		HttpStatus: statusByCode(code),
	}
	errors.AddErrors(c, publicErr)
}

func GetError(code string) *errors.PublicError {
	return &errors.PublicError{
		Title:      TitlesMap[code],
		Code:       code,
		HttpStatus: statusByCode(code),
	}
}

func AddErrorMeta(c *gin.Context, code string, meta interface{}) {
	publicErr := &errors.PublicError{
		Title:      TitlesMap[code],
		Code:       code,
		HttpStatus: statusByCode(code),
		Meta:       meta,
	}
	errors.AddErrors(c, publicErr)
}

func GetUsedCurrencyError(currencyCode string) *errors.PublicError {
	return &errors.PublicError{
		Title: fmt.Sprintf("%s currency is used in an account or card. Deactivating used currency is impossible", currencyCode),
		Code:  UsedCurrency,
		Meta: struct {
			Currency string `json:"currency"`
		}{currencyCode},
	}
}

func statusByCode(code string) int {
	status, ok := StatusCodes[code]
	if ok {
		return status
	}
	return http.StatusBadRequest
}
