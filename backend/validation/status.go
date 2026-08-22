package validation

import "fmt"

var allowed = map[string]bool{"online": true, "attention": true, "offline": true, "maintenance": true}

var allowedAlertStatuses = map[string]bool{"queued": true, "active": true, "paused": true, "closed": true}

func Status(value string) error {
	if !allowed[value] {
		return fmt.Errorf("status must be online, attention, offline, or maintenance")
	}
	return nil
}

// AlertStatus validates an alert work-order status.
func AlertStatus(value string) error {
	if !allowedAlertStatuses[value] {
		return fmt.Errorf("alert status must be queued, active, paused, or closed")
	}
	return nil
}
