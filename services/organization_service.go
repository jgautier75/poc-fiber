package services

import (
	"context"
	"poc-fiber/converters"
	"poc-fiber/dao"
	"poc-fiber/dtos"
	"poc-fiber/model"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type OrganizationService struct {
	tenantDao       dao.TenantDao
	organizationDao dao.OrganizationDao
	sectorDao       dao.SectorDao
	logger          zap.Logger
}

func NewOrganizationService(tDao dao.TenantDao, oDao dao.OrganizationDao, sDao dao.SectorDao, l zap.Logger) OrganizationService {
	orgService := OrganizationService{
		tenantDao:       tDao,
		organizationDao: oDao,
		sectorDao:       sDao,
		logger:          l,
	}
	return orgService
}

func (orgService *OrganizationService) CreateOrganization(tenantUuid string, code string, label string, otype string) (model.CompositeId, error) {

	// Find tenant
	tenant, errorTenant := orgService.tenantDao.FindByUuid(tenantUuid)
	if errorTenant != nil {
		orgService.logger.Error("error find tenant [%w]", zap.Error(errorTenant))
	}
	var nilComposite model.CompositeId

	// Get connection and init transaction
	conn, errConnect := orgService.tenantDao.DbPool.Acquire(context.Background())
	if errConnect != nil {
		return nilComposite, errConnect
	}
	tx, errTx := conn.BeginTx(context.Background(), pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead})
	if errTx != nil {
		return nilComposite, errTx
	}
	defer func() {
		if errTx != nil {
			errRbk := tx.Rollback(context.Background())
			if errRbk != nil {
				orgService.logger.Error("rollbak error [%w]", zap.Error(errRbk))
			}
		} else {
			errCmt := tx.Commit(context.Background())
			if errCmt != nil {
				orgService.logger.Error("commit error [%w]", zap.Error(errCmt))
			}
		}
	}()

	orgCid, errCreateOrg := orgService.organizationDao.WithTxCreateOrganization(tx, tenant.Id, code, label, otype)
	if errCreateOrg != nil {
		orgService.logger.Error("error creating organization [%w]", zap.Error(errCreateOrg))
	}

	sector := model.Sector{}
	sector.TenantId = tenant.Id
	sector.OrganizationId = orgCid.Id
	sector.Code = "root-" + code
	sector.Label = label
	sector.Depth = 0
	sector.HasParent = false

	sectorCid, errCreateSector := orgService.sectorDao.WithTxCreateSector(tx, sector)

	return sectorCid, errCreateSector
}

func (orgService *OrganizationService) FindAllOrganizations(tenantUuid string, logger zap.Logger) (dtos.OrgLightReponseList, error) {
	var orgsResponse = dtos.OrgLightReponseList{}
	logger.Info("find all organizations", zap.String("tenantUuid", tenantUuid))
	tenant, errTenant := orgService.tenantDao.FindByUuid(tenantUuid)
	if errTenant != nil {
		return orgsResponse, errTenant
	}
	orgs, errOrgs := orgService.organizationDao.FindAllByTenantId(tenant.Id)
	if errOrgs != nil {
		return orgsResponse, errOrgs
	}
	orgsResponse = converters.ConvertOrgEntityListToOrgLightList(orgs)
	return orgsResponse, nil
}
