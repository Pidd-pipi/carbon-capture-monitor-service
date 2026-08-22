package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type guardWriterB struct {
	headers   http.Header
	committed bool
}

func (w *guardWriterB) Header() http.Header {
	if w.committed {
		return http.Header{}
	}
	return w.headers
}
func (w *guardWriterB) WriteHeader(code int) { w.committed = true }
func (w *guardWriterB) Write(b []byte) (int, error) { w.committed = true; return len(b), nil }

func TestGuardLatencyHeaderImplicitStatus(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	w := &guardWriterB{headers: http.Header{}}
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/y", nil))
	if got := w.headers.Get("X-Operations-Latency-Ms"); got == "" {
		t.Fatal("missing latency header for implicit status")
	}
}
