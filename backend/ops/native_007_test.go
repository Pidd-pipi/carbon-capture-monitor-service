package ops

import (
	"testing"
	"time"
)

func TestGuardTimeoutContextEnforced(t *testing.T) {
	ctx, cancel := opsContext(nil, 100*time.Millisecond)
	defer cancel()
	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("opsContext returned no deadline")
	}
}
