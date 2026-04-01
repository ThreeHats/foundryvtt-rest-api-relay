package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// PendingRequestTypes lists all valid message types for the WS protocol.
var PendingRequestTypes = map[string]bool{
	"search": true, "entity": true, "structure": true, "contents": true,
	"create": true, "update": true, "delete": true,
	"rolls": true, "last-roll": true, "roll": true, "get-sheet": true,
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
	"sheet-screenshot": true,
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
	defer p.mu.Unlock()
	p.requests[id] = req
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
	defer p.mu.Unlock()
	delete(p.requests, id)
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
		log.Warn().Str("requestId", id).Msg("Response channel full or closed")
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

// CleanupStale removes pending requests older than the given duration.
func (p *PendingRequests) CleanupStale(maxAge time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	for id, req := range p.requests {
		if req.Timestamp.Before(cutoff) {
			close(req.ResponseCh)
			delete(p.requests, id)
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

// SanitizeResponse deep-removes sensitive keys from response data.
func SanitizeResponse(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			if key == "privateKey" || key == "apiKey" || key == "password" {
				continue
			}
			result[key] = SanitizeResponse(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = SanitizeResponse(val)
		}
		return result
	default:
		return data
	}
}

// WriteJSON writes a sanitized JSON response.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	sanitized := SanitizeResponse(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(sanitized)
}
