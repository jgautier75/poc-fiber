package services

import (
	"context"
	"errors"
	"poc-fiber/antlr/parser"
	"poc-fiber/commons"
	"poc-fiber/converters"
	"poc-fiber/dao"
	"poc-fiber/dtos"
	"poc-fiber/functions"
	"poc-fiber/logger"
	"poc-fiber/model"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

const LOGGER_NAME = "UserService"

type UserService struct {
	tenantFunctions functions.TenantFunctions
	orgsFunctions   functions.OrganizationsFunctions
	userDao         dao.UserDao
}

func NewUserService(tenantFunctions functions.TenantFunctions, orgsFunctions functions.OrganizationsFunctions, userDao dao.UserDao) UserService {
	userService := UserService{
		tenantFunctions: tenantFunctions,
		orgsFunctions:   orgsFunctions,
		userDao:         userDao,
	}
	return userService
}

func (userService UserService) CreateUser(tenantUuid string, orgUuid string, createUserReq dtos.CreateUserRequest, parentContext context.Context) (model.CompositeId, error) {
	var nilComposite model.CompositeId

	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "SERVICE-USER-CREATE")
	defer span.End()

	// Ensure tenant exists
	tenant, errFindTenant := userService.tenantFunctions.FindTenant(tenantUuid, c)
	if errFindTenant != nil {
		return nilComposite, errFindTenant
	}

	// Ensure organization exists
	org, errFindOrg := userService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, c)
	if errFindOrg != nil {
		return nilComposite, errFindOrg
	}

	// Ensure login not already in use
	loginExists, errLoginExists := userService.userDao.LoginExists(*createUserReq.Login, c)
	if errLoginExists != nil {
		return nilComposite, errLoginExists
	}
	if loginExists {
		return nilComposite, errors.New(commons.UserLoginAlreadyInUse)
	}

	// Ensure email not already in use
	emailExists, errMailExists := userService.userDao.EmailExists(*createUserReq.Login, c)
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
	cid, errCreate := userService.userDao.CreateUser(user, c)
	if errCreate != nil {
		return nilComposite, errCreate
	}
	return cid, nil
}

func (userService UserService) FindAllUsers(tenantUuid string, orgUuid string, parentContext context.Context) (dtos.UserListResponse, error) {
	var usersList = dtos.UserListResponse{}

	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "SERVICE-USER-LIST")
	defer span.End()
	logger.LogRecord(c, LOGGER_NAME, "find all users for tenant ["+tenantUuid+"] and org ["+orgUuid+"]")

	// Ensure tenant exists
	tenant, errFindTenant := userService.tenantFunctions.FindTenant(tenantUuid, c)
	if errFindTenant != nil {
		return usersList, errFindTenant
	}

	// Ensure organization exists
	org, errFindOrg := userService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, c)
	if errFindOrg != nil {
		return usersList, errFindOrg
	}

	users, errList := userService.userDao.FindAllByTenantAndOrganization(tenant.Id, org.Id, c)
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

func (userService UserService) FilterUsers(tenantUuid string, orgUuid string, expressions []parser.SearchExpression, pagination model.Pagination, parentContext context.Context) (dtos.UserListResponse, error) {
	var usersList = dtos.UserListResponse{}

	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "SERVICE-USER-FILTER")
	defer span.End()
	logger.LogRecord(c, LOGGER_NAME, "filter users for tenant ["+tenantUuid+"] and org ["+orgUuid+"]")

	// Ensure tenant exists
	tenant, errFindTenant := userService.tenantFunctions.FindTenant(tenantUuid, c)
	if errFindTenant != nil {
		return usersList, errFindTenant
	}

	// Ensure organization exists
	org, errFindOrg := userService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, c)
	if errFindOrg != nil {
		return usersList, errFindOrg
	}

	// Filter users
	cnt, users, errUsers := userService.userDao.FilterUsers(tenant.Id, org.Id, expressions, pagination, c)
	if errUsers != nil {
		return dtos.UserListResponse{}, errUsers
	}
	usersList.Total = cnt
	usersList.Pages = (cnt / pagination.RowsPerPage) + 1

	userArray := make([]dtos.UserResponse, len(users))
	for inc, usr := range users {
		userArray[inc] = converters.ConvertUserToResponse(usr)
	}
	usersList.Users = userArray
	return usersList, nil
}

func (userService UserService) DeleteUser(tenantUuid string, orgUuid string, userUuid string, parentContext context.Context) (bool, error) {

	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "SERVICE-USER-DELETE")
	defer span.End()

	// Ensure tenant exists
	tenant, errFindTenant := userService.tenantFunctions.FindTenant(tenantUuid, c)
	if errFindTenant != nil {
		return false, errFindTenant
	}

	// Ensure organization exists
	org, errFindOrg := userService.orgsFunctions.FindOrganization(tenant.Id, orgUuid, c)
	if errFindOrg != nil {
		return false, errFindOrg
	}

	userExists, err := userService.userDao.ExistsByUuid(tenant.Id, org.Id, userUuid, c)
	if err != nil {
		return false, err
	}

	return userExists, nil
}
