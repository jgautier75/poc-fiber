package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func PerformMigration(log *zap.Logger, pgAdminUrl string) {
	db, err := sql.Open("postgres", pgAdminUrl)
	if err != nil {
		panic(fmt.Errorf("error initializing database [%w]", err))
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("SQL instance [%w]", err))
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:migrate/files",
		"postgres", driver)
	if err != nil {
		panic(fmt.Errorf("database instance [%w]", err))
	}
	log.Info("Performing postgreSQL migration")
	errMig := m.Up()
	if errMig != nil {
		panic(errMig)
	}

}
