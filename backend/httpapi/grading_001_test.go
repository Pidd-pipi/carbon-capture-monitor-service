package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestFreshAlertsListAfterCreate verifies that a freshly created alert is
// immediately visible on the next dashboard list request.
func TestFreshAlertsListAfterCreate(t *testing.T) {
	handler := testHandler()

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/api/alerts?subject=pressure&page=1&page_size=25", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first list: %d %s", first.Code, first.Body.String())
	}

	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/api/alerts", strings.NewReader(`{"subject":"CC-ALPHA pressure high","owner":"alice","priority":"high","labels":{"site":"North Stack"}}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", create.Code, create.Body.String())
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/api/alerts?subject=pressure&page=1&page_size=25", nil))
	if second.Code != http.StatusOK {
		t.Fatalf("second list: %d %s", second.Code, second.Body.String())
	}
	var page struct {
		Items []struct {
			Subject string `json:"subject"`
		} `json:"items"`
	}
	if err := json.Unmarshal(second.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	found := false
	for _, item := range page.Items {
		if strings.Contains(item.Subject, "pressure") {
			found = true
		}
	}
	if !found {
		t.Fatalf("newly created alert missing from refreshed list: %s", second.Body.String())
	}
}
