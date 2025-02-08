package oracle

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	rootUser            = "SYS"
	defaultPassword     = "Qwerty123"
	defaultDatabaseName = "FREEPDB1"
)

// OracleContainer represents the Oracle container type used in the module
type OracleContainer struct {
	testcontainers.Container
	username string
	password string
	database string
}

// Run creates an instance of the Oracle container type
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*OracleContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        img,
		ExposedPorts: []string{"1521/tcp"},
		Env: map[string]string{
			"ORACLE_PASSWORD": defaultPassword,
		},
		WaitingFor: wait.ForLog("ALTER DATABASE OPEN"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, err
		}
	}

	username, ok := req.Env["APP_USER"]
	if !ok {
		username = rootUser
	}

	password, ok := req.Env["APP_USER_PASSWORD"]
	if !ok {
		password = defaultPassword
	}

	database, ok := req.Env["ORACLE_DATABASE"]
	if !ok {
		database = defaultDatabaseName
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)

	var c *OracleContainer
	if container != nil {
		c = &OracleContainer{
			Container: container,
			password:  password,
			username:  username,
			database:  database,
		}
	}

	if err != nil {
		return c, fmt.Errorf("could not create oracle container: %w", err)
	}

	return c, nil
}

// MustConnectionString panics if the address cannot be determined.
func (c *OracleContainer) MustConnectionString(ctx context.Context, args ...string) string {
	addr, err := c.ConnectionString(ctx, args...)
	if err != nil {
		panic(err)
	}
	return addr
}

func (c *OracleContainer) ConnectionString(ctx context.Context, args ...string) (string, error) {
	containerPort, err := c.MappedPort(ctx, "1521/tcp")
	if err != nil {
		return "", err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	extraArgs := ""
	if len(args) > 0 {
		extraArgs = strings.Join(args, "&")
	}
	if extraArgs != "" {
		extraArgs = "?" + extraArgs
	}
	// oracle://user:pass@server/service_name из дока на драйвер https://pkg.go.dev/github.com/sijms/go-ora
	connectionString := fmt.Sprintf("oracle://%s:%s@%s:%s/%s%s", c.username, c.password, host, containerPort.Port(), c.database, extraArgs)
	log.Printf("connection string: %s", connectionString)
	return connectionString, nil
}

func WithUsername(username string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["APP_USER"] = username

		return nil
	}
}

func WithPassword(password string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["APP_USER_PASSWORD"] = password

		return nil
	}
}

func WithDatabase(database string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["ORACLE_DATABASE"] = database

		return nil
	}
}

func WithScripts(scripts ...string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		var initScripts []testcontainers.ContainerFile
		for _, script := range scripts {
			cf := testcontainers.ContainerFile{
				HostFilePath:      script,
				ContainerFilePath: "/docker-entrypoint-initdb.d/" + filepath.Base(script),
				FileMode:          0o755,
			}
			initScripts = append(initScripts, cf)
		}
		req.Files = append(req.Files, initScripts...)

		return nil
	}
}
