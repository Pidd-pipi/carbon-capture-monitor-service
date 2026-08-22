package ops

import "testing"

// TestRewindPausedToQueued verifies a paused work order can be returned
// to the queued state.
func TestRewindPausedToQueued(t *testing.T) {
	m := NewStateMachine()
	if err := m.Move(OpsStatusPaused, OpsStatusQueued, "rewind"); err != nil {
		t.Fatalf("paused->queued transition should be allowed: %v", err)
	}
}
