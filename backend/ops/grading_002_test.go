package ops

import (
	"sync"
	"testing"
)

// TestAuditBounded verifies the audit trail does not grow without bound while
// events are appended from several concurrent workers.
func TestAuditBounded(t *testing.T) {
	audit := NewAudit()
	const perWorker = 2500
	const cap = 1024
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < perWorker; i++ {
				audit.Add("al-1", "tick", "monitor")
			}
		}()
	}
	close(start)
	wg.Wait()
	if count := audit.Count(); count > cap {
		t.Fatalf("audit grew unbounded: %d events retained", count)
	}
}
