package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	pag, err := FromParameters("25", "3", "+last_name,-first_name")
	assert.Nil(t, err, "error parsing pagination")
	assert.NotNil(t, pag, "pagination not null")
	assert.Equal(t, 2, len(pag.Sorting), "2 sort columns")
	assert.Equal(t, "last_name", pag.Sorting[0].Column, "last name sorting")
	assert.Equal(t, "asc", pag.Sorting[0].Order, "last name asc sorting")

	defPag, err := FromParameters("", "", "")
	assert.Nil(t, err, "error parsing pagination")
	assert.NotNil(t, defPag, "default parameters")
	assert.Equal(t, 1, defPag.Page, "page 1")
	assert.Equal(t, 10, defPag.RowsPerPage, "10 rows per page")
}
