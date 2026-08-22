package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/carbon-capture-monitor-service/config"
	"example.com/carbon-capture-monitor-service/httpapi"
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
	"example.com/carbon-capture-monitor-service/web"
)

type nilLabelsProvider struct{}

func (nilLabelsProvider) Labels() map[string]string { return nil }

// TestConfigureDefaultLabelsUsable verifies the wiring leaves a usable label
// map even when the provider reports no labels.
func TestConfigureDefaultLabelsUsable(t *testing.T) {
	configureDefaultLabels(nilLabelsProvider{})
	if ops.DefaultLabels() == nil {
		t.Fatal("configureDefaultLabels left a nil default label map")
	}
}

// TestDefaultConfigCreateAlertNoPanic drives the default-configuration alert
// creation path and asserts it does not crash.
func TestDefaultConfigCreateAlertNoPanic(t *testing.T) {
	t.Setenv("ALERT_DEFAULT_LABELS", "")
	configureDefaultLabels(config.NewLabelsProvider())

	alertSvc := ops.NewService(nil)
	readingStore := readings.NewStore(time.Hour, 100)
	handler := httpapi.NewHandler(store.New(), alertSvc, readingStore, web.FS)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"default config alert","owner":"alice","priority":"low"}`)))
	if rec.Code >= http.StatusInternalServerError {
		t.Fatalf("create alert on default config crashed the request: %d %s", rec.Code, rec.Body.String())
	}
}
