package main

import (
	"context"
	"sync"

	"example.com/carbon-capture-monitor-service/domain"
)

// cycleResult carries one unit's evaluation outcome back to the coordinator.
type cycleResult struct {
	unit     domain.CaptureUnit
	subject  string
	priority string
	trigger  bool
	err      error
}

// runCycle evaluates every capture unit concurrently and applies the outcomes.
func (m *Monitor) runCycle(ctx context.Context) error {
	units := m.units.List()
	results := make(chan cycleResult, len(units))

	var wg sync.WaitGroup
	for _, unit := range units {
		unit := unit
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := m.evaluateUnit(ctx, unit)
			select {
			case results <- result:
			case <-ctx.Done():
			}
		}()
	}
	wg.Wait()
	close(results)

	pruned := 0
	for result := range results {
		if result.err != nil {
			continue
		}
		if result.trigger {
			m.applyResult(ctx, result)
		}
	}
	now := m.clock()
	pruned = m.readings.Prune(now)
	_ = pruned
	m.recordCycle(nil)
	return nil
}
