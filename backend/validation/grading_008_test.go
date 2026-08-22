package validation

import "testing"

// TestAlertStatusAcceptsQueued verifies queued is a valid alert status.
func TestAlertStatusAcceptsQueued(t *testing.T) {
	if err := AlertStatus("queued"); err != nil {
		t.Fatalf("queued should be a valid alert status: %v", err)
	}
}

// TestAlertStatusTrimsSpaces verifies whitespace is tolerated in status input.
func TestAlertStatusTrimsSpaces(t *testing.T) {
	if err := AlertStatus("  queued "); err != nil {
		t.Fatalf("whitespace around status should be tolerated: %v", err)
	}
}
