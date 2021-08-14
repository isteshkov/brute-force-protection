package main

import (
	"log"
	"os"

	"gitlab.com/isteshkov/brute-force-protection/config"
	"gitlab.com/isteshkov/brute-force-protection/domain/database"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"gitlab.com/isteshkov/brute-force-protection/migrations"
	"gitlab.com/isteshkov/brute-force-protection/repositories"
	"gitlab.com/isteshkov/brute-force-protection/service"
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
		err = migrations.MigrateUp(cfg.DatabaseUrl)
		if err != nil {
			panic(err)
		}
		return
	}

	logger, err := logging.NewLogger(&logging.Config{
		LogLvl:        cfg.LogLevel,
		InstanceId:    cfg.InstanceId,
		ContainerId:   cfg.ContainerId,
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

	db, err := database.GetDatabase(database.Config{DatabaseURL: cfg.DatabaseUrl}, logger)
	if err != nil {
		panic(err)
	}

	subnetsRepository := repositories.NewSubnetListRepository(db, logger)

	server := service.NewService(&service.Config{
		ProfilingApiPort: cfg.ProfilingApiPort,
		TechnicalApiPort: cfg.TechnicalApiPort,
		RpcPort:          cfg.RpcPort,
	}, subnetsRepository, logger)

	server.ListenAndServe()
}
