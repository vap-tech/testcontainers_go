package testhelpers

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vap-tech/testcontainers_go/oracle"
	"time"
)

type OracleContainer struct {
	*oracle.OracleContainer
	ConnectionString string
}

func CreateOracleContainer(ctx context.Context) (*OracleContainer, error) {
	oraContainer, err := oracle.Run(ctx, "gvenzl/oracle-free",
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
	connStr, err := oraContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &OracleContainer{
		OracleContainer:  oraContainer,
		ConnectionString: connStr,
	}, nil
}
