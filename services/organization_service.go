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
)

const OTEL_TRACER_NAME = "otel-collector"

type OrganizationService struct {
	tenantDao       dao.TenantDao
	organizationDao dao.OrganizationDao
	sectorDao       dao.SectorDao
}

func NewOrganizationService(tDao dao.TenantDao, oDao dao.OrganizationDao, sDao dao.SectorDao) OrganizationService {
	orgService := OrganizationService{
		tenantDao:       tDao,
		organizationDao: oDao,
		sectorDao:       sDao,
	}
	return orgService
}

func (orgService *OrganizationService) CreateOrganization(tenantUuid string, orgCreateReq dtos.CreateOrgRequest, parentContext context.Context) (model.CompositeId, error) {
	var nilComposite model.CompositeId

	c, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "SERVICE-ORG-CREATE")
	defer span.End()

	// Find tenant
	tenant, errorTenant := orgService.tenantDao.FindByUuid(tenantUuid, context.Background())
	if errorTenant != nil {
		span.RecordError(errorTenant)
		return nilComposite, errorTenant
	}

	codeUsed, errOrg := orgService.organizationDao.ExistsByCode(*orgCreateReq.Code, c)
	if errOrg != nil {
		span.RecordError(errOrg)
		return nilComposite, errOrg
	}

	if codeUsed {
		return nilComposite, errors.New(commons.OrgAlreadyExistsByCode)
	}

	// Get connection and init transaction
	conn, errConnect := orgService.tenantDao.DbPool.Acquire(context.Background())
	if errConnect != nil {
		span.RecordError(errConnect)
		return nilComposite, errConnect
	}
	tx, errTx := conn.BeginTx(context.Background(), pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead})
	if errTx != nil {
		span.RecordError(errTx)
		return nilComposite, errTx
	}
	defer func() {
		if errTx != nil {
			errRbk := tx.Rollback(context.Background())
			if errRbk != nil {
				span.RecordError(errRbk)
			}
		} else {
			errCmt := tx.Commit(context.Background())
			if errCmt != nil {
				span.RecordError(errCmt)
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
	if errCreateSector != nil {
		span.RecordError(errCreateSector)
		return nilComposite, errCreateSector
	}

	return sectorCid, errCreateSector
}

func (orgService *OrganizationService) FindAllOrganizations(tenantUuid string, parentContext context.Context) (dtos.OrgLightReponseList, error) {
	var orgsResponse = dtos.OrgLightReponseList{}

	c, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "SERVICE-ORG-LIST")
	defer span.End()

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
