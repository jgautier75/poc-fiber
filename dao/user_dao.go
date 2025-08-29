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

func (udao UserDao) ExistsByUuid(tenantId int64, orgId int64, userUuid string, parentContext context.Context) (bool, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-USER-UUID_EXISTS")
	defer span.End()

	selStmt := viper.GetStringMapString(CONFIG_USERS)["existsbyuuid"]
	rows, e := udao.DbPool.Query(context.Background(), selStmt, tenantId, orgId, userUuid)
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

	logger.LogRecord(c, LOGGER_NAME, "user ["+userUuid+"] exists: "+strconv.FormatBool(exists))
	return exists, nil
}

func (udao UserDao) FindAllByTenantAndOrganization(tenantId int64, orgId int64, parentContext context.Context) ([]model.User, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-USER-LIST")
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

func (udao UserDao) FindByUuid(tenantId int64, orgId int64, uuid string) (model.User, error) {
	var nilUser model.User

	selStmt := viper.GetStringMapString(CONFIG_USERS)["findbyuuid"]

	rows, e := udao.DbPool.Query(context.Background(), selStmt, tenantId, orgId, uuid)
	if e != nil {
		return nilUser, e
	}
	defer rows.Close()

	user, errCollect := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return user, errCollect
}

func (udao UserDao) FilterUsers(tenantId int64, orgId int64, expressions []parser.SearchExpression, pagination model.Pagination, parentContext context.Context) (int, []model.User, error) {
	c, span := otel.Tracer(logger.OTEL_TRACER_NAME).Start(parentContext, "DAO-USER-FILTER")
	defer span.End()

	// Count number of results based on search expressions
	countStmt := viper.GetStringMapString(CONFIG_USERS)["countbytenantandorg"]
	fullQuery, errBuild, vals := BuildQueryFromExpressions(countStmt, expressions, tenantId, orgId)
	if errBuild != nil {
		return 0, nil, errBuild
	}
	rows, e := udao.DbPool.Query(context.Background(), fullQuery, vals...)
	if e != nil {
		return 0, nil, e
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		err := rows.Scan(&cnt)
		if err != nil {
			return 0, nil, err
		}
	}

	// Search based on expressions
	selStmt := viper.GetStringMapString(CONFIG_USERS)["findallbytenantandorg"]
	fullQuery, errBuild, svalues := BuildQueryFromExpressions(selStmt, expressions, tenantId, orgId)
	if errBuild != nil {
		return 0, nil, errBuild
	}
	qry := computePagination(fullQuery, pagination)
	searchRows, e := udao.DbPool.Query(context.Background(), qry, svalues...)
	if e != nil {
		return 0, nil, e
	}
	defer searchRows.Close()
	users, errCollect := pgx.CollectRows(searchRows, pgx.RowToStructByName[model.User])
	if errCollect != nil {
		return 0, nil, errCollect
	}

	logger.LogRecord(c, LOGGER_NAME, "nb of results ["+strconv.Itoa(len(users))+"]")
	return cnt, users, nil
}

func (udao UserDao) DeleteUser(userId int64) error {
	selStmt := viper.GetStringMapString(CONFIG_USERS)["delete"]
	rows, errQuery := udao.DbPool.Query(context.Background(), selStmt, userId)
	if errQuery != nil {
		return errQuery
	}
	defer rows.Close()
	return errQuery
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
		switch expr.Type {
		case parser.OpeningParenthesis:
			builder.Write([]byte(" ( "))
		case parser.ClosingParenthesis:
			builder.Write([]byte(" ) "))
		case parser.Comparison:
			sqlOp, errComp := convertComparisonExpressionToSql(expr.TextValue)
			if errComp != nil {
				return "", errComp, values
			}
			builder.Write([]byte(sqlOp))
		case parser.Negation:
			builder.Write([]byte(" not "))
		case parser.Property:
			ppty := strings.ReplaceAll(expr.TextValue, "'", "")
			builder.Write([]byte(ppty))
		case parser.Value:
			paramIndex := "$" + strconv.Itoa(inc)
			inc++
			values = append(values, strings.ReplaceAll(strings.Trim(expr.TextValue, "'"), "*", "%"))
			builder.Write([]byte(paramIndex))
		case parser.Operator:
			builder.Write([]byte(expr.TextValue))
		}
	}
	return builder.String(), nil, values
}

func convertComparisonExpressionToSql(exprOperator string) (string, error) {
	switch exprOperator {
	case "gt":
		return ">", nil
	case "ge":
		return ">=", nil
	case "lt":
		return "<", nil
	case "le":
		return "<=", nil
	case "eq":
		return "=", nil
	case "ne":
		return "!=", nil
	case "lk":
		return " like ", nil
	}
	return "", errors.New("invalid comparison [" + exprOperator + "]")
}

func computePagination(qry string, pagination model.Pagination) string {
	var builder strings.Builder
	builder.WriteString(qry)

	if len(pagination.Sorting) > 0 {
		builder.WriteString(" order by")
		for _, s := range pagination.Sorting {
			builder.WriteString(" ")
			builder.WriteString(s.Column)
			builder.WriteString(" ")
			builder.WriteString(s.Order)
		}
	}

	if pagination.Page > 1 {
		startPg := (pagination.Page - 1) * pagination.RowsPerPage
		builder.WriteString(" offset ")
		builder.WriteString(strconv.Itoa(startPg))
	}

	builder.WriteString(" limit ")
	builder.WriteString(strconv.Itoa(pagination.RowsPerPage))
	return builder.String()
}
