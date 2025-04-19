package dao

import (
	"context"
	"poc-fiber/logger"
	"poc-fiber/model"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

const CONFIG_USERS = "sql.users"

type UserDao struct {
	DbPool *pgxpool.Pool
}

func NewUserDaao(pool *pgxpool.Pool) UserDao {
	userDao := UserDao{}
	userDao.DbPool = pool
	return userDao
}

func (udao UserDao) CreateUser(user model.User, parentContext context.Context) (model.CompositeId, error) {
	_, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-USER-CREATE")
	defer span.End()

	insertStmt := viper.GetStringMapString(CONFIG_USERS)["create"]
	var id int64
	errQuery := udao.DbPool.QueryRow(context.Background(), insertStmt, user.TenantId, user.OrganizationId, user.Uuid, user.LastName, user.FirstName, user.Login, user.Email).Scan(&id)
	compId := model.CompositeId{
		Id:   id,
		Uuid: user.Uuid,
	}
	return compId, errQuery
}

func (udao UserDao) LoginExists(login string, parentContext context.Context) (bool, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-USER-LOGIN_EXISTS")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_USERS)["loginexists"]
	rows, e := udao.DbPool.Query(context.Background(), selStmt, login)
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

	logger.LogRecord(c, LOGGER_NAME, "login ["+login+"] already used: "+strconv.FormatBool(exists))
	return exists, nil
}

func (udao UserDao) EmailExists(email string, parentContext context.Context) (bool, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-USER-EMAIL_EXISTS")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_USERS)["emailexists"]
	rows, e := udao.DbPool.Query(context.Background(), selStmt, email)
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

	logger.LogRecord(c, LOGGER_NAME, "email ["+email+"] already used: "+strconv.FormatBool(exists))
	return exists, nil
}

func (udao UserDao) FindAllByTenantAndOrganization(tenantId int64, orgId int64, parentContext context.Context) ([]model.User, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "USER-LIST-DAO")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_USERS)["findallbytenantandorg"]
	rows, e := udao.DbPool.Query(context.Background(), selStmt, tenantId, orgId)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	users, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if errCollect != nil {
		return nil, errCollect
	}

	logger.LogRecord(c, LOGGER_NAME, "nb of results ["+strconv.Itoa(len(users))+"]")
	return users, nil
}
