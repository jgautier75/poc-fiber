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

const LOGGER_NAME = "OrganizationsFunctions"

type OrganizationsFunctions struct {
	organizationDao dao.OrganizationDao
	logger          zap.Logger
}

func NewOrganizationsFunctions(organizationDao dao.OrganizationDao, logger zap.Logger) OrganizationsFunctions {
	orgFunctions := OrganizationsFunctions{
		organizationDao: organizationDao,
		logger:          logger,
	}
	return orgFunctions
}

func (of *OrganizationsFunctions) FindOrganization(tenantId int64, uuid string, parentContext context.Context) (model.Organization, error) {
	var nilOrg model.Organization

	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "ORG-FIND-FUNC")
	defer span.End()
	logger.LogRecord(c, LOGGER_NAME, "find organization ["+uuid+"]")
	org, errFind := of.organizationDao.FindByTenantAndUuid(tenantId, uuid, c)
	if errFind != nil {
		return nilOrg, errFind
	}
	if org == nilOrg {
		return nilOrg, errors.New(commons.OrgNotFound)
	} else {
		return org, nil
	}
}
