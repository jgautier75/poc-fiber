package converters

import (
	"errors"
	"poc-fiber/commons"
	"poc-fiber/dtos"
)

func BuildSectorsLightHierarchy(sectorResponseList dtos.SectorResponseList) dtos.SectorLightResponse {
	var rootSector = sectorResponseList.Sectors
	return fetchLightRecursively(rootSector, sectorResponseList.Sectors.Children)
}

func fetchLightRecursively(rootSector dtos.SectorResponse, sectors []dtos.SectorResponse) dtos.SectorLightResponse {
	var c = make([]dtos.SectorLightResponse, 0)
	for _, sector := range sectors {
		c = append(c, fetchLightRecursively(sector, sector.Children))
	}
	return dtos.SectorLightResponse{
		Uuid:     rootSector.Uuid,
		Code:     rootSector.Code,
		Label:    rootSector.Label,
		Depth:    rootSector.Depth,
		Children: c,
	}
}

func BuildSectorsHierarchy(sectors []dtos.SectorResponse) (dtos.SectorResponse, error) {
	var nilSector dtos.SectorResponse
	var rootSector dtos.SectorResponse
	for _, sector := range sectors {
		if sector.Depth == 0 {
			rootSector = sector
			break
		}
	}
	if &rootSector == &nilSector {
		return rootSector, errors.New(commons.SectorRootNotFound)
	}
	return fetchRecursively(&rootSector, sectors), nil
}

func fetchRecursively(parentSector *dtos.SectorResponse, sectors []dtos.SectorResponse) dtos.SectorResponse {
	var c = make([]dtos.SectorResponse, 0)
	for _, sector := range sectors {
		if sector.ParentId.Int64 == parentSector.Id.Int64 {
			c = append(c, fetchRecursively(&sector, sectors))
		}
	}
	return dtos.SectorResponse{
		Id:       parentSector.Id,
		Uuid:     parentSector.Uuid,
		Code:     parentSector.Code,
		Label:    parentSector.Label,
		ParentId: parentSector.ParentId,
		Depth:    parentSector.Depth,
		Children: c,
	}
}
