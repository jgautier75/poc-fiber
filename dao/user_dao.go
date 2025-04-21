package dao

import (
	"context"
	"errors"
	"poc-fiber/antlr/parser"
	"poc-fiber/logger"
	"poc-fiber/model"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

const CONFIG_USERS = "sql.users"

type SqlSearchCriteria struct {
	Key   string
	Value string
}

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

func (udao UserDao) FilterUsers(tenantId int64, orgId int64, expressions []parser.SearchExpression, parentContext context.Context) ([]model.User, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "USER-FILTER-DAO")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_USERS)["findallbytenantandorg"]
	fullQuery, errBuild, vals := BuildQueryFromExpressions(selStmt, expressions, tenantId, orgId)
	if errBuild != nil {
		return nil, errBuild
	}

	rows, e := udao.DbPool.Query(context.Background(), fullQuery, vals...)
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

func BuildQueryFromExpressions(baseQuery string, expressions []parser.SearchExpression, tenantId int64, orgId int64) (query string, err error, params []interface{}) {
	var builder strings.Builder
	var values []interface{}
	builder.Write([]byte(baseQuery))

	if len(expressions) > 0 {
		builder.Write([]byte(" and "))
	}
	var inc = 1

	// Add default parameters
	// tenant id
	values = append(values, tenantId)
	inc++
	// organization id
	values = append(values, orgId)
	inc++

	for _, expr := range expressions {
		if expr.Type == parser.OpeningParenthesis {
			builder.Write([]byte(" ( "))
		} else if expr.Type == parser.ClosingParenthesis {
			builder.Write([]byte(" ) "))
		} else if expr.Type == parser.Comparison {
			sqlOp, errComp := convertComparisonExpressionToSql(expr.TextValue)
			if errComp != nil {
				return "", errComp, values
			}
			builder.Write([]byte(sqlOp))
		} else if expr.Type == parser.Negation {
			builder.Write([]byte(" not "))
		} else if expr.Type == parser.Property {
			ppty := strings.ReplaceAll(expr.TextValue, "'", "")
			builder.Write([]byte(ppty))
		} else if expr.Type == parser.Value {
			paramIndex := "$" + strconv.Itoa(inc)
			inc++
			values = append(values, strings.Trim(expr.TextValue, "'"))
			builder.Write([]byte(paramIndex))
		} else if expr.Type == parser.Operator {
			builder.Write([]byte(expr.TextValue))
		}
	}
	return builder.String(), nil, values
}

func convertComparisonExpressionToSql(exprOperator string) (string, error) {
	if exprOperator == "gt" {
		return ">", nil
	} else if exprOperator == "ge" {
		return ">=", nil
	} else if exprOperator == "lt" {
		return "<", nil
	} else if exprOperator == "le" {
		return "<=", nil
	} else if exprOperator == "eq" {
		return "=", nil
	} else if exprOperator == "ne" {
		return "!=", nil
	} else if exprOperator == "lk" {
		return " like ", nil
	}
	return "", errors.New("invalid comparison [" + exprOperator + "]")
}
