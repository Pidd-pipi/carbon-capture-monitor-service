package ops

import "testing"

func TestGuardNormalizeRecordWritableLabels(t *testing.T) {
	SetDefaultLabels(nil)
	record := NormalizeRecord(OpsRecord{ID: "al-1", Subject: "x", Owner: "a", Priority: OpsPriorityLow})
	record.Labels["source"] = "manual"
	if record.Labels["source"] != "manual" {
		t.Fatal("labels not writable")
	}
}

func TestGuardSetDefaultLabelsIsolated(t *testing.T) {
	original := map[string]string{"origin": "cfg"}
	SetDefaultLabels(original)
	original["tampered"] = "yes"
	if _, found := DefaultLabels()["tampered"]; found {
		t.Fatal("caller mutation leaked into defaults")
	}
}
