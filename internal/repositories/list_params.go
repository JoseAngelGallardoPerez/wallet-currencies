package repositories

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const AscSortCriteria bool = false
const DescSortCriteria bool = true

// Parameters for working with collecion of models
type ListParams struct {
	sortCreterias map[string]bool
}

// Initialize and returns new ListParams
func NewListParams(context *gin.Context) *ListParams {
	listParams := new(ListParams)
	listParams.setSortCriteriasString(context)
	return listParams
}

func (l *ListParams) setSortCriteriasString(context *gin.Context) {
	l.sortCreterias = make(map[string]bool)
	if sortParams := context.Query("sort"); sortParams != "" {
		var sortCriterias []string = strings.Split(sortParams, ",")
		for _, sortCriteria := range sortCriterias {
			if (string)(sortCriteria[0]) == "-" {
				l.sortCreterias[sortCriteria[1:]] = DescSortCriteria
			} else {
				l.sortCreterias[sortCriteria] = AscSortCriteria
			}
		}
	}
}

// Returns sql string can be passed for ordering
func (l *ListParams) GetSortCriteriasString() string {
	var newCriterias []string = make([]string, len(l.sortCreterias))
	i := 0
	for field, value := range l.sortCreterias {
		if value == AscSortCriteria {
			newCriterias[i] = field + " asc"
		} else {
			newCriterias[i] = field + " desc"
		}
	}
	i = i + 1
	return strings.Join(newCriterias, ", ")
}
