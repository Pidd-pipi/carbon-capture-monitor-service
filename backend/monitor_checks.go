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
// new alert should be raised for it.
func (m *Monitor) evaluateUnit(ctx context.Context, unit domain.CaptureUnit) cycleResult {
	result := cycleResult{unit: unit}
	if ctx.Err() != nil {
		result.err = ctx.Err()
		return result
	}
	latest, ok := m.readings.Latest(ctx, unit.ID)
	if !ok {
		return result
	}
	switch {
	case latest.PressureKPa >= pressureHighThreshold:
		result.trigger = true
		result.priority = "critical"
		result.subject = fmt.Sprintf("%s pressure high: %.1f kPa", unit.ID, latest.PressureKPa)
	case latest.CaptureRatePct < captureLowThreshold:
		result.trigger = true
		result.priority = "high"
		result.subject = fmt.Sprintf("%s capture efficiency low: %.1f%%", unit.ID, latest.CaptureRatePct)
	case latest.SolventLevel < solventLowThreshold:
		result.trigger = true
		result.priority = "normal"
		result.subject = fmt.Sprintf("%s solvent level low: %.1f%%", unit.ID, latest.SolventLevel)
	case strings.EqualFold(unit.Status, "attention") || strings.EqualFold(unit.Status, "offline"):
		result.trigger = true
		result.priority = "normal"
		result.subject = fmt.Sprintf("%s unit status %s", unit.ID, unit.Status)
	}
	return result
}
