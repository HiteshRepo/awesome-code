package tests

import (
	"github.com/hiteshrepo/awesome-code/integration_test/suites"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

const (
	TestTopic = "test.topic.v1"
)

type integrationTestSuite struct {
	suites.FullSuite
	kafkaSuite *suites.KafkaSuite
}

func TestIntegrationTestSuite(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	config := suites.FullSuiteConfig{
		RootDir:     workingDir,
		EnableKafka: true,
	}

	suite.Run(t, &integrationTestSuite{
		FullSuite:  suites.NewFullSuite(config),
	})
}

func (suite *integrationTestSuite) SetupSuite() {
	suite.FullSuite.SetupSuite()

	kafkaSuite := suite.Suites[suites.KafkaSuiteId]
	suite.kafkaSuite = kafkaSuite.(*suites.KafkaSuite)
	suite.kafkaSuite.CreateTopics([]string{TestTopic})
}

func (suite *integrationTestSuite) TearDownSuite() {
	suite.FullSuite.TearDownSuite()
}

func (suite *integrationTestSuite) TestSomething() {
	suite.Require().False(true)
}