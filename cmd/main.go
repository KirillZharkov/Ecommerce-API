package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/KirillZharkov/Ecommerce-API/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	cfg := Config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetEnv("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=ecommerce port=5432 sslmode=disable"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.db.dsn)
	api := API{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
