package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example.com/carbon-capture-monitor-service/domain"
	"example.com/carbon-capture-monitor-service/readings"
)

func (s *server) readingsPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		UnitID         string           `json:"unit_id"`
		CaptureRatePct float64          `json:"capture_rate_pct"`
		PressureKPa    float64          `json:"pressure_kpa"`
		SolventLevel   float64          `json:"solvent_level_pct"`
		Readings       []domain.Reading `json:"readings"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid readings payload")
		return
	}
	if strings.TrimSpace(body.UnitID) == "" {
		writeError(w, http.StatusBadRequest, "unit_id is required")
		return
	}
	if len(body.Readings) == 0 {
		sample := readings.Reading{
			UnitID:         body.UnitID,
			CaptureRatePct: body.CaptureRatePct,
			PressureKPa:    body.PressureKPa,
			SolventLevel:   body.SolventLevel,
			RecordedAt:     time.Now().UTC(),
		}
		if err := s.readings.Append(r.Context(), body.UnitID, sample); err != nil {
			writeError(w, http.StatusInternalServerError, "readings append failed")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"unit_id": body.UnitID, "count": 1})
		return
	}
	samples := make([]readings.Reading, 0, len(body.Readings))
	for _, item := range body.Readings {
		recordedAt, _ := item.At()
		samples = append(samples, readings.Reading{
			UnitID:         body.UnitID,
			CaptureRatePct: item.CaptureRatePct,
			PressureKPa:    item.PressureKPa,
			SolventLevel:   item.SolventLevel,
			RecordedAt:     recordedAt,
		})
	}
	if err := s.readings.AppendBatch(r.Context(), body.UnitID, samples); err != nil {
		writeError(w, http.StatusInternalServerError, "readings append failed")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"unit_id": body.UnitID, "count": len(samples)})
}

func (s *server) readingsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	unitID := r.URL.Query().Get("unit")
	if strings.TrimSpace(unitID) == "" {
		writeError(w, http.StatusBadRequest, "unit is required")
		return
	}
	from := parseTime(r.URL.Query().Get("from"))
	to := parseTime(r.URL.Query().Get("to"))
	limit := queryInt(r, "limit", 100)
	items, err := s.readings.List(r.Context(), unitID, from, to, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "readings query failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"unit_id": unitID, "items": items})
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed.UTC()
}
