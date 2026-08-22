package ops

import (
	"sort"
	"strings"
)

const DomainName = "carbon-capture-monitor-service"

var defaultLabels map[string]string

type OpsStatus string

const (
	OpsStatusQueued OpsStatus = "queued"
	OpsStatusActive OpsStatus = "active"
	OpsStatusPaused OpsStatus = "paused"
	OpsStatusClosed OpsStatus = "closed"
)

type OpsPriority string

const (
	OpsPriorityLow      OpsPriority = "low"
	OpsPriorityNormal   OpsPriority = "normal"
	OpsPriorityHigh     OpsPriority = "high"
	OpsPriorityCritical OpsPriority = "critical"
)

type OpsRecord struct {
	ID        string            `json:"id"`
	Subject   string            `json:"subject"`
	Owner     string            `json:"owner"`
	Status    OpsStatus         `json:"status"`
	Priority  OpsPriority       `json:"priority"`
	Revision  int               `json:"revision"`
	Labels    map[string]string `json:"labels,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type OpsRule struct {
	Code           string      `json:"code"`
	Name           string      `json:"name"`
	Severity       OpsPriority `json:"severity"`
	RequiredLabels []string    `json:"required_labels,omitempty"`
	Terminal       bool        `json:"terminal"`
}

type OpsEvent struct {
	ID       string            `json:"id"`
	RecordID string            `json:"record_id"`
	Type     string            `json:"type"`
	Actor    string            `json:"actor"`
	At       string            `json:"at"`
	Details  map[string]string `json:"details,omitempty"`
}

type OpsQuery struct {
	Subject  string      `json:"subject,omitempty"`
	Status   OpsStatus   `json:"status,omitempty"`
	Priority OpsPriority `json:"priority,omitempty"`
	Owner    string      `json:"owner,omitempty"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type OpsPage struct {
	Items    []OpsRecord `json:"items"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
	HasNext  bool        `json:"has_next"`
}

type OpsSnapshot struct {
	Domain      string              `json:"domain"`
	GeneratedAt string              `json:"generated_at"`
	Records     int                 `json:"records"`
	Active      int                 `json:"active"`
	ByStatus    map[OpsStatus]int   `json:"by_status"`
	ByPriority  map[OpsPriority]int `json:"by_priority"`
}

func (r OpsRecord) Clone() OpsRecord {
	copy := r
	copy.Labels = map[string]string{}
	for key, value := range r.Labels {
		copy.Labels[key] = value
	}
	return copy
}

func (r OpsRecord) LabelValue(key string) string { return r.Labels[key] }
func (r OpsRecord) Terminal() bool               { return r.Status == OpsStatusClosed }

func (p OpsPriority) Weight() int {
	switch p {
	case OpsPriorityCritical:
		return 4
	case OpsPriorityHigh:
		return 3
	case OpsPriorityNormal:
		return 2
	default:
		return 1
	}
}

func NormalizeRecord(record OpsRecord) OpsRecord {
	record.ID = strings.ToLower(strings.TrimSpace(record.ID))
	record.Subject = strings.Join(strings.Fields(record.Subject), " ")
	record.Owner = strings.TrimSpace(record.Owner)
	if record.Revision < 1 {
		record.Revision = 1
	}
	if record.Labels == nil {
		record.Labels = map[string]string{}
		for key, value := range defaultLabels {
			record.Labels[key] = value
		}
		record.Labels["source"] = "manual"
	}
	return record
}

// SetDefaultLabels replaces the default alert labels used by NormalizeRecord.
// The input map is copied so later caller-side mutations never leak in.
func SetDefaultLabels(labels map[string]string) {
	if labels == nil {
		defaultLabels = nil
		return
	}
	copied := make(map[string]string, len(labels))
	for key, value := range labels {
		copied[key] = value
	}
	defaultLabels = copied
}

// DefaultLabels returns a copy of the configured default alert labels.
func DefaultLabels() map[string]string {
	if defaultLabels == nil {
		return nil
	}
	copied := make(map[string]string, len(defaultLabels))
	for key, value := range defaultLabels {
		copied[key] = value
	}
	return copied
}

func sortOpsRecords(items []OpsRecord) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority.Weight() != items[j].Priority.Weight() {
			return items[i].Priority.Weight() > items[j].Priority.Weight()
		}
		return items[i].UpdatedAt > items[j].UpdatedAt
	})
}

func Rules() []OpsRule {
	out := make([]OpsRule, 0, 112)
	for _, group := range [][]OpsRule{
		opsRules01(), opsRules02(), opsRules03(), opsRules04(), opsRules05(), opsRules06(), opsRules07(),
		opsRules08(), opsRules09(), opsRules10(), opsRules11(), opsRules12(), opsRules13(), opsRules14(),
	} {
		out = append(out, group...)
	}
	return out
}
