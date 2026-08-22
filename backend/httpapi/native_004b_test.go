package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuardIllegalTransition422(t *testing.T) {
	handler := testHandler()
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"illegal","owner":"a","priority":"normal","labels":{"site":"North Stack"}}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create: %d", create.Code)
	}
	var out struct {
		ID string `json:"id"`
	}
	_ = json.NewDecoder(create.Body).Decode(&out)
	act := httptest.NewRecorder()
	handler.ServeHTTP(act, httptest.NewRequest(http.MethodPost, "/api/alerts/"+out.ID+"/transition", strings.NewReader(`{"expected_revision":1,"target_status":"active"}`)))
	if act.Code != http.StatusOK {
		t.Fatalf("activate: %d", act.Code)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts/"+out.ID+"/transition", strings.NewReader(`{"expected_revision":2,"target_status":"queued"}`)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}
