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
		// The latency header must be present on the wire the moment headers are
		// committed. Setting it via defer after next.ServeHTTP returns is too
		// late: the inner handler has already called WriteHeader/Write, after
		// which header mutations are silently dropped. Stamp it at the first
		// commit instead so every response — including panic recoveries and
		// streaming writes — carries the header.
		recorder := &latencyRecorder{ResponseWriter: w, start: start}
		next.ServeHTTP(recorder, r)
		recorder.stampLatency()
	})
}

// latencyRecorder writes the X-Operations-Latency-Ms header exactly once, at
// the point the response is first committed. It is safe to call stampLatency
// after the handler returns as well, which makes the header show up even when
// the inner handler never writes anything.
type latencyRecorder struct {
	http.ResponseWriter
	start   time.Time
	stamped bool
}

func (l *latencyRecorder) stampLatency() {
	if l.stamped {
		return
	}
	l.stamped = true
	l.ResponseWriter.Header().Set("X-Operations-Latency-Ms", formatOpsInt(int(time.Since(l.start).Milliseconds())))
}

func (l *latencyRecorder) WriteHeader(code int) {
	l.stampLatency()
	l.ResponseWriter.WriteHeader(code)
}

func (l *latencyRecorder) Write(p []byte) (int, error) {
	l.stampLatency()
	return l.ResponseWriter.Write(p)
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
