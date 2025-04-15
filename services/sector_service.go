package services

import (
	"errors"
	"poc-fiber/commons"
	"poc-fiber/dao"
	"poc-fiber/dtos"
	"poc-fiber/functions"

	"go.uber.org/zap"
)

type SectorService struct {
	tenantFunctions functions.TenantFunctions
	orgsFunctions   functions.OrganizationsFunctions
	sectorDao       dao.SectorDao
}

func NewSectorService(tenantFunctions functions.TenantFunctions, orgsFunctions functions.OrganizationsFunctions, sectorDao dao.SectorDao, l zap.Logger) SectorService {
	sectorService := SectorService{
		tenantFunctions: tenantFunctions,
		orgsFunctions:   orgsFunctions,
		sectorDao:       sectorDao,
	}
	return sectorService
}

func (sectorService *SectorService) FindSectorsByTenantAndOrganization(tenantUuid string, organizationUuid string, logger zap.Logger) (dtos.SectorResponseList, error) {
	var sectorsList = dtos.SectorResponseList{}

	// Ensure tenant exists
	tenant, errFindTenant := sectorService.tenantFunctions.FindTenant(tenantUuid, logger)
	if errFindTenant != nil {
		return sectorsList, errFindTenant
	}

	// Ensure organization exists
	org, errFindOrg := sectorService.orgsFunctions.FindOrganization(tenant.Id, organizationUuid, logger)
	if errFindOrg != nil {
		return sectorsList, errFindOrg
	}

	sectors, errSectors := sectorService.sectorDao.FindByTenantandOrganization(tenant.Id, org.Id)
	if errSectors != nil {
		return sectorsList, errSectors
	}

	// Convert to response objects
	sectorsResponseArray := make([]dtos.SectorResponse, len(sectors))
	for inc, s := range sectors {
		sgResponse := dtos.SectorResponse{
			Id:       &s.Id,
			Uuid:     &s.Uuid,
			Code:     &s.Code,
			Label:    &s.Label,
			Depth:    s.Depth,
			ParentId: &s.ParentId.Int64,
		}
		sectorsResponseArray[inc] = sgResponse
	}

	s, errHierarchy := buildSectorsHierarchy(sectorsResponseArray)
	if errHierarchy != nil {
		return sectorsList, errHierarchy
	}
	sectorsList.Sectors = s
	return sectorsList, nil
}

func buildSectorsHierarchy(sectors []dtos.SectorResponse) (dtos.SectorResponse, error) {
	var rootSector dtos.SectorResponse
	for _, sector := range sectors {
		if sector.Depth == 0 {
			rootSector = sector
			break
		}
	}
	if &rootSector == nil {
		return rootSector, errors.New(commons.SectorRootNotFound)
	}
	return fetchRecursively(&rootSector, sectors), nil
}

func fetchRecursively(parentSector *dtos.SectorResponse, sectors []dtos.SectorResponse) dtos.SectorResponse {
	var c = make([]dtos.SectorResponse, 0)
	for _, sector := range sectors {
		if sector.ParentId == parentSector.Id {
			c = append(c, fetchRecursively(&sector, sectors))
		}
	}
	return dtos.SectorResponse{
		Id:       parentSector.Id,
		Code:     parentSector.Code,
		Label:    parentSector.Label,
		ParentId: parentSector.ParentId,
		Depth:    parentSector.Depth,
		Children: c,
	}
}
