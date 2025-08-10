package setup

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/template/html/v2"
	"poc-fiber/exceptions"
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
	var defErrorHandler = func(c *fiber.Ctx, err error) error {
		var e *fiber.Error
		code := fiber.StatusInternalServerError
		if errors.As(err, &e) {
			code = e.Code
			if code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError {
				apiError := exceptions.ConvertToFunctionalError(err, code)
				return c.Status(code).JSON(apiError)
			} else {
				apiError := exceptions.ConvertToInternalError(err)
				return c.Status(code).JSON(apiError)
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(exceptions.ConvertToInternalError(err))
	}

	// load only the contents of the subfolder www
	engine := html.New("./www", ".html")
	engine.Delims("{{", "}}") // define delimiters to use in the templates

	fConfig := fiber.Config{
		AppName:           appName,
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: true,
		UnescapePath:      true,
		ErrorHandler:      defErrorHandler,
		Views:             engine,
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
