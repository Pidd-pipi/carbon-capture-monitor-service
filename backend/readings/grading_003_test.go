package readings

import (
	"context"
	"testing"
	"time"
)

// TestReadingsListIsolated verifies that a filtered list never corrupts the
// retained history of a unit.
func TestReadingsListIsolated(t *testing.T) {
	store := NewStore(time.Hour, 100)
	ctx := context.Background()
	base := time.Now().UTC().Add(-time.Minute)
	for i := 0; i < 3; i++ {
		store.Append(ctx, "CC-ALPHA", Reading{
			CaptureRatePct: float64(80 + i),
			PressureKPa:    170 + float64(i),
			SolventLevel:   70 - float64(i),
			RecordedAt:     base.Add(time.Duration(i) * time.Second),
		})
	}
	// A windowed query must not drop or duplicate readings in the store.
	_, err := store.List(ctx, "CC-ALPHA", base.Add(time.Second), time.Time{}, 0)
	if err != nil {
		t.Fatalf("windowed list: %v", err)
	}
	full, err := store.List(ctx, "CC-ALPHA", time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatalf("full list: %v", err)
	}
	if len(full) != 3 {
		t.Fatalf("history was corrupted by the windowed query: got %d readings, want 3", len(full))
	}
	seen := map[float64]bool{}
	for _, item := range full {
		seen[item.CaptureRatePct] = true
	}
	for _, want := range []float64{80, 81, 82} {
		if !seen[want] {
			t.Fatalf("reading %v missing after windowed query: %+v", want, full)
		}
	}
}

// TestFilterInPlaceDoesNotMutateInput verifies the filter returns a fresh
// slice and leaves the caller's backing array untouched.
func TestFilterInPlaceDoesNotMutateInput(t *testing.T) {
	base := time.Now().UTC()
	input := []Reading{
		{CaptureRatePct: 80, PressureKPa: 170, SolventLevel: 60, RecordedAt: base},
		{CaptureRatePct: 90, PressureKPa: 180, SolventLevel: 70, RecordedAt: base.Add(2 * time.Second)},
		{CaptureRatePct: 100, PressureKPa: 190, SolventLevel: 80, RecordedAt: base.Add(4 * time.Second)},
	}
	before := make([]Reading, len(input))
	copy(before, input)

	filtered := FilterInPlace(input, base.Add(time.Second), base.Add(3*time.Second))
	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered reading, got %d", len(filtered))
	}
	for i := range input {
		if input[i].CaptureRatePct != before[i].CaptureRatePct {
			t.Fatalf("FilterInPlace mutated the input backing array at index %d", i)
		}
	}
}

// TestAppendKeepsNewest verifies that when the per-unit cap is exceeded the
// newest readings are retained instead of the oldest.
func TestAppendKeepsNewest(t *testing.T) {
	store := NewStore(time.Hour, 3)
	ctx := context.Background()
	base := time.Now().UTC().Add(-time.Minute)
	for i := 0; i < 5; i++ {
		store.Append(ctx, "CC-GAMMA", Reading{
			CaptureRatePct: float64(70 + i),
			PressureKPa:    160 + float64(i),
			SolventLevel:   60 - float64(i),
			RecordedAt:     base.Add(time.Duration(i) * time.Second),
		})
	}
	items, err := store.List(ctx, "CC-GAMMA", time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 retained readings, got %d", len(items))
	}
	seen := map[float64]bool{}
	for _, item := range items {
		seen[item.CaptureRatePct] = true
	}
	for _, want := range []float64{72, 73, 74} {
		if !seen[want] {
			t.Fatalf("newest reading %v not retained: %+v", want, items)
		}
	}
}

// TestAggregateDoesNotMutateInput verifies aggregation has no side effects on
// the caller's slice.
func TestAggregateDoesNotMutateInput(t *testing.T) {
	base := time.Now().UTC()
	input := []Reading{
		{CaptureRatePct: 90, PressureKPa: 180, SolventLevel: 70, RecordedAt: base.Add(time.Second)},
		{CaptureRatePct: 80, PressureKPa: 170, SolventLevel: 60, RecordedAt: base},
		{CaptureRatePct: 100, PressureKPa: 190, SolventLevel: 80, RecordedAt: base.Add(2 * time.Second)},
	}
	before := make([]Reading, len(input))
	copy(before, input)

	Aggregate(input)

	for i := range input {
		if input[i].CaptureRatePct != before[i].CaptureRatePct {
			t.Fatalf("Aggregate mutated the input slice at index %d", i)
		}
	}
}
