package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

	// Wrap the request context with a cancellable context so CloseForClientID
	// can terminate all SSE connections for this client when the Foundry module
	// disconnects (important for multi-instance: ensures SSE clients reconnect
	// to whichever instance the module lands on after a reconnect).
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	unregister := sseMgr.RegisterSSECancel(clientID, cancel)
	defer unregister()
	r = r.WithContext(ctx)

	// Ownership check: the requesting API key must own this client.
	// This prevents one user from subscribing to another user's real-time events.
	reqCtx := helpers.GetRequestContext(r)
	if reqCtx != nil && reqCtx.MasterAPIKey != "" && client.APIKey() != reqCtx.MasterAPIKey {
		helpers.WriteError(w, http.StatusForbidden, "Access denied")
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
	case "hooks":
		hookFilters := []string{}
		if hf := r.URL.Query().Get("hooks"); hf != "" {
			for _, h := range strings.Split(hf, ",") {
				h = strings.TrimSpace(h)
				if h != "" {
					hookFilters = append(hookFilters, h)
				}
			}
		}

		conn := &helpers.GenericSSEConnection{
			W: w, Flusher: flusher, ClientID: clientID, HookFilters: hookFilters, Done: r.Context().Done(),
		}

		if !sseMgr.AddGenericSSE(conn) {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", mustJSON(map[string]string{"error": "Too many SSE connections"}))
			flusher.Flush()
			return
		}

		log.Info().Str("clientId", clientID).Int("filters", len(hookFilters)).Msg("Hooks SSE connection opened")
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
				sseMgr.RemoveGenericSSE(conn)
				log.Info().Str("clientId", clientID).Msg("Hooks SSE connection closed")
				return
			}
		}

	case "combat":
		encounterId := r.URL.Query().Get("encounterId")
		conn := &helpers.CombatSSEConnection{
			W: w, Flusher: flusher, ClientID: clientID, EncounterID: encounterId, Done: r.Context().Done(),
		}

		if !sseMgr.AddCombatSSE(conn) {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", mustJSON(map[string]string{"error": "Too many SSE connections"}))
			flusher.Flush()
			return
		}

		log.Info().Str("clientId", clientID).Msg("Combat SSE connection opened")
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
				sseMgr.RemoveCombatSSE(conn)
				log.Info().Str("clientId", clientID).Msg("Combat SSE connection closed")
				return
			}
		}

	case "actor":
		actorUuid := r.URL.Query().Get("actorUuid")
		if actorUuid == "" {
			helpers.WriteError(w, http.StatusBadRequest, "'actorUuid' is required for actor event subscriptions")
			return
		}
		conn := &helpers.ActorSSEConnection{
			W: w, Flusher: flusher, ClientID: clientID, ActorUUID: actorUuid, Done: r.Context().Done(),
		}

		if !sseMgr.AddActorSSE(conn) {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", mustJSON(map[string]string{"error": "Too many SSE connections"}))
			flusher.Flush()
			return
		}

		log.Info().Str("clientId", clientID).Str("actorUuid", actorUuid).Msg("Actor SSE connection opened")
		fmt.Fprintf(w, "event: connected\ndata: %s\n\n", mustJSON(map[string]string{"clientId": clientID, "actorUuid": actorUuid}))
		flusher.Flush()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, ": keepalive\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				sseMgr.RemoveActorSSE(conn)
				log.Info().Str("clientId", clientID).Msg("Actor SSE connection closed")
				return
			}
		}

	case "scene":
		sceneId := r.URL.Query().Get("sceneId")
		conn := &helpers.SceneSSEConnection{
			W: w, Flusher: flusher, ClientID: clientID, SceneID: sceneId, Done: r.Context().Done(),
		}

		if !sseMgr.AddSceneSSE(conn) {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", mustJSON(map[string]string{"error": "Too many SSE connections"}))
			flusher.Flush()
			return
		}

		log.Info().Str("clientId", clientID).Msg("Scene SSE connection opened")
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
				sseMgr.RemoveSceneSSE(conn)
				log.Info().Str("clientId", clientID).Msg("Scene SSE connection closed")
				return
			}
		}

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
