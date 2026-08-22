package main

import (
	"context"
	"testing"
	"time"
)

// TestMonitorRunStopsOnCancel verifies the background loop exits on shutdown.
func TestGuardMonitorRunStopsOnCancel(t *testing.T) {
	m, _, _ := newTestMonitor()
	m.interval = 5 * time.Millisecond
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
		t.Fatal("monitor run loop did not stop after cancel")
	}
}
