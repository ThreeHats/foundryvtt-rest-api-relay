package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
)

// Upload a file to Foundry's file system (handles both base64 and binary data)
//
// @tag FileSystem
// @param {string} clientId [query] Client ID for the Foundry world
// @param {string} path [query/body] The directory path to upload to
// @param {string} filename [query/body] The filename to save as
// @param {string} source [query/body] The source directory to use (data, systems, modules, etc.)
// @param {string} mimeType [query/body] The MIME type of the file
// @param {boolean} overwrite [query/body] Whether to overwrite an existing file
// @param {string} fileData [body] Base64 encoded file data (if sending as JSON) 250MB limit
// @param {string} userId [query/body] Foundry user ID or username to scope permissions (omit for GM-level access)
// @returns Result of the file upload operation
func uploadHandler(mgr *ws.ClientManager, pending *ws.PendingRequests) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		contentType := r.Header.Get("Content-Type")
		var filePath, filename, source, mimeType string
		var overwrite bool
		var fileData []byte

		// Enforce 250MB body size limit for file uploads
		r.Body = http.MaxBytesReader(w, r.Body, 250<<20)

		if strings.HasPrefix(contentType, "application/json") || contentType == "" {
			body, err := parseBody(r)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			filePath = bodyStr(body, "path")
			if filePath == "" {
				filePath = r.URL.Query().Get("path")
			}
			filename = bodyStr(body, "filename")
			if filename == "" {
				filename = r.URL.Query().Get("filename")
			}
			source = bodyStr(body, "source")
			if source == "" {
				source = "data"
			}
			mimeType = bodyStr(body, "mimeType")
			if ov, ok := body["overwrite"].(bool); ok {
				overwrite = ov
			}
			if fd, ok := body["fileData"].(string); ok {
				fileData = []byte(fd)
			}
		} else {
			filePath = r.URL.Query().Get("path")
			filename = r.URL.Query().Get("filename")
			source = r.URL.Query().Get("source")
			if source == "" {
				source = "data"
			}
			mimeType = r.URL.Query().Get("mimeType")
			overwrite = r.URL.Query().Get("overwrite") == "true"
			data, err := io.ReadAll(r.Body)
			if err != nil {
				helpers.WriteError(w, http.StatusBadRequest, "Failed to read upload data: "+err.Error())
				return
			}
			fileData = data
		}

		if filePath == "" || filename == "" {
			helpers.WriteError(w, http.StatusBadRequest, "'path' and 'filename' are required")
			return
		}

		clientID := r.URL.Query().Get("clientId")
		if reqCtx.ScopedKey != nil && reqCtx.ScopedKey.ScopedClientID != "" {
			clientID = reqCtx.ScopedKey.ScopedClientID
		}
		if clientID == "" {
			clients := mgr.GetConnectedClients(reqCtx.MasterAPIKey)
			if len(clients) == 1 {
				clientID = clients[0]
			} else {
				helpers.WriteError(w, http.StatusBadRequest, "'clientId' is required")
				return
			}
		}

		client := mgr.GetClient(clientID)
		if client == nil {
			helpers.WriteError(w, http.StatusNotFound, "Invalid client ID")
			return
		}

		requestID := fmt.Sprintf("upload-file_%d", time.Now().UnixMilli())
		responseCh := make(chan *ws.WSResponse, 1)
		pending.Store(requestID, &ws.PendingRequest{
			ResponseCh: responseCh, Type: "upload-file", ClientID: clientID, Timestamp: time.Now(),
		})

		payload := map[string]interface{}{
			"type": "upload-file", "requestId": requestID,
			"path": filePath, "filename": filename, "source": source,
			"mimeType": mimeType, "overwrite": overwrite, "fileData": string(fileData),
		}

		if !client.Send(payload) {
			pending.Delete(requestID)
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to send upload request")
			return
		}

		select {
		case resp := <-responseCh:
			if resp == nil {
				helpers.WriteError(w, http.StatusGatewayTimeout, "Upload timed out")
				return
			}
			if resp.RawData != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(resp.StatusCode)
				w.Write(resp.RawData)
			} else {
				helpers.WriteJSON(w, resp.StatusCode, resp.Data)
			}
		case <-time.After(30 * time.Second):
			pending.Delete(requestID)
			helpers.WriteError(w, http.StatusRequestTimeout, "Upload timed out")
		case <-r.Context().Done():
			pending.Delete(requestID)
		}
	}
}

// unused import guard
var _ = json.Marshal
