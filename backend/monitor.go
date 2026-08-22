package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
)

// Monitor periodically evaluates the latest telemetry of every capture unit,
// raises alert records through the ops service and prunes stale readings.
type Monitor struct {
	alerts   *ops.OpsService
	readings *readings.Store
	units    *store.Store
	interval time.Duration
	clock    func() time.Time

	mu      sync.Mutex
	cycles  int
	lastErr error
}

func newMonitor(alerts *ops.OpsService, readingStore *readings.Store, unitStore *store.Store, interval time.Duration) *Monitor {
	return &Monitor{
		alerts:   alerts,
		readings: readingStore,
		units:    unitStore,
		interval: interval,
		clock:    time.Now,
	}
}

// Start launches the background evaluation loop and returns a cancel function.
func (m *Monitor) Start(parent context.Context) context.CancelFunc {
	ctx, cancel := context.WithCancel(parent)
	go m.run(ctx)
	return cancel
}

func (m *Monitor) run(ctx context.Context) {
	for {
		timer := time.NewTimer(m.interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
		m.runCycleSafely(ctx)
	}
}

// runCycleSafely runs one evaluation pass and recovers from any panic so a
// single bad cycle never silently kills the background patrol loop. The panic
// is surfaced via LastError instead, and the next tick continues as scheduled.
func (m *Monitor) runCycleSafely(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			m.recordCycle(fmt.Errorf("monitor cycle panic: %v", r))
		}
	}()
	_ = m.runCycle(ctx)
}

func (m *Monitor) CycleCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cycles
}

func (m *Monitor) LastError() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastErr
}

func (m *Monitor) recordCycle(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cycles++
	m.lastErr = err
}
