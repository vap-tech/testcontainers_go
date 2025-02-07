package oracle

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"path/filepath"
)

const (
	defaultPassword = "oracle"
)

// DBContainer represents the oracle container type used in the module
type DBContainer struct {
	testcontainers.Container
	dbName       string
	user         string
	password     string
	snapshotName string
	// sqlDriverName is passed to sql.Open() to connect to the database when making or restoring snapshots.
	// This can be set if your app imports a different oracle driver
	sqlDriverName string
}

// MustConnectionString panics if the address cannot be determined.
func (c *DBContainer) MustConnectionString(ctx context.Context) string {
	addr, err := c.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}
	return addr
}

// ConnectionString returns the connection string for the oracle container, using the default 1521 port, and
// obtaining the host and exposed port from the container. It also accepts a variadic list of extra arguments
// which will be appended to the connection string.
func (c *DBContainer) ConnectionString(ctx context.Context) (string, error) {

	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	// oracle://user:pass@server/service_name
	connStr := fmt.Sprintf("oracle://%s:%s@%s/%s", c.user, c.password, host, c.dbName)
	return connStr, nil
}

// WithInitScripts sets the init scripts to be run when the container starts
func WithInitScripts(scripts ...string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		initScripts := []testcontainers.ContainerFile{}
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

// WithPassword sets the initial password of the user to be created when the container starts
// It is required for you to use the Oracle DB image. It must not be empty or undefined.
// This environment variable sets the superuser password for Oracle DB.
func WithPassword(password string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["ORACLE_PASSWORD"] = password

		return nil
	}
}

// Deprecated: use Run instead
// RunContainer creates an instance of the Oracle container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*DBContainer, error) {
	return Run(ctx, "gvenzl/oracle-free:latest", opts...)
}

// Run creates an instance of the Oracle container type
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*DBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: img,
		Env: map[string]string{
			"ORACLE_PASSWORD": defaultPassword, //This variable is mandatory for the first container startup and specifies the password for the Oracle Database SYS and SYSTEM users.
		},
		ExposedPorts: []string{"1521/tcp"},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	// Gather all config options (defaults and then apply provided options)
	settings := defaultOptions()
	for _, opt := range opts {
		if apply, ok := opt.(Option); ok {
			apply(&settings)
		}
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, err
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	var c *DBContainer
	if container != nil {
		c = &DBContainer{
			Container:     container,
			password:      req.Env["ORACLE_PASSWORD"],
			sqlDriverName: settings.SQLDriverName,
			snapshotName:  settings.Snapshot,
		}
	}

	if err != nil {
		return c, fmt.Errorf("generic container: %w", err)
	}

	return c, nil
}
