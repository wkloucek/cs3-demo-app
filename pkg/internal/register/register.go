package register

import (
	"context"
	"fmt"

	registryv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
)

func Register(ctx context.Context, gwc gatewayv1beta1.GatewayAPIClient) error {

	req := &registryv1beta1.AddAppProviderRequest{
		Provider: &registryv1beta1.ProviderInfo{
			Name:        "demoapp",
			Description: "this is an demo app",
			Icon:        "",
			Address:     "127.0.0.1:5678",
			MimeTypes: []string{
				"image/png",
			},
		},
	}

	resp, err := gwc.AddAppProvider(ctx, req)
	if err != nil {
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		return fmt.Errorf("") // TODO: transform error
	}

	return nil
}
