package cs3demoapp

import (
	"context"
	"time"

	"github.com/wkloucek/cs3-demo-app/pkg/internal/app"
)

func Start() error {
	ctx := context.Background()

	app := app.New()

	if err := app.GetCS3apiClient(); err != nil {
		return err
	}

	if err := app.RegisterDemoApp(ctx); err != nil {
		return err
	}

	if err := app.GRPCServer(ctx); err != nil {
		return err
	}

	if err := app.HTTPServer(ctx); err != nil {
		return err
	}

	for {
		time.Sleep(1 * time.Second)
	}

	return nil
}
