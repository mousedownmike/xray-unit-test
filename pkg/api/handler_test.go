package api

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMyHandler_Get(t *testing.T) {
	var h MyHandler
	var w *httptest.ResponseRecorder
	var r *http.Request

	mockCtx, xs := xray.BeginSegment(context.Background(), t.Name())
	defer xs.Close(nil)
	before := func() {
		h = MyHandler{}
		w = httptest.NewRecorder()
		r = &http.Request{}
		r = r.WithContext(mockCtx)
	}

	t.Run("success", func(t *testing.T) {
		before()
		h.succeed = true
		h.Get(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "GET", xs.Annotations["method"])
	})
	t.Run("fail", func(t *testing.T) {
		before()
		h.succeed = false
		h.Get(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "GET", xs.Annotations["method"])
	})
}

func TestMyHandler_Post(t *testing.T) {
	var h MyHandler
	var w *httptest.ResponseRecorder
	var r *http.Request

	mockCtx, xs := xray.BeginSegment(context.Background(), t.Name())
	defer xs.Close(nil)
	before := func() {
		h = MyHandler{}
		w = httptest.NewRecorder()
		r = &http.Request{}
		r = r.WithContext(mockCtx)
	}

	t.Run("success", func(t *testing.T) {
		before()
		h.succeed = true
		h.Post(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "POST", xs.Annotations["method"])
	})
	t.Run("fail", func(t *testing.T) {
		before()
		h.succeed = false
		h.Post(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "POST", xs.Annotations["method"])
	})
}
