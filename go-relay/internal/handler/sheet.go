package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
)

// Get actor sheet HTML
//
// This endpoint retrieves the HTML for an actor sheet based on the provided UUID or selected actor.
// Only works on Foundry version 12.
// @tag Sheet
// @param {string} uuid [query] The UUID of the entity to get the sheet for
// @param {boolean} selected [query] Whether to get the sheet for the selected entity
// @param {boolean} actor [query] Whether to get the sheet for the selected token's actor if selected is true
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} format [query] The format to return the sheet in (html, json)
// @param {number} scale [query] The initial scale of the sheet
// @param {number} tab [query] The active tab index to open
// @param {boolean} darkMode [query] Whether to use dark mode for the sheet
// @param {string} userId [query] Foundry user ID or username to scope permissions (omit for GM-level access)
// @returns The sheet HTML or data depending on format requested
func sheetGetHandler(mgr *ws.ClientManager, pending *ws.PendingRequests) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		query := r.URL.Query()
		uuid := query.Get("uuid")
		selectedStr := query.Get("selected")
		format := query.Get("format")
		if format == "" {
			format = "html"
		}
		if uuid == "" && selectedStr == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Either 'uuid' or 'selected' is required")
			return
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

		requestID := fmt.Sprintf("get-sheet_%d", time.Now().UnixMilli())
		responseCh := make(chan *ws.WSResponse, 1)
		pendingReq := &ws.PendingRequest{ResponseCh: responseCh, Type: "get-sheet", ClientID: clientID, Format: format, Timestamp: time.Now()}
		if tab := query.Get("tab"); tab != "" {
			if v, err := strconv.Atoi(tab); err == nil {
				pendingReq.ActiveTab = &v
			}
		}
		if query.Get("darkMode") == "true" {
			pendingReq.DarkMode = true
		}
		if scale := query.Get("scale"); scale != "" {
			if v, err := strconv.ParseFloat(scale, 64); err == nil {
				pendingReq.InitScale = &v
			}
		}
		pending.Store(requestID, pendingReq)

		payload := map[string]interface{}{"type": "get-sheet", "requestId": requestID, "uuid": uuid, "data": map[string]interface{}{}}
		if selectedStr != "" {
			payload["selected"] = selectedStr == "true"
		}
		if query.Get("actor") != "" {
			payload["actor"] = query.Get("actor") == "true"
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
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(resp.StatusCode)
				w.Write(resp.RawData)
			} else {
				if resp.RawData != nil { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(resp.StatusCode); w.Write(resp.RawData) } else { helpers.WriteJSON(w, resp.StatusCode, resp.Data) }
			}
		case <-time.After(20 * time.Second):
			pending.Delete(requestID)
			helpers.WriteError(w, http.StatusRequestTimeout, "Request timed out")
		case <-r.Context().Done():
			pending.Delete(requestID)
		}
	}
}

// Get actor sheet as screenshot image
//
// Captures the rendered actor sheet using html2canvas and returns it as a PNG or JPEG image.
// Works on both Foundry v12 and v13+.
// @tag SheetImage
// @param {string} uuid [query] The UUID of the entity to screenshot
// @param {boolean} selected [query] Whether to screenshot the selected entity's sheet
// @param {boolean} actor [query] Whether to use the selected token's actor if selected is true
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} format [query] Image format: png or jpeg (default: png)
// @param {number} quality [query] Image quality 0-1 for JPEG (default: 0.9)
// @param {number} scale [query] Capture scale factor (default: 1)
// @param {string} userId [query] Foundry user ID or username to scope permissions (omit for GM-level access)
// @returns The sheet screenshot as an image
func sheetImageHandler(mgr *ws.ClientManager, pending *ws.PendingRequests) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		query := r.URL.Query()
		uuid := query.Get("uuid")
		selectedStr := query.Get("selected")
		if uuid == "" && selectedStr == "" {
			helpers.WriteError(w, http.StatusBadRequest, "Either 'uuid' or 'selected' is required")
			return
		}

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
		scale := 1.0
		if s := query.Get("scale"); s != "" {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				scale = v
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

		requestID := fmt.Sprintf("sheet-screenshot_%d", time.Now().UnixMilli())
		responseCh := make(chan *ws.WSResponse, 1)
		pending.Store(requestID, &ws.PendingRequest{ResponseCh: responseCh, Type: "sheet-screenshot", ClientID: clientID, Format: format, Timestamp: time.Now()})

		payload := map[string]interface{}{"type": "sheet-screenshot", "requestId": requestID, "uuid": uuid, "format": format, "quality": quality, "scale": scale, "data": map[string]interface{}{}}
		if selectedStr != "" {
			payload["selected"] = selectedStr == "true"
		}
		if query.Get("actor") != "" {
			payload["actor"] = query.Get("actor") == "true"
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
				if resp.RawData != nil { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(resp.StatusCode); w.Write(resp.RawData) } else { helpers.WriteJSON(w, resp.StatusCode, resp.Data) }
			}
		case <-time.After(20 * time.Second):
			pending.Delete(requestID)
			helpers.WriteError(w, http.StatusRequestTimeout, "Request timed out")
		case <-r.Context().Done():
			pending.Delete(requestID)
		}
	}
}
