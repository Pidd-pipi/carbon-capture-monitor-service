package main

import (
	"context"
	"testing"
	"time"

	"example.com/carbon-capture-monitor-service/domain"
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
)

func newTestMonitor() (*Monitor, *ops.OpsService, *readings.Store) {
	alertSvc := ops.NewService(nil)
	readingStore := readings.NewStore(time.Hour, 100)
	unitStore := store.New()
	m := newMonitor(alertSvc, readingStore, unitStore, time.Hour)
	m.clock = func() time.Time { return time.Now().UTC() }
	return m, alertSvc, readingStore
}

func TestEvaluateUnitHighPressure(t *testing.T) {
	m, _, readingStore := newTestMonitor()
	readingStore.Append(context.Background(), "CC-ALPHA", readings.Reading{
		CaptureRatePct: 90, PressureKPa: 210, SolventLevel: 70, RecordedAt: time.Now().UTC(),
	})
	result := m.evaluateUnit(context.Background(), domain.CaptureUnit{ID: "CC-ALPHA", Facility: "North Stack"})
	if !result.trigger || result.priority != "critical" {
		t.Fatalf("expected critical pressure alert, got %+v", result)
	}
}

func TestEvaluateUnitNoAlertForHealthy(t *testing.T) {
	m, _, readingStore := newTestMonitor()
	readingStore.Append(context.Background(), "CC-BETA", readings.Reading{
		CaptureRatePct: 92, PressureKPa: 175, SolventLevel: 75, RecordedAt: time.Now().UTC(),
	})
	result := m.evaluateUnit(context.Background(), domain.CaptureUnit{ID: "CC-BETA", Facility: "River Plant", Status: "online"})
	if result.trigger {
		t.Fatalf("did not expect an alert for healthy unit: %+v", result)
	}
}

func TestApplyResultCreatesAlert(t *testing.T) {
	m, alertSvc, _ := newTestMonitor()
	result := cycleResult{
		unit:     domain.CaptureUnit{ID: "CC-ALPHA", Facility: "North Stack"},
		subject:  "CC-ALPHA solvent level low: 55.0%",
		priority: "normal",
		trigger:  true,
	}
	m.applyResult(context.Background(), result)
	page, err := alertSvc.Search(context.Background(), ops.OpsQuery{Subject: "solvent", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if page.Total != 1 {
		t.Fatalf("expected 1 alert, got %d", page.Total)
	}
}
