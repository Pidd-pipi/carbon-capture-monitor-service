package ops

import (
	"testing"
	"time"
)

// TestTimeoutContextEnforced verifies the timeout argument produces a real
// deadline instead of being silently dropped.
func TestTimeoutContextEnforced(t *testing.T) {
	ctx, cancel := opsContext(nil, 100*time.Millisecond)
	defer cancel()
	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("opsContext returned a context without a deadline")
	}
}
