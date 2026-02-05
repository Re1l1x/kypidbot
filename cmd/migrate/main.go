package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jus1d/kypidbot/internal/config"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	"github.com/pressly/goose/v3"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	c := config.MustLoad()
	p := c.Postgres
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.Name, p.ModeSSL)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("postgresql: failed to connect", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("postgresql: failed to set dialect", sl.Err(err))
		os.Exit(1)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Error("postgresql: failed to apply up migrations", sl.Err(err))
		os.Exit(1)
	}
}
