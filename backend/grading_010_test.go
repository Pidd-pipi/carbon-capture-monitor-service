package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// committedWriter mimics net/http's committed-response semantics: headers
// written after WriteHeader are dropped from the wire.
type committedWriter struct {
	headers   http.Header
	committed bool
}

func (w *committedWriter) Header() http.Header {
	if w.committed {
		return http.Header{}
	}
	return w.headers
}

func (w *committedWriter) WriteHeader(code int) { w.committed = true }

func (w *committedWriter) Write(b []byte) (int, error) {
	w.committed = true
	return len(b), nil
}

// TestServerTimeoutsConfigured verifies the enterprise server keeps resource
// protection timeouts in place instead of dropping them.
func TestServerTimeoutsConfigured(t *testing.T) {
	srv := newEnterpriseServer(":0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if srv.IdleTimeout <= 0 {
		t.Fatalf("IdleTimeout must be positive, got %v", srv.IdleTimeout)
	}
	if srv.ReadTimeout <= 0 || srv.WriteTimeout <= 0 {
		t.Fatalf("read/write timeouts must be positive: %v/%v", srv.ReadTimeout, srv.WriteTimeout)
	}
}

// TestOpsLatencyHeaderPresent verifies every response carries the operations
// latency header even though the handler commits its status early.
func TestOpsLatencyHeaderPresent(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	w := &committedWriter{headers: http.Header{}}
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/x", nil))
	if got := w.headers.Get("X-Operations-Latency-Ms"); got == "" {
		t.Fatal("response missing X-Operations-Latency-Ms header")
	}
}

// TestOpsLatencyHeaderForEmptyResponse verifies the latency header is also set for
// responses that never call WriteHeader explicitly.
func TestOpsLatencyHeaderForEmptyResponse(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("no-status"))
	}))
	w := &committedWriter{headers: http.Header{}}
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/y", nil))
	if got := w.headers.Get("X-Operations-Latency-Ms"); got == "" {
		t.Fatal("response missing latency header for implicit status")
	}
}
