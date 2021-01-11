package api

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AlwaysSamplingStrategy struct{}

func (s *AlwaysSamplingStrategy) ShouldTrace(_ *sampling.Request) *sampling.Decision {
	return &sampling.Decision{Sample: true}
}

type MockEmitter struct {
	xray.Emitter
	segments []*xray.Segment
}

func (m *MockEmitter) Emit(s *xray.Segment) {
	m.segments = append(m.segments, s)
}

func TestMyHandler_Get(t *testing.T) {
	var h MyHandler
	var w *httptest.ResponseRecorder
	var r *http.Request
	me := &MockEmitter{}

	ctx, _ := xray.ContextWithConfig(context.Background(), xray.Config{
		SamplingStrategy: &AlwaysSamplingStrategy{},
		Emitter:          me,
	})
	ctx, xs := xray.BeginSegment(ctx, t.Name())
	defer xs.Close(nil)
	before := func() {
		h = MyHandler{}
		w = httptest.NewRecorder()
		r = &http.Request{Method: http.MethodGet}
		r = r.WithContext(ctx)
	}

	t.Run("success", func(t *testing.T) {
		before()
		h.succeed = true
		h.Get(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "GET", xs.Annotations["parent_annotation"])
		assert.Len(t, me.segments, 1)
	})
	t.Run("fail", func(t *testing.T) {
		before()
		h.succeed = false
		h.Get(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		xs.Close(nil)
		assert.Equal(t, "GET", xs.Annotations["parent_annotation"])
	})
}

func TestMyHandler_Post(t *testing.T) {
	var h MyHandler
	var w *httptest.ResponseRecorder
	var r *http.Request

	ctx, _ := xray.ContextWithConfig(context.Background(), xray.Config{SamplingStrategy: &AlwaysSamplingStrategy{}})
	ctx, xs := xray.BeginSegment(ctx, t.Name())
	defer xs.Close(nil)
	before := func() {
		h = MyHandler{}
		w = httptest.NewRecorder()
		r = &http.Request{Method: http.MethodPost}
		r = r.WithContext(ctx)
	}

	t.Run("success", func(t *testing.T) {
		before()
		h.succeed = true
		h.Post(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		s := xray.GetSegment(ctx)
		assert.Equal(t, "POST", s.Annotations["method"])
		assert.Equal(t, "POST", xs.Annotations["parent_annotation"])
	})
	t.Run("fail", func(t *testing.T) {
		before()
		h.succeed = false
		h.Post(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "POST", xs.Annotations["parent_annotation"])
	})
}
