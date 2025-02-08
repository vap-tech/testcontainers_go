package testhelpers

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vap-tech/testcontainers_go/oracle"
	"time"
)

type PostgresContainer struct {
	*oracle.OracleContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := oracle.Run(ctx, "gvenzl/oracle-free",
		//oracle.WithScripts(filepath.Join("..", "testdata", "init-db.sql")), // Пока похоже эфекта не имеет
		//oracle.WithDatabase("FREEPDB2"), // Работает
		//oracle.WithUsername("test-user"), // Контейнер останавливается с ненулевым кодом
		//oracle.WithPassword("test-password"), // Контейнер останавливается с ненулевым кодом
		testcontainers.WithWaitStrategy(
			wait.ForLog("DATABASE IS READY TO USE!").
				WithOccurrence(1).WithStartupTimeout(5*time.Minute)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		OracleContainer:  pgContainer,
		ConnectionString: connStr,
	}, nil
}
