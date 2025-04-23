package model

import (
	"strconv"
	"strings"
)

type OrderBy struct {
	Column string
	Order  string
}

type Pagination struct {
	RowsPerPage int
	Page        int
	Sorting     []OrderBy
}

func FromParameters(rowsPerPage string, page string, sorting string) (p Pagination, e error) {
	var nilPagination Pagination
	var rpp = 10
	var pg = 1
	if rowsPerPage != "" {
		rpp, e = strconv.Atoi(rowsPerPage)
		if e != nil {
			return nilPagination, e
		}
	}
	if page != "" {
		pg, e = strconv.Atoi(page)
		if e != nil {
			return nilPagination, e
		}
	}

	orderList := make([]OrderBy, 0)
	if sorting != "" {
		strArray := strings.Split(sorting, ",")
		for _, sortStr := range strArray {
			if strings.HasPrefix(sortStr, "+") {
				sortColumn := OrderBy{
					Column: sortStr[1:],
					Order:  "asc",
				}
				orderList = append(orderList, sortColumn)
			}
			if strings.HasPrefix(sortStr, "-") {
				sortColumn := OrderBy{
					Column: sortStr[1:],
					Order:  "desc",
				}
				orderList = append(orderList, sortColumn)
			}
		}
	}
	return Pagination{
		RowsPerPage: rpp,
		Page:        pg,
		Sorting:     orderList,
	}, nil
}
