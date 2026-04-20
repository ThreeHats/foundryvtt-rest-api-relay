package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics — exposed via /metrics endpoint.
// Labels are designed to be PII-safe: user IDs (integers), endpoint patterns,
// and HTTP status codes only. NEVER use email or API key as a label value.

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests by method, path pattern, and status",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency by method and path pattern",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	WSConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ws_connections_active",
			Help: "Currently active WebSocket connections",
		},
	)

	WSPendingRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ws_pending_requests",
			Help: "WebSocket requests awaiting a response from a Foundry module",
		},
	)

	RateLimitHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_hits_total",
			Help: "Total rate limit rejections by limiter name",
		},
		[]string{"limiter"},
	)

	HeadlessSessionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "headless_sessions_active",
			Help: "Currently active headless ChromeDP sessions",
		},
	)

	SSEConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sse_connections_active",
			Help: "Currently active Server-Sent Event subscriptions",
		},
	)

	AdminLoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "admin_logins_total",
			Help: "Total admin login attempts by outcome (success|failure|locked)",
		},
		[]string{"outcome"},
	)

	// AuthFailuresTotal counts API auth rejections by reason.
	// Reasons: bad_key, expired, disabled, account_issue, rate_limited
	AuthFailuresTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_failures_total",
			Help: "API authentication rejections by reason",
		},
		[]string{"reason"},
	)

	// WSRoundTripTimeouts counts WS requests that timed out waiting for Foundry module response.
	WSRoundTripTimeouts = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ws_roundtrip_timeouts_total",
			Help: "WebSocket requests that timed out waiting for Foundry module response",
		},
	)
)
