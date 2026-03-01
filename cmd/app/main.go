package main

import (
	"log/slog"
	"os"

	"github.com/nakle1ka/piggy/internal/app"
	"github.com/nakle1ka/piggy/internal/config"
	"github.com/nakle1ka/piggy/internal/database/postgres"
	"github.com/nakle1ka/piggy/internal/database/redis"
	"github.com/nakle1ka/piggy/internal/logger"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.SetupLogger(config.App.Env)

	db, err := postgres.NewPostgresDB(
		postgres.WithDbname(config.Database.Name),
		postgres.WithHost(config.Database.Host),
		postgres.WithPort(config.Database.Port),
		postgres.WithUser(config.Database.Username),
		postgres.WithPassword(config.Database.Password),
		postgres.WithSslmode(config.Database.SSLMode),
	)
	if err != nil {
		slog.Error("failed connect to postgres", "err", err)
		os.Exit(1)
	}

	rdb, err := redis.ConnectRedis(
		redis.WithDb(config.Redis.DB),
		redis.WithHost(config.Redis.Host),
		redis.WithPort(config.Redis.Port),
		redis.WithPassword(config.Redis.Password),
	)
	if err != nil {
		slog.Error("failed connect to redis", "err", err)
		os.Exit(1)
	}
	defer rdb.Close()

	app := app.NewApp(config, db, rdb)

	if err := app.Run(); err != nil {
		slog.Error("failed connect to database", "err", err)
		os.Exit(1)
	}
}
