package readings

import (
	"testing"
	"time"
)

func TestGuardAggregateKeepsInputOrder(t *testing.T) {
	base := time.Now().UTC()
	input := []Reading{
		{CaptureRatePct: 90, PressureKPa: 180, SolventLevel: 70, RecordedAt: base.Add(time.Second)},
		{CaptureRatePct: 80, PressureKPa: 170, SolventLevel: 60, RecordedAt: base},
	}
	before := make([]Reading, len(input))
	copy(before, input)
	Aggregate(input)
	for i := range input {
		if input[i].CaptureRatePct != before[i].CaptureRatePct {
			t.Fatal("Aggregate mutated the input slice")
		}
	}
}

func TestGuardFilterInPlaceKeepsInput(t *testing.T) {
	base := time.Now().UTC()
	input := []Reading{
		{CaptureRatePct: 80, PressureKPa: 170, SolventLevel: 60, RecordedAt: base},
		{CaptureRatePct: 100, PressureKPa: 190, SolventLevel: 80, RecordedAt: base.Add(4 * time.Second)},
	}
	before := make([]Reading, len(input))
	copy(before, input)
	_ = FilterInPlace(input, base.Add(time.Second), base.Add(3*time.Second))
	for i := range input {
		if input[i].CaptureRatePct != before[i].CaptureRatePct {
			t.Fatal("FilterInPlace mutated the input backing array")
		}
	}
}
