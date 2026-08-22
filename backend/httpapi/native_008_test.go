package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardQueuedStatusView(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/alerts/status/queued", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
