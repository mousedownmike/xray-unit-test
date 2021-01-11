package api

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"net/http"
)

type MyHandler struct {
	succeed bool
}

func (h MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, xs := xray.BeginSubsegment(r.Context(), "MyHandler.Get")
	h.useContext(ctx)
	_ = xs.ParentSegment.AddAnnotation("parent_annotation", "GET")
	_ = xs.AddAnnotation("method", "GET")
	if h.succeed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, e := w.Write([]byte("GET MyHandler"))
	xs.Close(e)
}

func (h MyHandler) Post(w http.ResponseWriter, r *http.Request) {
	ctx, xs := xray.BeginSubsegment(r.Context(), "MyHandler.Post")
	h.useContext(ctx)
	_ = xs.ParentSegment.AddAnnotation("parent_annotation", "POST")
	_ = xs.AddAnnotation("method", "POST")
	if h.succeed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, e := w.Write([]byte("POST MyHandler"))
	xs.Close(e)
}

func (h MyHandler) useContext(_ context.Context) {
	// make calls somewhere with the context
}
