package middleware

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

var accessLogFieldsPool = sync.Pool{
	New: func() interface{} { return &AccessLogFields{} },
}

// AccessLogFields is a mutable holder injected into the request context by
// MetricsMiddleware BEFORE calling next, so AuthMiddleware (which runs
// downstream) can backfill the userId and keyPrefix fields. Using a pointer
// means auth can mutate the struct without needing to propagate a new context
// back upstream.
type AccessLogFields struct {
	UserID    int64
	KeyPrefix string
}

type accessLogFieldsKeyType struct{}

// AccessLogFieldsKey is exported so AuthMiddleware (same package) can use it.
var AccessLogFieldsKey = accessLogFieldsKeyType{}

// GetAccessLogFields retrieves the access-log holder from a request context.
// Returns nil if MetricsMiddleware is not in the chain.
func GetAccessLogFields(r *http.Request) *AccessLogFields {
	f, _ := r.Context().Value(AccessLogFieldsKey).(*AccessLogFields)
	return f
}

// statusRecorder wraps http.ResponseWriter to capture the status code and the
// first 512 bytes of error response bodies (for structured access logging).
//
// It must also implement http.Hijacker so that WebSocket upgrades work when
// this middleware is in the chain. Embedding http.ResponseWriter does NOT
// automatically forward Hijack() — the embedded interface only exposes the
// methods declared on http.ResponseWriter (Header, Write, WriteHeader). We
// need to explicitly implement Hijack() and delegate to the underlying writer
// if it implements http.Hijacker.
type statusRecorder struct {
	http.ResponseWriter
	status    int
	errorBody []byte // up to 512 bytes of non-2xx response body for log enrichment
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Write captures the first 512 bytes of error responses for log enrichment,
// then delegates to the underlying writer.
func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status >= 400 && len(r.errorBody) < 512 {
		remaining := 512 - len(r.errorBody)
		if len(b) <= remaining {
			r.errorBody = append(r.errorBody, b...)
		} else {
			r.errorBody = append(r.errorBody, b[:remaining]...)
		}
	}
	return r.ResponseWriter.Write(b)
}

// Hijack implements http.Hijacker so WebSocket upgrades work through this middleware.
// We delegate to the underlying ResponseWriter if it supports hijacking.
func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
	}
	return hijacker.Hijack()
}

// Flush implements http.Flusher so SSE streaming works through this middleware.
func (r *statusRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// MetricsMiddleware records Prometheus metrics and emits a structured access
// log line for each HTTP request. Uses the chi route pattern (not the actual
// URL) as the label/log field so cardinality stays bounded.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}

		// Inject a mutable holder BEFORE calling next so downstream middleware
		// (AuthMiddleware) can backfill userId and keyPrefix into it.
		logFields := accessLogFieldsPool.Get().(*AccessLogFields)
		*logFields = AccessLogFields{} // zero before reuse
		defer accessLogFieldsPool.Put(logFields)
		r = r.WithContext(context.WithValue(r.Context(), AccessLogFieldsKey, logFields))

		next.ServeHTTP(rec, r)

		path := chi.RouteContext(r.Context()).RoutePattern()
		if path == "" {
			path = "unknown"
		}
		// Skip successful WebSocket upgrades — gorilla hijacks the connection and
		// writes the 101 response directly to the wire without calling WriteHeader,
		// so status stays at 200 with userId=0. The real events are the
		// "Client WS connected/disconnected" log lines emitted by the WS handlers.
		// Rejected WS connections (401, 403, etc.) still show in the access log.
		if rec.status == 200 && strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			return
		}

		// Skip /metrics itself to avoid noise
		if strings.HasPrefix(path, "/metrics") {
			return
		}

		duration := time.Since(start)
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration.Seconds())
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rec.status)).Inc()
		metrics.Global.Record(logFields.UserID, path, rec.status)

		// Emit structured access log. WebSocket upgrades (101) have their own
		// connect/disconnect logs — skip them here to avoid duplication.
		if rec.status == 101 {
			return
		}

		// Extract error message from response body for 4xx/5xx log lines.
		var errMsg string
		if rec.status >= 400 && len(rec.errorBody) > 0 {
			var errJSON map[string]interface{}
			if jsonErr := json.Unmarshal(rec.errorBody, &errJSON); jsonErr == nil {
				if msg, ok := errJSON["error"].(string); ok {
					errMsg = msg
				}
			}
		}

		// Dashboard GET requests (session auth → no API key prefix) polling endpoints
		// like /auth/connection-tokens every 10 s are pure noise at Info level.
		// API key callers (keyPrefix set) stay at Info so their reads remain visible.
		isDashboardRead := r.Method == "GET" && rec.status < 400 && logFields.KeyPrefix == ""

		logEvent := log.Info()
		if rec.status >= 500 {
			logEvent = log.Error()
		} else if rec.status >= 400 {
			logEvent = log.Warn()
		} else if isDashboardRead {
			logEvent = log.Debug()
		}
		logEvent = logEvent.
			Str("method", r.Method).
			Str("path", path).
			Int("status", rec.status).
			Int64("latencyMs", duration.Milliseconds()).
			Int64("userId", logFields.UserID).
			Str("keyPrefix", logFields.KeyPrefix)
		if errMsg != "" {
			logEvent = logEvent.Str("error", errMsg)
		}
		logEvent.Msg("request")
	})
}
