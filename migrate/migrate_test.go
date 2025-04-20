package migrate

import (
	"context"
	"fmt"
	"poc-fiber/logger"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

type TestLogConsumer struct {
	Msgs []string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	g.Msgs = append(g.Msgs, string(l.Content))
}

func TestOrganizationDao(t *testing.T) {
	ctx := context.Background()
	consoleLogger := logger.ConfigureConsoleLogger()

	logConsumer := TestLogConsumer{
		Msgs: []string{},
	}

	// Init postgreSQL container with testcontainers
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:17.4",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			consoleLogger.Error("error terminating container", zap.Error(err))
		}
	}(dbContainer, ctx)

	errLogProd := dbContainer.StartLogProducer(ctx)
	if errLogProd != nil {
		fmt.Printf("Error on log producer: [%v]", errLogProd)
	}
	dbContainer.FollowOutput(&logConsumer)

	// Retrieve postgreSQL container host and port
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")
	fmt.Printf("postgreSQL started on [%s]:[%s] \n", host, port)
	pgUrl := fmt.Sprintf("postgres://postgres:postgres@%s:%d/testdb?sslmode=disable", host, port.Int())

	consoleLogger.Info("Postgresql", zap.String("pgUrl", pgUrl))

	errMigrate := PerformMigration(*consoleLogger, pgUrl, "../migrate/files")
	assert.Nil(t, errMigrate, "no errors performing db migration")

}
