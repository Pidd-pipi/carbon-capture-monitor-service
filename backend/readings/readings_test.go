package readings

import (
	"context"
	"testing"
	"time"
)

func TestAppendListLatest(t *testing.T) {
	store := NewStore(time.Hour, 10)
	ctx := context.Background()
	base := time.Now().UTC().Add(-time.Minute)
	for i := 0; i < 3; i++ {
		err := store.Append(ctx, "CC-ALPHA", Reading{
			CaptureRatePct: float64(80 + i),
			PressureKPa:    170 + float64(i),
			SolventLevel:   70 - float64(i),
			RecordedAt:     base.Add(time.Duration(i) * time.Second),
		})
		if err != nil {
			t.Fatalf("append: %v", err)
		}
	}
	items, err := store.List(ctx, "CC-ALPHA", time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 readings, got %d", len(items))
	}
	// Newest first
	if items[0].CaptureRatePct != 82 {
		t.Fatalf("expected newest first, got %+v", items[0])
	}
	latest, ok := store.Latest(ctx, "CC-ALPHA")
	if !ok || latest.CaptureRatePct != 82 {
		t.Fatalf("unexpected latest: %+v ok=%v", latest, ok)
	}
}

func TestRetentionPrunesOldReadings(t *testing.T) {
	store := NewStore(10*time.Second, 100)
	ctx := context.Background()
	now := time.Now().UTC()
	store.Append(ctx, "CC-BETA", Reading{CaptureRatePct: 90, PressureKPa: 170, SolventLevel: 70, RecordedAt: now.Add(-time.Minute)})
	store.Append(ctx, "CC-BETA", Reading{CaptureRatePct: 91, PressureKPa: 171, SolventLevel: 71, RecordedAt: now})
	if count := store.Count("CC-BETA"); count != 1 {
		t.Fatalf("expected 1 retained reading, got %d", count)
	}
}

func TestAggregate(t *testing.T) {
	items := []Reading{
		{CaptureRatePct: 80, PressureKPa: 170, SolventLevel: 60},
		{CaptureRatePct: 90, PressureKPa: 180, SolventLevel: 70},
		{CaptureRatePct: 100, PressureKPa: 190, SolventLevel: 80},
	}
	stats := Aggregate(items)
	if stats.Count != 3 || stats.AvgCapture != 90 || stats.MinCapture != 80 || stats.MaxCapture != 100 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
	if stats.MinPressure != 170 || stats.MaxPressure != 190 || stats.AvgSolvent != 70 {
		t.Fatalf("unexpected pressure/solvent stats: %+v", stats)
	}
	if empty := Aggregate(nil); empty.Count != 0 {
		t.Fatalf("expected empty stats, got %+v", empty)
	}
}

func TestSummaryWindow(t *testing.T) {
	store := NewStore(time.Hour, 100)
	ctx := context.Background()
	now := time.Now().UTC()
	store.Append(ctx, "CC-GAMMA", Reading{CaptureRatePct: 88, PressureKPa: 175, SolventLevel: 72, RecordedAt: now.Add(-5 * time.Second)})
	summary, err := store.Summary(ctx, "CC-GAMMA", time.Minute, now)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Stats.Count != 1 || summary.Stats.AvgCapture != 88 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}
