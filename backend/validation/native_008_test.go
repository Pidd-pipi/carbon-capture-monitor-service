package validation

import "testing"

func TestGuardAlertStatusAcceptsQueued(t *testing.T) {
	if err := AlertStatus("queued"); err != nil {
		t.Fatalf("queued invalid: %v", err)
	}
}
