package main

import (
	"log/slog"

	"github.com/spf13/pflag"

	"AvailabilityLinks/internal/config"
	"AvailabilityLinks/internal/handlers"
	"AvailabilityLinks/internal/storage"
	"AvilabilityLinks/internal/router"
)

func main() {

	configPath := pafseFlags()

	cfg, err := config.MustLoadConfig(configPath)
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
	}

	db, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		slog.Error("Error initializing storage", slog.Any("error", err))
	}
	defer db.Close()

	h := handlers.NewHandler(db)
	r := router.InitRouter(h)
	_ = r

	//TODO: start server
}

func pafseFlags() string {

	var configPath string

	pflag.StringVar(&configPath, "config", "configs/local.yaml", "Path to config file")
	pflag.Parse()

	return configPath
}
