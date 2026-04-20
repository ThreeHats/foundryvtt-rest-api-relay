package middleware

import (
	"context"
	"strings"
	"sync"
	"time"
)

// AdminLoginRateLimiter rate-limits admin login attempts per IP.
// Initialized with defaults; reinitialized via InitAdminRateLimiters() after .env loads.
var AdminLoginRateLimiter *RateLimiter

func init() {
	ctx := context.Background()
	AdminLoginRateLimiter = NewRateLimiter(ctx, 15*time.Minute, getEnvInt("ADMIN_LOGIN_RATE_LIMIT", 10),
		"Too many admin login attempts from this IP, please try again after 15 minutes")
}

// InitAdminRateLimiters re-initializes admin rate limiters after .env load.
func InitAdminRateLimiters(ctx context.Context) {
	AdminLoginRateLimiter = NewRateLimiter(ctx, 15*time.Minute, getEnvInt("ADMIN_LOGIN_RATE_LIMIT", 10),
		"Too many admin login attempts from this IP, please try again after 15 minutes")
}

// AdminLockoutTracker tracks failed admin login attempts per email
// and locks accounts after exceeding the threshold.
type AdminLockoutTracker struct {
	mu         sync.Mutex
	state      map[string]*lockoutState
	threshold  int
	lockWindow time.Duration
}

type lockoutState struct {
	failures    int
	lockedUntil time.Time
}

// NewAdminLockoutTracker creates a tracker that locks accounts after `threshold`
// failed attempts for `lockWindow` duration.
func NewAdminLockoutTracker(threshold int, lockWindow time.Duration) *AdminLockoutTracker {
	t := &AdminLockoutTracker{
		state:      make(map[string]*lockoutState),
		threshold:  threshold,
		lockWindow: lockWindow,
	}
	go t.cleanupLoop()
	return t
}

// IsLocked reports whether the given email is currently locked out.
func (t *AdminLockoutTracker) IsLocked(email string) bool {
	if t == nil || t.threshold <= 0 {
		return false
	}
	key := strings.ToLower(email)
	t.mu.Lock()
	defer t.mu.Unlock()
	st, ok := t.state[key]
	if !ok {
		return false
	}
	if !st.lockedUntil.IsZero() && time.Now().Before(st.lockedUntil) {
		return true
	}
	// Lock window passed — reset
	if !st.lockedUntil.IsZero() && time.Now().After(st.lockedUntil) {
		delete(t.state, key)
	}
	return false
}

// RecordFailure records a failed login attempt and locks the account if the
// threshold is exceeded.
func (t *AdminLockoutTracker) RecordFailure(email string) {
	if t == nil || t.threshold <= 0 {
		return
	}
	key := strings.ToLower(email)
	t.mu.Lock()
	defer t.mu.Unlock()
	st, ok := t.state[key]
	if !ok {
		st = &lockoutState{}
		t.state[key] = st
	}
	st.failures++
	if st.failures >= t.threshold {
		st.lockedUntil = time.Now().Add(t.lockWindow)
	}
}

// RecordSuccess clears any failed-attempt state for the email.
func (t *AdminLockoutTracker) RecordSuccess(email string) {
	if t == nil {
		return
	}
	key := strings.ToLower(email)
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.state, key)
}

// LockedUntil returns the time at which the lock expires (zero if not locked).
func (t *AdminLockoutTracker) LockedUntil(email string) time.Time {
	if t == nil {
		return time.Time{}
	}
	key := strings.ToLower(email)
	t.mu.Lock()
	defer t.mu.Unlock()
	if st, ok := t.state[key]; ok {
		return st.lockedUntil
	}
	return time.Time{}
}

func (t *AdminLockoutTracker) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		t.mu.Lock()
		now := time.Now()
		for k, st := range t.state {
			if !st.lockedUntil.IsZero() && now.After(st.lockedUntil) {
				delete(t.state, k)
			}
		}
		t.mu.Unlock()
	}
}

// AdminLockout is the package-level lockout tracker for admin logins.
// Initialized with defaults; reinitialized via InitAdminRateLimiters().
var AdminLockout = NewAdminLockoutTracker(10, 30*time.Minute)

// InitAdminLockout re-initializes the admin lockout tracker from env vars.
func InitAdminLockout() {
	AdminLockout = NewAdminLockoutTracker(
		getEnvInt("ADMIN_LOCKOUT_THRESHOLD", 10),
		time.Duration(getEnvInt("ADMIN_LOCKOUT_MINUTES", 30))*time.Minute,
	)
}
