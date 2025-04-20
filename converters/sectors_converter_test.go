package converters

import (
	"database/sql"
	"poc-fiber/dtos"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	DEFAULT_TENANT = 1
	DEFAULT_ORG    = 1
)

func TestSectorsConverter(t *testing.T) {
	sectors := make([]dtos.SectorResponse, 4)
	rootUuid := uuid.New().String()
	rootSector := dtos.SectorResponse{
		Id:    sql.NullInt64{Int64: 1, Valid: true},
		Uuid:  rootUuid,
		Code:  "root",
		Label: "root",
		Depth: 0,
	}
	sectors[0] = rootSector

	northUuid := uuid.New().String()
	northSector := dtos.SectorResponse{
		Id:       sql.NullInt64{Int64: 2, Valid: true},
		Uuid:     northUuid,
		Code:     "north",
		Label:    "north",
		Depth:    1,
		ParentId: sql.NullInt64{Int64: 1, Valid: true},
	}
	sectors[1] = northSector

	southUuid := uuid.New().String()
	southSector := dtos.SectorResponse{
		Id:       sql.NullInt64{Int64: 3, Valid: true},
		Uuid:     southUuid,
		Code:     "south",
		Label:    "south",
		Depth:    1,
		ParentId: sql.NullInt64{Int64: 1, Valid: true},
	}
	sectors[2] = southSector

	northEastUuid := uuid.New().String()
	northEastSector := dtos.SectorResponse{
		Id:       sql.NullInt64{Int64: 4, Valid: true},
		Uuid:     northEastUuid,
		Code:     "north-east",
		Label:    "north-east",
		Depth:    2,
		ParentId: sql.NullInt64{Int64: 2, Valid: true},
	}
	sectors[3] = northEastSector

	rootResponse, errHierarchy := BuildSectorsHierarchy(sectors)
	assert.Nil(t, errHierarchy, "error building hierarchy")
	assert.NotNil(t, rootResponse, "root sector not nul")

	assert.Equal(t, 2, len(rootResponse.Children), "2 children under root sector")
	assert.Equal(t, "north", rootResponse.Children[0].Code, "north sector found")
	assert.Equal(t, "south", rootResponse.Children[1].Code, "south sector found")
	assert.Equal(t, 1, len(rootResponse.Children[0].Children), "1 child for north sector")
	assert.Equal(t, "north-east", rootResponse.Children[0].Children[0].Code, "north-east sector as north sector child")
}
