package ops

import (
	"context"
	"testing"
)

func TestServiceCreateSearchTransition(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()

	created, err := svc.Create(ctx, OpsRecord{
		Subject:  "CC-ALPHA pressure high",
		Owner:    "alice",
		Priority: OpsPriorityHigh,
		Labels:   map[string]string{"site": "North Stack"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" || created.Status != OpsStatusQueued || created.Revision != 1 {
		t.Fatalf("unexpected created record: %+v", created)
	}

	page, err := svc.Search(ctx, OpsQuery{Subject: "pressure", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("search total=%d items=%d", page.Total, len(page.Items))
	}

	updated, err := svc.Transition(ctx, created.ID, 1, OpsStatusActive, "alice")
	if err != nil {
		t.Fatalf("transition: %v", err)
	}
	if updated.Status != OpsStatusActive || updated.Revision != 2 {
		t.Fatalf("unexpected updated record: %+v", updated)
	}

	audit := svc.Audit(created.ID)
	if len(audit) != 2 {
		t.Fatalf("expected 2 audit events, got %d", len(audit))
	}
}

func TestServicePolicyRejectsMissingOwner(t *testing.T) {
	svc := NewService(nil)
	_, err := svc.Create(context.Background(), OpsRecord{
		Subject:  "no owner alert",
		Priority: OpsPriorityLow,
		Labels:   map[string]string{"site": "North Stack"},
	})
	if err == nil {
		t.Fatal("expected policy error for missing owner")
	}
}

func TestStateMachineTransitionRules(t *testing.T) {
	m := NewStateMachine()
	if err := m.Move(OpsStatusQueued, OpsStatusActive, "start"); err != nil {
		t.Fatalf("queued->active should be allowed: %v", err)
	}
	if err := m.Move(OpsStatusQueued, OpsStatusClosed, "close"); err != nil {
		t.Fatalf("queued->closed should be allowed: %v", err)
	}
	if err := m.Move(OpsStatusActive, OpsStatusQueued, "rewind"); err == nil {
		t.Fatal("active->queued should be rejected")
	}
	if err := m.Move(OpsStatusClosed, OpsStatusActive, "reopen"); err == nil {
		t.Fatal("closed->active should be rejected")
	}
	if last, ok := m.Last(); !ok || last.From != OpsStatusQueued {
		t.Fatal("expected a recorded transition")
	}
}

func TestRulesCatalog(t *testing.T) {
	rules := Rules()
	if len(rules) < 100 {
		t.Fatalf("expected at least 100 rules, got %d", len(rules))
	}
}
