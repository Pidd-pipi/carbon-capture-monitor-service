package ops

import (
	"context"
	"sync"
)

// OpsStore is an in-memory, revision-tracked store for operation records.
// The backing slice is append-only under a write lock; every read path hands
// out independent copies so callers can never mutate stored state.
type OpsStore struct {
	mu    sync.RWMutex
	items []OpsRecord
}

func NewStore(seed []OpsRecord) *OpsStore {
	s := &OpsStore{items: []OpsRecord{}}
	for _, item := range seed {
		item = NormalizeRecord(item)
		s.items = append(s.items, item)
	}
	return s
}

func (s *OpsStore) Get(ctx context.Context, id string) (OpsRecord, error) {
	select {
	case <-ctx.Done():
		return OpsRecord{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.items {
		if item.ID == id {
			return item.Clone(), nil
		}
	}
	return OpsRecord{}, ErrOpsNotFound
}

func (s *OpsStore) List(ctx context.Context) ([]OpsRecord, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]OpsRecord, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item.Clone())
	}
	return out, nil
}

func (s *OpsStore) Put(ctx context.Context, item OpsRecord) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	item = NormalizeRecord(item)
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.items {
		if existing.ID == item.ID {
			return ErrOpsConflict
		}
	}
	s.items = append(s.items, item)
	return nil
}

func (s *OpsStore) Update(ctx context.Context, item OpsRecord, expected int) (OpsRecord, error) {
	select {
	case <-ctx.Done():
		return OpsRecord{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID != item.ID {
			continue
		}
		if expected > 0 && s.items[i].Revision != expected {
			return OpsRecord{}, ErrOpsConflict
		}
		item.Revision = s.items[i].Revision + 1
		item.UpdatedAt = timeNowOps()
		s.items[i] = item.Clone()
		return s.items[i].Clone(), nil
	}
	return OpsRecord{}, ErrOpsNotFound
}

func (s *OpsStore) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}
	return ErrOpsNotFound
}

func (s *OpsStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
