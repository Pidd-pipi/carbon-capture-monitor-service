package readings

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Reading is one telemetry sample captured from a capture unit.
type Reading struct {
	UnitID         string
	CaptureRatePct float64
	PressureKPa    float64
	SolventLevel   float64
	RecordedAt     time.Time
}

// Store keeps a bounded in-memory history of readings per capture unit.
// Older samples beyond the retention window are pruned as new ones arrive.
type Store struct {
	mu         sync.RWMutex
	perUnit    map[string][]Reading
	window     time.Duration
	maxPerUnit int
}

// NewStore creates a readings store that retains at most maxPerUnit samples
// per unit and prunes anything older than the retention window.
func NewStore(window time.Duration, maxPerUnit int) *Store {
	return &Store{
		perUnit:    map[string][]Reading{},
		window:     window,
		maxPerUnit: maxPerUnit,
	}
}

// Append records a single reading for a unit, pruning the oldest samples
// that fall outside the retention window.
func (s *Store) Append(ctx context.Context, unitID string, r Reading) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r.UnitID = unitID
	if r.RecordedAt.IsZero() {
		r.RecordedAt = time.Now().UTC()
	}
	items := s.perUnit[unitID]
	items = append(items, r)
	cutoff := r.RecordedAt.Add(-s.window)
	keep := items
	for i, item := range items {
		if !item.RecordedAt.Before(cutoff) {
			keep = items[i:]
			break
		}
	}
	if len(keep) > s.maxPerUnit {
		keep = keep[len(keep)-s.maxPerUnit:]
	}
	s.perUnit[unitID] = keep
	return nil
}

// AppendBatch records several readings at once for a unit.
func (s *Store) AppendBatch(ctx context.Context, unitID string, samples []Reading) error {
	for _, sample := range samples {
		if err := s.Append(ctx, unitID, sample); err != nil {
			return err
		}
	}
	return nil
}

// List returns the readings of a unit within [from, to], newest first.
func (s *Store) List(ctx context.Context, unitID string, from, to time.Time, limit int) ([]Reading, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.perUnit[unitID]
	out := make([]Reading, 0, len(items))
	for _, item := range items {
		if (from.IsZero() || !item.RecordedAt.Before(from)) && (to.IsZero() || !item.RecordedAt.After(to)) {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RecordedAt.After(out[j].RecordedAt) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// Latest returns the most recent reading for a unit.
func (s *Store) Latest(ctx context.Context, unitID string) (Reading, bool) {
	select {
	case <-ctx.Done():
		return Reading{}, false
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.perUnit[unitID]
	if len(items) == 0 {
		return Reading{}, false
	}
	latest := items[0]
	for _, item := range items {
		if item.RecordedAt.After(latest.RecordedAt) {
			latest = item
		}
	}
	return latest, true
}

// Units returns the identifiers of every unit that has at least one reading.
func (s *Store) Units() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.perUnit))
	for id := range s.perUnit {
		if len(s.perUnit[id]) > 0 {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// Count returns how many readings are currently retained for a unit.
func (s *Store) Count(unitID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.perUnit[unitID])
}

// Prune removes readings older than the retention window for every unit.
// It is used by the background maintenance loop.
func (s *Store) Prune(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	pruned := 0
	cutoff := now.Add(-s.window)
	for id, items := range s.perUnit {
		keep := items
		for i, item := range items {
			if !item.RecordedAt.Before(cutoff) {
				keep = items[i:]
				break
			}
		}
		if len(keep) < len(items) {
			pruned += len(items) - len(keep)
		}
		if len(keep) > s.maxPerUnit {
			keep = keep[len(keep)-s.maxPerUnit:]
		}
		s.perUnit[id] = keep
	}
	return pruned
}
