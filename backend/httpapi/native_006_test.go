package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuardSummaryPropagatesClientCancel(t *testing.T) {
	handler := testHandler()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/readings/summary?unit=CC-ALPHA&window_minutes=60", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 499 {
		t.Fatalf("expected 499, got %d", rec.Code)
	}
}
