package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/Infamous003/GoChat/internal/data"
	"github.com/Infamous003/GoChat/internal/db"
	_ "github.com/lib/pq"
)

type config struct {
	port int

	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	cfg    config
	logger *slog.Logger
	models data.Models
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var cfg config

	flag.IntVar(&cfg.port, "port", 9090, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GOCHAT_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 10*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	database, err := db.OpenDB(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("successfully connected to the database")

	app := &application{
		cfg:    cfg,
		logger: logger,
		models: data.NewModels(database),
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
