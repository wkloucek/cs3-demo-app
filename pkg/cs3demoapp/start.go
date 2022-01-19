package cs3demoapp

import (
	"context"
	"time"

	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/register"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/server/grpc"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/server/http"
)

func Start() error {
	ctx := context.Background()

	gwc, err := pool.GetGatewayServiceClient("localhost:9142")
	if err != nil {
		return err
	}

	err = register.Register(ctx, gwc)
	if err != nil {
		return err
	}

	grpcServer, err := grpc.Server(ctx)
	if err != nil {
		return err
	}
	defer grpcServer.GracefulStop()

	_, err = http.Server(ctx, gwc)
	if err != nil {
		return err
	}

	for {
		time.Sleep(1 * time.Second)
	}

	return nil
}
