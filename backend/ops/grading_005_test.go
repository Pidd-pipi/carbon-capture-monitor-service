package ops

import "testing"

// TestSetDefaultLabelsIsolated verifies the default labels are copied so the
// caller's map cannot be mutated through the stored reference.
func TestSetDefaultLabelsIsolated(t *testing.T) {
	original := map[string]string{"origin": "cfg"}
	SetDefaultLabels(original)
	original["tampered"] = "yes"

	got := DefaultLabels()
	if _, found := got["tampered"]; found {
		t.Fatal("mutating the source map leaked into the stored default labels")
	}
	if got["origin"] != "cfg" {
		t.Fatalf("expected origin=cfg, got %v", got)
	}
}

// TestNormalizeRecordWritableLabels verifies records without labels get a
// writable map even when no defaults are configured.
func TestNormalizeRecordWritableLabels(t *testing.T) {
	SetDefaultLabels(nil)
	record := NormalizeRecord(OpsRecord{ID: "al-nil", Subject: "no labels", Owner: "alice", Priority: OpsPriorityLow})
	record.Labels["source"] = "manual"
	if record.Labels["source"] != "manual" {
		t.Fatalf("normalized labels not writable: %v", record.Labels)
	}
}
