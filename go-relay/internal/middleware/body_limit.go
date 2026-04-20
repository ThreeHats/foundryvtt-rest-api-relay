package middleware

import (
	"encoding/json"
	"net/http"
)

// BodyLimit returns middleware that rejects requests whose body exceeds maxBytes.
// It checks Content-Length first (fast path) and wraps the body with
// http.MaxBytesReader as a safety net for chunked/streaming requests.
func BodyLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				json.NewEncoder(w).Encode(map[string]string{"error": "Request body too large"}) //nolint:errcheck
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
