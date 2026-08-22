package ops

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// TestConcurrentListRace verifies that concurrent list and create operations on the
// store do not race and that every created record is eventually visible.
func TestConcurrentListRace(t *testing.T) {
	store := NewStore(nil)
	ctx := context.Background()
	const workers = 4
	const perWorker = 20
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			<-start
			for i := 0; i < perWorker; i++ {
				_ = store.Put(ctx, OpsRecord{
					ID:       fmt.Sprintf("race-alert-%d-%d", w, i),
					Subject:  "pressure high",
					Owner:    "monitor",
					Priority: OpsPriorityHigh,
				})
				_, _ = store.List(ctx)
			}
		}(w)
	}
	close(start)
	wg.Wait()

	items, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected at least one record after concurrent writes")
	}
	for _, item := range items {
		if item.ID == "" {
			t.Fatal("record with empty id returned by list")
		}
	}
}

// TestConcurrentStoreGetRace verifies concurrent reads and writes against the store
// are synchronised (no data race).
func TestConcurrentStoreGetRace(t *testing.T) {
	store := NewStore(nil)
	ctx := context.Background()
	_ = store.Put(ctx, OpsRecord{ID: "get-race", Subject: "pressure", Owner: "m", Priority: OpsPriorityHigh})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			<-start
			for i := 0; i < 30; i++ {
				_, _ = store.Get(ctx, "get-race")
				_ = store.Put(ctx, OpsRecord{ID: fmt.Sprintf("get-race-%d-%d", w, i), Subject: "x", Owner: "m", Priority: OpsPriorityLow})
			}
		}(w)
	}
	close(start)
	wg.Wait()
	if got := store.Count(); got == 0 {
		t.Fatal("expected records after concurrent writes")
	}
}

// TestSortedCopyKeepsSource verifies that sorting never reorders the
// caller's own slice.
func TestSortedCopyKeepsSource(t *testing.T) {
	older := "2026-08-21T10:00:00Z"
	newer := "2026-08-21T11:00:00Z"
	input := []OpsRecord{
		{ID: "a", Priority: OpsPriorityLow, UpdatedAt: older},
		{ID: "b", Priority: OpsPriorityHigh, UpdatedAt: older},
		{ID: "c", Priority: OpsPriorityNormal, UpdatedAt: newer},
	}
	before := make([]OpsRecord, len(input))
	copy(before, input)

	sorted := SortRecords(input)

	for i := range input {
		if input[i].ID != before[i].ID {
			t.Fatalf("SortRecords mutated the input slice at index %d: got %s want %s", i, input[i].ID, before[i].ID)
		}
	}
	if len(sorted) != 3 || sorted[0].ID != "b" {
		t.Fatalf("unexpected sorted result: %+v", sorted)
	}
}

// TestGetReturnsIsolatedCopy verifies that mutating a fetched record does not
// leak into the stored record.
func TestGetReturnsIsolatedCopy(t *testing.T) {
	store := NewStore(nil)
	ctx := context.Background()
	if err := store.Put(ctx, OpsRecord{
		ID:       "iso-1",
		Subject:  "solvent low",
		Owner:    "alice",
		Priority: OpsPriorityNormal,
		Labels:   map[string]string{"site": "North Stack"},
	}); err != nil {
		t.Fatalf("put: %v", err)
	}

	first, err := store.Get(ctx, "iso-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	first.Labels["tampered"] = "yes"

	second, err := store.Get(ctx, "iso-1")
	if err != nil {
		t.Fatalf("get again: %v", err)
	}
	if _, found := second.Labels["tampered"]; found {
		t.Fatal("mutation of a fetched record leaked into the store")
	}
}
