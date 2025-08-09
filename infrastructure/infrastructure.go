package infrastructure

import (
	"context"

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
