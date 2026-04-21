package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/rs/zerolog/log"
)

// forwardClient is a shared HTTP client for inter-instance forwarding.
// Using a shared client enables TCP connection pooling to Fly.io instances.
var forwardClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	},
}

// sseForwardClient is used for SSE proxy forwarding with no overall timeout
// (SSE streams must remain open as long as the subscriber is connected).
var sseForwardClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       0, // no idle timeout for long-lived streams
		ResponseHeaderTimeout: 15 * time.Second,
	},
	Timeout: 0, // no deadline — bound to request context
}

// RequestForwarder forwards requests to the correct Fly.io instance when a client
// is connected to a different instance than the one receiving the request.
type RequestForwarder struct {
	redis      *config.RedisClient
	instanceID string
	appName    string
	internalPort string
}

// NewRequestForwarder creates a new request forwarder middleware.
func NewRequestForwarder(redis *config.RedisClient, cfg *config.Config) *RequestForwarder {
	return &RequestForwarder{
		redis:        redis,
		instanceID:   cfg.InstanceID(),
		appName:      cfg.AppName,
		internalPort: cfg.FlyInternalPort,
	}
}

// Middleware returns the HTTP middleware that checks Redis and forwards if needed.
func (f *RequestForwarder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if no Redis or local mode
		if f.redis == nil || !f.redis.IsConnected() || f.appName == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Skip health checks and static assets
		path := r.URL.Path
		if path == "/health" || strings.HasPrefix(path, "/static") || path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("X-API-Key")
		}
		if apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		var targetInstanceID string

		// Check client instance
		clientID := r.URL.Query().Get("clientId")
		if clientID != "" {
			instanceID, err := f.redis.SafeGet(ctx, fmt.Sprintf("client:%s:instance", clientID))
			if err == nil && instanceID != "" && instanceID != f.instanceID {
				targetInstanceID = instanceID
				log.Info().
					Str("clientId", clientID).
					Str("targetInstance", targetInstanceID).
					Str("thisInstance", f.instanceID).
					Msg("Client on different instance")
			}
		}

		// Check API key instance
		if targetInstanceID == "" {
			instanceID, err := f.redis.SafeGet(ctx, fmt.Sprintf("apikey:%s:instance", apiKey))
			if err == nil && instanceID != "" && instanceID != f.instanceID {
				targetInstanceID = instanceID
			}
		}

		// No forwarding needed
		if targetInstanceID == "" {
			next.ServeHTTP(w, r)
			return
		}

		// SSE subscriptions require a streaming proxy with no fixed timeout.
		// A standard 20s timeout would kill the connection mid-stream.
		if r.Header.Get("Accept") == "text/event-stream" {
			f.forwardSSE(w, r, targetInstanceID)
			return
		}

		// Forward the request
		f.forwardRequest(w, r, targetInstanceID)
	})
}

func (f *RequestForwarder) forwardRequest(w http.ResponseWriter, r *http.Request, targetInstanceID string) {
	port := f.internalPort
	if port == "" {
		port = "3010"
	}
	targetURL := fmt.Sprintf("http://%s.vm.%s.internal:%s%s", targetInstanceID, f.appName, port, r.RequestURI)
	log.Info().Str("url", targetURL).Msg("Forwarding request")

	// Determine timeout
	isFileOp := r.URL.Path == "/upload" || r.URL.Path == "/download"
	timeout := 20 * time.Second
	if isFileOp {
		timeout = 45 * time.Second
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	// Create forwarded request
	proxyReq, err := http.NewRequestWithContext(ctx, r.Method, targetURL, r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create proxy request")
		http.Error(w, "Bad Gateway: failed to forward request", http.StatusBadGateway)
		return
	}

	// Copy headers (skip host)
	for key, values := range r.Header {
		if strings.EqualFold(key, "host") {
			continue
		}
		for _, v := range values {
			proxyReq.Header.Add(key, v)
		}
	}

	// Execute request using the shared client (context carries the timeout)
	resp, err := forwardClient.Do(proxyReq)
	if err != nil {
		log.Error().Err(err).Str("targetInstance", targetInstanceID).Msg("Forward request failed")
		http.Error(w, "Bad Gateway: upstream instance unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers (filter problematic ones)
	for key, values := range resp.Header {
		lower := strings.ToLower(key)
		if lower == "connection" || lower == "content-length" || lower == "transfer-encoding" {
			continue
		}
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	// Copy response, capped at 256 MB to match the JSON body limit and prevent memory exhaustion.
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, io.LimitReader(resp.Body, 256<<20))
}

// forwardSSE proxies an SSE subscription to the target instance with no fixed
// timeout. The connection lives as long as either side closes it (bounded only
// by the original request's context — the browser or SDK controls the lifetime).
// Events are flushed immediately so they reach clients without buffering delay.
func (f *RequestForwarder) forwardSSE(w http.ResponseWriter, r *http.Request, targetInstanceID string) {
	port := f.internalPort
	if port == "" {
		port = "3010"
	}
	targetURL := fmt.Sprintf("http://%s.vm.%s.internal:%s%s", targetInstanceID, f.appName, port, r.RequestURI)
	log.Info().Str("url", targetURL).Msg("Forwarding SSE request")

	// Use the request's own context — no additional timeout.
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create SSE proxy request")
		http.Error(w, "Bad Gateway: failed to forward SSE", http.StatusBadGateway)
		return
	}
	for key, values := range r.Header {
		if strings.EqualFold(key, "host") {
			continue
		}
		for _, v := range values {
			proxyReq.Header.Add(key, v)
		}
	}

	resp, err := sseForwardClient.Do(proxyReq)
	if err != nil {
		log.Error().Err(err).Str("targetInstance", targetInstanceID).Msg("SSE forward failed")
		http.Error(w, "Bad Gateway: upstream SSE unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		lower := strings.ToLower(key)
		if lower == "connection" || lower == "content-length" || lower == "transfer-encoding" {
			continue
		}
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Stream with per-write flushing so SSE events reach clients immediately.
	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			break
		}
	}
}
