package counter

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

// UserStore is the subset of model.UserStore needed by RequestCounter.
type UserStore interface {
	IncrementRequestsBy(ctx context.Context, id int64, count int) error
}

// RequestCounter batches per-user request counts in memory and flushes them to
// the DB periodically. This avoids one DB write per request (which serializes on
// SQLite's single writer lock) by coalescing many increments into a single UPDATE.
type RequestCounter struct {
	mu      sync.Mutex
	pending map[int64]int // userID -> count since last flush
}

// New creates a new RequestCounter.
func New() *RequestCounter {
	return &RequestCounter{
		pending: make(map[int64]int),
	}
}

// Add records one request for the given user. It is safe to call concurrently.
func (c *RequestCounter) Add(userID int64) {
	c.mu.Lock()
	c.pending[userID]++
	c.mu.Unlock()
}

// Flush writes all pending counts to the DB and resets the in-memory state.
// It is safe to call concurrently with Add.
func (c *RequestCounter) Flush(ctx context.Context, store UserStore) {
	c.mu.Lock()
	if len(c.pending) == 0 {
		c.mu.Unlock()
		return
	}
	snapshot := c.pending
	c.pending = make(map[int64]int, len(snapshot))
	c.mu.Unlock()

	for id, count := range snapshot {
		if err := store.IncrementRequestsBy(ctx, id, count); err != nil {
			log.Warn().Err(err).Int64("userId", id).Int("count", count).Msg("Failed to flush request counter")
		}
	}
}
