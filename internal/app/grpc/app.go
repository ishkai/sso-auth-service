package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
	"sso/internal/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log          *slog.Logger
	gRPCServer   *grpc.Server
	healthServer *health.Server
	port         int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()))

	authgrpc.Register(gRPCServer, authService)

	healthServer := health.NewServer()

	healthgrpc.RegisterHealthServer(gRPCServer, healthServer)

	healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("auth.Auth", healthgrpc.HealthCheckResponse_SERVING)

	reflection.Register(gRPCServer)

	return &App{
		log:          log,
		gRPCServer:   gRPCServer,
		healthServer: healthServer,
		port:         port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
	a.healthServer.SetServingStatus("auth.Auth", healthgrpc.HealthCheckResponse_NOT_SERVING)

	a.gRPCServer.GracefulStop()
}
