package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/rs/zerolog/log"
)

// AdminIPAllowlist returns middleware that restricts access to a list of IPs/CIDRs.
// An empty allowlist means all IPs are allowed (open mode — useful in development).
//
// Relies solely on r.RemoteAddr, which chi's RealIP middleware (applied earlier in the
// chain) has already normalized from X-Forwarded-For. Do NOT check X-Forwarded-For here
// directly — doing so would allow spoofing by injecting arbitrary headers.
func AdminIPAllowlist(allowedCSV string) func(http.Handler) http.Handler {
	allowed := parseAllowlist(allowedCSV)
	return func(next http.Handler) http.Handler {
		if len(allowed.exact) == 0 && len(allowed.cidrs) == 0 {
			// Open mode — pass through
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if ip == nil {
				helpers.WriteError(w, http.StatusForbidden, "Access denied")
				return
			}
			if !allowed.contains(ip) {
				log.Warn().Str("ip", ip.String()).Str("path", r.URL.Path).Msg("admin: blocked by IP allowlist")
				helpers.WriteError(w, http.StatusForbidden, "Access denied")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type ipAllowlist struct {
	exact map[string]struct{}
	cidrs []*net.IPNet
}

func (a *ipAllowlist) contains(ip net.IP) bool {
	if _, ok := a.exact[ip.String()]; ok {
		return true
	}
	for _, c := range a.cidrs {
		if c.Contains(ip) {
			return true
		}
	}
	return false
}

func parseAllowlist(csv string) *ipAllowlist {
	result := &ipAllowlist{exact: make(map[string]struct{})}
	if csv == "" {
		return result
	}
	for _, raw := range strings.Split(csv, ",") {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			_, cidr, err := net.ParseCIDR(entry)
			if err != nil {
				log.Warn().Str("entry", entry).Err(err).Msg("admin: invalid CIDR in ADMIN_ALLOWED_IPS")
				continue
			}
			result.cidrs = append(result.cidrs, cidr)
			continue
		}
		ip := net.ParseIP(entry)
		if ip == nil {
			log.Warn().Str("entry", entry).Msg("admin: invalid IP in ADMIN_ALLOWED_IPS")
			continue
		}
		result.exact[ip.String()] = struct{}{}
		// If any loopback address is listed, accept both IPv4 and IPv6 loopback.
		// Node.js resolves "localhost" to ::1 while browsers typically use 127.0.0.1.
		if ip.IsLoopback() {
			result.exact["127.0.0.1"] = struct{}{}
			result.exact["::1"] = struct{}{}
		}
	}
	return result
}

func clientIP(r *http.Request) net.IP {
	// Use only RemoteAddr — chi's RealIP middleware has already resolved the
	// true client IP from X-Forwarded-For before this middleware runs.
	// Reading X-Forwarded-For again here would allow IP spoofing.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if ip := net.ParseIP(r.RemoteAddr); ip != nil {
			return ip
		}
		return nil
	}
	return net.ParseIP(host)
}
