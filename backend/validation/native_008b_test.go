package validation

import "testing"

func TestGuardAlertStatusTrimsSpaces(t *testing.T) {
	if err := AlertStatus("  queued "); err != nil {
		t.Fatalf("whitespace not tolerated: %v", err)
	}
}
