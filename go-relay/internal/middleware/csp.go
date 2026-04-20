package middleware

import "net/http"

// AdminCSP applies a strict Content-Security-Policy and related security headers
// to admin dashboard responses.
//
// Notes:
// - 'unsafe-inline' is included for style-src because Svelte components emit
//   scoped <style> blocks at build time. If a future build uses CSS extraction
//   and a nonce, that can be tightened.
// - 'self' is the only allowed origin for scripts/connections — no CDN, no inline scripts.
// - 'frame-ancestors none' prevents the admin pages being embedded in an iframe.
func AdminCSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data:; "+
				"font-src 'self'; "+
				"connect-src 'self'; "+
				"frame-ancestors 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
