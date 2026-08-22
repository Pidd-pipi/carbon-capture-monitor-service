package main

import (
	"context"
	"testing"

	"example.com/carbon-capture-monitor-service/domain"
	"example.com/carbon-capture-monitor-service/ops"
)

func TestGuardApplyResultHonorsCancel(t *testing.T) {
	m, alertSvc, _ := newTestMonitor()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	m.applyResult(ctx, cycleResult{unit: domain.CaptureUnit{ID: "CC-ALPHA", Facility: "North Stack"}, subject: "CC-ALPHA pressure high", priority: "critical", trigger: true})
	page, _ := alertSvc.Search(context.Background(), ops.OpsQuery{Subject: "pressure", Page: 1, PageSize: 10})
	if page.Total != 0 {
		t.Fatalf("alert created despite cancelled ctx: %d", page.Total)
	}
}
