package dao

import (
	"context"
	"fmt"
	"poc-fiber/logger"
	"poc-fiber/migrate"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

func TestOrganizationDao(t *testing.T) {
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

	orgDao := NewOrganizationDao(dbPool)
	cid, errCreate := orgDao.CreateOrganization(1, "test-org-code", "test-org-label", "community", ctx)
	assert.Nil(t, errCreate, "no error creating organization")
	assert.NotNil(t, cid, "composite id not null")
	consoleLogger.Info("created uuid", zap.String("uuid", cid.Uuid))

	org, errFind := orgDao.FindByTenantAndUuid(1, cid.Uuid, ctx)
	assert.Nil(t, errFind, "no error finding organization")
	assert.NotNil(t, org, "organization not null")

	orgExists, errExistsCode := orgDao.ExistsByCode("test-org-code", ctx)
	assert.Nil(t, errExistsCode, "no error finding organization by code")
	assert.True(t, orgExists, "organization exists")

	orgList, errList := orgDao.FindAllByTenantId(DEFAULT_TENANT, ctx)
	assert.Nil(t, errList, "no error listing organizations")
	assert.NotNil(t, orgList, "organizations list not nil")
	assert.Equal(t, 1, len(orgList), "1 organization found")
}
