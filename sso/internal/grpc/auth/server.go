package auth

import (
	"context"
	ssopb "github.com/FIZZI-77/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssopb.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssopb.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, req *ssopb.LoginRequest) (*ssopb.LoginResponse, error) {
	panic("implement me")
}

func (s *serverApi) Register(ctx context.Context, req *ssopb.RegisterRequest) (*ssopb.RegisterResponse, error) {
	panic("implement me")
}
