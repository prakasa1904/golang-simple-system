package member

import (
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"
)

/**
Converter HTTP request query to repository structure
*/

func ConvertQueryToLimit(url *url.URL) int {
	var limit = 6

	limitQuery := url.Query().Get("limit")

	if url.Query().Get("limit") != "" {
		limitQueryInt, err := strconv.Atoi(limitQuery)
		if err != nil {
			return limit
		}

		return limitQueryInt
	}

	return limit
}

func ConvertQueryToFilter(url *url.URL) map[string]string {
	var filter = make(map[string]string)

	for key, val := range url.Query() {
		for _, allowKey := range AllowedFilterQuery {
			if allowKey == key {
				filter[key+" LIKE ?"] = "%" + val[0] + "%"
			}
		}
	}

	return filter
}

func ConvertQueryToOrder(url *url.URL) clause.OrderByColumn {
	var isDescOrder = true
	var orderBy = "id"

	order := strings.ToLower(url.Query().Get("order"))
	by := strings.ToLower(url.Query().Get("by"))

	if order == "asc" {
		isDescOrder = false
	}

	if by != "" {
		orderBy = by
	}

	return clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: isDescOrder}
}
