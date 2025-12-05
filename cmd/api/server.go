package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (app *application) serve() error {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server...", "port", app.cfg.port)
	return s.ListenAndServe()
}
