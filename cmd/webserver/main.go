package main

import (
	"log/slog"

	"github.com/spf13/pflag"

	"AvailabilityLinks/internal/config"
	"AvailabilityLinks/internal/storage"
)

func main() {

	configPath := pafseFlags()

	cfg, err := config.MustLoadConfig(configPath)
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
	}

	if err := storage.NewStorage(cfg.StoragePath); err != nil {
		slog.Error("Error initializing storage", slog.Any("error", err))
	}

	//TODO: init route: chi, "chi render"

	//TODO: start server
}

func pafseFlags() string {

	var configPath string

	pflag.StringVar(&configPath, "config", "configs/local.yaml", "Path to config file")
	pflag.Parse()

	return configPath
}
