package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/models"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/test/suites"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

const (
	pgPassword = "mySecretPassword"
	pgDatabase = "new_books_db"
	pgUsername = "postgres"
)

type NewBookAPISuite struct {
	suites.FullSuite
	ctx          context.Context
	rootDir      string
	dbPort       string
	dbConn       *gorm.DB
	router       *mux.Router
	appPort      string
	dbHost       string
	appContainer testcontainers.Container
	port         nat.Port
	host         string
	appUrl       string
}

func TestNewBooksAPI(t *testing.T) {
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	rootDir := path.Join(workingDirectory, "../../")

	cnf := suites.FullSuiteConfig{
		RootDir:        rootDir,
		EnablePostgres: true,
	}

	suite.Run(t, &NewBookAPISuite{
		FullSuite: suites.NewFullSuite(cnf),
		rootDir:   rootDir,
	})
}

func (suite *NewBookAPISuite) SetupTest() {
	suite.ctx = context.Background()
	suite.appPort = "8010"
	suite.SetupSuite()

	//redisSuite := suite.Suites[suites.RedisSuiteId].(*suites.RedisSuite)
	//suite.redisSuite = fmt.Sprintf("%s:%s", redisSuite.GetContainerHost(), redisSuite.GetContainerMappedPort().Port())

	pgSuite := suite.Suites[suites.PostgresSuiteId].(*suites.PostgresSuite)
	suite.dbPort = pgSuite.GetContainerMappedPort().Port()
	suite.dbHost = pgSuite.GetContainerHost()

	suite.startApp()
}

func (suite *NewBookAPISuite) startApp() {
	request := testcontainers.ContainerRequest{
		Name: "new_books_api",
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       suite.Config.RootDir,
			Dockerfile:    "build/package/Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{
			"APP_PORT":    suite.appPort,
			"DB_USER":     pgUsername,
			"DB_NAME":     pgDatabase,
			"DB_PASSWORD": pgPassword,
			"DB_HOST":     suite.dbHost,
			"DB_PORT":     suite.dbPort,
		},
		ExposedPorts: []string{"5500"},
		Networks:     []string{suite.Network.Name},
		Cmd:          []string{"/bin/sh", "-c", "/wait && ./new_books_api -configFile default.yaml"},
		WaitingFor:   wait.ForLog("Starting server"),
	}

	app, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
	})
	suite.Require().NoError(err)

	suite.appContainer = app

	err = suite.appContainer.Start(suite.ctx)
	suite.Require().NoError(err)

	suite.port, err = suite.appContainer.MappedPort(suite.ctx, "5500")
	suite.Require().NoError(err)

	suite.host, err = suite.appContainer.Host(suite.ctx)
	suite.Require().NoError(err)
	suite.appUrl = fmt.Sprintf("%s:%s", suite.host, suite.port.Port())
}

func (suite *NewBookAPISuite) TestGetAllBooks() {
	response, err := suite.makeRequest("GET", "/books", nil)
	suite.Require().NoError(err)
	defer response.Body.Close()
	suite.Assert().Equal(http.StatusOK, response.StatusCode)

	var booksResponse []*models.Book
	data, err := ioutil.ReadAll(response.Body)
	suite.Require().NoError(err)

	err = json.Unmarshal(data, &booksResponse)
	suite.Require().NoError(err)

	expectedBooksResponse := suite.getExpectedBooks()
	suite.Assert().Equal(expectedBooksResponse, booksResponse, "books response does not match")
}

func (suite *NewBookAPISuite) makeRequest(method, url string, body io.Reader) (*http.Response, error) {

	formattedUrl := fmt.Sprintf("http://%s%s", suite.appUrl, url)

	client := &http.Client{}
	req, err := http.NewRequest(method, formattedUrl, body)
	req.Header.Set("Content-Type", "application/json")
	suite.Require().NoError(err)
	resp, err := client.Do(req)
	return resp, err
}

func (suite *NewBookAPISuite) getExpectedBooks() []*models.Book {
	return []*models.Book{
		{ISBN: 12345, Name: "da vinci code", Author: "dan brown"},
		{ISBN: 12346, Name: "the best laid plans", Author: "leslie stewart"},
		{ISBN: 12347, Name: "atomic habits", Author: "james clear"},
	}
}
