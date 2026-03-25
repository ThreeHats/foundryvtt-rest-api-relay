// src/routes/websocket.ts
import { WebSocketServer, WebSocket } from "ws";
import { log } from "../utils/logger";
import { ClientManager } from "../core/ClientManager";
import { validateHeadlessSession } from "../workers/headlessSessions";
import { validateApiKeyDetailed } from "../utils/validateApiKey";

// Read ping interval from environment variable, default to 20 seconds
const WEBSOCKET_PING_INTERVAL_MS = parseInt(process.env.WEBSOCKET_PING_INTERVAL_MS || '20000', 10);
// Read client cleanup interval from environment variable, default to 15 seconds
const CLIENT_CLEANUP_INTERVAL_MS = parseInt(process.env.CLIENT_CLEANUP_INTERVAL_MS || '15000', 10);

export const wsRoutes = (wss: WebSocketServer): void => {
  wss.on("connection", async (ws, req) => {
    try {
      // Parse URL parameters
      const url = new URL(req.url || "", `http://${req.headers.host}`);
      const id = url.searchParams.get("id");
      const token = url.searchParams.get("token");
      const worldId = url.searchParams.get("worldId");
      const worldTitle = url.searchParams.get("worldTitle");
      const foundryVersion = url.searchParams.get("foundryVersion");
      const systemId = url.searchParams.get("systemId");
      const systemTitle = url.searchParams.get("systemTitle");
      const systemVersion = url.searchParams.get("systemVersion");
      const customName = url.searchParams.get("customName");

      if (!id || !token) {
        log.warn("Rejecting WebSocket connection: missing id or token");
        ws.close(1008, "Missing client ID or token");
        return;
      }

      // Validate API key for ALL clients before accepting connection
      const keyValidation = await validateApiKeyDetailed(token);
      if (!keyValidation.valid) {
        log.warn(`Rejecting WebSocket connection: invalid API key for client ${id}`);
        ws.close(1008, "Invalid API key");
        return;
      }

      // For scoped keys, validate scopedClientId constraint and use master key for registration
      let registrationToken = token;
      if (keyValidation.masterApiKey && keyValidation.masterApiKey !== token) {
        // This is a scoped key — validate clientId constraint if set
        if (keyValidation.scopedClientId && keyValidation.scopedClientId !== id) {
          log.warn(`Rejecting WebSocket connection: scoped key clientId mismatch (expected ${keyValidation.scopedClientId}, got ${id})`);
          ws.close(1008, "Client ID does not match scoped API key");
          return;
        }
        // Register under the parent's master key so it appears in tokenGroups
        registrationToken = keyValidation.masterApiKey;
      }

      // Additional validation for headless sessions
      const isValid = await validateHeadlessSession(id, token);
      if (!isValid) {
        log.warn(`Rejecting invalid headless client: ${id}`);
        ws.close(1008, "Invalid headless session");
        return;
      }

      // Register client (use master key for scoped keys so it appears in correct tokenGroups)
      const client = await ClientManager.addClient(ws, id, registrationToken, worldId, worldTitle, foundryVersion, systemId, systemTitle, systemVersion, customName);
      if (!client) return; // Connection already rejected

      // Add protocol-level ping/pong to keep the TCP connection active
      const pingInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.ping(Buffer.from("keepalive"));
          log.debug(`Sent WebSocket protocol ping to ${id}`);
        }
      }, WEBSOCKET_PING_INTERVAL_MS); // Use configured interval

      // Handle disconnection
      ws.on("close", () => {
        clearInterval(pingInterval);
        ClientManager.removeClient(id);
      });

      // Handle pong responses to update client activity
      ws.on("pong", () => {
        // Update the client's last seen timestamp
        client.updateLastSeen();
      });

      // Handle errors
      ws.on("error", (error) => {
        clearInterval(pingInterval);
        log.error(`WebSocket error for client ${id}: ${error}`);
        ClientManager.removeClient(id);
      });
    } catch (error) {
      log.error(`WebSocket connection error: ${error}`);
      try {
        ws.close(1011, "Server error");
      } catch (e) {
        // Ignore errors closing socket
      }
    }
  });

  // Set up periodic cleanup
  setInterval(() => {
    ClientManager.cleanupInactiveClients();
  }, CLIENT_CLEANUP_INTERVAL_MS); // Use configured interval
};

// Export the ClientManager for usage in API routes
export { ClientManager };
