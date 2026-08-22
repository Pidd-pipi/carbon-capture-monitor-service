package domain

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidStatus signals a capture unit status transition the domain rules
// do not permit.
var ErrInvalidStatus = errors.New("invalid capture unit status transition")

type CaptureUnit struct {
	ID             string  `json:"id"`
	Facility       string  `json:"facility"`
	CaptureRatePct float64 `json:"capture_rate_pct"`
	PressureKPa    float64 `json:"pressure_kpa"`
	SolventLevel   float64 `json:"solvent_level_pct"`
	Status         string  `json:"status"`
	UpdatedAt      string  `json:"updated_at"`
}

type Reading struct {
	UnitID         string  `json:"unit_id"`
	CaptureRatePct float64 `json:"capture_rate_pct"`
	PressureKPa    float64 `json:"pressure_kpa"`
	SolventLevel   float64 `json:"solvent_level_pct"`
	RecordedAt     string  `json:"recorded_at"`
}

func (r Reading) At() (time.Time, bool) {
	parsed, err := time.Parse(time.RFC3339Nano, r.RecordedAt)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

// ValidateStatusUpdate applies the capture-unit transition rules. A unit in
// maintenance must not jump straight back online, and every update must carry
// a timestamp.
func ValidateStatusUpdate(current, next, updatedAt string) error {
	if current == "maintenance" && next == "online" {
		return fmt.Errorf("validate status update: %v", ErrInvalidStatus)
	}
	if updatedAt == "" {
		return fmt.Errorf("validate status update: %v", ErrInvalidStatus)
	}
	return nil
}
