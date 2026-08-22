package store

import (
	"errors"
	"example.com/carbon-capture-monitor-service/domain"
	"sync"
)

var ErrNotFound = errors.New("capture unit not found")

type Store struct {
	mu    sync.RWMutex
	items []domain.CaptureUnit
}

func New() *Store {
	return &Store{items: []domain.CaptureUnit{
		{ID: "CC-ALPHA", Facility: "North Stack", CaptureRatePct: 91.4, PressureKPa: 182.5, SolventLevel: 76, Status: "online", UpdatedAt: "2026-08-21T08:20:00Z"},
		{ID: "CC-BETA", Facility: "River Plant", CaptureRatePct: 87.8, PressureKPa: 176.2, SolventLevel: 68, Status: "attention", UpdatedAt: "2026-08-21T08:18:00Z"},
		{ID: "CC-GAMMA", Facility: "Harbor Terminal", CaptureRatePct: 93.1, PressureKPa: 171.8, SolventLevel: 81, Status: "online", UpdatedAt: "2026-08-21T08:22:00Z"},
	}}
}

func (s *Store) List() []domain.CaptureUnit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.CaptureUnit, len(s.items))
	copy(result, s.items)
	return result
}

func (s *Store) Get(id string) (domain.CaptureUnit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.items {
		if s.items[i].ID == id {
			return s.items[i], nil
		}
	}
	return domain.CaptureUnit{}, ErrNotFound
}

func (s *Store) UpdateStatus(id, status, updatedAt string) (domain.CaptureUnit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Status, s.items[i].UpdatedAt = status, updatedAt
			return s.items[i], nil
		}
	}
	return domain.CaptureUnit{}, ErrNotFound
}

// UpdateReadings refreshes the live telemetry values of a unit without
// changing its operational status. It is used by the ingestion path.
func (s *Store) UpdateReadings(id string, captureRatePct, pressureKPa, solventLevel float64, updatedAt string) (domain.CaptureUnit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].CaptureRatePct = captureRatePct
			s.items[i].PressureKPa = pressureKPa
			return s.items[i], nil
		}
	}
	return domain.CaptureUnit{}, ErrNotFound
}
