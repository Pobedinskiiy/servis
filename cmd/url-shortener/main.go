package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"servis/internal/lib/logger/sl"
	"servis/internal/lib/logger/slogpretty"
	"servis/internal/storage/sqlite"

	"servis/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	err = storage.DeleteURL("google")
	if err != nil {
		log.Error("Failed to delete url", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("https://google.com", "google")
	if err != nil {
		log.Error("Failed to save url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("save url", slog.Int64("id", id))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//router.Use(middleware.New(log))

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := slogpretty.NewPrettyHandler(os.Stdout, opts)
		logger = slog.New(handler)
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
