package health

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ErrDependencyDown reports that a service dependency failed its health probe.
var ErrDependencyDown = errors.New("service dependency down")

// dependencyOK is the shared dependency flag maintained by the service wiring.
var dependencyOK = true

func dependencyHealth() error {
	if !dependencyOK {
		return ErrDependencyDown
	}
	return nil
}

func Handler(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		storeStatus := "ok"
		if err := dependencyHealth(); err != nil {
			if errors.Is(err, ErrDependencyDown) {
				storeStatus = "degraded"
			} else {
				storeStatus = "unknown"
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": service, "store": storeStatus})
	}
}
