package main

import (
	"fmt"
	"log/slog"

	"github.com/spf13/pflag"

	"AvailabilityLinks/internal/config"
)

func main() {

	configPath := pafseFlags()

	cfg, err := config.MustLoadConfig(configPath)
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
	}

	if cfg == nil {
		slog.Error("Failed to load config")
		return
	}

	fmt.Println(cfg)

	//TODO: init config: cleanenv

	//TODO: init logger: log/slog

	//TODO: init storage: sqlite

	//TODO: init route: chi, "chi render"

	//TODO: start server
}

func pafseFlags() string {
	var configPath string
	pflag.StringVar(&configPath, "config", "configs/local.yaml", "Path to config file")
	pflag.Parse()
	return configPath
}
