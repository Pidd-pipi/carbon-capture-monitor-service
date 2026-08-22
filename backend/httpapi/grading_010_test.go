package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAppJServed verifies the browser asset is served by the embedded page.
func TestAppJServed(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/app.js", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for /app.js, got %d", rec.Code)
	}
}
