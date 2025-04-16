package services

import (
	"database/sql"
	"errors"
	"poc-fiber/converters"
	"poc-fiber/dao"
	"poc-fiber/dtos"
	"poc-fiber/exceptions"
	"poc-fiber/functions"
	"poc-fiber/model"

	"github.com/google/uuid"

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

	sectors, errSectors := sectorService.sectorDao.FindAllByTenantAndOrganization(tenant.Id, org.Id)
	if errSectors != nil {
		return sectorsList, errSectors
	}

	// Convert to response objects
	sectorsResponseArray := make([]dtos.SectorResponse, len(sectors))
	for inc, s := range sectors {
		sgResponse := dtos.SectorResponse{
			Id:       sql.NullInt64{Int64: s.Id},
			Uuid:     &s.Uuid,
			Code:     &s.Code,
			Label:    &s.Label,
			Depth:    s.Depth,
			ParentId: s.ParentId,
		}
		sectorsResponseArray[inc] = sgResponse
	}

	s, errHierarchy := converters.BuildSectorsHierarchy(sectorsResponseArray)
	if errHierarchy != nil {
		return sectorsList, errHierarchy
	}
	sectorsList.Sectors = s
	return sectorsList, nil
}

func (sectorService *SectorService) CreateSector(tenantUuid string, orgUuid string, sectorReq dtos.SectorCreateRequest, logger zap.Logger) (model.CompositeId, error) {
	var nilComposite model.CompositeId
	var nilSector model.Sector

	// Ensure tenant exists
	tenant, errFindTenant := sectorService.tenantFunctions.FindTenant(tenantUuid, logger)
	if errFindTenant != nil {
		return nilComposite, errFindTenant
	}

	// Ensure organization exists
	org, errFindOrg := sectorService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, logger)
	if errFindOrg != nil {
		return nilComposite, errFindOrg
	}

	parentSector, errParent := sectorService.sectorDao.FindByUuid(*sectorReq.ParentUuid)
	if errParent != nil {
		return nilComposite, errParent
	}

	if parentSector == nilSector {
		return nilComposite, errors.New(exceptions.SECTOR_NOT_FOUND)
	}

	// Create sector
	parentId := sql.NullInt64{
		Int64: parentSector.Id,
		Valid: true,
	}
	nuuid := uuid.New().String()

	sector := model.Sector{
		TenantId:       tenant.Id,
		OrganizationId: org.Id,
		Code:           *sectorReq.Code,
		Label:          *sectorReq.Label,
		ParentId:       parentId,
		Uuid:           nuuid,
	}

	sectorId, errCreate := sectorService.sectorDao.CreateSector(sector)
	if errCreate != nil {
		return nilComposite, errCreate
	}

	return sectorId, nil
}
