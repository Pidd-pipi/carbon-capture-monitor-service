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

func TestGuardDefaultConfigCreateAlertNoPanic(t *testing.T) {
	t.Setenv("ALERT_DEFAULT_LABELS", "")
	provider := config.NewLabelsProvider()
	labels := provider.Labels()
	ops.SetDefaultLabels(labels)
	alertSvc := ops.NewService(nil)
	handler := httpapi.NewHandler(store.New(), alertSvc, readings.NewStore(time.Hour, 100), web.FS)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"x","owner":"a","priority":"low"}`)))
	if rec.Code >= http.StatusInternalServerError {
		t.Fatalf("request crashed: %d %s", rec.Code, rec.Body.String())
	}
}
