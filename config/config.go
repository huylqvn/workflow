package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env     string `default:"development" envconfig:"ENV"`
	Port    string `default:"8080" envconfig:"PORT"`
	Version string `default:"v1" envconfig:"VERSION"`
	APP     string `default:"prototype" envconfig:"APP"`

	Database Database

	NumberOfWorkers int `default:"2" envconfig:"NUMBER_OF_WORKERS"`
	BatchSize       int `default:"10" envconfig:"BATCH_SIZE"`

	ElasticApmServerURL   string `default:"http://localhost:8200" envconfig:"ELASTIC_APM_SERVER_URL"`
	ElasticApmSecretToken string `default:"token" envconfig:"ELASTIC_APM_SECRET_TOKEN"`
	ElasticApmServiceName string `default:"service" envconfig:"ELASTIC_APM_SERVICE_NAME"`
	ElasticApmEnvironment string `default:"env" envconfig:"ELASTIC_APM_ENVIRONMENT"`

	GroupID int64 `default:"4266372848" envconfig:"GROUP_ID"`
}

type Database struct {
	Type     string `default:"postgres" envconfig:"DATABASE_TYPE"`
	Host     string `default:"103.82.38.155" envconfig:"DATABASE_HOST"`
	Port     string `default:"5432" envconfig:"DATABASE_PORT"`
	Username string `default:"postgres" envconfig:"DATABASE_USER"`
	Password string `default:"0Canpass!!!" envconfig:"DATABASE_PASSWORD"`
	DBName   string `default:"postgres" envconfig:"DATABASE_NAME"`
	SSL      string `default:"disable" envconfig:"DATABASE_SSL"`
}

var Conf *Config = &Config{}
var once sync.Once

func Load() *Config {
	once.Do(func() {
		godotenv.Load()
		if err := envconfig.Process("", Conf); err != nil {
			panic(err)
		}
	})

	return Conf
}

func (c Config) IsAPIServer() bool {
	return c.APP == "api" || c.APP == "prototype"
}

func (c Config) IsWorker() bool {
	return c.APP == "worker" || c.APP == "prototype"
}

func (c Config) IsScheduler() bool {
	return c.APP == "scheduler" || c.APP == "prototype"
}
