package suite

import (
	"context"
	ssopb "github.com/FIZZI-77/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssopb.AuthClient
}

const grpcHost = "localhost"

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	cfg := config.MustLoadByPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(), grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssopb.NewAuthClient(cc),
	}

}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
