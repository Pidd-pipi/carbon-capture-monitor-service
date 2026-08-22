package readings

import (
	"sort"
	"time"
)

// Stats aggregates a set of readings into count/avg/min/max summaries.
type Stats struct {
	Count       int     `json:"count"`
	AvgCapture  float64 `json:"avg_capture_rate_pct"`
	MinCapture  float64 `json:"min_capture_rate_pct"`
	MaxCapture  float64 `json:"max_capture_rate_pct"`
	AvgPressure float64 `json:"avg_pressure_kpa"`
	MinPressure float64 `json:"min_pressure_kpa"`
	MaxPressure float64 `json:"max_pressure_kpa"`
	AvgSolvent  float64 `json:"avg_solvent_level_pct"`
	MinSolvent  float64 `json:"min_solvent_level_pct"`
	MaxSolvent  float64 `json:"max_solvent_level_pct"`
}

// Aggregate computes summary statistics over the given samples.
// A nil or empty input produces a zero-valued Stats with Count 0.
// The input slice is never mutated.
func Aggregate(items []Reading) Stats {
	var stats Stats
	if len(items) == 0 {
		return stats
	}
	sorted := make([]Reading, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].RecordedAt.Before(sorted[j].RecordedAt) })
	stats.Count = len(sorted)
	stats.MinCapture = sorted[0].CaptureRatePct
	stats.MaxCapture = sorted[0].CaptureRatePct
	stats.MinPressure = sorted[0].PressureKPa
	stats.MaxPressure = sorted[0].PressureKPa
	stats.MinSolvent = sorted[0].SolventLevel
	stats.MaxSolvent = sorted[0].SolventLevel
	for _, item := range sorted {
		stats.AvgCapture += item.CaptureRatePct
		stats.AvgPressure += item.PressureKPa
		stats.AvgSolvent += item.SolventLevel
		if item.CaptureRatePct < stats.MinCapture {
			stats.MinCapture = item.CaptureRatePct
		}
		if item.CaptureRatePct > stats.MaxCapture {
			stats.MaxCapture = item.CaptureRatePct
		}
		if item.PressureKPa < stats.MinPressure {
			stats.MinPressure = item.PressureKPa
		}
		if item.PressureKPa > stats.MaxPressure {
			stats.MaxPressure = item.PressureKPa
		}
		if item.SolventLevel < stats.MinSolvent {
			stats.MinSolvent = item.SolventLevel
		}
		if item.SolventLevel > stats.MaxSolvent {
			stats.MaxSolvent = item.SolventLevel
		}
	}
	stats.AvgCapture /= float64(len(sorted))
	stats.AvgPressure /= float64(len(sorted))
	stats.AvgSolvent /= float64(len(sorted))
	return stats
}

// FilterInPlace returns the readings whose timestamps fall inside [from, to].
// Despite the name, the returned slice does not alias the input backing array,
// so the caller's slice is left intact.
func FilterInPlace(items []Reading, from, to time.Time) []Reading {
	out := make([]Reading, 0, len(items))
	for _, item := range items {
		if (from.IsZero() || !item.RecordedAt.Before(from)) && (to.IsZero() || !item.RecordedAt.After(to)) {
			out = append(out, item)
		}
	}
	return out
}
