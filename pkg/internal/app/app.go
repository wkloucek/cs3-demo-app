package app

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	registryv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/disintegration/imaging"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/helpers"
	"google.golang.org/grpc"
)

type demoApp struct {
	gwc        gatewayv1beta1.GatewayAPIClient
	grpcServer *grpc.Server
}

func New() *demoApp {
	return &demoApp{}
}

func (app *demoApp) GetCS3apiClient() error {
	// establish a connection to the cs3 api endpoint
	// in this case a REVA gateway, started by oCIS
	gwc, err := pool.GetGatewayServiceClient("localhost:9142")
	if err != nil {
		return err
	}
	app.gwc = gwc

	return nil
}

func (app *demoApp) RegisterDemoApp(ctx context.Context) error {
	req := &registryv1beta1.AddAppProviderRequest{
		Provider: &registryv1beta1.ProviderInfo{
			Name:        "demoapp",
			Description: "this is an demo app",
			Icon:        "image-edit",
			Address:     "127.0.0.1:5678", // address of the grpc server we start in this demo app
			MimeTypes: []string{
				// supported mime types
				"image/png",
				"image/jpeg",
				"image/gif",
			},
		},
	}

	resp, err := app.gwc.AddAppProvider(ctx, req)
	if err != nil {
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		return errors.New("status code != CODE_OK")
	}

	return nil
}

func (app *demoApp) GRPCServer(ctx context.Context) error {
	opts := []grpc.ServerOption{}
	app.grpcServer = grpc.NewServer(opts...)

	// register the app provider interface / OpenInApp call
	appproviderv1beta1.RegisterProviderAPIServer(app.grpcServer, app)

	l, err := net.Listen("tcp", "localhost:5678")
	if err != nil {
		return err
	}
	go app.grpcServer.Serve(l)

	return nil
}

func (app *demoApp) OpenInApp(ctx context.Context, req *appproviderv1beta1.OpenInAppRequest) (*appproviderv1beta1.OpenInAppResponse, error) {
	return &appproviderv1beta1.OpenInAppResponse{
		Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
		AppUrl: &appproviderv1beta1.OpenInAppURL{
			AppUrl: "http://localhost:6789",
			Method: "POST",
			FormParameters: map[string]string{
				// these parameters will be passed to the web server by the app provider application
				"access_token": req.AccessToken,
				"storage_id":   req.ResourceInfo.Id.StorageId,
				"opaque_id":    req.ResourceInfo.Id.OpaqueId,
			},
		},
	}, nil
}

func (app *demoApp) HTTPServer(ctx context.Context) error {
	// start a simple web server that will get requests from
	// app provider client, eg. ownCloud Web
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		PictureHandler(app, w, r)
	})

	if err := http.ListenAndServe("localhost:6789", r); err != nil {
		return err
	}
	return nil
}

func PictureHandler(app *demoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get form parameters
	if err := r.ParseForm(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	token := r.Form.Get("access_token")
	storageID := r.Form.Get("storage_id")
	opaqueID := r.Form.Get("opaque_id")
	if token == "" || storageID == "" || opaqueID == "" {
		render.Status(r, http.StatusBadRequest)
		return
	}

	// download the image
	resp, err := helpers.DownloadFile(
		ctx,
		&providerv1beta1.Reference{
			ResourceId: &providerv1beta1.ResourceId{
				StorageId: storageID,
				OpaqueId:  opaqueID,
			},
			Path: ".",
		},
		app.gwc,
		token,
	)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	// read the image from the body
	defer resp.Body.Close()
	bodyImg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	// decode the image
	img, _, err := image.Decode(bytes.NewReader(bodyImg))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	// do whatever needed with the image
	img = imaging.Rotate(img, 180, color.Black)

	// convert back to a file
	buf := new(bytes.Buffer)
	if err = png.Encode(buf, img); err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	// normally you would return some proper html
	// but we will just return the image here
	w.Write(buf.Bytes())
}
