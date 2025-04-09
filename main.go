package main

import (
	"fmt"
	"os"
	"os/signal"
	"poc-fiber/certificates"
	"poc-fiber/dao"
	"poc-fiber/endpoints"
	"poc-fiber/infrastructure"
	"poc-fiber/logger"
	"poc-fiber/migrate"
	"poc-fiber/services"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found : [%w]", err))
		} else {
			panic(fmt.Errorf("error regind config : [%w]", err))
		}
	}

	logger := logger.ConfigureLogger(viper.GetString("app.logFile"), true, true)

	// Generate certificates
	certificates.GenerateSelfSignedCerts(logger)

	// Perform SQL Migration
	migrate.PerformMigration(logger, viper.GetString("app.pgUrl"))

	//Setup rdbms connection pool
	dbPool, poolErr := infrastructure.SetupCnxPool(viper.GetString("app.pgUrl"), viper.GetInt32("app.pgPoolMin"), viper.GetInt32("app.pgPoolMax"), logger)
	if poolErr != nil {
		panic(poolErr)
	}
	defer dbPool.Close()

	//Setup Dao & Services
	var tenantDao = dao.NewTenantDao(dbPool)
	var orgDao = dao.NewOrganizationDao(dbPool)
	var sectorDao = dao.NewSectorDao(dbPool)
	var orgService = services.NewOrganizationService(tenantDao, orgDao, sectorDao, logger)

	appContext := viper.GetString("app.server.context")
	apiBaseUri := viper.GetString("app.server.api")
	var fullApiUri = "/" + appContext + "/" + apiBaseUri

	fConfig := endpoints.BuildFiberConfig(viper.GetString("app.name"))
	logger.Info("Application -> Setup")
	app := fiber.New(fConfig)
	app.Get(fullApiUri+"/tenants/:tenantUuid/organizations", endpoints.MakeOrgFindAll(orgService))

	go func() {
		logger.Info("Application -> Listen TLS")
		if errTls := app.ListenTLS(":"+viper.GetString("app.server.port"), "cert.pem", "key.pem"); errTls != nil {
			panic(errTls)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
}
