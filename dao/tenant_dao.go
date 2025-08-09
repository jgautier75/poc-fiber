package dao

import (
	"context"
	"poc-fiber/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

const CONFIG_KEY = "sql.tenants"

type TenantDao struct {
	DbPool *pgxpool.Pool
}

func NewTenantDao(pool *pgxpool.Pool) TenantDao {
	tenantDao := TenantDao{}
	tenantDao.DbPool = pool
	return tenantDao
}

func (tenantDao TenantDao) FindByCode(code string, parentContext context.Context) (model.Tenant, error) {
	var nilTenant model.Tenant

	_, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "DAO-TENANT-FIND_BY_CODE")
	defer span.End()

	sqlTenantsMaps := viper.GetStringMapString(CONFIG_KEY)
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

func (tenantDao TenantDao) FindByUuid(uuid string, parentContext context.Context) (model.Tenant, error) {
	var nilTenant model.Tenant

	_, span := otel.Tracer(OTEL_TRACER_NAME).Start(parentContext, "DAO-TENANT-FIND_BY_UUID")
	defer span.End()

	sqlTenantsMaps := viper.GetStringMapString(CONFIG_KEY)
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
	sqlTenantsMaps := viper.GetStringMapString(CONFIG_KEY)
	_, errQuery := tenantDao.DbPool.Exec(context.Background(), sqlTenantsMaps["updatelabelbyuuid"], newLabel, uuid)
	return errQuery
}
