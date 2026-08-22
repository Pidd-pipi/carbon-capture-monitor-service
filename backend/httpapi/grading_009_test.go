package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUnknownUnitLookup verifies updating an unknown unit returns
// 404 rather than a server error.
func TestUnknownUnitLookup(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-NOPE","status":"online"}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown unit, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestMaintenanceToOnlineRejected verifies the maintenance->online
// transition is rejected as a bad request instead of a server error.
func TestMaintenanceToOnlineRejected(t *testing.T) {
	handler := testHandler()
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-ALPHA","status":"maintenance"}`)))
	if first.Code != http.StatusOK {
		t.Fatalf("set maintenance: %d %s", first.Code, first.Body.String())
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-ALPHA","status":"online"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for maintenance->online, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestCaptureUnitValidUpdateStillWorks verifies a legitimate transition still
// succeeds end to end.
func TestCaptureUnitValidUpdateStillWorks(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-BETA","status":"online"}`)))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"status":"online"`) {
		t.Fatalf("valid update: %d %s", rec.Code, rec.Body.String())
	}
}
