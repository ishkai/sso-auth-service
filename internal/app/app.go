package app

import (
	"context"
	"fmt"
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	metricsapp "sso/internal/app/metrics"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCServer    *grpcapp.App
	MetricsServer *metricsapp.App
}

type storageProvider interface {
	auth.UserSaver
	auth.UserProvider
	auth.AppProvider
}

func New(log *slog.Logger, grpcPort int, metricsPort int, storageType string, storagepath string, storageDSN string, tokenTTL time.Duration) *App {
	const op = "app.New"

	var st storageProvider
	var err error

	switch storageType {
	case "sqlite":
		st, err = sqlite.New(storagepath)
	case "postgres":
		st, err = postgres.New(context.Background(), storageDSN)
	default:
		panic(fmt.Sprintf("%s: unknown storage type %s", op, storageType))
	}

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, st, st, st, tokenTTL)

	metricsApp := metricsapp.New(log, metricsPort)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer:    grpcApp,
		MetricsServer: metricsApp,
	}
}
