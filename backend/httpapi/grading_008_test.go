package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestQueuedStatusView verifies the status-scoped alert listing accepts
// the queued status and returns a page instead of an error.
func TestQueuedStatusView(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/alerts/status/queued", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for queued status list, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestEmptyStatusRejected verifies an empty status path is rejected
// instead of silently listing every alert.
func TestEmptyStatusRejected(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/alerts/status/", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty status, got %d: %s", rec.Code, rec.Body.String())
	}
}
