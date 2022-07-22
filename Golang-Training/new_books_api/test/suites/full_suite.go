package suites

import (
	"context"
	"sync"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type SuiteId int

const (
	RedisSuiteId SuiteId = iota
	PostgresSuiteId
)

type IntegrationSuite interface {
	suite.TestingSuite
	suite.SetupAllSuite
	suite.TearDownAllSuite

	GetContainerHost() string
	GetContainerMappedPort() nat.Port

	SetCtx(ctx context.Context)
	SetNetwork(network *testcontainers.DockerNetwork)
	SetRootDir(rootDir string)
}

type FullSuiteConfig struct {
	RootDir        string
	EnablePostgres bool
	//EnableRedis bool
}

type FullSuite struct {
	suite.Suite

	Ctx     context.Context
	Config  FullSuiteConfig
	Network *testcontainers.DockerNetwork
	Suites  map[SuiteId]IntegrationSuite
}

func NewFullSuite(config FullSuiteConfig) FullSuite {
	return FullSuite{Config: config}
}

func (f *FullSuite) SetupSuite() {
	var err error

	f.Ctx = context.Background()
	f.Network, err = CreateNetwork(f.Ctx, uuid.NewString())
	f.Require().NoError(err)

	f.Suites = make(map[SuiteId]IntegrationSuite)

	f.addSuiteIfEnabled(f.Config.EnablePostgres, PostgresSuiteId, &PostgresSuite{})
	// f.addSuiteIfEnabled(f.Config.EnablePostgres, RedisSuiteId, &RedisSuite{})

	var wg sync.WaitGroup

	for _, intSuite := range f.Suites {
		wg.Add(1)

		is := intSuite
		go func() {
			defer wg.Done()

			is.SetCtx(f.Ctx)
			is.SetT(f.T())
			is.SetNetwork(f.Network)
			is.SetRootDir(f.Config.RootDir)
			is.SetupSuite()
		}()
	}

	wg.Wait()
}

func (f *FullSuite) TearDownSuite() {
	for _, testSuite := range f.Suites {
		testSuite.TearDownSuite()
	}

	if f.Network != nil {
		_ = f.Network.Remove(f.Ctx)
	}
}

func (f *FullSuite) addSuiteIfEnabled(enabled bool, id SuiteId, suite IntegrationSuite) {
	if enabled {
		f.Suites[id] = suite
	}
}

func CreateNetwork(ctx context.Context, clusterNetworkName string) (*testcontainers.DockerNetwork, error) {

	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: clusterNetworkName},
	})

	if err != nil {
		return nil, err
	}

	return network.(*testcontainers.DockerNetwork), nil
}
