package ws

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
)

// RegisterMessageHandlers sets up the WS message type handlers that route
// Foundry responses back to pending HTTP/WS requests. This is the critical
// bridge between the WS protocol and the REST API.
// EventFanout is called when chat or roll events arrive from Foundry.
// Set by the server after SSEManager is created.
type EventFanout struct {
	OnChatEvent func(clientID string, data map[string]interface{})
	OnRollData  func(clientID string, data map[string]interface{})
}

// RegisterMessageHandlers sets up all WS response handlers.
// fanout may be nil if SSE is not yet configured.
func RegisterMessageHandlers(manager *ClientManager, pending *PendingRequests, fanout *EventFanout, sheetSessions *SheetSessionManager) {
	// For every pending request type, register a "<type>-result" handler
	// that resolves the matching pending request.
	for reqType := range PendingRequestTypes {
		rt := reqType // capture for closure

		// Skip types that have special response handlers
		if rt == "get-sheet" || rt == "download-file" || rt == "sheet-screenshot" {
			continue
		}

		manager.OnRawMessageType(rt+"-result", func(client *Client, data map[string]interface{}, rawBytes []byte) {
			requestID, _ := data["requestId"].(string)
			if requestID == "" {
				return
			}

			// Determine status code
			statusCode := 200
			if _, hasError := data["error"]; hasError {
				statusCode = 400
			}

			// Send raw bytes directly — avoids re-marshaling the entire response
			pending.ResolveRaw(requestID, statusCode, rawBytes)
		})
	}

	// Special handler: get-sheet-response (actor sheet HTML)
	manager.OnMessageType("get-sheet-response", func(client *Client, data map[string]interface{}) {
		requestID, _ := data["requestId"].(string)
		if requestID == "" {
			return
		}

		log.Info().Str("requestId", requestID).Msg("Received actor sheet response")

		req, ok := pending.Load(requestID)
		if !ok {
			return
		}

		// Check for errors
		errMsg, _ := data["error"].(string)
		if errMsg == "" {
			if d, ok := data["data"].(map[string]interface{}); ok {
				errMsg, _ = d["error"].(string)
			}
		}

		if errMsg != "" {
			pending.Resolve(requestID, 400, map[string]interface{}{
				"requestId": requestID,
				"clientId":  req.ClientID,
				"error":     errMsg,
			})
			return
		}

		// Extract HTML, CSS, and UUID from either data root or data.data
		html, _ := data["html"].(string)
		css, _ := data["css"].(string)
		responseUUID, _ := data["uuid"].(string)
		if responseUUID == "" {
			if d, ok := data["data"].(map[string]interface{}); ok {
				if u, ok := d["uuid"].(string); ok {
					responseUUID = u
				}
			}
		}
		if html == "" {
			if d, ok := data["data"].(map[string]interface{}); ok {
				if h, ok := d["html"].(string); ok {
					html = h
				}
				if c, ok := d["css"].(string); ok {
					css = c
				}
			}
		}

		if req.Format == "json" || req.WSCallback != nil {
			// JSON format
			resp := map[string]interface{}{
				"requestId": requestID,
				"clientId":  req.ClientID,
				"html":      html,
				"css":       css,
			}
			if responseUUID != "" {
				resp["uuid"] = responseUUID
			}
			pending.Resolve(requestID, 200, resp)
		} else {
			// HTML format — send raw HTML bytes
			// TODO: Full HTML template rendering (Phase 12)
			// For now, return the raw HTML
			pending.ResolveRaw(requestID, 200, []byte(html))
		}
	})

	// Special handler: download-file-result (binary file download)
	manager.OnMessageType("download-file-result", func(client *Client, data map[string]interface{}) {
		requestID, _ := data["requestId"].(string)
		if requestID == "" {
			return
		}

		log.Info().Str("requestId", requestID).Msg("Received file download response")

		req, ok := pending.Load(requestID)
		if !ok {
			return
		}

		// Check for error
		if errMsg, ok := data["error"].(string); ok && errMsg != "" {
			pending.Resolve(requestID, 400, map[string]interface{}{
				"requestId": requestID,
				"error":     errMsg,
			})
			return
		}

		// For WS callbacks or base64 format, send JSON with base64 data
		if req.WSCallback != nil || req.Format == "base64" {
			pending.Resolve(requestID, 200, data)
			return
		}

		// For HTTP binary format, decode base64 to binary
		fileDataStr, _ := data["fileData"].(string)
		if fileDataStr == "" {
			pending.Resolve(requestID, 200, data)
			return
		}

		// Handle data URL format: "data:mime;base64,XXXX"
		b64 := fileDataStr
		if idx := strings.Index(b64, ","); idx >= 0 {
			b64 = b64[idx+1:]
		}

		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			// Fall back to JSON
			pending.Resolve(requestID, 200, data)
			return
		}

		mimeType, _ := data["mimeType"].(string)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		fileName, _ := data["fileName"].(string)
		if fileName == "" {
			fileName, _ = data["filename"].(string)
		}

		resp := &WSResponse{
			StatusCode: 200,
			RawData:    decoded,
			Data: map[string]interface{}{
				"_contentType":        mimeType,
				"_contentDisposition": "attachment; filename=\"" + fileName + "\"",
				"_contentLength":      len(decoded),
			},
		}
		pendReq, _ := pending.Load(requestID)
		if pendReq != nil {
			pending.Delete(requestID)
			select {
			case pendReq.ResponseCh <- resp:
			default:
			}
		}
	})

	// Special handler: sheet-screenshot-result
	manager.OnMessageType("sheet-screenshot-result", func(client *Client, data map[string]interface{}) {
		requestID, _ := data["requestId"].(string)
		if requestID == "" {
			return
		}

		log.Info().Str("requestId", requestID).Msg("Received sheet screenshot response")

		req, ok := pending.Load(requestID)
		if !ok {
			return
		}

		if errMsg, ok := data["error"].(string); ok && errMsg != "" {
			pending.Resolve(requestID, 400, map[string]interface{}{
				"requestId": requestID,
				"error":     errMsg,
			})
			return
		}

		// For WS callbacks, send JSON
		if req.WSCallback != nil {
			pending.Resolve(requestID, 200, data)
			return
		}

		// For HTTP, decode base64 to binary image
		imageDataStr, _ := data["imageData"].(string)
		if imageDataStr == "" {
			if d, ok := data["data"].(map[string]interface{}); ok {
				imageDataStr, _ = d["imageData"].(string)
			}
		}
		if imageDataStr == "" {
			pending.Resolve(requestID, 200, data)
			return
		}

		b64 := imageDataStr
		if idx := strings.Index(b64, ","); idx >= 0 {
			b64 = b64[idx+1:]
		}

		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			pending.Resolve(requestID, 200, data)
			return
		}

		mimeType, _ := data["mimeType"].(string)
		if mimeType == "" {
			mimeType = "image/png"
		}

		resp := &WSResponse{
			StatusCode: 200,
			RawData:    decoded,
			Data: map[string]interface{}{
				"_contentType": mimeType,
				"width":        data["width"],
				"height":       data["height"],
			},
		}
		pendReq, _ := pending.Load(requestID)
		if pendReq != nil {
			pending.Delete(requestID)
			select {
			case pendReq.ResponseCh <- resp:
			default:
			}
		}
	})

	// Chat event handler — fan out to SSE and WS event connections
	manager.OnMessageType("chat-event", func(client *Client, data map[string]interface{}) {
		log.Debug().Str("clientId", client.ID()).Msg("Chat event received")
		if fanout != nil && fanout.OnChatEvent != nil {
			fanout.OnChatEvent(client.ID(), data)
		}
	})

	// Roll data handler — fan out to SSE and WS event connections
	manager.OnMessageType("roll-data", func(client *Client, data map[string]interface{}) {
		log.Debug().Str("clientId", client.ID()).Msg("Roll data received")
		if fanout != nil && fanout.OnRollData != nil {
			fanout.OnRollData(client.ID(), data)
		}
	})

	// Sheet session response handlers — forward Foundry responses to consumer WS
	if sheetSessions != nil {
		manager.OnMessageType("sheet-session-started", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := sheetSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			sheetSessions.ActivateSession(sessionID)
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "sheet-session-started", "sessionId": sessionID,
				"imageData": data["imageData"], "mimeType": data["mimeType"],
				"width": data["width"], "height": data["height"],
			})
			session.ConsumerConn.WriteMessage(1, msg) // 1 = TextMessage
		})

		manager.OnMessageType("sheet-frame", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := sheetSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			sheetSessions.UpdateActivity(sessionID)
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "sheet-frame", "sessionId": sessionID,
				"imageData": data["imageData"], "mimeType": data["mimeType"],
				"width": data["width"], "height": data["height"], "trigger": data["trigger"],
			})
			session.ConsumerConn.WriteMessage(1, msg)
		})

		manager.OnMessageType("sheet-session-error", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := sheetSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "sheet-session-error", "sessionId": sessionID, "error": data["error"],
			})
			session.ConsumerConn.WriteMessage(1, msg)
			if fatal, ok := data["fatal"].(bool); ok && fatal {
				sheetSessions.EndSession(sessionID)
			}
		})

		manager.OnMessageType("sheet-session-ended", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := sheetSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "sheet-session-ended", "sessionId": sessionID, "reason": data["reason"],
			})
			session.ConsumerConn.WriteMessage(1, msg)
			sheetSessions.EndSession(sessionID)
		})
	}

	log.Info().Int("count", len(PendingRequestTypes)).Msg("Registered WS message handlers")
}
