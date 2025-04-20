package migrate

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func PerformMigration(log zap.Logger, pgAdminUrl string, migrationFiles string) error {
	db, errOpen := sql.Open("postgres", pgAdminUrl)
	if errOpen != nil {
		return errOpen
	}
	driver, errPosgres := postgres.WithInstance(db, &postgres.Config{})
	if errPosgres != nil {
		return errPosgres
	}

	var migrationUri = "file:" + migrationFiles
	m, errInstantiate := migrate.NewWithDatabaseInstance(
		migrationUri,
		"postgres", driver)
	if errInstantiate != nil {
		return errInstantiate
	}
	log.Info("Performing postgreSQL migration")
	return m.Up()
}
