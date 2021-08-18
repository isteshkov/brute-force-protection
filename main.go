package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"gitlab.com/isteshkov/brute-force-protection/config"
	"gitlab.com/isteshkov/brute-force-protection/domain/database"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"gitlab.com/isteshkov/brute-force-protection/migrations"
	"gitlab.com/isteshkov/brute-force-protection/repositories"
	"gitlab.com/isteshkov/brute-force-protection/service"
	ratelimiter "gitlab.com/isteshkov/brute-force-protection/service/rate_limiter"
)

func main() {
	file := ""
	args := os.Args[1:]
	if len(args) > 0 {
		file = args[0]
	}

	cfg, err := config.LoadConfig(file)
	if err != nil {
		log.Println("Error loading config, using default")
		panic(err)
	}

	if len(args) > 1 && args[1] == "migrate" {
		err = migrations.MigrateUp(cfg.DatabaseURL)
		if err != nil {
			panic(err)
		}
		return
	}

	logger, err := logging.NewLogger(&logging.Config{
		LogLvl:        cfg.LogLevel,
		InstanceID:    cfg.InstanceID,
		ContainerID:   cfg.ContainerID,
		ContainerName: cfg.ContainerName,
		EnvName:       cfg.EnvName,
		Version:       cfg.Version,
		Release:       cfg.Release,
		CommitSha:     cfg.CommitSha,
		ServiceName:   cfg.ServiceName,
	})
	if err != nil {
		panic(err)
	}

	db, err := database.GetDatabase(database.Config{DatabaseURL: cfg.DatabaseURL}, logger)
	if err != nil {
		panic(err)
	}

	subnetsRepository := repositories.NewSubnetListRepository(db, logger)
	rateLimiter := ratelimiter.NewRateLim(
		cfg.LoginAttemptsPerMinuteCount,
		cfg.PasswordAttemptsPerMinuteCount,
		cfg.IPAttemptsPerMinuteCount,
		logger,
	)

	server := service.NewService(&service.Config{
		ProfilingAPIPort: cfg.ProfilingAPIPort,
		TechnicalAPIPort: cfg.TechnicalAPIPort,
		RPCPort:          cfg.RPCPort,
	}, subnetsRepository, rateLimiter, logger)

	server.ListenAndServe()
}
