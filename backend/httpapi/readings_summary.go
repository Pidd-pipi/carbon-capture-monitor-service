package httpapi

import (
	"net/http"
	"strings"
	"time"
)

func (s *server) readingsSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	unitID := r.URL.Query().Get("unit")
	if strings.TrimSpace(unitID) == "" {
		writeError(w, http.StatusBadRequest, "unit is required")
		return
	}
	window := time.Duration(queryInt(r, "window_minutes", 60)) * time.Minute
	if window <= 0 {
		window = time.Minute
	}
	summary, err := s.readings.Summary(r.Context(), unitID, window, time.Now().UTC())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "readings summary failed")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
