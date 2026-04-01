package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// sseSubscribeHandler handles SSE subscribe for both roll and chat events.
func sseSubscribeHandler(w http.ResponseWriter, r *http.Request, mgr *ws.ClientManager, sseMgr *helpers.SSEManager, eventType string) {
	clientID := r.URL.Query().Get("clientId")
	if clientID == "" {
		helpers.WriteError(w, http.StatusBadRequest, "'clientId' is required")
		return
	}

	client := mgr.GetClient(clientID)
	if client == nil {
		helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch eventType {
	case "roll":
		filters := helpers.RollSSEFilters{}
		if uid := r.URL.Query().Get("userId"); uid != "" {
			filters.UserID = uid
		}

		conn := &helpers.RollSSEConnection{
			W: w, Flusher: flusher, ClientID: clientID, Filters: filters, Done: r.Context().Done(),
		}

		if !sseMgr.AddRollSSE(conn) {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", mustJSON(map[string]string{"error": "Too many SSE connections"}))
			flusher.Flush()
			return
		}

		log.Info().Str("clientId", clientID).Msg("Roll SSE connection opened")
		fmt.Fprintf(w, "event: connected\ndata: %s\n\n", mustJSON(map[string]string{"clientId": clientID}))
		flusher.Flush()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, ": keepalive\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				sseMgr.RemoveRollSSE(conn)
				log.Info().Str("clientId", clientID).Msg("Roll SSE connection closed")
				return
			}
		}

	case "chat":
		filters := helpers.SSEFilters{}
		if s := r.URL.Query().Get("speaker"); s != "" {
			filters.Speaker = s
		}
		if t := r.URL.Query().Get("type"); t != "" {
			if v, err := strconv.Atoi(t); err == nil {
				filters.Type = &v
			}
		}
		if r.URL.Query().Get("whisperOnly") == "true" {
			filters.WhisperOnly = true
		}
		if uid := r.URL.Query().Get("userId"); uid != "" {
			filters.UserID = uid
		}

		conn := &helpers.SSEConnection{
			W: w, Flusher: flusher, ClientID: clientID, Filters: filters, Done: r.Context().Done(),
		}

		if !sseMgr.AddChatSSE(conn) {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", mustJSON(map[string]string{"error": "Too many SSE connections"}))
			flusher.Flush()
			return
		}

		log.Info().Str("clientId", clientID).Msg("Chat SSE connection opened")
		fmt.Fprintf(w, "event: connected\ndata: %s\n\n", mustJSON(map[string]string{"clientId": clientID}))
		flusher.Flush()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, ": keepalive\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				sseMgr.RemoveChatSSE(conn)
				log.Info().Str("clientId", clientID).Msg("Chat SSE connection closed")
				return
			}
		}
	}
}

func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
