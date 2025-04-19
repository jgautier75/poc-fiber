package services

import (
	"context"
	"errors"
	"poc-fiber/commons"
	"poc-fiber/converters"
	"poc-fiber/dao"
	"poc-fiber/dtos"
	"poc-fiber/model"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const OTEL_TRACER_NAME = "otel-collector"

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

func (orgService *OrganizationService) CreateOrganization(tenantUuid string, orgCreateReq dtos.CreateOrgRequest, parentContext context.Context) (model.CompositeId, error) {
	var nilComposite model.CompositeId

	c, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "ORG-LIST-SERVICE")
	defer span.End()

	// Find tenant
	tenant, errorTenant := orgService.tenantDao.FindByUuid(tenantUuid, context.Background())
	if errorTenant != nil {
		orgService.logger.Error("error find tenant [%w]", zap.Error(errorTenant))
		return nilComposite, errorTenant
	}

	codeUsed, errOrg := orgService.organizationDao.ExistsByCode(*orgCreateReq.Code, c)
	if errOrg != nil {
		orgService.logger.Error("error find organization [%w]", zap.Error(errOrg))
		return nilComposite, errOrg
	}

	if codeUsed {
		return nilComposite, errors.New(commons.OrgAlreadyExistsByCode)
	}

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
				span.RecordError(errRbk)
			}
		} else {
			errCmt := tx.Commit(context.Background())
			if errCmt != nil {
				span.RecordError(errCmt)
				orgService.logger.Error("commit error [%w]", zap.Error(errCmt))
			}
		}
	}()

	orgCid, errCreateOrg := orgService.organizationDao.WithTxCreateOrganization(tx, tenant.Id, *orgCreateReq.Code, *orgCreateReq.Label, *orgCreateReq.Type, c)
	if errCreateOrg != nil {
		span.RecordError(errCreateOrg)
		return nilComposite, errCreateOrg
	}

	sector := model.Sector{}
	sector.TenantId = tenant.Id
	sector.OrganizationId = orgCid.Id
	sector.Code = "root-" + *orgCreateReq.Code
	sector.Label = *orgCreateReq.Label
	sector.Depth = 0
	sector.HasParent = false

	sectorCid, errCreateSector := orgService.sectorDao.WithTxCreateSector(tx, sector, c)

	return sectorCid, errCreateSector
}

func (orgService *OrganizationService) FindAllOrganizations(tenantUuid string, logger zap.Logger, parentContext context.Context) (dtos.OrgLightReponseList, error) {
	var orgsResponse = dtos.OrgLightReponseList{}

	c, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "ORG-LIST-SERVICE")
	defer span.End()

	logger.Info("find all organizations", zap.String("tenantUuid", tenantUuid))
	tenant, errTenant := orgService.tenantDao.FindByUuid(tenantUuid, c)
	if errTenant != nil {
		span.RecordError(errTenant)
		return orgsResponse, errTenant
	}
	orgs, errOrgs := orgService.organizationDao.FindAllByTenantId(tenant.Id, c)
	if errOrgs != nil {
		span.RecordError(errOrgs)
		return orgsResponse, errOrgs
	}
	orgsResponse = converters.ConvertOrgEntityListToOrgLightList(orgs)
	return orgsResponse, nil
}
