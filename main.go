package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"os/signal"
	"poc-fiber/certificates"
	"poc-fiber/dao"
	"poc-fiber/endpoints"
	"poc-fiber/functions"
	"poc-fiber/infrastructure"
	"poc-fiber/logger"
	"poc-fiber/middleware"
	"poc-fiber/migrate"
	"poc-fiber/oauth"
	"poc-fiber/opentelemetry"
	"poc-fiber/services"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	openbao "github.com/openbao/openbao/api/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func main() {

	gob.Register(time.Time{})
	gob.Register(oauth2.Token{})

	// Read main config
	viper.SetEnvPrefix("EV")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found : [%w]", err))
		} else {
			panic(fmt.Errorf("error reading config : [%w]", err))
		}
	}

	// Read sql queries config
	viper.SetConfigName("sql")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.MergeInConfig()

	logger := logger.ConfigureLogger(viper.GetString("app.logFile"), true, true)

	// OpenBao
	config := openbao.DefaultConfig()
	config.Address = viper.GetString("vault.endpoint")
	client, err := openbao.NewClient(config)
	if err != nil {
		panic(fmt.Errorf("unable to initialize OpenBao client: %v", err))
	}
	client.SetToken(viper.GetString("vault.token"))

	secret, err := client.KVv2(viper.GetString("vault.path")).Get(context.Background(), viper.GetString("vault.data"))
	if err != nil {
		panic(fmt.Errorf("unable to read secret: %v", err))
	}
	vaultData := secret.Data

	// Generate certificates
	certificates.GenerateSelfSignedCerts(logger)

	// Perform SQL Migration
	pgUrl := vaultData["pgUrl"].(string)
	errMig := migrate.PerformMigration(logger, pgUrl, "migrate/files")
	if errMig != nil && errMig.Error() != "no change" {
		panic(errMig)
	}

	//Setup rdbms connection pool
	dbPool, poolErr := infrastructure.SetupCnxPool(pgUrl, viper.GetInt32("app.pgPoolMin"), viper.GetInt32("app.pgPoolMax"), logger)
	if poolErr != nil {
		panic(poolErr)
	}
	defer dbPool.Close()

	clientId := vaultData["clientId"].(string)
	clientSecret := vaultData["clientSecret"].(string)
	appContext := viper.GetString("app.server.context")

	// Opentelemetry
	logger.Info("OpenTelemetry > Setup")
	otelShutdownFuncs, otmMetrics, errOtel := opentelemetry.Setup(context.Background())
	if errOtel != nil {
		panic(fmt.Errorf("error setting up opentelemetry [%w]", errOtel))
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		for _, shutdownFn := range otelShutdownFuncs {
			errShutdown := shutdownFn(context.Background())
			logger.Error("error shutting down opentelemetry", zap.Error(errShutdown))
		}
	}()

	//Setup Dao & Services
	var tenantDao = dao.NewTenantDao(dbPool)
	var orgDao = dao.NewOrganizationDao(dbPool)
	var sectorDao = dao.NewSectorDao(dbPool)
	var userDao = dao.NewUserDaao(dbPool)
	var tenantFunctions = functions.NewTenantFunctions(tenantDao, logger)
	var orgFunctions = functions.NewOrganizationsFunctions(orgDao, logger)

	var orgService = services.NewOrganizationService(tenantDao, orgDao, sectorDao)
	var sectorService = services.NewSectorService(tenantFunctions, orgFunctions, sectorDao)
	var userService = services.NewUserService(tenantFunctions, orgFunctions, userDao)

	apiBaseUri := viper.GetString("app.server.api")
	var fullApiUri = "/" + appContext + "/" + apiBaseUri
	var versionsApi = fullApiUri + "/versions"
	var apiOrgsPrefix = fullApiUri + "/tenants/:tenantUuid/organizations"
	var apiSectorsPrefix = apiOrgsPrefix + "/:organizationUuid/sectors"
	var apiUsersPrefix = apiOrgsPrefix + "/:organizationUuid/users"

	// Redis setup (session storage)
	defCfg := session.ConfigDefault
	redisStorage := endpoints.ConfigureRedisStorage(viper.GetString("redis.host"), viper.GetInt("redis.port"))
	defCfg.Storage = redisStorage
	store := session.New(defCfg)

	// Fiber endpoints
	fConfig := endpoints.BuildFiberConfig(viper.GetString("app.name"))
	logger.Info("Application -> Setup")
	app := fiber.New(fConfig)

	// Fetch OIDC .well-known url
	logger.Info("OIDC -> Fetch .well-known url [" + viper.GetString("oauth2.issuer") + "]")
	fOAuth := oauth.NewOAuthManager()
	authMgr, errFetch := fOAuth.InitOAuthManager(context.Background(), logger, clientId, clientSecret)
	if errFetch != nil {
		panic(fmt.Errorf("error fetching .well-known issuer: [%w]", errFetch))
	}

	logger.Info("Middleware -> Setup")
	app.Use(middleware.InitOidcMiddleware(authMgr, fullApiUri, versionsApi, store, clientId, clientSecret))
	app.Use(middleware.HttpMiddleWareStats(otmMetrics))

	logger.Info("Endpoints -> Setup")
	app.Get("/"+appContext+"/home", endpoints.MakeIndex(authMgr.OAuthConfig, store))
	app.Get(versionsApi, endpoints.MakeVersions(viper.GetString("app.version")))
	app.Delete(fullApiUri+"/sessions", endpoints.DeleteSession(clientId, clientSecret, store, authMgr.OAuthEndpoints, logger))

	// OIDC
	app.Get(authMgr.OAuthCallBackUri, endpoints.MakeOAuthCallback(authMgr.OAuthConfig, store, authMgr.Verifier))

	// Organizations
	app.Get(apiOrgsPrefix, endpoints.MakeOrgFindAll(orgService))
	app.Post(apiOrgsPrefix, endpoints.MakeOrgCreate(orgService))

	// Sectors
	app.Get(apiSectorsPrefix, endpoints.MakeSectorsFindAll(sectorService))
	app.Post(apiSectorsPrefix, endpoints.MakeSectorCreate(sectorService))

	// Users
	app.Get(apiUsersPrefix, endpoints.MakeUsersList(userService))
	app.Post(apiUsersPrefix, endpoints.MakeUserCreate(userService))
	app.Get(apiUsersPrefix+"/filter", endpoints.MakeUsersFilter(userService))

	go func() {
		logger.Info("Application -> Listen TLS")
		if errTls := app.ListenTLS(":"+viper.GetString("app.server.port"), "cert.pem", "key.pem"); errTls != nil {
			panic(errTls)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()
}
