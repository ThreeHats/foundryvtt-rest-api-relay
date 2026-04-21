package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/metrics"
	"github.com/rs/zerolog/log"
)

// PendingRequestTypes lists all valid message types for the WS protocol.
var PendingRequestTypes = map[string]bool{
	"search": true, "entity": true, "structure": true, "contents": true,
	"create": true, "update": true, "delete": true,
	"rolls": true, "last-roll": true, "roll": true,
	"macro-execute": true, "macros": true,
	"encounters": true, "start-encounter": true, "next-turn": true, "next-round": true,
	"last-turn": true, "last-round": true, "end-encounter": true,
	"add-to-encounter": true, "remove-from-encounter": true,
	"kill": true, "decrease": true, "increase": true, "give": true, "remove": true,
	"execute-js": true, "select": true, "selected": true,
	"file-system": true, "upload-file": true, "download-file": true,
	"get-actor-details": true, "modify-item-charges": true,
	"use-ability": true, "use-feature": true, "use-spell": true, "use-item": true,
	"modify-experience": true, "add-item": true, "remove-item": true,
	"get-folder": true, "create-folder": true, "delete-folder": true,
	"players": true,
	"get-scene": true, "create-scene": true, "update-scene": true,
	"delete-scene": true, "switch-scene": true,
	"get-canvas-documents": true, "create-canvas-document": true,
	"update-canvas-document": true, "delete-canvas-document": true,
	"chat-messages": true, "chat-send": true, "chat-delete": true, "chat-flush": true,
	"short-rest": true, "long-rest": true, "skill-check": true,
	"ability-save": true, "ability-check": true, "death-save": true,
	"get-effects": true, "add-effect": true, "remove-effect": true,
	"get-status-effects": true,
	"sheet-screenshot": true,
	"move-token": true, "measure-distance": true,
	"get-playlists": true, "playlist-play": true, "playlist-stop": true,
	"playlist-next": true, "playlist-volume": true, "play-sound": true, "stop-sound": true,
	"world-info": true,
	"get-concentration": true, "break-concentration": true, "concentration-save": true,
	"equip-item": true, "attune-item": true, "transfer-currency": true,
	"scene-screenshot": true, "scene-raw-image": true,
	"get-users": true, "get-user": true, "create-user": true,
	"update-user": true, "delete-user": true,
}

// WSResponse is the response received from a Foundry client via WebSocket.
type WSResponse struct {
	StatusCode int
	Data       map[string]interface{}
	RawData    []byte // For binary responses (files, screenshots)
}

// PendingRequest tracks an in-flight request waiting for a WS response.
type PendingRequest struct {
	ResponseCh chan *WSResponse
	Type       string
	ClientID   string
	UUID       string
	Path       string
	Query      string
	Filter     string
	Format     string // "json", "html", "binary"
	ActiveTab  *int
	DarkMode   bool
	InitScale  *float64
	Timestamp  time.Time
	MaxAge     time.Duration // 0 = use the global cleanup default

	// For HTTP responses
	Writer  http.ResponseWriter
	Flusher http.Flusher

	// For WS callback responses
	WSCallback func(statusCode int, data interface{})
}

// PendingRequests is a thread-safe map of request ID to pending request.
type PendingRequests struct {
	mu       sync.RWMutex
	requests map[string]*PendingRequest
}

// NewPendingRequests creates a new PendingRequests tracker.
func NewPendingRequests() *PendingRequests {
	return &PendingRequests{
		requests: make(map[string]*PendingRequest),
	}
}

// Store adds a pending request.
func (p *PendingRequests) Store(id string, req *PendingRequest) {
	p.mu.Lock()
	p.requests[id] = req
	count := len(p.requests)
	p.mu.Unlock()
	metrics.WSPendingRequests.Set(float64(count))
}

// Count returns the number of pending requests currently tracked.
func (p *PendingRequests) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.requests)
}

// Load retrieves a pending request.
func (p *PendingRequests) Load(id string) (*PendingRequest, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	req, ok := p.requests[id]
	return req, ok
}

// Delete removes a pending request.
func (p *PendingRequests) Delete(id string) {
	p.mu.Lock()
	delete(p.requests, id)
	count := len(p.requests)
	p.mu.Unlock()
	metrics.WSPendingRequests.Set(float64(count))
}

// Resolve sends a response to a pending request and removes it.
func (p *PendingRequests) Resolve(id string, statusCode int, data map[string]interface{}) {
	p.mu.Lock()
	req, ok := p.requests[id]
	if ok {
		delete(p.requests, id)
	}
	p.mu.Unlock()

	if !ok {
		log.Warn().Str("requestId", id).Msg("No pending request found")
		return
	}

	resp := &WSResponse{StatusCode: statusCode, Data: data}

	select {
	case req.ResponseCh <- resp:
	default:
		log.Warn().
			Str("requestId", id).
			Str("clientId", req.ClientID).
			Str("type", req.Type).
			Msg("Response channel full or closed")
	}
}

// ResolveRaw sends raw bytes to a pending request (for binary responses).
func (p *PendingRequests) ResolveRaw(id string, statusCode int, raw []byte) {
	p.mu.Lock()
	req, ok := p.requests[id]
	if ok {
		delete(p.requests, id)
	}
	p.mu.Unlock()

	if !ok {
		return
	}

	resp := &WSResponse{StatusCode: statusCode, RawData: raw}
	select {
	case req.ResponseCh <- resp:
	default:
	}
}

// CleanupStale removes pending requests older than the given duration,
// sending each a 408 timeout response so waiting handlers exit cleanly.
// We never close ResponseCh — closing is not needed (goroutines use select
// with ctx.Done) and avoids any theoretical panic if a send races the close.
func (p *PendingRequests) CleanupStale(maxAge time.Duration) {
	p.mu.Lock()
	now := time.Now()
	var stale []*PendingRequest
	for id, req := range p.requests {
		age := maxAge
		if req.MaxAge > 0 {
			age = req.MaxAge
		}
		if req.Timestamp.Before(now.Add(-age)) {
			stale = append(stale, req)
			delete(p.requests, id)
		}
	}
	count := len(p.requests)
	p.mu.Unlock()

	metrics.WSPendingRequests.Set(float64(count))

	// Send timeout responses outside the lock so slow receivers don't block cleanup.
	for _, req := range stale {
		select {
		case req.ResponseCh <- &WSResponse{StatusCode: 408, Data: map[string]interface{}{"error": "request timeout"}}:
		default:
		}
	}
}

// StartCleanupLoop periodically removes stale pending requests.
func (p *PendingRequests) StartCleanupLoop(ctx interface{ Done() <-chan struct{} }, interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.CleanupStale(maxAge)
			case <-ctx.Done():
				return
			}
		}
	}()
}

