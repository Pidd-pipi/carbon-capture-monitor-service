package health

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHealthReportsStoreOK verifies the health endpoint reflects a healthy
// store dependency instead of an unknown state.
func TestHealthReportsStoreOK(t *testing.T) {
	dependencyOK = true
	handler := Handler("carbon-capture-monitor-service")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("health: %d %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"store":"ok"`) {
		t.Fatalf("health body missing store=ok: %s", rec.Body.String())
	}
}
