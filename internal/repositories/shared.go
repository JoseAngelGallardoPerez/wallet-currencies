package repositories

import (
	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/jinzhu/gorm"
)

// processDbQuery returns new db query with applied options from ListParams
func processDbQuery(params *list_params.ListParams, query *gorm.DB) *gorm.DB {
	str, arguments := params.GetWhereCondition()
	resut := query.Where(str, arguments...)

	resut = resut.Order(params.GetOrderByString())

	if limit := params.GetLimit(); limit != 0 {
		resut = resut.Limit(limit)
	}
	resut = resut.Offset(params.GetOffset())

	resut = resut.Joins(params.GetJoinCondition())

	for _, preloadName := range params.GetPreloads() {
		resut = query.Preload(preloadName)
	}
	return resut
}
