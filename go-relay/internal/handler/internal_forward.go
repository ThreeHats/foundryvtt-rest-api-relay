package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
)

type internalForwardBody struct {
	TargetClientID string                 `json:"targetClientId"`
	Action         string                 `json:"action"`
	Payload        map[string]interface{} `json:"payload"`
}

// InternalForwardActionHandler handles cross-instance remote-request forwarding.
// Called by a peer relay instance when the target Foundry client is local to this instance.
// Only reachable on Fly.io private network — no external auth required.
func InternalForwardActionHandler(manager *ws.ClientManager, pending *ws.PendingRequests) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResult := func(success bool, errMsg string, data map[string]interface{}) {
			w.Header().Set("Content-Type", "application/json")
			out := map[string]interface{}{"success": success}
			if errMsg != "" {
				out["error"] = errMsg
			}
			if data != nil {
				out["data"] = data
			}
			json.NewEncoder(w).Encode(out)
		}

		var body internalForwardBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeResult(false, "bad request body", nil)
			return
		}
		if body.TargetClientID == "" || body.Action == "" {
			writeResult(false, "targetClientId and action are required", nil)
			return
		}

		target := manager.GetClient(body.TargetClientID)
		if target == nil {
			writeResult(false, "target not connected on this instance", nil)
			return
		}

		relayRequestID := fmt.Sprintf("irr_%d_%s", time.Now().UnixNano(), body.Action)
		responseCh := make(chan *ws.WSResponse, 1)
		pending.Store(relayRequestID, &ws.PendingRequest{
			ResponseCh: responseCh,
			Type:       body.Action,
			ClientID:   target.ID(),
			Timestamp:  time.Now(),
			MaxAge:     90 * time.Second,
		})

		msg := make(map[string]interface{}, len(body.Payload)+2)
		for k, v := range body.Payload {
			msg[k] = v
		}
		msg["type"] = body.Action
		msg["requestId"] = relayRequestID

		if !target.Send(msg) {
			pending.Delete(relayRequestID)
			writeResult(false, "failed to send to target client", nil)
			return
		}

		select {
		case resp := <-responseCh:
			if resp == nil {
				writeResult(false, "nil response from target", nil)
				return
			}
			var result map[string]interface{}
			if resp.RawData != nil {
				json.Unmarshal(resp.RawData, &result)
			} else if resp.Data != nil {
				result = resp.Data
			}
			if result == nil {
				result = map[string]interface{}{}
			}
			delete(result, "requestId")
			writeResult(true, "", result)
		case <-time.After(60 * time.Second):
			pending.Delete(relayRequestID)
			writeResult(false, "request timed out", nil)
		}
	}
}
