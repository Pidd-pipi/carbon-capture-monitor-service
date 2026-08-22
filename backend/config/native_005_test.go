package config

import "testing"

func TestGuardDefaultProviderUsable(t *testing.T) {
	t.Setenv("ALERT_DEFAULT_LABELS", "")
	if labels := NewLabelsProvider().Labels(); labels == nil {
		t.Fatal("default provider returned nil labels")
	}
}
