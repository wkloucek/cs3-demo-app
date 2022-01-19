package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"google.golang.org/grpc/metadata"
)

func DownloadFile(ctx context.Context, ref *providerv1beta1.Reference, gwc gatewayv1beta1.GatewayAPIClient, token string) (http.Response, error) {
	req := &providerv1beta1.InitiateFileDownloadRequest{
		Ref: ref,
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "x-access-token", token)

	resp, err := gwc.InitiateFileDownload(ctx, req)
	if err != nil {
		return http.Response{}, err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		return http.Response{}, errors.New("status code != CODE_OK")
	}

	downloadEndpoint := ""
	downloadToken := ""

	for _, proto := range resp.Protocols {
		if proto.Protocol == "simple" || proto.Protocol == "spaces" {
			downloadEndpoint = proto.DownloadEndpoint
			downloadToken = proto.Token
		}
	}

	if downloadEndpoint == "" || downloadToken == "" {
		return http.Response{}, errors.New("download endpoint or token is missing")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	httpReq, err := http.NewRequest(http.MethodGet, downloadEndpoint, bytes.NewReader([]byte("")))
	if err != nil {
		return http.Response{}, err
	}
	httpReq.Header.Add("X-Reva-Transfer", downloadToken)
	httpReq.Header.Add("x-access-token", token)

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return http.Response{}, err
	}

	return *httpResp, nil
}
