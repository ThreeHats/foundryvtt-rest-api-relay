package service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// RemoteRequestBatcher batches cross-world remote-request notifications so
// that a single entity transfer (which may comprise 20-50 individual
// remote-request calls) triggers ONE summary notification instead of 20-50
// individual ones.
//
// Grouping key: (userID, sourceTokenID, targetClientID). Window defaults to
// 5 minutes from the first event in the batch, but can be set per-user via
// SetUserWindow. When the window expires, the batcher composes a summary and
// dispatches it through the Dispatcher.
type RemoteRequestBatcher struct {
	mu         sync.Mutex
	pending    map[batchKey]*batchEntry
	dispatcher *Dispatcher
	window     time.Duration
	done       chan struct{}
	// userWindow overrides the global window for specific users.
	// Key: userID, Value: per-user batch window duration.
	userWindow map[int64]time.Duration
}

type batchKey struct {
	UserID        int64
	SourceTokenID int64
	TargetClient  string
}

type batchEntry struct {
	ScopeCounts map[string]int // scope → count
	FailCount   int
	TotalCount  int
	FirstAt     time.Time
	LastAt      time.Time
	Timer       *time.Timer
	TokenName   string // for the notification body
}

// NewRemoteRequestBatcher creates a batcher with the given notification window.
func NewRemoteRequestBatcher(dispatcher *Dispatcher, window time.Duration) *RemoteRequestBatcher {
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &RemoteRequestBatcher{
		pending:    make(map[batchKey]*batchEntry),
		dispatcher: dispatcher,
		window:     window,
		done:       make(chan struct{}),
		userWindow: make(map[int64]time.Duration),
	}
}

// SetUserWindow configures a per-user batch window, overriding the global default.
// Pass 0 to revert to the global default. Safe to call concurrently.
func (b *RemoteRequestBatcher) SetUserWindow(userID int64, d time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if d <= 0 {
		delete(b.userWindow, userID)
	} else {
		b.userWindow[userID] = d
	}
}

// Add records a remote-request event. If this is the first event for the
// batch key, a timer is started. When the timer fires, the batch is flushed.
func (b *RemoteRequestBatcher) Add(userID int64, sourceTokenID int64, targetClientID, action, tokenName string, success bool) {
	key := batchKey{
		UserID:        userID,
		SourceTokenID: sourceTokenID,
		TargetClient:  targetClientID,
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	entry, exists := b.pending[key]
	if !exists {
		entry = &batchEntry{
			ScopeCounts: make(map[string]int),
			FirstAt:     time.Now(),
			TokenName:   tokenName,
		}
		// Use per-user window if configured, otherwise fall back to global default.
		window := b.window
		if uw, ok := b.userWindow[userID]; ok && uw > 0 {
			window = uw
		}
		// Start a timer for this batch. When it fires, flush and dispatch.
		entry.Timer = time.AfterFunc(window, func() {
			b.flush(key)
		})
		b.pending[key] = entry
	}

	// Use the action name directly as the grouping key. The scope mapping
	// lives in model/scopes.go which we can't import from the service
	// package (cycle risk). The action name is descriptive enough for the
	// notification summary (e.g., "create × 12, upload × 8").
	scope := action

	entry.ScopeCounts[scope]++
	entry.TotalCount++
	if !success {
		entry.FailCount++
	}
	entry.LastAt = time.Now()
}

// flush dispatches a summary notification for the given batch key and removes
// the entry from the pending map.
func (b *RemoteRequestBatcher) flush(key batchKey) {
	b.mu.Lock()
	entry, exists := b.pending[key]
	if !exists {
		b.mu.Unlock()
		return
	}
	delete(b.pending, key)
	b.mu.Unlock()

	// Build the summary description
	var lines []string
	for scope, count := range entry.ScopeCounts {
		lines = append(lines, fmt.Sprintf("%s × %d", scope, count))
	}
	failNote := ""
	if entry.FailCount > 0 {
		failNote = fmt.Sprintf(" (%d failed)", entry.FailCount)
	}
	desc := fmt.Sprintf(
		"Cross-world activity: %d action(s)%s in %s\n%s",
		entry.TotalCount,
		failNote,
		time.Since(entry.FirstAt).Round(time.Second),
		strings.Join(lines, "\n"),
	)
	if entry.TokenName != "" {
		desc = fmt.Sprintf("Token \"%s\" → %s\n%s", entry.TokenName, key.TargetClient, desc)
	}

	log.Info().
		Int64("userId", key.UserID).
		Int64("sourceTokenId", key.SourceTokenID).
		Str("targetClient", key.TargetClient).
		Int("totalActions", entry.TotalCount).
		Int("failedActions", entry.FailCount).
		Msg("Flushing batched remote-request notification")

	b.dispatcher.Dispatch(NotificationContext{
		Event:       EventRemoteRequest,
		UserID:      key.UserID,
		ClientID:    key.TargetClient,
		Description: desc,
		Severity:    "info",
	})
}

// Shutdown flushes all pending batches immediately and prevents new ones.
func (b *RemoteRequestBatcher) Shutdown() {
	close(b.done)
	b.mu.Lock()
	keys := make([]batchKey, 0, len(b.pending))
	for k, entry := range b.pending {
		entry.Timer.Stop()
		keys = append(keys, k)
	}
	b.mu.Unlock()

	for _, k := range keys {
		b.flush(k)
	}
}
