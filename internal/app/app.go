package app

import (
	grpcapp "exchanger-microservice/internal/app/grpc"
	"exchanger-microservice/internal/config"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcCfg config.GRPCConfig, storageCfg config.StorageConfig) *App {
	grpcApp := grpcapp.New(log, grpcCfg, storageCfg)

	return &App{
		GRPCServer: grpcApp,
	}
}
