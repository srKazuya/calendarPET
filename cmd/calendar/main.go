package main

import (
	"calendar/internal/config"
	"calendar/internal/event"
	"calendar/internal/infrastructure/http/handlers"
	"calendar/internal/infrastructure/http/middleware"
	"calendar/internal/infrastructure/storage/in_memory"
	"calendar/pkg/sl_logger/sl"
	"errors"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//entry point

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	storage := inmem.New()
	service := event.NewService(storage)

	mux := http.NewServeMux()

	// root := middleware.NewMWLogger(log)(middleware.RequestID(mux))
	mux.Handle("/create_event",
		middleware.NewMWLogger(log)(
			middleware.RequestID(
				http.HandlerFunc(handlers.NewAddEventHandler(log, service)),
			),
		),
	)
	mux.Handle("/events_for_day",
		middleware.NewMWLogger(log)(
			middleware.RequestID(
				http.HandlerFunc(handlers.NewEventsForDayHandler(log, service)),
			),
		),
	)
	mux.Handle("/events_for_month",
		middleware.NewMWLogger(log)(
			middleware.RequestID(
				http.HandlerFunc(handlers.NewEventsForMonthHandler(log, service)),
			),
		),
	)
	mux.Handle("/events_for_week",
		middleware.NewMWLogger(log)(
			middleware.RequestID(
				http.HandlerFunc(handlers.NewEventsForWeekHandler(log, service)),
			),
		),
	)
	mux.Handle("/delete_event",
		middleware.NewMWLogger(log)(
			middleware.RequestID(
				http.HandlerFunc(handlers.NewDeleteEventHandler(log, service)),
			),
		),
	)
	mux.Handle("/update_event",
		middleware.NewMWLogger(log)(
			middleware.RequestID(
				http.HandlerFunc(handlers.NewUpdateEventHandler(log, service)),
			),
		),
	)


	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting HTTP server",
		slog.String("address", cfg.Address),
		slog.Duration("read_timeout", cfg.HTTPServer.Timeout),
		slog.Duration("write_timeout", cfg.HTTPServer.Timeout),
		slog.Duration("idle_timeout", cfg.IdleTimeout),
	)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start server", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Debug("Debug level enabled")
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Debug("Debug level enabled")
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		log.Info("Info level enabled")
	}

	return log
}
