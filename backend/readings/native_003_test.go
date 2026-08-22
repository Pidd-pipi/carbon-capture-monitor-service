package readings

import (
	"context"
	"testing"
	"time"
)

func TestGuardReadingsListIsolated(t *testing.T) {
	store := NewStore(time.Hour, 100)
	ctx := context.Background()
	base := time.Now().UTC().Add(-time.Minute)
	for i := 0; i < 3; i++ {
		store.Append(ctx, "CC-ALPHA", Reading{CaptureRatePct: float64(80 + i), PressureKPa: 170 + float64(i), SolventLevel: 70 - float64(i), RecordedAt: base.Add(time.Duration(i) * time.Second)})
	}
	_, _ = store.List(ctx, "CC-ALPHA", base.Add(time.Second), time.Time{}, 0)
	full, _ := store.List(ctx, "CC-ALPHA", time.Time{}, time.Time{}, 0)
	if len(full) != 3 {
		t.Fatalf("windowed query corrupted history: %d", len(full))
	}
}

func TestGuardAppendKeepsNewest(t *testing.T) {
	store := NewStore(time.Hour, 3)
	ctx := context.Background()
	base := time.Now().UTC().Add(-time.Minute)
	for i := 0; i < 5; i++ {
		store.Append(ctx, "CC-GAMMA", Reading{CaptureRatePct: float64(70 + i), PressureKPa: 160 + float64(i), SolventLevel: 60 - float64(i), RecordedAt: base.Add(time.Duration(i) * time.Second)})
	}
	items, _ := store.List(ctx, "CC-GAMMA", time.Time{}, time.Time{}, 0)
	if len(items) != 3 {
		t.Fatalf("expected 3 retained readings, got %d", len(items))
	}
	if items[0].CaptureRatePct != 74 {
		t.Fatalf("newest reading not retained: %+v", items[0])
	}
}
