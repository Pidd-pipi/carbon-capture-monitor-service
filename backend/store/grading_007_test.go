package store

import "testing"

// TestUpdateReadingsComplete verifies every telemetry field is refreshed by
// the ingestion update.
func TestUpdateReadingsComplete(t *testing.T) {
	s := New()
	updated, err := s.UpdateReadings("CC-ALPHA", 80.5, 201.0, 55.0, "2026-08-22T12:00:00Z")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.CaptureRatePct != 80.5 || updated.PressureKPa != 201.0 {
		t.Fatalf("capture/pressure not updated: %+v", updated)
	}
	if updated.SolventLevel != 55.0 {
		t.Fatalf("solvent level not updated: %+v", updated)
	}
	if updated.UpdatedAt != "2026-08-22T12:00:00Z" {
		t.Fatalf("updated_at not refreshed: %+v", updated)
	}
}
