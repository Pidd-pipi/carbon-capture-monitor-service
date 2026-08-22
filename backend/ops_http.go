package main

import (
	"encoding/json"
	"example.com/carbon-capture-monitor-service/ops"
	"net/http"
	"strings"
	"time"
)

func opsEnterpriseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("X-Operations-Domain", ops.DomainName)
		if strings.TrimSpace(r.Header.Get("X-Request-ID")) == "" {
			w.Header().Set("X-Operations-Request", "generated")
		} else {
			w.Header().Set("X-Operations-Request", "provided")
		}
		lw := &latencyWriter{ResponseWriter: w, started: start}
		next.ServeHTTP(lw, r)
	})
}

// latencyWriter injects the operations latency header at the moment the
// response is committed, so it is never lost to a deferred write.
type latencyWriter struct {
	http.ResponseWriter
	started time.Time
	status  int
}

func (w *latencyWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
		w.Header().Set("X-Operations-Latency-Ms", formatOpsInt(int(time.Since(w.started).Milliseconds())))
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *latencyWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
func formatOpsInt(value int) string {
	if value == 0 {
		return "0"
	}
	out := ""
	for value > 0 {
		out = string(rune('0'+value%10)) + out
		value /= 10
	}
	return out
}
func opsJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func opsAllowed(method string, allowed ...string) bool {
	for _, candidate := range allowed {
		if method == candidate {
			return true
		}
	}
	return false
}
func opsPathID(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.Trim(strings.TrimPrefix(path, prefix), "/")
}
func opsActorFromRequest(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("X-Operator"))
	if value == "" {
		return "web"
	}
	return value
}
func opsNoStore(w http.ResponseWriter)    { w.Header().Set("Cache-Control", "no-store") }
func opsRequestID(r *http.Request) string { return r.Header.Get("X-Request-ID") }
