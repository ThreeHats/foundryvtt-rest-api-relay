import express, { Request, Response } from 'express';
import WebSocket from 'ws';
import { log } from '../utils/logger';
import { ClientManager } from '../core/ClientManager';

// Extracted from api.ts
export function sanitizeResponse(response: any): any {
    if (response === null || response === undefined) {
      return response;
    }
    
    if (typeof response !== 'object') {
      return response;
    }
    
    // Custom deep clone and key removal
    function removeSensitiveKeys(obj: any): any {
      if (obj === null || typeof obj !== 'object') {
        return obj;
      }
      
      if (Array.isArray(obj)) {
        return obj.map(item => removeSensitiveKeys(item));
      }
      
      const newObj: any = {};
      for (const key in obj) {
        if (key !== 'privateKey' && key !== 'apiKey' && key !== 'password') {
          newObj[key] = removeSensitiveKeys(obj[key]);
        }
      }
      return newObj;
    }
    
    return removeSensitiveKeys(response);
}
  
export function safeResponse(res: Response, statusCode: number, data: any): void {
    if (res.headersSent) {
      log.warn(`Headers already sent for request. Cannot send response:`, data);
      return;
    }
    const sanitizedData = sanitizeResponse(data);
    res.status(statusCode).json(sanitizedData);
}

export const PENDING_REQUEST_TYPES = [
    'search', 'entity', 'structure', 'contents', 'create', 'update', 'delete',
    'rolls', 'last-roll', 'roll', 'get-sheet', 'macro-execute', 'macros',
    'encounters', 'start-encounter', 'next-turn', 'next-round', 'last-turn', 'last-round',
    'end-encounter', 'add-to-encounter', 'remove-from-encounter', 'kill', 'decrease', 'increase', 'give', 'remove', 'execute-js',
    'select', 'selected', 'file-system', 'upload-file', 'download-file',
    'get-actor-details', 'modify-item-charges', 'use-ability', 'use-feature', 'use-spell', 'use-item', 'modify-experience', 'add-item', 'remove-item',
    'get-folder', 'create-folder', 'delete-folder',
    'players',
    'get-scene', 'create-scene', 'update-scene', 'delete-scene', 'switch-scene',
    'get-canvas-documents', 'create-canvas-document', 'update-canvas-document', 'delete-canvas-document',
    'chat-messages', 'chat-send', 'chat-delete', 'chat-flush',
    'short-rest', 'long-rest', 'skill-check', 'ability-save', 'ability-check', 'death-save',
    'get-effects', 'add-effect', 'remove-effect',
    'sheet-screenshot'
] as const;
  
export type PendingRequestType = typeof PENDING_REQUEST_TYPES[number];

export interface PendingRequest {
    res?: express.Response;
    wsCallback?: (statusCode: number, data: any) => void;
    type: PendingRequestType;
    clientId?: string;
    uuid?: string;
    path?: string;
    query?: string;
    filter?: string;
    timestamp: number;
    format?: string;
    initialScale?: number | null;
    activeTab?: number | null;
    darkMode?: boolean;
}

export const pendingRequests = new Map<string, PendingRequest>();

/**
 * Resolve a pending request via either HTTP response or WS callback.
 * Replaces direct safeResponse(pending.res, ...) calls in message handlers.
 */
export function resolvePendingRequest(requestId: string, statusCode: number, data: any): void {
    const pending = pendingRequests.get(requestId);
    if (!pending) {
        log.warn(`resolvePendingRequest: no pending request for ${requestId}`);
        return;
    }

    if (pending.wsCallback) {
        const sanitizedData = sanitizeResponse(data);
        pending.wsCallback(statusCode, sanitizedData);
    } else if (pending.res) {
        safeResponse(pending.res, statusCode, data);
    } else {
        log.warn(`resolvePendingRequest: pending request ${requestId} has neither res nor wsCallback`);
    }

    pendingRequests.delete(requestId);
}

// SSE connection tracking for chat events
export interface SSEConnection {
    res: express.Response;
    clientId: string;
    filters: {
        speaker?: string;
        type?: number;
        whisperOnly?: boolean;
        userId?: string;
    };
}

const sseConnections = new Map<string, Set<SSEConnection>>();
const MAX_SSE_CONNECTIONS_PER_CLIENT = 10;

export function addSSEConnection(clientId: string, connection: SSEConnection): boolean {
    if (!sseConnections.has(clientId)) {
        sseConnections.set(clientId, new Set());
    }
    const connections = sseConnections.get(clientId)!;
    if (connections.size >= MAX_SSE_CONNECTIONS_PER_CLIENT) {
        return false;
    }
    connections.add(connection);
    return true;
}

export function removeSSEConnection(clientId: string, connection: SSEConnection): void {
    const connections = sseConnections.get(clientId);
    if (connections) {
        connections.delete(connection);
        if (connections.size === 0) {
            sseConnections.delete(clientId);
        }
    }
}

export function getSSEConnections(clientId: string): Set<SSEConnection> | undefined {
    return sseConnections.get(clientId);
}

// SSE connection tracking for roll events
export interface RollSSEConnection {
    res: express.Response;
    clientId: string;
    filters: {
        userId?: string;
    };
}

const rollSSEConnections = new Map<string, Set<RollSSEConnection>>();

