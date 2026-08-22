package ops

import (
	"errors"
	"fmt"
	"testing"
)

// TestOpsErrorCodeMap verifies wrapped sentinel errors keep their
// semantic code instead of falling back to "internal".
func TestOpsErrorCodeMap(t *testing.T) {
	cases := []struct {
		err  error
		want string
	}{
		{fmt.Errorf("transition al-1: %w", ErrOpsNotFound), "not_found"},
		{fmt.Errorf("transition al-1: %w", ErrOpsConflict), "conflict"},
		{fmt.Errorf("transition al-1: %w", ErrOpsTransition), "transition"},
		{fmt.Errorf("transition al-1: %w", ErrOpsPolicy), "policy"},
		{wrapOps("create", "store.put", errors.New("boom")), "create"},
	}
	for _, c := range cases {
		if got := ErrorCode(c.err); got != c.want {
			t.Fatalf("ErrorCode(%v) = %q, want %q", c.err, got, c.want)
		}
	}
}
