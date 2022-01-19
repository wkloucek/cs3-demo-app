package http

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/disintegration/imaging"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/wkloucek/cs3-demo-app/pkg/internal/helpers"
)

var gwc gatewayv1beta1.GatewayAPIClient

func Server(ctx context.Context, gatewayclient gatewayv1beta1.GatewayAPIClient) (http.Handler, error) {

	gwc = gatewayclient

	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		PictureHandler(w, r)
	})

	err := http.ListenAndServe("localhost:6789", r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func PictureHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseForm()
	if err != nil {
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

	resp, err := helpers.DownloadFile(
		ctx,
		&providerv1beta1.Reference{
			ResourceId: &providerv1beta1.ResourceId{
				StorageId: storageID,
				OpaqueId:  opaqueID,
			},
			Path: ".",
		},
		gwc,
		token,
	)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	bodyImg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	img, _, _ := image.Decode(bytes.NewReader(bodyImg))
	img = imaging.Rotate(img, 180, color.Black)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	imgByte := buf.Bytes()

	w.Write(imgByte)
}
