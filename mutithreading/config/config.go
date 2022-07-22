package config

import (
	"flag"
	"log"
	"os"
	"path"
	"sync"
)

var (
	once sync.Once
	cachedConfig *Config
)

const (
	configFileKey     = "fName"
	defaultConfigFile = "Downloads/text.txt"
	configFileUsage   = "this is config file path"

	noOfWorkersKey        = "M"
	defaultNoOfWorkersKey = 2
	noOfWorkersKeyUsage   = "this is number of workers"
)

type Config struct {
	FileName    string
	NoOfWorkers int64
}

func ProvideInputConfigs() (c *Config) {
	once.Do(func() {
		var configFile string
		var noOfWorkers int64
		flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
		flag.Int64Var(&noOfWorkers, noOfWorkersKey, defaultNoOfWorkersKey, noOfWorkersKeyUsage)
		flag.Parse()

		if configFile == defaultConfigFile {
			log.Println("reading default file")
			wd, err := os.UserHomeDir()
			check(err)
			configFile = path.Join(wd, configFile)
		}

		c = &Config{
			FileName:    configFile,
			NoOfWorkers: noOfWorkers,
		}

		cachedConfig = c
	})
	return cachedConfig
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}