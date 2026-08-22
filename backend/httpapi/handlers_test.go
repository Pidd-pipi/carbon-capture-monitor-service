package httpapi

import (
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
	"example.com/carbon-capture-monitor-service/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testHandler() http.Handler {
	return NewHandler(store.New(), ops.NewService(nil), readings.NewStore(time.Hour, 100), web.FS)
}

func TestCaptureRoutes(t *testing.T) {
	handler := testHandler()
	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/capture-units", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), "CC-ALPHA") {
		t.Fatalf("collection: %d %s", get.Code, get.Body.String())
	}
	post := httptest.NewRecorder()
	handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-BETA","status":"online"}`)))
	if post.Code != http.StatusOK || !strings.Contains(post.Body.String(), `"status":"online"`) {
		t.Fatalf("update: %d %s", post.Code, post.Body.String())
	}
}

func TestCaptureRejectsUnknownUnit(t *testing.T) {
	handler := testHandler()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/capture-units/status", strings.NewReader(`{"id":"CC-NOPE","status":"online"}`)))
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}

func TestAlertLifecycle(t *testing.T) {
	handler := testHandler()
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"CC-ALPHA pressure high","owner":"alice","priority":"high","labels":{"site":"North Stack"}}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create alert: %d %s", create.Code, create.Body.String())
	}
	if !strings.Contains(create.Body.String(), `"subject":"CC-ALPHA pressure high"`) {
		t.Fatalf("create alert body: %s", create.Body.String())
	}
	list := httptest.NewRecorder()
	handler.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/api/alerts?page=1&page_size=10", nil))
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "CC-ALPHA pressure high") {
		t.Fatalf("list alerts: %d %s", list.Code, list.Body.String())
	}
}

func TestReadingsRoundTrip(t *testing.T) {
	handler := testHandler()
	post := httptest.NewRecorder()
	handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/readings", strings.NewReader(`{"unit_id":"CC-ALPHA","capture_rate_pct":88.1,"pressure_kpa":179.0,"solvent_level_pct":70}`)))
	if post.Code != http.StatusCreated {
		t.Fatalf("post reading: %d %s", post.Code, post.Body.String())
	}
	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/readings?unit=CC-ALPHA", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), "88.1") {
		t.Fatalf("list readings: %d %s", get.Code, get.Body.String())
	}
}
