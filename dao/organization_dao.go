package dao

import (
	"context"
	"poc-fiber/logger"
	"poc-fiber/model"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

const LOGGER_NAME = "OrganizationDao"
const CONFIG_ORGS = "sql.organizations"
const OTEL_TRACER_NAME = "otel-collector"

type OrganizationDao struct {
	DbPool *pgxpool.Pool
}

func NewOrganizationDao(pool *pgxpool.Pool) OrganizationDao {
	orgDao := OrganizationDao{}
	orgDao.DbPool = pool
	return orgDao
}

func (orgDao *OrganizationDao) CreateOrganization(tenantId int64, code string, label string, otype string, parentContext context.Context) (model.CompositeId, error) {
	_, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-ORG-CREATE")
	defer span.End()

	insertStmt := viper.GetStringMapString(CONFIG_ORGS)["create"]
	nuuid := uuid.New().String()
	var id int64
	errQuery := orgDao.DbPool.QueryRow(context.Background(), insertStmt, tenantId, nuuid, code, label, otype).Scan(&id)
	compId := model.CompositeId{
		Id:   id,
		Uuid: nuuid,
	}
	return compId, errQuery
}

func (orgDao *OrganizationDao) WithTxCreateOrganization(tx pgx.Tx, tenantId int64, code string, label string, otype string, parentContext context.Context) (model.CompositeId, error) {
	_, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-ORG-CREATE_TX")
	defer span.End()

	insertStmt := viper.GetStringMapString(CONFIG_ORGS)["create"]
	nuuid := uuid.New().String()
	var id int64
	errQuery := tx.QueryRow(context.Background(), insertStmt, tenantId, nuuid, code, label, otype).Scan(&id)
	compId := model.CompositeId{
		Id:   id,
		Uuid: nuuid,
	}
	return compId, errQuery
}

func (orgDao *OrganizationDao) updateLabel(uuid string, nlabel string) error {
	updateStmt := viper.GetStringMapString(CONFIG_ORGS)["updatelabelbyuuid"]
	_, errQuery := orgDao.DbPool.Exec(context.Background(), updateStmt, nlabel, uuid)
	return errQuery
}

func (orgDao *OrganizationDao) FindAllByTenantId(tenantId int64, parentContext context.Context) ([]model.Organization, error) {
	var nilOrg []model.Organization
	c, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "DAO-ORG-LIST")
	defer span.End()
	selStmt := viper.GetStringMapString(CONFIG_ORGS)["findalldisplay"]
	logger.LogRecord(c, LOGGER_NAME, "select statement ["+selStmt+"]")
	rows, errQry := orgDao.DbPool.Query(context.Background(), selStmt, tenantId)
	if errQry != nil {
		span.RecordError(errQry)
		return nil, errQry
	}
	defer rows.Close()
	orgs, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.Organization])
	if errCollect != nil {
		span.RecordError(errCollect)
		return nilOrg, errCollect
	}
	logger.LogRecord(c, LOGGER_NAME, "nb of results ["+strconv.Itoa(len(orgs))+"]")
	return orgs, nil
}

func (orgDao *OrganizationDao) FindByTenantAndUuid(tenantId int64, orgUuid string, parentContext context.Context) (model.Organization, error) {
	var nilOrg model.Organization

	_, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "DAO-ORG-FIND_BY_UUID")
	defer span.End()

	sqlOrgsMaps := viper.GetStringMapString(CONFIG_ORGS)
	rows, e := orgDao.DbPool.Query(context.Background(), sqlOrgsMaps["findbytenantanduuid"], tenantId, orgUuid)
	if e != nil {
		return nilOrg, e
	}
	defer rows.Close()
	org, errCollect := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Organization])
	if errCollect != nil {
		if errCollect.Error() == "no rows in result set" {
			return nilOrg, nil
		}
		return nilOrg, errCollect
	}
	return org, nil
}

func (orgDao *OrganizationDao) ExistsByCode(code string, parentContext context.Context) (bool, error) {
	c, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "DAO-ORG-EXISTS_BY_CODE")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_ORGS)["existsbycode"]
	rows, e := orgDao.DbPool.Query(context.Background(), selStmt, code)
	if e != nil {
		return false, e
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		err := rows.Scan(&cnt)
		if err != nil {
			return false, err
		}
	}

	var exists = false
	if cnt > 0 {
		exists = true
	}
	logger.LogRecord(c, LOGGER_NAME, "exists by code ["+code+"]: "+strconv.FormatBool(exists))
	return exists, nil
}
