package metrics

import (
	"sort"
	"strings"
	"sync"
	"time"
)

// Rolling tracks per-minute, per-hour, and per-day request counts in memory.
// Ephemeral by design — resets on server restart.
//
// Buckets are time-aligned: minute buckets start at second 0, hour buckets at
// minute 0, day buckets at hour 0 UTC.
type Rolling struct {
	mu sync.RWMutex

	// minute buckets (last 60 minutes)
	minuteBuckets [60]int
	minuteIdx     int
	lastMinute    time.Time

	// hour buckets (last 24 hours)
	hourBuckets [24]int
	hourIdx     int
	lastHour    time.Time

	// day buckets (last 30 days)
	dayBuckets [30]int
	dayIdx     int
	lastDay    time.Time

	// per-user tally (lifetime since startup; cleared on reset)
	byUser map[int64]int

	// per-endpoint tally
	byEndpoint map[string]int

	// errors (5xx) total
	errorsTotal int
}

// NewRolling creates an in-memory rolling metrics tracker.
func NewRolling() *Rolling {
	now := time.Now()
	return &Rolling{
		lastMinute: now.Truncate(time.Minute),
		lastHour:   now.Truncate(time.Hour),
		lastDay:    now.Truncate(24 * time.Hour),
		byUser:     make(map[int64]int),
		byEndpoint: make(map[string]int),
	}
}

// Global is the package-level rolling metrics instance, initialized at server startup.
var Global = NewRolling()

// Record adds a request to the rolling counters.
func (r *Rolling) Record(userID int64, endpoint string, status int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.advance(now)

	r.minuteBuckets[r.minuteIdx]++
	r.hourBuckets[r.hourIdx]++
	r.dayBuckets[r.dayIdx]++

	if userID > 0 {
		r.byUser[userID]++
	}
	if isFoundryAPIPath(endpoint) {
		r.byEndpoint[endpoint]++
	}
	if status >= 500 {
		r.errorsTotal++
	}
}

// advance rolls forward the bucket cursors based on elapsed time, zeroing
// out skipped buckets so stale data doesn't pollute the rolling window.
func (r *Rolling) advance(now time.Time) {
	// Minutes
	currentMinute := now.Truncate(time.Minute)
	minutesElapsed := int(currentMinute.Sub(r.lastMinute) / time.Minute)
	if minutesElapsed > 0 {
		for i := 0; i < minutesElapsed && i < 60; i++ {
			r.minuteIdx = (r.minuteIdx + 1) % 60
			r.minuteBuckets[r.minuteIdx] = 0
		}
		if minutesElapsed >= 60 {
			for i := 0; i < 60; i++ {
				r.minuteBuckets[i] = 0
			}
		}
		r.lastMinute = currentMinute
	}

	currentHour := now.Truncate(time.Hour)
	hoursElapsed := int(currentHour.Sub(r.lastHour) / time.Hour)
	if hoursElapsed > 0 {
		for i := 0; i < hoursElapsed && i < 24; i++ {
			r.hourIdx = (r.hourIdx + 1) % 24
			r.hourBuckets[r.hourIdx] = 0
		}
		if hoursElapsed >= 24 {
			for i := 0; i < 24; i++ {
				r.hourBuckets[i] = 0
			}
		}
		r.lastHour = currentHour
	}

	currentDay := now.Truncate(24 * time.Hour)
	daysElapsed := int(currentDay.Sub(r.lastDay) / (24 * time.Hour))
	if daysElapsed > 0 {
		for i := 0; i < daysElapsed && i < 30; i++ {
			r.dayIdx = (r.dayIdx + 1) % 30
			r.dayBuckets[r.dayIdx] = 0
		}
		if daysElapsed >= 30 {
			for i := 0; i < 30; i++ {
				r.dayBuckets[i] = 0
			}
		}
		r.lastDay = currentDay
	}
}

// Overview returns summary stats: requests in the last min/hour/day + error count.
func (r *Rolling) Overview() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return map[string]interface{}{
		"requestsPerMinute": sum(r.minuteBuckets[:]),
		"requestsPerHour":   sum(r.hourBuckets[:]),
		"requestsPerDay":    sum(r.dayBuckets[:]),
		"errorsTotal":       r.errorsTotal,
	}
}

// ByEndpoint returns a sorted breakdown of request counts by endpoint pattern.
func (r *Rolling) ByEndpoint() []map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]map[string]interface{}, 0, len(r.byEndpoint))
	for ep, count := range r.byEndpoint {
		out = append(out, map[string]interface{}{
			"path":  ep,
			"count": count,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i]["count"].(int) > out[j]["count"].(int)
	})
	return out
}

// TopConsumers returns the top N users by request count.
func (r *Rolling) TopConsumers(n int) []map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	type entry struct {
		userID int64
		count  int
	}
	all := make([]entry, 0, len(r.byUser))
	for id, c := range r.byUser {
		all = append(all, entry{id, c})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].count > all[j].count })
	if n > 0 && len(all) > n {
		all = all[:n]
	}
	out := make([]map[string]interface{}, 0, len(all))
	for _, e := range all {
		out = append(out, map[string]interface{}{
			"userId": e.userID,
			"count":  e.count,
		})
	}
	return out
}

// Errors returns the total error (5xx) count since startup.
func (r *Rolling) Errors() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.errorsTotal
}

func sum(slice []int) int {
	total := 0
	for _, v := range slice {
		total += v
	}
	return total
}

// isFoundryAPIPath returns true for Foundry VTT API endpoint paths (the routes
// registered in RegisterAPIRoutes that external apps call via API key). It
// excludes dashboard, admin, auth, and WebSocket paths so the endpoint
// breakdown only shows operationally useful API traffic.
func isFoundryAPIPath(path string) bool {
	for _, prefix := range []string{"/api", "/auth", "/admin", "/ws", "/relay", "/_", "/docs", "/privacy", "/static"} {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}
	return path != "" && path != "unknown"
}

// Export returns a deep copy of the cumulative metrics state for persistence.
// Rolling time-window buckets are intentionally excluded — they are real-time
// gauges that reset on restart by design.
func (r *Rolling) Export() (byEndpoint map[string]int, byUser map[int64]int, errorsTotal int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ep := make(map[string]int, len(r.byEndpoint))
	for k, v := range r.byEndpoint {
		ep[k] = v
	}
	usr := make(map[int64]int, len(r.byUser))
	for k, v := range r.byUser {
		usr[k] = v
	}
	return ep, usr, r.errorsTotal
}

// Import restores cumulative metrics state from a persisted snapshot.
// Called once at startup before the server begins accepting requests.
func (r *Rolling) Import(byEndpoint map[string]int, byUser map[int64]int, errorsTotal int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byEndpoint = byEndpoint
	r.byUser = byUser
	r.errorsTotal = errorsTotal
}
