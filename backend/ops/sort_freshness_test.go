package ops

import "testing"

// TestSortRecordsKeepsInputOrder guards the sorting helper against mutating
// the caller's slice.
func TestSortRecordsKeepsInputOrder(t *testing.T) {
	input := []OpsRecord{
		{ID: "a", Priority: OpsPriorityLow, UpdatedAt: "2026-08-21T10:00:00Z"},
		{ID: "b", Priority: OpsPriorityHigh, UpdatedAt: "2026-08-21T10:00:00Z"},
	}
	before := make([]OpsRecord, len(input))
	copy(before, input)
	_ = SortRecords(input)
	for i := range input {
		if input[i].ID != before[i].ID {
			t.Fatalf("SortRecords mutated the input slice at index %d", i)
		}
	}
}
