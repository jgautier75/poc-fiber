package functions

import (
	"errors"
	"poc-fiber/dao"
	"poc-fiber/exceptions"
	"poc-fiber/model"

	"go.uber.org/zap"
)

type TenantFunctions struct {
	tenantDao dao.TenantDao
	logger    zap.Logger
}

func NewTenantFunctions(tenantDao dao.TenantDao, logger zap.Logger) TenantFunctions {
	tenantFunctions := TenantFunctions{
		tenantDao: tenantDao,
		logger:    logger,
	}
	return tenantFunctions
}

func (tf *TenantFunctions) FindTenant(uuid string, logger zap.Logger) (model.Tenant, error) {
	var nilTenant model.Tenant
	logger.Info("find tenant", zap.String("uuid", uuid))
	tenant, errFind := tf.tenantDao.FindByUuid(uuid)
	if errFind != nil {
		return nilTenant, errFind
	}
	if tenant == nilTenant {
		return nilTenant, errors.New(exceptions.TENANT_NOT_FOUND)
	} else {
		return tenant, nil
	}
}