export function addRollSSEConnection(clientId: string, connection: RollSSEConnection): boolean {
    if (!rollSSEConnections.has(clientId)) {
        rollSSEConnections.set(clientId, new Set());
    }
    const connections = rollSSEConnections.get(clientId)!;
    if (connections.size >= MAX_SSE_CONNECTIONS_PER_CLIENT) {
        return false;
    }
    connections.add(connection);
    return true;
}

export function removeRollSSEConnection(clientId: string, connection: RollSSEConnection): void {
    const connections = rollSSEConnections.get(clientId);
    if (connections) {
        connections.delete(connection);
        if (connections.size === 0) {
            rollSSEConnections.delete(clientId);
        }
    }
}

export function getRollSSEConnections(clientId: string): Set<RollSSEConnection> | undefined {
    return rollSSEConnections.get(clientId);
}

// WebSocket event connection tracking (parallel to SSE)
export interface WSEventConnection {
    ws: WebSocket;
    clientId: string;
    channel: 'chat-events' | 'roll-events';
    filters: {
        speaker?: string;
        type?: number;
        whisperOnly?: boolean;
        userId?: string;
    };
}

const wsEventConnections = new Map<string, Set<WSEventConnection>>();
const MAX_WS_EVENT_CONNECTIONS_PER_CLIENT = 10;

export function addWSEventConnection(clientId: string, connection: WSEventConnection): boolean {
    if (!wsEventConnections.has(clientId)) {
        wsEventConnections.set(clientId, new Set());
    }
    const connections = wsEventConnections.get(clientId)!;
    if (connections.size >= MAX_WS_EVENT_CONNECTIONS_PER_CLIENT) {
        return false;
    }
    connections.add(connection);
    return true;
}

export function removeWSEventConnection(clientId: string, connection: WSEventConnection): void {
    const connections = wsEventConnections.get(clientId);
    if (connections) {
        connections.delete(connection);
        if (connections.size === 0) {
            wsEventConnections.delete(clientId);
        }
    }
}

export function getWSEventConnections(clientId: string): Set<WSEventConnection> | undefined {
    return wsEventConnections.get(clientId);
}

/**
 * Resolve clientId for manual route handlers (non-createApiRoute endpoints).
 * Applies scoped key enforcement and auto-resolution.
 * @returns resolved clientId or null if an error response was sent
 */
export async function resolveClientId(req: Request, res: Response, rawClientId?: string): Promise<string | null> {
  let clientId = rawClientId;

  // Scoped key enforcement
  if (req.scopedKey?.scopedClientId) {
    clientId = req.scopedKey.scopedClientId;
  }

  // Auto-resolve if missing
  if (!clientId) {
    const masterKey = req.masterApiKey || req.user?.apiKey ||
      (req.user?.getDataValue ? req.user.getDataValue('apiKey') : undefined);
    if (masterKey) {
      const clients = await ClientManager.getConnectedClients(masterKey);
      if (clients.length === 1) {
        clientId = clients[0];
      } else if (clients.length === 0) {
        // No clients connected — try auto-starting a headless session
        const autoClientId = await tryAutoStartForScopedKey(req);
        if (autoClientId) {
          clientId = autoClientId;
        } else {
          safeResponse(res, 404, { error: 'No connected Foundry clients found' });
          return null;
        }
      } else {
        safeResponse(res, 400, {
          error: 'Multiple clients connected. Please specify clientId.',
          connectedClients: clients
        });
        return null;
      }
    } else {
      safeResponse(res, 400, { error: 'clientId is required' });
      return null;
    }
  }

  // Verify the resolved client actually exists; if not, try auto-start for scoped keys
  const client = await ClientManager.getClient(clientId);
  if (!client && req.scopedKey) {
    const autoClientId = await tryAutoStartForScopedKey(req);
    if (autoClientId) return autoClientId;
  }

  return clientId;
}

/**
 * Try to auto-start a headless session for a scoped key with stored credentials.
 * Uses lazy imports to avoid circular dependencies.
 */
export async function tryAutoStartForScopedKey(req: Request): Promise<string | null> {
  if (!req.scopedKey?.id || !req.masterApiKey) return null;

  try {
    const { ApiKey } = await import('../models/apiKey');
    const scopedKeyRecord = await ApiKey.findOne({ where: { id: req.scopedKey.id } });
    if (!scopedKeyRecord) return null;

    const get = (f: string) => scopedKeyRecord.getDataValue ? scopedKeyRecord.getDataValue(f) : (scopedKeyRecord as any)[f];
    const hasCredentials = get('encryptedFoundryPassword') && get('foundryUrl') && get('foundryUsername');
    if (!hasCredentials) return null;

    log.info(`Auto-starting headless session for scoped key ${req.scopedKey.id}`);
    const { startHeadlessWithStoredCredentials } = await import('../utils/headlessSessionStarter');
    return await startHeadlessWithStoredCredentials(req.scopedKey.id, req.masterApiKey);
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);
    log.error(`Auto-start headless session failed: ${msg}`);
    return null;
  }
}

/**
 * Resolve userId for manual route handlers.
 * Applies scoped key userId enforcement.
 */
export function resolveScopedUserId(req: Request, rawUserId?: string): string | undefined {
  if (req.scopedKey?.scopedUserId) {
    return req.scopedKey.scopedUserId;
  }
  return rawUserId;
}