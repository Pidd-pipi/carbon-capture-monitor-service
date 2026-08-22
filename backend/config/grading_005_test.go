package config

import "testing"

// TestDefaultProviderUsable verifies the default configuration path always
// yields a usable label map instead of a nil trap.
func TestDefaultProviderUsable(t *testing.T) {
	t.Setenv("ALERT_DEFAULT_LABELS", "")
	provider := NewLabelsProvider()
	if provider == nil {
		t.Fatal("default provider must not be nil")
	}
	labels := provider.Labels()
	if labels == nil {
		t.Fatal("default provider returned a nil label map")
	}
}
