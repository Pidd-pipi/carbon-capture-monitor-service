package ops

import "testing"

func TestGuardRewindPausedToQueued(t *testing.T) {
	m := NewStateMachine()
	if err := m.Move(OpsStatusPaused, OpsStatusQueued, "rewind"); err != nil {
		t.Fatalf("paused->queued should be allowed: %v", err)
	}
}
