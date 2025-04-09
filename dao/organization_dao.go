package dao

import (
	"context"
	"poc-fiber/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

const CONFIG_ORGS = "sql.organizations"

type OrganizationDao struct {
	DbPool *pgxpool.Pool
}

func NewOrganizationDao(pool *pgxpool.Pool) OrganizationDao {
	orgDao := OrganizationDao{}
	orgDao.DbPool = pool
	return orgDao
}

func (orgDao *OrganizationDao) CreateOrganization(tenantId int64, code string, label string, otype string) (model.CompositeId, error) {
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

func (orgDao *OrganizationDao) WithTxCreateOrganization(tx pgx.Tx, tenantId int64, code string, label string, otype string) (model.CompositeId, error) {
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

func (orgDao *OrganizationDao) FindAllByTenantId(tenantId int64) ([]model.Organization, error) {
	var nilOrg []model.Organization
	selStmt := viper.GetStringMapString(CONFIG_ORGS)["findalldisplay"]
	rows, errQry := orgDao.DbPool.Query(context.Background(), selStmt, tenantId)
	if errQry != nil {
		return nil, errQry
	}
	defer rows.Close()
	orgs, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.Organization])
	if errCollect != nil {
		return nilOrg, errCollect
	}
	return orgs, nil
}
