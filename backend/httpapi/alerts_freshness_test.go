package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAlertsListShowsFreshCreate guards the dashboard refresh behaviour: a
// newly created alert must be visible on the very next list request.
func TestAlertsListShowsFreshCreate(t *testing.T) {
	handler := testHandler()
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/api/alerts?page=1&page_size=50", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first list: %d %s", first.Code, first.Body.String())
	}
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"fresh alert","owner":"alice","priority":"normal","labels":{"site":"North Stack"}}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", create.Code, create.Body.String())
	}
	list := httptest.NewRecorder()
	handler.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/api/alerts?page=1&page_size=50", nil))
	if list.Code != http.StatusOK {
		t.Fatalf("list: %d %s", list.Code, list.Body.String())
	}
	if !strings.Contains(list.Body.String(), "fresh alert") {
		t.Fatalf("freshly created alert missing from refreshed list: %s", list.Body.String())
	}
}
