package main

import (
	"context"
	"testing"
)

func TestGuardMonitorCycleSurfacesUnitErrors(t *testing.T) {
	m, _, _ := newTestMonitor()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := m.runCycle(ctx)
	if err == nil && m.LastError() == nil {
		t.Fatal("cycle swallowed the cancelled-context unit failure")
	}
}
