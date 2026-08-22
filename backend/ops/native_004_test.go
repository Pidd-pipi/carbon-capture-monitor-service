package ops

import (
	"errors"
	"fmt"
	"testing"
)

func TestGuardErrorCodeClassifies(t *testing.T) {
	if got := ErrorCode(fmt.Errorf("wrap: %w", ErrOpsNotFound)); got != "not_found" {
		t.Fatalf("ErrorCode = %q, want not_found", got)
	}
	if got := ErrorCode(wrapOps("create", "x", errors.New("boom"))); got != "create" {
		t.Fatalf("ErrorCode = %q, want create", got)
	}
}
