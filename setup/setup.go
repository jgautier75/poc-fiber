package setup

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/template/html/v2"
	"net/http"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func SetupCnxPool(pgUrl string, minConns int32, maxConns int32, zapLogger zap.Logger) (*pgxpool.Pool, error) {
	zapLogger.Info("CnxPool -> Parse configuration")
	dbConfig, errDbCfg := pgxpool.ParseConfig(pgUrl)
	if errDbCfg != nil {
		panic(errDbCfg)
	}
	dbConfig.MinConns = minConns
	dbConfig.MaxConns = maxConns
	zapLogger.Info("Connection Pool -> Initialize pool")
	return pgxpool.NewWithConfig(context.Background(), dbConfig)
}

func BuildFiberConfig(appName string) fiber.Config {

	var defaultErrorHandler = func(c fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError
		var e *fiber.Error
		matched := errors.As(err, &e)
		if matched && e != nil {
			code = e.Code
		}
		message := http.StatusText(code)
		if err != nil && !(matched && e == nil) {
			message = err.Error()
		}

		// Set Content-Type: text/plain; charset=utf-8
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

		// Return status code with error message
		return c.Status(code).SendString(message)
	}

	// load only the contents of the subfolder www
	engine := html.New("./www", ".html")
	engine.Delims("{{", "}}") // define delimiters to use in the templates

	fConfig := fiber.Config{
		AppName:       appName,
		CaseSensitive: true,
		StrictRouting: true,
		UnescapePath:  true,
		ErrorHandler:  defaultErrorHandler,
		Views:         engine,
	}
	return fConfig
}

func ConfigureRedisStorage(redisHost string, redisPort int) *redis.Storage {
	return redis.New(redis.Config{
		Host:      redisHost,
		Port:      redisPort,
		Username:  "",
		Password:  "",
		URL:       "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	},
	)
}
