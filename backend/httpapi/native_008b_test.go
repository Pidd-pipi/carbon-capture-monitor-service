package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardEmptyStatusRejected(t *testing.T) {
	rec := httptest.NewRecorder()
	testHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/alerts/status/", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
