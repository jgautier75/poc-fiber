package dao

import (
	"context"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	PG_VERSION     = "postgres:17.5"
	PG_PORT        = "5432/tcp"
	PG_DB          = "testdb"
	PG_USER        = "postgres"
	PG_PASS        = "postgres"
	DEFAULT_TENANT = 1
)

type TestLogConsumer struct {
	Msgs []string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	g.Msgs = append(g.Msgs, string(l.Content))
}

func CreatePgContainerTest(ctx context.Context, logger zap.Logger) (*testcontainers.Container, string, nat.Port, error) {
	logConsumer := TestLogConsumer{
		Msgs: []string{},
	}
	// Init postgreSQL container with testcontainers
	containerReq := testcontainers.ContainerRequest{
		Image:        PG_VERSION,
		ExposedPorts: []string{PG_PORT},
		WaitingFor:   wait.ForListeningPort(PG_PORT),
		Env: map[string]string{
			"POSTGRES_DB":       PG_DB,
			"POSTGRES_PASSWORD": PG_USER,
			"POSTGRES_USER":     PG_PASS,
		},
	}
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	errLogProd := dbContainer.StartLogProducer(ctx)
	if errLogProd != nil {
		logger.Error("container error log producer", zap.Error(errLogProd))
	}
	dbContainer.FollowOutput(&logConsumer)
	host, errHost := dbContainer.Host(ctx)
	var nilHost string
	var nilPort nat.Port
	if errHost != nil {
		return nil, nilHost, nilPort, errHost
	}
	port, errPort := dbContainer.MappedPort(ctx, "5432")
	return &dbContainer, host, port, errPort
}
