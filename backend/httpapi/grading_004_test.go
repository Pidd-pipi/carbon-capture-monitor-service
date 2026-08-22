package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func decodeBody(rec *httptest.ResponseRecorder, out any) error {
	return json.Unmarshal(rec.Body.Bytes(), out)
}

func createAlert(t *testing.T, handler http.Handler, subject string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(
		`{"subject":"`+subject+`","owner":"alice","priority":"normal","labels":{"site":"North Stack"}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var out struct {
		ID string `json:"id"`
	}
	_ = decodeBody(rec, &out)
	return out.ID
}

func TestTransitionUnknownAlert(t *testing.T) {
	handler := testHandler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts/al-does-not-exist/transition", strings.NewReader(`{"expected_revision":1,"target_status":"active"}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown alert, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConflictStaleRevision(t *testing.T) {
	handler := testHandler()
	id := createAlert(t, handler, "conflict alert")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts/"+id+"/transition", strings.NewReader(`{"expected_revision":99,"target_status":"active"}`)))
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 for stale revision, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestIllegalTransition(t *testing.T) {
	handler := testHandler()
	id := createAlert(t, handler, "illegal transition alert")
	act := httptest.NewRecorder()
	handler.ServeHTTP(act, httptest.NewRequest(http.MethodPost, "/api/alerts/"+id+"/transition", strings.NewReader(`{"expected_revision":1,"target_status":"active"}`)))
	if act.Code != http.StatusOK {
		t.Fatalf("activate: %d %s", act.Code, act.Body.String())
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/alerts/"+id+"/transition", strings.NewReader(`{"expected_revision":2,"target_status":"queued"}`)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for illegal transition, got %d: %s", rec.Code, rec.Body.String())
	}
}
