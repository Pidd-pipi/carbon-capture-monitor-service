package main

import (
	"context"
	"testing"

	"example.com/carbon-capture-monitor-service/domain"
)

// TestEvaluateUnitStatusWithoutReadings verifies a degraded unit still raises
// a status alert when no telemetry readings have arrived yet.
func TestEvaluateUnitStatusWithoutReadings(t *testing.T) {
	m, _, _ := newTestMonitor()
	result := m.evaluateUnit(context.Background(), domain.CaptureUnit{ID: "CC-BETA", Facility: "River Plant", Status: "attention"})
	if !result.trigger {
		t.Fatalf("expected a status alert for attention unit without readings, got %+v", result)
	}
	if result.priority != "normal" {
		t.Fatalf("expected normal priority, got %q", result.priority)
	}
}
