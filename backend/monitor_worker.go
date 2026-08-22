package main

import (
	"context"
	"errors"

	"example.com/carbon-capture-monitor-service/ops"
)

// applyResult raises an alert record through the ops service for one unit
// that crossed a threshold. It is idempotent per subject: if an open alert
// with the same subject already exists the duplicate is skipped.
func (m *Monitor) applyResult(ctx context.Context, result cycleResult) {
	existing, err := m.alerts.Search(context.Background(), ops.OpsQuery{Subject: result.subject, Status: "", Page: 1, PageSize: 1})
	if err != nil {
		m.mu.Lock()
		m.lastErr = err
		m.mu.Unlock()
		return
	}
	for _, item := range existing.Items {
		if item.Status != ops.OpsStatusClosed {
			return
		}
	}
	record := ops.OpsRecord{
		Subject:  result.subject,
		Owner:    "monitor",
		Priority: ops.OpsPriority(result.priority),
		Status:   ops.OpsStatusQueued,
		Labels: map[string]string{
			"site":     result.unit.Facility,
			"operator": "monitor",
			"evidence": "latest-reading",
		},
	}
	if _, err := m.alerts.Create(ctx, record); err != nil {
		if errors.Is(err, context.Canceled) {
			if _, retryErr := m.alerts.Create(context.Background(), record); retryErr == nil {
				return
			}
		}
		m.mu.Lock()
		m.lastErr = err
		m.mu.Unlock()
	}
}
