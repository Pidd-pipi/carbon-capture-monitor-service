package ops

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

var alertSequence uint64

func newAlertID() string { return fmt.Sprintf("AL-%06d", atomic.AddUint64(&alertSequence, 1)) }

type OpsPolicy struct {
	RequireOwner  bool
	RequiredLabel string
	MaxActive     int
}
type OpsService struct {
	store  *OpsStore
	audit  *OpsAudit
	state  *OpsStateMachine
	policy OpsPolicy
	clock  OpsClock
}

func NewService(seed []OpsRecord) *OpsService {
	return &OpsService{store: NewStore(seed), audit: NewAudit(), state: NewStateMachine(), policy: OpsPolicy{RequireOwner: true, RequiredLabel: "site", MaxActive: 1000}, clock: NewClock()}
}
func (p OpsPolicy) Check(record OpsRecord) error {
	if p.RequireOwner && strings.TrimSpace(record.Owner) == "" {
		return fmt.Errorf("%w: owner required", ErrOpsPolicy)
	}
	if p.RequiredLabel != "" && record.LabelValue(p.RequiredLabel) == "" {
		return fmt.Errorf("%w: %s label required", ErrOpsPolicy, p.RequiredLabel)
	}
	if record.Priority == "" {
		return fmt.Errorf("%w: priority required", ErrOpsPolicy)
	}
	return nil
}
func (s *OpsService) Create(ctx context.Context, record OpsRecord) (OpsRecord, error) {
	if record.ID == "" {
		record.ID = newAlertID()
	}
	record = NormalizeRecord(record)
	if record.Status == "" {
		record.Status = OpsStatusQueued
	}
	if err := s.policy.Check(record); err != nil {
		return OpsRecord{}, err
	}
	record.CreatedAt = s.clock.Stamp()
	record.UpdatedAt = record.CreatedAt
	if err := s.store.Put(ctx, record); err != nil {
		return OpsRecord{}, wrapOps("create", "store.put", err)
	}
	s.audit.Add(record.ID, "created", record.Owner)
	return record, nil
}
func (s *OpsService) Get(ctx context.Context, id string) (OpsRecord, error) {
	return s.store.Get(ctx, id)
}
func (s *OpsService) Search(ctx context.Context, q OpsQuery) (OpsPage, error) {
	items, err := s.store.List(ctx)
	if err != nil {
		return OpsPage{}, err
	}
	filtered := make([]OpsRecord, 0, len(items))
	for _, item := range items {
		if opsMatch(item, q) {
			filtered = append(filtered, item)
		}
	}
	q = opsQueryDefaults(q)
	sorted := SortRecords(filtered)
	start, end := opsBounds(len(sorted), q.Page, q.PageSize)
	return OpsPage{Items: sorted[start:end], Page: q.Page, PageSize: q.PageSize, Total: len(sorted), HasNext: end < len(sorted)}, nil
}
func (s *OpsService) Transition(ctx context.Context, id string, expected int, target OpsStatus, actor string) (OpsRecord, error) {
	ctx, cancel := opsContext(ctx, 3*time.Second)
	defer cancel()
	record, err := s.store.Get(ctx, id)
	if err != nil {
		return OpsRecord{}, fmt.Errorf("transition %s: %v", id, err)
	}
	if expected > 0 && expected != record.Revision {
		return OpsRecord{}, fmt.Errorf("transition %s: %v", id, ErrOpsConflict)
	}
	if err := s.state.Move(record.Status, target, "operator update"); err != nil {
		return OpsRecord{}, fmt.Errorf("transition %s: %v", id, err)
	}
	record.Status = target
	updated, err := s.store.Update(ctx, record, expected)
	if err != nil {
		return OpsRecord{}, fmt.Errorf("transition %s: %v", id, err)
	}
	s.audit.Add(updated.ID, "status_changed", actor)
	return updated, nil
}
func (s *OpsService) Audit(id string) []OpsEvent { return s.audit.For(id) }
func (s *OpsService) Snapshot() OpsSnapshot {
	items, _ := s.store.List(context.Background())
	out := OpsSnapshot{Domain: DomainName, GeneratedAt: s.clock.Stamp(), ByStatus: map[OpsStatus]int{}, ByPriority: map[OpsPriority]int{}}
	for _, i := range items {
		out.Records++
		out.ByStatus[i.Status]++
		out.ByPriority[i.Priority]++
		if i.Status == OpsStatusActive {
			out.Active++
		}
	}
	return out
}
func (s *OpsService) Domain() string { return DomainName }
func (s *OpsService) Count() int     { return s.store.Count() }
func timeNowOps() string             { return time.Now().UTC().Format(time.RFC3339Nano) }
