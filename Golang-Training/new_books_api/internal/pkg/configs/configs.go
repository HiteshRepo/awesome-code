package configs

import (
	"flag"
	"github.com/hashicorp/go-multierror"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/pkg/configs/db"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/pkg/configs/server"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"os"
	"sync"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = ""
	configFileUsage   = "this is config file path"
)

var (
	once         sync.Once
	cachedConfig *AppConfig
)

type AppConfig struct {
	ServerConfig server.ServerConfig `mapstructure:"app"`
	DBConfig     db.DatabaseConfig   `mapstructure:"db"`
}

func ProvideAppConfig() (c *AppConfig, err error) {
	once.Do(func() {
		var configFile string
		flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
		flag.Parse()

		var configReader io.ReadCloser
		configReader, err = os.Open(configFile)
		defer configReader.Close() //nolint:staticcheck

		if err != nil {
			return
		}

		c, err = LoadConfig(configReader)
		if err != nil {
			return
		}

		cachedConfig = c
	})

	return cachedConfig, err
}

func LoadConfig(reader io.Reader) (*AppConfig, error) {
	var appConfig AppConfig

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	keysToEnvironmentVariables := map[string]string{
		"app.port":    "APP_PORT",
		"db.name":     "DB_NAME",
		"db.user":     "DB_USER",
		"db.host":     "DB_HOST",
		"db.port":     "DB_PORT",
		"db.password": "DB_PASSWORD",
	}

	err := bind(keysToEnvironmentVariables)
	if err != nil {
		return nil, err
	}

	if err := viper.ReadConfig(reader); err != nil {
		return nil, errors.Wrap(err, "Failed to load app config file")
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, errors.Wrap(err, "Unable to parse app config file")
	}

	return &appConfig, nil
}

func bind(keysToEnvironmentVariables map[string]string) error {
	var bindErrors error

	for key, environmentVariable := range keysToEnvironmentVariables {
		if err := viper.BindEnv(key, environmentVariable); err != nil {
			bindErrors = multierror.Append(bindErrors, err)
		}
	}

	return bindErrors
}
