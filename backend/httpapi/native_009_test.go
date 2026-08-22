package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuardUnknownUnitLookup(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-NOPE","status":"online"}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGuardMaintenanceToOnlineRejected(t *testing.T) {
	handler := testHandler()
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-ALPHA","status":"maintenance"}`)))
	if first.Code != http.StatusOK {
		t.Fatalf("maintenance: %d", first.Code)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-ALPHA","status":"online"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
