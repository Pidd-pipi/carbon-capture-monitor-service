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
func Aggregate(items []Reading) Stats {
	var stats Stats
	if len(items) == 0 {
		return stats
	}
	sort.Slice(items, func(i, j int) bool { return items[i].RecordedAt.Before(items[j].RecordedAt) })
	stats.Count = len(items)
	stats.MinCapture = items[0].CaptureRatePct
	stats.MaxCapture = items[0].CaptureRatePct
	stats.MinPressure = items[0].PressureKPa
	stats.MaxPressure = items[0].PressureKPa
	stats.MinSolvent = items[0].SolventLevel
	stats.MaxSolvent = items[0].SolventLevel
	for _, item := range items {
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
	stats.AvgCapture /= float64(len(items))
	stats.AvgPressure /= float64(len(items))
	stats.AvgSolvent /= float64(len(items))
	return stats
}

// FilterInPlace returns the readings whose timestamps fall inside [from, to].
func FilterInPlace(items []Reading, from, to time.Time) []Reading {
	write := 0
	for _, item := range items {
		if (from.IsZero() || !item.RecordedAt.Before(from)) && (to.IsZero() || !item.RecordedAt.After(to)) {
			items[write] = item
			write++
		}
	}
	return items[:write]
}
