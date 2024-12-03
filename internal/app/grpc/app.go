package grpc

import (
	"context"
	"exchanger-microservice/internal/config"
	"exchanger-microservice/internal/domain/exchanger"
	"exchanger-microservice/internal/domain/exchanger/db"
	exchange "exchanger-microservice/internal/grpc/exchange"
	"exchanger-microservice/pkg/client/psql"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(logger *slog.Logger, grpcCfg config.GRPCConfig, storageCfg config.StorageConfig) *App {
	const op = "app.grpc.New"
	log := logger.With(slog.String("op", op))
	log.Info("Initializing GRPC server")

	gRPCServer := grpc.NewServer()
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		storageCfg.StorageUser,
		storageCfg.StoragePass,
		storageCfg.StorageHost,
		storageCfg.StoragePort,
		storageCfg.StorageDatabase,
	)
	client, err := psql.NewClient(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(client, log)
	service := exchanger.NewService(log, storage)
	exchange.Register(gRPCServer, service)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       grpcCfg.Port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port))

	lis, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
