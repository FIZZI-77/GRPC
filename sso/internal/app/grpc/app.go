package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func (s *App) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func NewApp(log *slog.Logger, port int) *App {
	grpcServer := grpc.NewServer()
	authgrpc.Register(grpcServer)
	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}
func (s *App) Run() error {
	const op = "grpcapp.Run"

	log := s.log.With(slog.String("op", op), slog.Int("port", s.port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", lis.Addr().String()))

	if err := s.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *App) Stop() {
	const op = "grpcapp.Stop"

	s.log.With(slog.String("op", op)).Info("grpc server is stopping", slog.Int("port", s.port))

	s.gRPCServer.GracefulStop()
}
