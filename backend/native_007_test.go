package main

import (
	"context"
	"testing"

	"example.com/carbon-capture-monitor-service/domain"
)

func TestGuardEvaluateUnitStatusWithoutReadings(t *testing.T) {
	m, _, _ := newTestMonitor()
	result := m.evaluateUnit(context.Background(), domain.CaptureUnit{ID: "CC-BETA", Facility: "River Plant", Status: "attention"})
	if !result.trigger {
		t.Fatalf("expected status alert, got %+v", result)
	}
}
