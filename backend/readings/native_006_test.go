package readings

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGuardSummaryHonorsCancel(t *testing.T) {
	store := NewStore(time.Hour, 100)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := store.Summary(ctx, "CC-ALPHA", time.Minute, time.Now().UTC())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context error, got %v", err)
	}
}
