package controllers

import (
	"fmt"

	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/gin-gonic/gin"
)

func getListParams(ctx *gin.Context) *list_params.ListParams {
	params := list_params.NewListParamsFromQuery(ctx.Request.URL.RawQuery, models.Currency{})
	params.AllowPagination()
	addFilters(params)
	addSortings(params)
	return params
}

func addFilters(params *list_params.ListParams) {
	params.AllowFilters([]string{"active", "type"})
	params.AddCustomFilter("active", list_params.BoolFilter("active"))
}

func addSortings(params *list_params.ListParams) {
	params.AllowSortings([]string{"type"})
	params.AddCustomSortings("type", sortByType)
}

func sortByType(direction string, params *list_params.ListParams) (string, error) {
	return fmt.Sprintf("currencies.type != 'fiat', currencies.type != 'crypto', currencies.type != 'other', currencies.code %s", direction), nil
}
