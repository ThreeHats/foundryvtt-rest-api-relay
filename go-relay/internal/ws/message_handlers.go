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
	OnChatEvent   func(clientID string, data map[string]interface{})
	OnRollData    func(clientID string, data map[string]interface{})
	OnHookEvent   func(clientID string, data map[string]interface{})
	OnCombatEvent func(clientID string, data map[string]interface{})
	OnActorEvent  func(clientID string, data map[string]interface{})
	OnSceneEvent  func(clientID string, data map[string]interface{})
}

// RegisterMessageHandlers sets up all WS response handlers.
// fanout may be nil if SSE is not yet configured.
func RegisterMessageHandlers(manager *ClientManager, pending *PendingRequests, fanout *EventFanout, interactiveSessions *InteractiveSessionManager) {
	// For every pending request type, register a "<type>-result" handler
	// that resolves the matching pending request.
	for reqType := range PendingRequestTypes {
		rt := reqType // capture for closure

		// Skip types that have special response handlers
		if rt == "download-file" || rt == "sheet-screenshot" || rt == "scene-screenshot" {
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

	// Special handler: download-file-result (binary file download)
	manager.OnMessageType("download-file-result", func(client *Client, data map[string]interface{}) {
		requestID, _ := data["requestId"].(string)
		if requestID == "" {
			return
		}

		log.Debug().Str("requestId", requestID).Msg("Received file download response")

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

		log.Debug().Str("requestId", requestID).Msg("Received sheet screenshot response")

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

	// Special handler: scene-screenshot-result (same pattern as sheet-screenshot)
	manager.OnMessageType("scene-screenshot-result", func(client *Client, data map[string]interface{}) {
		requestID, _ := data["requestId"].(string)
		if requestID == "" {
			return
		}

		log.Debug().Str("requestId", requestID).Msg("Received scene screenshot response")

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

		if req.WSCallback != nil {
			pending.Resolve(requestID, 200, data)
			return
		}

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

	// Hook event handler — generic Foundry hook forwarding (firehose)
	manager.OnMessageType("hook-event", func(client *Client, data map[string]interface{}) {
		log.Debug().Str("clientId", client.ID()).Msg("Hook event received")
		if fanout != nil && fanout.OnHookEvent != nil {
			fanout.OnHookEvent(client.ID(), data)
		}
	})

	// Combat event handler
	manager.OnMessageType("combat-event", func(client *Client, data map[string]interface{}) {
		log.Debug().Str("clientId", client.ID()).Msg("Combat event received")
		if fanout != nil && fanout.OnCombatEvent != nil {
			fanout.OnCombatEvent(client.ID(), data)
		}
	})

	// Actor event handler
	manager.OnMessageType("actor-event", func(client *Client, data map[string]interface{}) {
		log.Debug().Str("clientId", client.ID()).Msg("Actor event received")
		if fanout != nil && fanout.OnActorEvent != nil {
			fanout.OnActorEvent(client.ID(), data)
		}
	})

	// Scene event handler
	manager.OnMessageType("scene-event", func(client *Client, data map[string]interface{}) {
		log.Debug().Str("clientId", client.ID()).Msg("Scene event received")
		if fanout != nil && fanout.OnSceneEvent != nil {
			fanout.OnSceneEvent(client.ID(), data)
		}
	})

	// Interactive session response handlers — forward Foundry responses to consumer WS
	if interactiveSessions != nil {
		manager.OnMessageType("interactive-session-started", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := interactiveSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			interactiveSessions.ActivateSession(sessionID)
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "interactive-session-started", "sessionId": sessionID,
				"imageData": data["imageData"], "mimeType": data["mimeType"],
				"width": data["width"], "height": data["height"],
			})
			session.ConsumerConn.WriteMessage(1, msg) // 1 = TextMessage
		})

		manager.OnMessageType("interactive-frame", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := interactiveSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			interactiveSessions.UpdateActivity(sessionID)
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "interactive-frame", "sessionId": sessionID,
				"imageData": data["imageData"], "mimeType": data["mimeType"],
				"width": data["width"], "height": data["height"], "trigger": data["trigger"],
			})
			session.ConsumerConn.WriteMessage(1, msg)
		})

		manager.OnMessageType("interactive-session-error", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := interactiveSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "interactive-session-error", "sessionId": sessionID, "error": data["error"],
			})
			session.ConsumerConn.WriteMessage(1, msg)
			if fatal, ok := data["fatal"].(bool); ok && fatal {
				interactiveSessions.EndSession(sessionID)
			}
		})

		manager.OnMessageType("interactive-session-ended", func(client *Client, data map[string]interface{}) {
			sessionID, _ := data["sessionId"].(string)
			session := interactiveSessions.GetSession(sessionID)
			if session == nil {
				return
			}
			msg, _ := json.Marshal(map[string]interface{}{
				"type": "interactive-session-ended", "sessionId": sessionID, "reason": data["reason"],
			})
			session.ConsumerConn.WriteMessage(1, msg)
			interactiveSessions.EndSession(sessionID)
		})
	}

	log.Info().Int("count", len(PendingRequestTypes)).Msg("Registered WS message handlers")
}
