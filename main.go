package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/signal"
	"poc-fiber/authentik"
	"poc-fiber/certificates"
	"poc-fiber/dao"
	"poc-fiber/endpoints"
	"poc-fiber/functions"
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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	gob.Register(time.Time{})
	gob.Register(oauth2.Token{})

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
	var otelServiceName = semconv.ServiceNameKey.String(viper.GetString("app.name"))
	var otelServiceVersion = semconv.ServiceVersionKey.String(viper.GetString("app.version"))

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

	clientId := viper.GetString("oauth2.clientId")
	clientSecret := viper.GetString("oauth2.clientSecret")
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

	var oauthConfig *authentik.OauthConfiguration
	asyncOAuthConfig := make(chan *authentik.OauthConfiguration)
	go func() {
		authCfg := authentik.FetchOAuthConfiguration(viper.GetString("oauth2.issuer"), logger)
		asyncOAuthConfig <- authCfg
	}()
	oauthConfig = <-asyncOAuthConfig

	oauth2Config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  oauthCallBackFull,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "offline_access"},
	}
	tokenVerifier := provider.Verifier(&oidc.Config{ClientID: clientId})

	// Opentelemetry
	logger.Info("OpenTelemetry > Setup")
	otelResource, errResource := resource.New(context.Background(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			otelServiceName,
			// The service version used to display traces in backends
			otelServiceVersion,
		),
	)
	if errResource != nil {
		panic(fmt.Errorf("error setting up opentelemetry resource [%w]", errResource))
	}

	grpcCLientCon, errGrpcCon := initGrpcConn(viper.GetString("otel.endpoint"))
	if errGrpcCon != nil {
		panic(fmt.Errorf("error setting up opentelemetry grpc connection [%w]", errGrpcCon))
	}

	shutdownTracerProvider, err := initTracerProvider(context.Background(), otelResource, grpcCLientCon)
	if err != nil {
		panic(fmt.Errorf("error setting up opentelemetry tracer provider [%w]", errGrpcCon))
	}
	defer func() {
		if err := shutdownTracerProvider(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	//Setup Dao & Services
	var tenantDao = dao.NewTenantDao(dbPool)
	var orgDao = dao.NewOrganizationDao(dbPool)
	var sectorDao = dao.NewSectorDao(dbPool)
	var userDao = dao.NewUserDaao(dbPool)
	var tenantFunctions = functions.NewTenantFunctions(tenantDao, logger)
	var orgFunctions = functions.NewOrganizationsFunctions(orgDao, logger)

	var orgService = services.NewOrganizationService(tenantDao, orgDao, sectorDao, logger)
	var sectorService = services.NewSectorService(tenantFunctions, orgFunctions, sectorDao, logger)
	var userService = services.NewUserService(tenantFunctions, orgFunctions, userDao, logger)

	apiBaseUri := viper.GetString("app.server.api")
	var fullApiUri = "/" + appContext + "/" + apiBaseUri
	var versionsApi = fullApiUri + "/versions"
	var apiOrgsPrefix = fullApiUri + "/tenants/:tenantUuid/organizations"
	var apiSectorsPrefix = fullApiUri + apiOrgsPrefix + "/:organizationUuid/sectors"
	var apiUsersPrefix = fullApiUri + apiOrgsPrefix + "/:organizationUuid/users"

	// Redis setup (session storage)
	defCfg := session.ConfigDefault
	redisStorage := endpoints.ConfigureRedisStorage(viper.GetString("redis.host"), viper.GetInt("redis.port"))
	defCfg.Storage = redisStorage
	store := session.New(defCfg)

	// Fiber endpoints
	fConfig := endpoints.BuildFiberConfig(viper.GetString("app.name"))
	logger.Info("Application -> Setup")
	app := fiber.New(fConfig)

	logger.Info("Middleware -> Setup")
	app.Use(middleware.NewApiOidcHandler(fullApiUri, versionsApi, provider, tokenVerifier, store, clientId, clientSecret))

	logger.Info("Endpoints -> Setup")
	app.Get("/"+appContext+"/home", endpoints.MakeIndex(oauth2Config, store))
	app.Get(apiOrgsPrefix, endpoints.MakeOrgFindAll(orgService, logger))
	app.Post(apiOrgsPrefix, endpoints.MakeOrgCreate(orgService, logger))
	app.Get(oauthCallBackUri, endpoints.MakeOAuthCallback(oauth2Config, store, tokenVerifier))
	app.Get(versionsApi, endpoints.MakeVersions(viper.GetString("app.version")))
	app.Get(apiSectorsPrefix, endpoints.MakeSectorsFindAll(sectorService, logger))
	app.Post(apiSectorsPrefix, endpoints.MakeSectorCreate(sectorService, logger))
	app.Delete(fullApiUri+"/sessions", endpoints.DeleteSession(clientId, clientSecret, store, oauthConfig, logger))
	app.Get(apiUsersPrefix, endpoints.MakUsersList(userService, logger))
	app.Post(apiUsersPrefix, endpoints.MakeUserCreate(userService, logger))

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

func initGrpcConn(grpcEndpoint string) (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	// You may replace `localhost:4317` with your endpoint.
	conn, err := grpc.NewClient(grpcEndpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}
