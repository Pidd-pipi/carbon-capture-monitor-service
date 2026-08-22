package main

import (
	"testing"

	"example.com/carbon-capture-monitor-service/ops"
)

type guardNilLabelsProvider struct{}

func (guardNilLabelsProvider) Labels() map[string]string { return nil }

func TestGuardConfigureDefaultLabels(t *testing.T) {
	configureDefaultLabels(guardNilLabelsProvider{})
	if ops.DefaultLabels() == nil {
		t.Fatal("configureDefaultLabels left a nil default label map")
	}
}
