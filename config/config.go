package config

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	DefaultServiceName   = "brute-force-protection-service"
	DefaultProfilingPort = ":6060"
	DefaultTechnicalPort = ":8001"
)

var ErrorProducerLoading = errors.NewProducer("LOADING_ERROR")

type Config struct {
	ServiceName   string `env:"SERVICE_NAME"`
	InstanceId    string `env:"INSTANCE_ID"`
	ContainerId   string `env:"CONTAINER_ID"`
	ContainerName string `env:"CONTAINER_NAME"`
	EnvName       string `env:"ENV_NAME"`
	Version       string `env:"VERSION"`
	Release       string `env:"RELEASE"`
	CommitSha     string `env:"COMMIT_SHA"`

	RpcPort          string `env:"RPC_PORT,required"`
	TechnicalApiPort string `env:"TECHNICAL_API_PORT"`
	ProfilingApiPort string `env:"PROFILING_API_PORT"`
	LogLevel         string `env:"LOG_LEVEL"`
	DatabaseUrl      string `env:"DATABASE_URL,required"`
}

func LoadConfig(fileName ...string) (cfg *Config, err error) {
	if len(fileName) == 1 {
		err = loadEnvFromFile(fileName[0])
	}
	cfg = &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return
	}

	fillDefault(cfg)

	return
}

func loadEnvFromFile(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		if err := godotenv.Load(filename); err != nil {
			return ErrorProducerLoading.WrapF(err, fmt.Sprintf("error loading file: %s", filename))
		}
		log.Printf("Config file %s is using.\n", filename)
	} else {
		log.Printf("Missed config %s. Using env only.\n", filename)
	}

	return nil
}

func fillDefault(cfg *Config) {
	if cfg.ServiceName == "" {
		cfg.ServiceName = DefaultServiceName
	}
	if cfg.ProfilingApiPort == "" {
		cfg.ProfilingApiPort = DefaultProfilingPort
	}
	if cfg.TechnicalApiPort == "" {
		cfg.TechnicalApiPort = DefaultTechnicalPort
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = logging.LevelDebug
	}
}
