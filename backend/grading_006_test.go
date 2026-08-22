package main

import (
	"context"
	"testing"

	"example.com/carbon-capture-monitor-service/domain"
	"example.com/carbon-capture-monitor-service/ops"
)

// TestApplyResultHonorsCancel verifies a cancelled monitor context stops the
// alert creation instead of continuing to write.
func TestApplyResultHonorsCancel(t *testing.T) {
	m, alertSvc, _ := newTestMonitor()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := cycleResult{
		unit:     domain.CaptureUnit{ID: "CC-ALPHA", Facility: "North Stack"},
		subject:  "CC-ALPHA pressure high: 210.0 kPa",
		priority: "critical",
		trigger:  true,
	}
	m.applyResult(ctx, result)

	page, err := alertSvc.Search(context.Background(), ops.OpsQuery{Subject: "pressure", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if page.Total != 0 {
		t.Fatalf("alert was created despite cancelled context: %d alerts", page.Total)
	}
}
