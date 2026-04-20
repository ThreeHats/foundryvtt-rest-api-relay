package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
)

// Get a rendered screenshot of a scene
//
// Captures the full rendered canvas of a scene including all visible layers
// (tokens, lights, walls, etc.) as an image. The scene can be specified by ID
// or defaults to the active scene.
// @tag Scene
// @param {string} sceneId [query] Scene ID (defaults to viewed/active scene)
// @param {boolean} active [query] If true, explicitly use the player-facing active scene instead of the viewed scene
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} format [query] Image format: png or jpeg (default: png)
// @param {number} quality [query] Image quality 0-1 for JPEG (default: 0.9)
// @param {boolean} viewport [query] If true, capture exactly what the browser currently shows instead of the full scene
// @param {number} width [query] Output image width in pixels (default: scene width)
// @param {number} height [query] Output image height in pixels (default: scene height)
// @param {boolean} showGrid [query] Include grid lines in capture (default: false)
// @param {boolean} hideOverlays [query] Hide fog of war, weather, vision, and UI overlays (default: false)
// @param {string} userId [query] Foundry user ID or username
// @returns The scene screenshot as an image
func sceneImageHandler(mgr *ws.ClientManager, pending *ws.PendingRequests) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		query := r.URL.Query()

		format := query.Get("format")
		if format == "" {
			format = "png"
		}
		quality := 0.9
		if q := query.Get("quality"); q != "" {
			if v, err := strconv.ParseFloat(q, 64); err == nil {
				quality = v
			}
		}

		clientID := query.Get("clientId")
		if reqCtx != nil && reqCtx.ScopedKey != nil && reqCtx.ScopedKey.ScopedClientID != "" {
			clientID = reqCtx.ScopedKey.ScopedClientID
		}
		if clientID == "" && reqCtx != nil {
			clients := mgr.GetConnectedClients(reqCtx.MasterAPIKey)
			if len(clients) == 1 {
				clientID = clients[0]
			}
		}
		if clientID == "" {
			helpers.WriteError(w, http.StatusBadRequest, "'clientId' is required")
			return
		}
		client := mgr.GetClient(clientID)
		if client == nil {
			helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
			return
		}

		requestID := fmt.Sprintf("scene-screenshot_%d", time.Now().UnixMilli())
		responseCh := make(chan *ws.WSResponse, 1)
		pending.Store(requestID, &ws.PendingRequest{ResponseCh: responseCh, Type: "scene-screenshot", ClientID: clientID, Format: format, Timestamp: time.Now()})

		payload := map[string]interface{}{
			"type": "scene-screenshot", "requestId": requestID,
			"format": format, "quality": quality,
			"data": map[string]interface{}{},
		}
		if sceneId := query.Get("sceneId"); sceneId != "" {
			payload["sceneId"] = sceneId
		}
		if query.Get("active") == "true" {
			payload["active"] = true
		}
		if query.Get("viewport") == "true" {
			payload["viewport"] = true
		}
		if query.Get("showGrid") == "true" {
			payload["showGrid"] = true
		}
		if query.Get("hideOverlays") == "true" {
			payload["hideOverlays"] = true
		}
		if w := query.Get("width"); w != "" {
			if v, err := strconv.Atoi(w); err == nil {
				payload["width"] = v
			}
		}
		if h := query.Get("height"); h != "" {
			if v, err := strconv.Atoi(h); err == nil {
				payload["height"] = v
			}
		}
		if userId := query.Get("userId"); userId != "" {
			payload["userId"] = userId
		}

		if !client.Send(payload) {
			pending.Delete(requestID)
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to send request")
			return
		}

		select {
		case resp := <-responseCh:
			if resp == nil {
				helpers.WriteError(w, http.StatusGatewayTimeout, "Request timed out")
				return
			}
			if resp.RawData != nil {
				mimeType := "image/png"
				if format == "jpeg" {
					mimeType = "image/jpeg"
				}
				w.Header().Set("Content-Type", mimeType)
				if resp.Data != nil {
					if width, ok := resp.Data["width"]; ok {
						w.Header().Set("X-Image-Width", fmt.Sprintf("%v", width))
					}
					if height, ok := resp.Data["height"]; ok {
						w.Header().Set("X-Image-Height", fmt.Sprintf("%v", height))
					}
				}
				w.WriteHeader(http.StatusOK)
				w.Write(resp.RawData)
			} else {
				helpers.WriteJSON(w, resp.StatusCode, resp.Data)
			}
		case <-time.After(30 * time.Second):
			pending.Delete(requestID)
			helpers.WriteError(w, http.StatusRequestTimeout, "Request timed out")
		case <-r.Context().Done():
			pending.Delete(requestID)
		}
	}
}
