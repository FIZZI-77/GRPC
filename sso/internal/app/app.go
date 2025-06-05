package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, tokenTTL time.Duration) *App {

	db, err := postgres.StorageConnect()
	if err != nil {
		log.Error("failed to connect to postgres", err)
	}

	storage := postgres.NewStorage(db)

	authService := auth.NewAuth(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
