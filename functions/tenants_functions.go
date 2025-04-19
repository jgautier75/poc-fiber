package functions

import (
	"context"
	"errors"
	"poc-fiber/commons"
	"poc-fiber/dao"
	"poc-fiber/logger"
	"poc-fiber/model"

	"go.opentelemetry.io/otel"
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

func (tf *TenantFunctions) FindTenant(uuid string, parentContext context.Context) (model.Tenant, error) {
	var nilTenant model.Tenant

	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "TENANT-FIND-FUNC")
	defer span.End()

	logger.LogRecord(c, "TenantsFunctions", "find tenant ["+uuid+"]")
	tenant, errFind := tf.tenantDao.FindByUuid(uuid, c)
	if errFind != nil {
		return nilTenant, errFind
	}
	if tenant == nilTenant {
		return nilTenant, errors.New(commons.TenantNotFound)
	} else {
		return tenant, nil
	}
}
