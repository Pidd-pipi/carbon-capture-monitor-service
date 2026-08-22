package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardAppJServed(t *testing.T) {
	rec := httptest.NewRecorder()
	testHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/app.js", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for /app.js, got %d", rec.Code)
	}
}
