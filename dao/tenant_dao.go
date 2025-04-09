package dao

import (
	"context"
	"poc-fiber/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type TenantDao struct {
	DbPool *pgxpool.Pool
}

func NewTenantDao(pool *pgxpool.Pool) TenantDao {
	tenantDao := TenantDao{}
	tenantDao.DbPool = pool
	return tenantDao
}

func (tenantDao *TenantDao) FindByCode(code string) (model.Tenant, error) {
	var nilTenant model.Tenant
	sqlTenantsMaps := viper.GetStringMapString("sql.tenants")
	rows, e := tenantDao.DbPool.Query(context.Background(), sqlTenantsMaps["findbycode"], code)
	if e != nil {
		return nilTenant, e
	}
	defer rows.Close()
	tenant, errCollect := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Tenant])
	if errCollect != nil {
		return nilTenant, errCollect
	}
	return tenant, nil
}

func (tenantDao *TenantDao) FindByUuid(uuid string) (model.Tenant, error) {
	var nilTenant model.Tenant
	sqlTenantsMaps := viper.GetStringMapString("sql.tenants")
	rows, e := tenantDao.DbPool.Query(context.Background(), sqlTenantsMaps["findbyuuid"], uuid)
	if e != nil {
		return nilTenant, e
	}
	defer rows.Close()
	tenant, errCollect := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Tenant])
	if errCollect != nil {
		return nilTenant, errCollect
	}
	return tenant, nil
}

func (tenantDao TenantDao) UpdateLabel(uuid string, newLabel string) error {
	sqlTenantsMaps := viper.GetStringMapString("sql.tenants")
	_, errQuery := tenantDao.DbPool.Exec(context.Background(), sqlTenantsMaps["updatelabelbyuuid"], newLabel, uuid)
	return errQuery
}
