package main

import (
	"log"
	"os"

	"gitlab.com/isteshkov/brute-force-protection/config"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
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

	server := service.NewService(&service.Config{
		ProfilingApiPort: cfg.ProfilingApiPort,
		TechnicalApiPort: cfg.TechnicalApiPort,
		RpcPort:          cfg.RpcPort,
	}, logger)

	server.ListenAndServe()
}
