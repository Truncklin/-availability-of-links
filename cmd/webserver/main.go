package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"

	"AvailabilityLinks/internal/config"
	"AvailabilityLinks/internal/handlers"
	"AvailabilityLinks/internal/router"
	"AvailabilityLinks/internal/storage"
)

func main() {
	slogInit()

	configPath := pafseFlags()

	cfg, err := config.MustLoadConfig(configPath)
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
	}

	db, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		slog.Error("Error initializing storage", slog.Any("error", err))
	}

	if err := db.Ping(); err != nil {
		slog.Error("DB ping failed:", "error", err)
	}

	h := handlers.NewHandler(db)
	r := router.InitRouter(h)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Host,
		Handler:      r,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	go func() {
		slog.Debug("Server start on", slog.Any("debugMessage", cfg.HttpServer.Host))
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("Server Failed", slog.Any("error", err))
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Debug("Shuting server :(")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

}

func slogInit() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}

func pafseFlags() string {

	var configPath string

	pflag.StringVar(&configPath, "config", "configs/local.yaml", "Path to config file")
	pflag.Parse()

	return configPath
}
