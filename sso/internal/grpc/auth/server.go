package auth

import (
	"context"
	ssopb "github.com/FIZZI-77/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
type serverApi struct {
	ssopb.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

func (s *serverApi) Login(ctx context.Context, req *ssopb.LoginRequest) (*ssopb.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssopb.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(ctx context.Context, req *ssopb.RegisterRequest) (*ssopb.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssopb.RegisterResponse{
		UserId: userID,
	}, nil

}

func (s *serverApi) IsAdmin(ctx context.Context, req *ssopb.IsAdminRequest) (*ssopb.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssopb.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssopb.LoginRequest) error {

	_, err := mail.ParseAddress(req.GetEmail())

	if req.GetEmail() == " " && err == nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == " " {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "appId is required")
	}
	return nil
}

func validateRegister(req *ssopb.RegisterRequest) error {
	_, err := mail.ParseAddress(req.GetEmail())

	if req.GetEmail() == " " && err == nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == " " {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateIsAdmin(req *ssopb.IsAdminRequest) error {

	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
