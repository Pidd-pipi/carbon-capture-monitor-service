package health

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuardHealthReportsStoreOK(t *testing.T) {
	dependencyOK = true
	rec := httptest.NewRecorder()
	Handler("carbon-capture-monitor-service").ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if !strings.Contains(rec.Body.String(), `"store":"ok"`) {
		t.Fatalf("health body: %s", rec.Body.String())
	}
}
