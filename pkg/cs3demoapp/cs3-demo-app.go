package cs3demoapp

import (
	"context"

	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/register"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/server"
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

	err = server.Server()
	if err != nil {
		return err
	}

	return nil
}
