package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func Port() string {
	value := os.Getenv("PORT")
	port, err := strconv.Atoi(value)
	if err == nil && port > 0 && port < 65536 {
		return value
	}
	return "8080"
}

func MonitorInterval() time.Duration {
	value := os.Getenv("MONITOR_INTERVAL_SECONDS")
	seconds, err := strconv.Atoi(value)
	if err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	return 30 * time.Second
}

func ReadingRetention() time.Duration {
	value := os.Getenv("READING_RETENTION_HOURS")
	hours, err := strconv.Atoi(value)
	if err == nil && hours > 0 {
		return time.Duration(hours) * time.Hour
	}
	return 24 * time.Hour
}

func MaxReadingsPerUnit() int {
	value := os.Getenv("MAX_READINGS_PER_UNIT")
	limit, err := strconv.Atoi(value)
	if err == nil && limit > 0 {
		return limit
	}
	return 2000
}

// AlertRulesEnabled toggles whether the monitor consults the ops rule set
// when evaluating a unit's readings. Disabled in tests to keep fixtures small.
func AlertRulesEnabled() bool {
	value := os.Getenv("ALERT_RULES_ENABLED")
	return value != "false"
}

// LabelsProvider supplies the default labels attached to manually created
// alert records.
type LabelsProvider interface {
	Labels() map[string]string
}

type envLabelsProvider struct {
	labels map[string]string
}

func (p *envLabelsProvider) Labels() map[string]string {
	if p == nil {
		return nil
	}
	return p.labels
}

// NewLabelsProvider builds the default alert labels from ALERT_DEFAULT_LABELS
// (a comma separated key=value list). When the variable is unset or empty the
// provider still returns a usable (non-nil) empty label map so the service
// runs without any environment configuration.
func NewLabelsProvider() LabelsProvider {
	out := map[string]string{}
	raw := strings.TrimSpace(os.Getenv("ALERT_DEFAULT_LABELS"))
	if raw != "" {
		for _, part := range strings.Split(raw, ",") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	}
	return &envLabelsProvider{labels: out}
}
