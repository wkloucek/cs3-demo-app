package app

import (
	"context"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"google.golang.org/grpc/metadata"
)

type DemoApp struct {
}

func (a DemoApp) OpenInApp(
	ctx context.Context,
	req *appproviderv1beta1.OpenInAppRequest,
) (*appproviderv1beta1.OpenInAppResponse, error) {
	authTokens := []string{}
	token := ""

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		authTokens = md.Get("x-access-token")
	}
	if len(authTokens) > 0 {
		token = authTokens[0]
	}

	if token == "" {
		return &appproviderv1beta1.OpenInAppResponse{
			Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_UNAUTHENTICATED},
			AppUrl: &appproviderv1beta1.OpenInAppURL{},
		}, nil
	}

	fp := make(map[string]string)
	fp["access_token"] = token
	fp["storage_id"] = req.ResourceInfo.Id.StorageId
	fp["opaque_id"] = req.ResourceInfo.Id.OpaqueId

	resp := &appproviderv1beta1.OpenInAppResponse{
		Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
		AppUrl: &appproviderv1beta1.OpenInAppURL{
			AppUrl:         "http://localhost:6789",
			Method:         "POST",
			FormParameters: fp,
		},
	}
	return resp, nil
}
