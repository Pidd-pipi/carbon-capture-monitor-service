package main

import (
	"context"
	"fmt"
	"strings"

	"example.com/carbon-capture-monitor-service/domain"
)

const (
	pressureHighThreshold = 190.0
	captureLowThreshold   = 85.0
	solventLowThreshold   = 60.0
)

// evaluateUnit inspects the latest reading of one unit and decides whether a
// new alert should be raised for it. The returned subject describes the
// defect condition rather than the raw reading, so a unit with a sustained
// fault produces one stable subject across jittering measurements instead of
// a new one every cycle.
func (m *Monitor) evaluateUnit(ctx context.Context, unit domain.CaptureUnit) cycleResult {
	result := cycleResult{unit: unit}
	if ctx.Err() != nil {
		result.err = ctx.Err()
		return result
	}
	latest, ok := m.readings.Latest(ctx, unit.ID)
	if ok {
		switch {
		case latest.PressureKPa >= pressureHighThreshold:
			result.trigger = true
			result.priority = "critical"
			result.subject = fmt.Sprintf("%s pressure high", unit.ID)
		case latest.CaptureRatePct < captureLowThreshold:
			result.trigger = true
			result.priority = "high"
			result.subject = fmt.Sprintf("%s capture efficiency low", unit.ID)
		case latest.SolventLevel < solventLowThreshold:
			result.trigger = true
			result.priority = "normal"
			result.subject = fmt.Sprintf("%s solvent level low", unit.ID)
		}
	}
	// A unit flagged attention/offline needs surfacing even before any
	// telemetry arrives, so operators learn about a misbehaving unit rather
	// than silently waiting for the first reading.
	if !result.trigger && (strings.EqualFold(unit.Status, "attention") || strings.EqualFold(unit.Status, "offline")) {
		result.trigger = true
		result.priority = "normal"
		result.subject = fmt.Sprintf("%s unit status %s", unit.ID, unit.Status)
	}
	return result
}
