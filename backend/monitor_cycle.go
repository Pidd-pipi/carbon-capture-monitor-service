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
	wg.Add(len(units))
	for _, unit := range units {
		unit := unit
		go func() {
			defer wg.Done()
			results <- m.evaluateUnit(ctx, unit)
		}()
	}
	wg.Wait()
	close(results)

	var firstErr error
	for result := range results {
		if result.err != nil && firstErr == nil {
			firstErr = result.err
			continue
		}
		if result.trigger {
			m.applyResult(ctx, result)
		}
	}
	m.readings.Prune(m.clock())
	m.recordCycle(firstErr)
	return firstErr
}
