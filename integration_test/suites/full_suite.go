package suites

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"sync"
)

type SuiteId int

const (
	KafkaSuiteId SuiteId = iota
)

type IntegrationSuite interface {
	suite.TestingSuite
	suite.SetupAllSuite
	suite.TearDownAllSuite

	GetContainerHost() string
	GetContainerMappedPort() nat.Port
	SetNetwork(network *testcontainers.DockerNetwork)
	SetCtx(ctx context.Context)
	SetRootDir(rootDir string)
}

type FullSuiteConfig struct {
	RootDir string
	EnableKafka bool
}

type FullSuite struct {
	suite.Suite

	ctx context.Context
	config FullSuiteConfig
	Network *testcontainers.DockerNetwork
	Suites map[SuiteId]IntegrationSuite
}

func NewFullSuite(config FullSuiteConfig) FullSuite {
	return FullSuite{config: config}
}

func (suite *FullSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()
	suite.Network, err = CreateNetwork(suite.ctx, uuid.NewString())
	suite.Require().Nil(err)

	suite.Suites = make(map[SuiteId]IntegrationSuite)

	suite.addSuiteIfEnabled(suite.config.EnableKafka, KafkaSuiteId, &KafkaSuite{})

	var wg sync.WaitGroup

	for _, inSuite := range suite.Suites {
		wg.Add(1)

		is := inSuite
		go func(){
			defer wg.Done()

			is.SetCtx(suite.ctx)
			is.SetNetwork(suite.Network)
			is.SetRootDir(suite.config.RootDir)
			is.SetupSuite()
		}()
	}
}

func CreateNetwork(ctx context.Context, clusterName string) (*testcontainers.DockerNetwork, error) {
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: clusterName},
	})

	if err != nil {
		return nil, err
	}

	return network.(*testcontainers.DockerNetwork), nil
}

func (suite *FullSuite) TearDownSuite() {
	for _, testSuite := range suite.Suites {
		testSuite.TearDownSuite()
	}

	if suite.Network != nil {
		err := suite.Network.Remove(suite.ctx)
		suite.Require().NoError(err)
	}
}

func (suite *FullSuite) addSuiteIfEnabled(enabled bool, id SuiteId, is IntegrationSuite) {
	if enabled {
		suite.Suites[id] = is
	}
}
