package dao

import (
	"context"
	"fmt"
	"poc-fiber/logger"
	"poc-fiber/migrate"
	"poc-fiber/model"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

func TestUserDao(t *testing.T) {
	ctx := context.Background()
	consoleLogger := logger.ConfigureConsoleLogger()

	dbContainer, host, port, errCreate := CreatePgContainerTest(ctx, *consoleLogger)
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			consoleLogger.Error("error terminating container", zap.Error(err))
		}
	}(*dbContainer, ctx)
	if errCreate != nil {
		panic(errCreate)
	}

	// Retrieve postgreSQL container host and port
	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", PG_USER, PG_PASS, host, port.Int(), PG_DB)
	consoleLogger.Info("Postgresql", zap.String("pgUrl", pgUrl))

	errMigration := migrate.PerformMigration(*consoleLogger, pgUrl, "../migrate/files")
	if errMigration != nil {
		panic(errMigration)
	}

	dbConfig, errDbCfg := pgxpool.ParseConfig(pgUrl)
	if errDbCfg != nil {
		panic(errDbCfg)
	}
	dbPool, poolErr := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if poolErr != nil {
		panic(poolErr)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found : [%w]", err))
		} else {
			panic(fmt.Errorf("error reading config : [%w]", err))
		}
	}

	userDao := UserDao{
		DbPool: dbPool,
	}
	orgDao := OrganizationDao{
		DbPool: dbPool,
	}

	orgCid, errOrg := orgDao.CreateOrganization(DEFAULT_TENANT, "test-org-code", "Test Org Code", "community", ctx)
	if errOrg != nil {
		assert.Nil(t, errOrg, "no error creating organization")
	}
	user := model.User{
		TenantId:       DEFAULT_TENANT,
		Uuid:           uuid.New().String(),
		OrganizationId: orgCid.Id,
		LastName:       "test-last-name",
		FirstName:      "test-user-name",
		Login:          "test-login",
		Email:          "test@test.fr",
	}
	cidUser, errCreateUser := userDao.CreateUser(user, ctx)
	if errCreateUser != nil {
		assert.Nil(t, errCreateUser, "no error creating user")
	}
	assert.NotNil(t, cidUser, "generated user composite id not null")
	consoleLogger.Info("user", zap.String("uuid", user.Uuid))

	emailUsed, errEmail := userDao.EmailExists("test@test.fr", ctx)
	assert.Nil(t, errEmail, "no error user exists by email")
	assert.True(t, emailUsed, "email already in use")

	loginUsed, errLogin := userDao.LoginExists("test-login", ctx)
	assert.Nil(t, errLogin, "no error user exists by login")
	assert.True(t, loginUsed, "login already in use")

	users, errList := userDao.FindAllByTenantAndOrganization(DEFAULT_TENANT, orgCid.Id, ctx)
	assert.Nil(t, errList, "no error listing users")
	assert.Equal(t, 1, len(users), "1 user in result")
}
