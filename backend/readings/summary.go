package readings

import (
	"context"
	"fmt"
	"time"
)

// Summary is the aggregate view of a unit's readings over a window.
type Summary struct {
	UnitID    string `json:"unit_id"`
	Window    string `json:"window"`
	Stats     Stats  `json:"stats"`
	Generated string `json:"generated_at"`
}

// Summary computes aggregate statistics for one unit over the trailing window.
// If the unit has no readings inside the window it returns a zero-value Stats
// with Count 0 rather than an error.
//
// The caller's ctx is honored so a per-request timeout or cancellation aborts
// the work in flight instead of leaving orphaned work running on a detached
// context. Each request computes its own deadline independently.
func (s *Store) Summary(ctx context.Context, unitID string, window time.Duration, now time.Time) (Summary, error) {
	if err := ctx.Err(); err != nil {
		return Summary{}, fmt.Errorf("summary: %w", err)
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	from := now.Add(-window)
	items, err := s.List(ctx, unitID, from, now, 0)
	if err != nil {
		return Summary{}, fmt.Errorf("summary: %w", err)
	}
	return Summary{
		UnitID:    unitID,
		Window:    window.String(),
		Stats:     Aggregate(items),
		Generated: now.UTC().Format(time.RFC3339Nano),
	}, nil
}
