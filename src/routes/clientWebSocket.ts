import { WebSocketServer, WebSocket } from 'ws';
import { IncomingMessage } from 'http';
import { log } from '../utils/logger';
import { validateApiKeyDetailed, trackWsApiUsage } from '../utils/validateApiKey';
import { ClientManager } from '../core/ClientManager';
import {
  pendingRequests,
  PENDING_REQUEST_TYPES,
  PendingRequestType,
  WSEventConnection,
  addWSEventConnection,
  removeWSEventConnection,
} from './shared';
import { SheetSessionManager } from '../core/SheetSessionManager';
import { ApiKey } from '../models/apiKey';
import { startHeadlessWithStoredCredentials } from '../utils/headlessSessionStarter';

const PING_INTERVAL_MS = parseInt(process.env.WEBSOCKET_PING_INTERVAL_MS || '20000', 10);

interface ClientWSState {
  apiKey: string;
  masterApiKey: string; // Master key for usage tracking (same as apiKey for non-scoped keys)
  clientId: string;
  scopedUserId?: string | null;
  subscriptions: Set<WSEventConnection>;
  pendingRequestIds: Set<string>;
  ws: WebSocket;
}

/**
 * Handle a new client-facing WebSocket connection on /ws/api.
 * Authenticates via ?token=<apiKey>&clientId=<id> query params.
 */
export function setupClientWebSocket(wss: WebSocketServer): void {
  // Initialize the sheet session manager cleanup interval
  SheetSessionManager.init();

  wss.on('connection', async (ws: WebSocket, req: IncomingMessage) => {
    let state: ClientWSState | null = null;

    try {
      const url = new URL(req.url || '', `http://${req.headers.host}`);
      const token = url.searchParams.get('token');
      let clientId = url.searchParams.get('clientId');

      if (!token) {
        ws.close(1008, 'Missing token query parameter');
        return;
      }

      const keyInfo = await validateApiKeyDetailed(token);
      if (!keyInfo.valid) {
        ws.close(1008, 'Invalid API key');
        return;
      }

      // The API key used to match clients — master key for scoped keys
      const matchKey = keyInfo.masterApiKey || token;

      // Auto-resolve clientId if not provided
      if (!clientId) {
        const clients = await ClientManager.getConnectedClients(matchKey);
        if (clients.length === 1) {
          clientId = clients[0];
          log.info(`WS auto-resolved clientId: ${clientId}`);
        } else if (clients.length > 1) {
          // If scoped key has a scopedClientId, use that
          if (keyInfo.scopedClientId) {
            clientId = keyInfo.scopedClientId;
          } else {
            ws.close(1008, 'Multiple clients connected. Please specify clientId.');
            return;
          }
        }
      }

      // If still no clientId, try to auto-start a headless session for scoped keys with stored creds
      if (!clientId && keyInfo.masterApiKey) {
        try {
          // Look up the scoped key to check for stored credentials
          const scopedKey = await ApiKey.findByKey(token);
          if (scopedKey) {
            const get = (f: string) => scopedKey.getDataValue ? scopedKey.getDataValue(f) : (scopedKey as any)[f];
            const scopedKeyId = get('id');
            const hasCredentials = get('encryptedFoundryPassword') && get('foundryUrl') && get('foundryUsername');

            if (hasCredentials) {
              ws.send(JSON.stringify({ type: 'status', message: 'Starting headless Foundry session...' }));
              log.info(`WS: Auto-starting headless session for scoped key ${token.substring(0, 8)}...`);
              clientId = await startHeadlessWithStoredCredentials(scopedKeyId, keyInfo.masterApiKey);
            }
          }
        } catch (err) {
          const msg = err instanceof Error ? err.message : String(err);
          log.error(`WS: Failed to auto-start headless session: ${msg}`);
          ws.close(1008, `Failed to start headless session: ${msg}`);
          return;
        }
      }

      if (!clientId) {
        ws.close(1008, 'No clientId provided and no connected Foundry client found');
        return;
      }

      // Verify the clientId exists and belongs to this API key
      const foundryClient = await ClientManager.getClient(clientId);
      if (!foundryClient) {
        ws.close(1008, 'Invalid clientId — no connected Foundry instance with that ID');
        return;
      }
      if (foundryClient.getApiKey() !== matchKey) {
        ws.close(1008, 'API key does not match the specified clientId');
        return;
      }

      state = {
        apiKey: token,
        masterApiKey: keyInfo.masterApiKey || token,
        clientId,
        scopedUserId: keyInfo.scopedUserId,
        subscriptions: new Set(),
        pendingRequestIds: new Set(),
        ws,
      };

      log.info(`Client WS connected: clientId=${clientId}, apiKey=${token.substring(0, 8)}...`);

      // Send welcome message
      ws.send(JSON.stringify({
        type: 'connected',
        clientId,
        supportedTypes: [...PENDING_REQUEST_TYPES],
        eventChannels: ['chat-events', 'roll-events'],
      }));

      // Ping/pong keepalive
      const pingInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.ping();
        }
      }, PING_INTERVAL_MS);

      ws.on('message', async (raw: Buffer | string) => {
        if (!state) return;

        let msg: any;
        try {
          msg = JSON.parse(typeof raw === 'string' ? raw : raw.toString('utf8'));
        } catch {
          ws.send(JSON.stringify({ type: 'error', error: 'Invalid JSON' }));
          return;
        }

        // Track usage per message (use master key for user lookup)
        const usage = await trackWsApiUsage(state.masterApiKey);
        if (!usage.allowed) {
          ws.send(JSON.stringify({ type: 'error', error: usage.error, requestId: msg.requestId }));
          return;
        }

        await handleClientMessage(ws, state, msg);
      });

      ws.on('close', () => {
        clearInterval(pingInterval);
        if (state) {
          cleanupClientState(state);
          log.info(`Client WS disconnected: clientId=${state.clientId}`);
        }
      });

      ws.on('error', (err) => {
        clearInterval(pingInterval);
        log.error(`Client WS error: ${err}`);
        if (state) {
          cleanupClientState(state);
        }
      });
    } catch (error) {
      log.error(`Client WS connection error: ${error}`);
      try {
        ws.close(1011, 'Server error');
      } catch {
        // Ignore close errors
      }
    }
  });
}

