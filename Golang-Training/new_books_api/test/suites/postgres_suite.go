package suites

import (
	"context"
	"fmt"
	"path"
	"time"

	// this is required for the driver to load
	_ "github.com/lib/pq"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	pgConnectionString      = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	pgPort                  = "5432"
	pgUsername              = "postgres"
	pgPassword              = "mySecretPassword"
	pgDatabase              = "new_books_db"
	pgContainerSetupTimeout = 60
)

type PostgresSuite struct {
	suite.Suite

	ctx     context.Context
	network *testcontainers.DockerNetwork
	rootDir string

	pgContainer testcontainers.Container
}

func (p *PostgresSuite) SetupSuite() {
	p.createPgContainer()

	err := p.pgContainer.Start(p.ctx)
	p.Require().NoError(err)
}

func (p *PostgresSuite) TearDownSuite() {
	p.terminateContainer(p.pgContainer)
}

func (p *PostgresSuite) GetContainerHost() string {
	host, err := p.pgContainer.Host(p.ctx)
	p.Require().NoError(err)

	return host
}

func (p *PostgresSuite) GetContainerMappedPort() nat.Port {
	port, err := p.pgContainer.MappedPort(p.ctx, pgPort)
	p.Require().NoError(err)

	return port
}

func (p *PostgresSuite) SetCtx(ctx context.Context) {
	p.ctx = ctx
}

func (p *PostgresSuite) SetNetwork(network *testcontainers.DockerNetwork) {
	p.network = network
}

func (p *PostgresSuite) SetRootDir(rootDir string) {
	p.rootDir = rootDir
}

func (p *PostgresSuite) createPgContainer() {
	getDBUrl := func(port nat.Port) string {
		username := pgUsername
		password := pgPassword
		database := pgDatabase
		return fmt.Sprintf(pgConnectionString, username, password, "localhost", port.Port(), database)
	}

	pathToSeedSql := path.Join(p.rootDir, "/test/scripts")

	request := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{pgPort + "/tcp"},
		WaitingFor:   wait.ForSQL(pgPort+"/tcp", "postgres", getDBUrl).Timeout(pgContainerSetupTimeout * time.Second),
		BindMounts:   map[string]string{"/docker-entrypoint-initdb.d": pathToSeedSql},
		Env: map[string]string{
			"POSTGRES_DB":       pgDatabase,
			"POSTGRES_USER":     pgUsername,
			"POSTGRES_PASSWORD": pgPassword,
		},
		Networks:     []string{p.network.Name},
	}
	container, err := testcontainers.GenericContainer(p.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
	})

	p.Require().NoError(err)
	p.pgContainer = container
}

func (p *PostgresSuite) terminateContainer(container testcontainers.Container) {
	if container != nil {
		if err := container.Terminate(p.ctx); err != nil {
			p.Require().NoError(err)
		}
	}
}

func (p *PostgresSuite) GetConnectionString(host, port string) string {
	return fmt.Sprintf(pgConnectionString, pgUsername, pgPassword, host, port, pgDatabase)
}
