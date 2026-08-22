package main

import (
	"context"
	"runtime"
	"testing"
	"time"

	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
)

// TestMonitorRunStopsOnCancel verifies the monitor loop exits promptly when
// the service is shutting down instead of leaking its goroutine.
func TestMonitorRunStopsOnCancel(t *testing.T) {
	m, _, readingStore := newTestMonitor()
	m.interval = 5 * time.Millisecond
	_ = readingStore
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		m.run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("monitor run loop did not stop after context cancel")
	}
}

// TestMonitorCycleCompletes drives a full evaluation cycle and asserts the
// triggered alert is raised without any worker panicking on a closed channel.
func TestMonitorCycleCompletes(t *testing.T) {
	old := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(old)

	m, alertSvc, readingStore := newTestMonitor()
	ctx := context.Background()
	readingStore.Append(ctx, "CC-ALPHA", readings.Reading{
		CaptureRatePct: 90, PressureKPa: 210, SolventLevel: 70, RecordedAt: time.Now().UTC(),
	})
	err := m.runCycle(ctx)
	if err != nil {
		t.Fatalf("runCycle: %v", err)
	}
	// Yield so any mis-synchronised worker goroutine actually runs before the
	// test process exits.
	time.Sleep(100 * time.Millisecond)

	page, err := alertSvc.Search(ctx, ops.OpsQuery{Subject: "pressure", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if page.Total == 0 {
		t.Fatal("expected a pressure alert to be created by the cycle")
	}
}

// TestMonitorCycleSurfacesUnitErrors verifies a failed unit evaluation is
// recorded instead of being silently swallowed.
func TestMonitorCycleSurfacesUnitErrors(t *testing.T) {
	m, _, _ := newTestMonitor()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := m.runCycle(ctx)
	if err == nil && m.LastError() == nil {
		t.Fatal("expected the cycle to surface the cancelled-context unit failure")
	}
}
