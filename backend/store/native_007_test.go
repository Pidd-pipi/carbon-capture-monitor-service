package store

import "testing"

func TestGuardUpdateReadingsComplete(t *testing.T) {
	s := New()
	updated, err := s.UpdateReadings("CC-ALPHA", 80.5, 201.0, 55.0, "2026-08-22T12:00:00Z")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.SolventLevel != 55.0 || updated.UpdatedAt != "2026-08-22T12:00:00Z" {
		t.Fatalf("fields not fully updated: %+v", updated)
	}
}
