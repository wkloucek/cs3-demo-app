package app

import (
	"context"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
)

type DemoApp struct {
}

func (a DemoApp) OpenInApp(
	ctx context.Context,
	req *appproviderv1beta1.OpenInAppRequest,
) (*appproviderv1beta1.OpenInAppResponse, error) {

	fp := make(map[string]string)
	fp["access_token"] = "fooo"

	resp := &appproviderv1beta1.OpenInAppResponse{
		Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
		AppUrl: &appproviderv1beta1.OpenInAppURL{
			AppUrl:         "https://lorem-ipsum",
			Method:         "POST",
			FormParameters: fp,
		},
	}
	return resp, nil
}
