package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardServerTimeoutsConfigured(t *testing.T) {
	srv := newEnterpriseServer(":0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if srv.IdleTimeout <= 0 {
		t.Fatalf("IdleTimeout must be positive, got %v", srv.IdleTimeout)
	}
}

type guardCommittedWriter struct {
	headers   http.Header
	committed bool
}

func (w *guardCommittedWriter) Header() http.Header {
	if w.committed {
		return http.Header{}
	}
	return w.headers
}
func (w *guardCommittedWriter) WriteHeader(code int) { w.committed = true }
func (w *guardCommittedWriter) Write(b []byte) (int, error) { w.committed = true; return len(b), nil }

func TestGuardOpsLatencyHeaderPresent(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	w := &guardCommittedWriter{headers: http.Header{}}
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/x", nil))
	if got := w.headers.Get("X-Operations-Latency-Ms"); got == "" {
		t.Fatal("missing latency header")
	}
}
