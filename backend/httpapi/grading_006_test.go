package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/carbon-capture-monitor-service/readings"
)

// TestSummaryHonorsCancel verifies the summary path returns the context error
// instead of silently completing after cancellation.
func TestSummaryHonorsCancel(t *testing.T) {
	store := readings.NewStore(time.Hour, 100)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := store.Summary(ctx, "CC-ALPHA", time.Minute, time.Now().UTC())
	if err == nil {
		t.Fatal("Summary completed despite a cancelled context")
	}
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected a context error, got %v", err)
	}
}

// TestSummaryPropagatesClientCancel verifies the summary endpoint honours the
// request context so a cancelled client stops the aggregation.
func TestSummaryPropagatesClientCancel(t *testing.T) {
	handler := testHandler()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/readings/summary?unit=CC-ALPHA&window_minutes=60", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 499 {
		t.Fatalf("expected 499 for cancelled client, got %d: %s", rec.Code, rec.Body.String())
	}
}
