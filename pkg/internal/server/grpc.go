package server

import (
	"net"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/app"
	"google.golang.org/grpc"
)

func Server() error {
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	app := app.DemoApp{}

	appproviderv1beta1.RegisterProviderAPIServer(s, app)

	l, err := net.Listen("tcp", "localhost:5678")
	if err != nil {
		return err
	}

	s.Serve(l)

	return nil
}
