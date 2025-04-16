package functions

import (
	"errors"
	"poc-fiber/commons"
	"poc-fiber/dao"
	"poc-fiber/model"

	"go.uber.org/zap"
)

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

func (of *OrganizationsFunctions) FindOrganization(tenantId int64, uuid string, logger zap.Logger) (model.Organization, error) {
	var nilOrg model.Organization
	logger.Info("find organization", zap.String("uuid", uuid))
	org, errFind := of.organizationDao.FindByTenantAndUuid(tenantId, uuid)
	if errFind != nil {
		return nilOrg, errFind
	}
	if org == nilOrg {
		return nilOrg, errors.New(commons.OrgNotFound)
	} else {
		return org, nil
	}
}
