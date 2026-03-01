package logger

import (
	"log/slog"
	"os"
)

const (
	envDev  string = "prod"
	envProd string = "dev"
)

func SetupLogger(env string) {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
		log.Warn("Unknown environment, defaulting to development")
	}

	slog.SetDefault(log)
}
