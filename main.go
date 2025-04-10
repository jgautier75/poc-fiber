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
	"poc-fiber/infrastructure"
	"poc-fiber/logger"
	"poc-fiber/middleware"
	"poc-fiber/migrate"
	"poc-fiber/services"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func main() {

	gob.Register(time.Time{})

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found : [%w]", err))
		} else {
			panic(fmt.Errorf("error reading config : [%w]", err))
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

	appBase := viper.GetString("app.server.base")
	appContext := viper.GetString("app.server.context")
	appPort := viper.GetString("app.server.port")
	oauthCallBackUri := "/" + appContext + "/oauth2/callback"
	oauthCallBackFull := appBase + ":" + appPort + oauthCallBackUri

	// Setup OIDC - Fetch .well-known endpoint  asynchronously
	var provider *oidc.Provider
	asyncOidcIssuer := make(chan *oidc.Provider, 1)
	go func() {
		oidcprov, oidcError := oidc.NewProvider(context.Background(), viper.GetString("oauth2.issuer"))
		if oidcError != nil {
			panic(fmt.Errorf("oidc prodider error: [%w]", oidcError))
		}
		asyncOidcIssuer <- oidcprov
	}()

	provider = <-asyncOidcIssuer

	oauth2Config := oauth2.Config{
		ClientID:     viper.GetString("oauth2.clientId"),
		ClientSecret: viper.GetString("oauth2.clientSecret"),
		RedirectURL:  oauthCallBackFull,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "offline_access"},
	}
	tokenVerifier := provider.Verifier(&oidc.Config{ClientID: viper.GetString("oauth2.clientId")})

	//Setup Dao & Services
	var tenantDao = dao.NewTenantDao(dbPool)
	var orgDao = dao.NewOrganizationDao(dbPool)
	var sectorDao = dao.NewSectorDao(dbPool)
	var orgService = services.NewOrganizationService(tenantDao, orgDao, sectorDao, logger)

	apiBaseUri := viper.GetString("app.server.api")
	var fullApiUri = "/" + appContext + "/" + apiBaseUri

	// Redis setup (session storage)
	defCfg := session.ConfigDefault
	redisStorage := endpoints.ConfigureRedisStorage(viper.GetString("redis.host"), viper.GetInt("redis.port"))
	defCfg.Storage = redisStorage
	store := session.New(defCfg)

	// Fiber endpoints
	fConfig := endpoints.BuildFiberConfig(viper.GetString("app.name"))
	logger.Info("Application -> Setup")
	app := fiber.New(fConfig)
	app.Use(middleware.NewApiOidcHandler(fullApiUri, tokenVerifier))
	app.Get(fullApiUri+"/tenants/:tenantUuid/organizations", endpoints.MakeOrgFindAll(orgService))
	app.Get("/"+appContext+"/home", endpoints.MakeIndex(oauth2Config, store))
	app.Get(oauthCallBackUri, endpoints.MakeOAuthCallback(oauth2Config, store, tokenVerifier))

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
	_ = app.Shutdown()
}
