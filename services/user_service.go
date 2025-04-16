package services

import (
	"errors"
	"poc-fiber/commons"
	"poc-fiber/converters"
	"poc-fiber/dao"
	"poc-fiber/dtos"
	"poc-fiber/functions"
	"poc-fiber/model"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type UserService struct {
	tenantFunctions functions.TenantFunctions
	orgsFunctions   functions.OrganizationsFunctions
	userDao         dao.UserDao
}

func NewUserService(tenantFunctions functions.TenantFunctions, orgsFunctions functions.OrganizationsFunctions, userDao dao.UserDao, l zap.Logger) UserService {
	userService := UserService{
		tenantFunctions: tenantFunctions,
		orgsFunctions:   orgsFunctions,
		userDao:         userDao,
	}
	return userService
}

func (userService UserService) CreateUser(tenantUuid string, orgUuid string, createUserReq dtos.CreateUserRequest, logger zap.Logger) (model.CompositeId, error) {
	var nilComposite model.CompositeId
	var nilTenant model.Tenant
	var nilOrg model.Organization

	// Ensure tenant exists
	tenant, errFindTenant := userService.tenantFunctions.FindTenant(tenantUuid, logger)
	if errFindTenant != nil {
		return nilComposite, errFindTenant
	}
	if tenant == nilTenant {
		return nilComposite, errors.New(commons.TenantNotFound)
	}

	// Ensure organization exists
	org, errFindOrg := userService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, logger)
	if errFindOrg != nil {
		return nilComposite, errFindOrg
	}
	if org == nilOrg {
		return nilComposite, errors.New(commons.OrgNotFound)
	}

	// Ensure login not already in use
	loginExists, errLoginExists := userService.userDao.LoginExists(*createUserReq.Login)
	if errLoginExists != nil {
		return nilComposite, errLoginExists
	}
	if loginExists {
		return nilComposite, errors.New(commons.UserLoginAlreadyInUse)
	}

	// Ensure email not already in use
	emailExists, errMailExists := userService.userDao.EmailExists(*createUserReq.Login)
	if errMailExists != nil {
		return nilComposite, errMailExists
	}
	if emailExists {
		return nilComposite, errors.New(commons.UserEmailAlreadyInUse)
	}

	nuuid := uuid.New().String()
	var user = model.User{
		TenantId:       tenant.Id,
		OrganizationId: org.Id,
		Uuid:           nuuid,
		FirstName:      *createUserReq.FirstName,
		LastName:       *createUserReq.LastName,
		Login:          *createUserReq.Login,
		Email:          *createUserReq.Email,
	}
	cid, errCreate := userService.userDao.CreateUser(user)
	if errCreate != nil {
		return nilComposite, errCreate
	}
	return cid, nil
}

func (userService UserService) FindAllUsers(tenantUuid string, orgUuid string, logger zap.Logger) (dtos.UserListResponse, error) {
	var nilTenant model.Tenant
	var nilOrg model.Organization
	var usersList = dtos.UserListResponse{}

	// Ensure tenant exists
	tenant, errFindTenant := userService.tenantFunctions.FindTenant(tenantUuid, logger)
	if errFindTenant != nil {
		return usersList, errFindTenant
	}
	if tenant == nilTenant {
		return usersList, errors.New(commons.TenantNotFound)
	}

	// Ensure organization exists
	org, errFindOrg := userService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, logger)
	if errFindOrg != nil {
		return usersList, errFindOrg
	}
	if org == nilOrg {
		return usersList, errors.New(commons.OrgNotFound)
	}

	users, errList := userService.userDao.FindAllByTenantAndOrganization(tenant.Id, org.Id)
	if errList != nil {
		return usersList, errList
	}

	userArray := make([]dtos.UserResponse, len(users))
	for inc, usr := range users {
		userArray[inc] = converters.ConvertUserToResponse(usr)
	}
	usersList.Users = userArray
	return usersList, nil
}
