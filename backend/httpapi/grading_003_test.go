package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestReadingsLimitHonored verifies the list endpoint returns at most the
// requested number of readings.
func TestReadingsLimitHonored(t *testing.T) {
	handler := testHandler()
	post := httptest.NewRecorder()
	handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/readings", strings.NewReader(
		`{"unit_id":"CC-ALPHA","readings":[
			{"capture_rate_pct":81,"pressure_kpa":171,"solvent_level_pct":71,"recorded_at":"2026-08-22T10:00:00Z"},
			{"capture_rate_pct":82,"pressure_kpa":172,"solvent_level_pct":72,"recorded_at":"2026-08-22T10:01:00Z"},
			{"capture_rate_pct":83,"pressure_kpa":173,"solvent_level_pct":73,"recorded_at":"2026-08-22T10:02:00Z"},
			{"capture_rate_pct":84,"pressure_kpa":174,"solvent_level_pct":74,"recorded_at":"2026-08-22T10:03:00Z"},
			{"capture_rate_pct":85,"pressure_kpa":175,"solvent_level_pct":75,"recorded_at":"2026-08-22T10:04:00Z"}
		]}`)))
	if post.Code != http.StatusCreated {
		t.Fatalf("post: %d %s", post.Code, post.Body.String())
	}

	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/readings?unit=CC-ALPHA&limit=2", nil))
	if get.Code != http.StatusOK {
		t.Fatalf("get: %d %s", get.Code, get.Body.String())
	}
	var body struct {
		Items []struct {
			CaptureRatePct float64 `json:"capture_rate_pct"`
		} `json:"items"`
	}
	if err := json.Unmarshal(get.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Items) != 2 {
		t.Fatalf("expected exactly 2 readings with limit=2, got %d: %s", len(body.Items), get.Body.String())
	}
}
