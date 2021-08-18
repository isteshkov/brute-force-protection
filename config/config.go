package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
)

const (
	DefaultServiceName   = "brute-force-protection-service"
	DefaultProfilingPort = ":6060"
	DefaultTechnicalPort = ":8001"
)

var ErrorProducerLoading = errors.NewProducer("LOADING_ERROR")

type Config struct {
	ServiceName   string `env:"SERVICE_NAME"`
	InstanceID    string `env:"INSTANCE_ID"`
	ContainerID   string `env:"CONTAINER_ID"`
	ContainerName string `env:"CONTAINER_NAME"`
	EnvName       string `env:"ENV_NAME"`
	Version       string `env:"VERSION"`
	Release       string `env:"RELEASE"`
	CommitSha     string `env:"COMMIT_SHA"`

	RPCPort          string `env:"RPC_PORT,required"`
	TechnicalAPIPort string `env:"TECHNICAL_API_PORT"`
	ProfilingAPIPort string `env:"PROFILING_API_PORT"`
	LogLevel         string `env:"LOG_LEVEL"`
	DatabaseURL      string `env:"DATABASE_URL,required"`

	LoginAttemptsPerMinuteCount    int `env:"LOGIN_ATTEMPTS_PER_MINUTE_COUNT" binding:"max=10"`
	PasswordAttemptsPerMinuteCount int `env:"PASSWORD_ATTEMPTS_PER_MINUTE_COUNT" binding:"max=100"`
	IPAttemptsPerMinuteCount       int `env:"IP_ATTEMPTS_PER_MINUTE_COUNT" binding:"max=1000"`
}

func LoadConfig(fileName string) (cfg *Config, err error) {
	cfg = &Config{}

	if len(fileName) > 0 {
		err = loadEnvFromFile(fileName)
		if err != nil {
			err = env.Parse(cfg)
			if err != nil {
				return
			}
		}
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
	if cfg.ProfilingAPIPort == "" {
		cfg.ProfilingAPIPort = DefaultProfilingPort
	}
	if cfg.TechnicalAPIPort == "" {
		cfg.TechnicalAPIPort = DefaultTechnicalPort
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = logging.LevelDebug
	}
}
