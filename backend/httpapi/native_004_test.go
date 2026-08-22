package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuardTransitionUnknownKeeps404(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts/al-nope/transition", strings.NewReader(`{"expected_revision":1,"target_status":"active"}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGuardTransitionStaleRevisionKeeps409(t *testing.T) {
	handler := testHandler()
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"conflict","owner":"a","priority":"normal","labels":{"site":"North Stack"}}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create: %d", create.Code)
	}
	var out struct {
		ID string `json:"id"`
	}
	_ = json.NewDecoder(create.Body).Decode(&out)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts/"+out.ID+"/transition", strings.NewReader(`{"expected_revision":99,"target_status":"active"}`)))
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}