async function handleClientMessage(ws: WebSocket, state: ClientWSState, msg: any): Promise<void> {
  const { type, requestId } = msg;

  if (!type) {
    ws.send(JSON.stringify({ type: 'error', error: 'Missing "type" field', requestId }));
    return;
  }

  // Handle subscribe/unsubscribe
  if (type === 'subscribe') {
    handleSubscribe(ws, state, msg);
    return;
  }
  if (type === 'unsubscribe') {
    handleUnsubscribe(ws, state, msg);
    return;
  }

  // Handle ping
  if (type === 'ping') {
    ws.send(JSON.stringify({ type: 'pong', requestId }));
    return;
  }

  // Handle sheet session messages
  if (type === 'sheet-session-start') {
    await handleSheetSessionStart(ws, state, msg);
    return;
  }
  if (type === 'sheet-input') {
    await handleSheetInput(ws, state, msg);
    return;
  }
  if (type === 'sheet-session-end') {
    await handleSheetSessionEnd(ws, state, msg);
    return;
  }

  // Validate message type is a known request type
  if (!PENDING_REQUEST_TYPES.includes(type as PendingRequestType)) {
    ws.send(JSON.stringify({
      type: 'error',
      error: `Unknown message type: "${type}". Supported types: ${PENDING_REQUEST_TYPES.join(', ')}`,
      requestId,
    }));
    return;
  }

  if (!requestId) {
    ws.send(JSON.stringify({ type: 'error', error: 'Missing "requestId" field for request messages' }));
    return;
  }

  // Get the Foundry client
  const foundryClient = await ClientManager.getClient(state.clientId);
  if (!foundryClient) {
    ws.send(JSON.stringify({
      type: `${type}-result`,
      requestId,
      error: 'Foundry client is no longer connected',
    }));
    return;
  }

  const internalRequestId = `ws_${type}_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;

  // Register pending request with wsCallback
  pendingRequests.set(internalRequestId, {
    type: type as PendingRequestType,
    clientId: state.clientId,
    timestamp: Date.now(),
    format: msg.format || 'json',
    wsCallback: (statusCode: number, data: any) => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          ...data,
          type: `${type}-result`,
          requestId, // Return client's requestId, not internal one
        }));
      }
      state.pendingRequestIds.delete(internalRequestId);
    },
  });
  state.pendingRequestIds.add(internalRequestId);

  // Build payload — strip type and requestId, forward everything else
  const { type: _type, requestId: _rid, format: _fmt, ...payload } = msg;

  // Inject scopedUserId if set — caller cannot override
  if (state.scopedUserId) {
    payload.userId = state.scopedUserId;
  }

  const sent = foundryClient.send({
    type,
    requestId: internalRequestId,
    ...payload,
    data: {
      ...payload.data,
    },
  });

  if (!sent) {
    pendingRequests.delete(internalRequestId);
    state.pendingRequestIds.delete(internalRequestId);
    ws.send(JSON.stringify({
      type: `${type}-result`,
      requestId,
      error: 'Failed to send request to Foundry client',
    }));
    return;
  }

  // Timeout
  setTimeout(() => {
    if (pendingRequests.has(internalRequestId)) {
      pendingRequests.delete(internalRequestId);
      state.pendingRequestIds.delete(internalRequestId);
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: `${type}-result`,
          requestId,
          error: 'Request timed out',
        }));
      }
    }
  }, 30000);
}

function handleSubscribe(ws: WebSocket, state: ClientWSState, msg: any): void {
  const { channel, filters = {}, requestId } = msg;

  if (channel !== 'chat-events' && channel !== 'roll-events') {
    ws.send(JSON.stringify({
      type: 'error',
      error: `Invalid channel: "${channel}". Supported: chat-events, roll-events`,
      requestId,
    }));
    return;
  }

  const connection: WSEventConnection = {
    ws,
    clientId: state.clientId,
    channel,
    filters: {
      speaker: filters.speaker,
      type: filters.type !== undefined ? Number(filters.type) : undefined,
      whisperOnly: filters.whisperOnly,
      userId: filters.userId,
    },
  };

  const added = addWSEventConnection(state.clientId, connection);
  if (!added) {
    ws.send(JSON.stringify({
      type: 'error',
      error: 'Maximum event subscriptions reached for this client',
      requestId,
    }));
    return;
  }

  state.subscriptions.add(connection);
  ws.send(JSON.stringify({ type: 'subscribed', channel, requestId }));
  log.info(`Client WS subscribed to ${channel} for clientId=${state.clientId}`);
}

function handleUnsubscribe(ws: WebSocket, state: ClientWSState, msg: any): void {
  const { channel, requestId } = msg;

  let removed = 0;
  for (const sub of state.subscriptions) {
    if (!channel || sub.channel === channel) {
      removeWSEventConnection(state.clientId, sub);
      state.subscriptions.delete(sub);
      removed++;
    }
  }

  ws.send(JSON.stringify({ type: 'unsubscribed', channel: channel || 'all', removed, requestId }));
  log.info(`Client WS unsubscribed from ${channel || 'all'} for clientId=${state.clientId} (${removed} removed)`);
}

function cleanupClientState(state: ClientWSState): void {
  // Clean up event subscriptions
  for (const sub of state.subscriptions) {
    removeWSEventConnection(state.clientId, sub);
  }
  state.subscriptions.clear();

  // Clean up pending requests
  for (const reqId of state.pendingRequestIds) {
    pendingRequests.delete(reqId);
  }
  state.pendingRequestIds.clear();

  // Clean up sheet sessions - notify Foundry to close them
  const sessionIds = SheetSessionManager.terminateSessionsForConsumer(state.ws);
  if (sessionIds.length > 0) {
    ClientManager.getClient(state.clientId).then(foundryClient => {
      if (foundryClient) {
        for (const sessionId of sessionIds) {
          foundryClient.send({ type: 'sheet-session-end', sessionId });
        }
      }
    }).catch(() => { /* ignore */ });
  }
}

async function handleSheetSessionStart(ws: WebSocket, state: ClientWSState, msg: any): Promise<void> {
  const foundryClient = await ClientManager.getClient(state.clientId);
  if (!foundryClient) {
    ws.send(JSON.stringify({ type: 'sheet-session-error', error: 'Foundry client is no longer connected' }));
    return;
  }

  const result = SheetSessionManager.createSession(
    state.clientId,
    state.apiKey,
    ws,
    { uuid: msg.uuid, quality: msg.quality, scale: msg.scale }
  );

  if ('error' in result) {
    ws.send(JSON.stringify({ type: 'sheet-session-error', error: result.error }));
    return;
  }

  const sent = foundryClient.send({
    type: 'sheet-session-start',
    sessionId: result.sessionId,
    uuid: msg.uuid,
    selected: msg.selected,
    actor: msg.actor,
    quality: msg.quality,
    scale: msg.scale,
    userId: state.scopedUserId || msg.userId,
  });

  if (!sent) {
    SheetSessionManager.endSession(result.sessionId);
    ws.send(JSON.stringify({ type: 'sheet-session-error', error: 'Failed to send session start to Foundry' }));
  }
}

async function handleSheetInput(ws: WebSocket, state: ClientWSState, msg: any): Promise<void> {
  const { sessionId } = msg;

  const session = SheetSessionManager.getSession(sessionId);
  if (!session || session.consumerWs !== ws) {
    ws.send(JSON.stringify({ type: 'sheet-session-error', sessionId, error: 'Invalid session' }));
    return;
  }

  SheetSessionManager.updateActivity(sessionId);

  const foundryClient = await ClientManager.getClient(state.clientId);
  if (!foundryClient) {
    ws.send(JSON.stringify({ type: 'sheet-session-error', sessionId, error: 'Foundry client disconnected' }));
    SheetSessionManager.endSession(sessionId);
    return;
  }

  foundryClient.send({
    type: 'sheet-input',
    sessionId,
    action: msg.action,
    x: msg.x,
    y: msg.y,
    button: msg.button,
    deltaX: msg.deltaX,
    deltaY: msg.deltaY,
    key: msg.key,
    code: msg.code,
    modifiers: msg.modifiers,
  });
}

async function handleSheetSessionEnd(ws: WebSocket, state: ClientWSState, msg: any): Promise<void> {
  const { sessionId } = msg;

  SheetSessionManager.endSession(sessionId);

  const foundryClient = await ClientManager.getClient(state.clientId);
  if (foundryClient) {
    foundryClient.send({ type: 'sheet-session-end', sessionId });
  }
}
