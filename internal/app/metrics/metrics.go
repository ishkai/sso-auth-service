package metricsapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	port   int
}

func New(log *slog.Logger, port int) *App {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}

	return &App{
		log:    log,
		server: server,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "metricsapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))
	log.Info("starting metrics app", slog.String("address", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "metricsapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping metrics server")

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
