package ops

import "testing"

// TestAuditBounded verifies the audit trail does not grow without bound.
func TestGuardAuditBounded(t *testing.T) {
	audit := NewAudit()
	for i := 0; i < 5000; i++ {
		audit.Add("al-1", "tick", "monitor")
	}
	if count := audit.Count(); count > 1024 {
		t.Fatalf("audit grew unbounded: %d", count)
	}
}
