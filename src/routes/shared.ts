import express, { Response } from 'express';
import { log } from '../utils/logger';

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
    'chat-messages', 'chat-send', 'chat-delete', 'chat-flush'
] as const;
  
export type PendingRequestType = typeof PENDING_REQUEST_TYPES[number];

export interface PendingRequest {
    res: express.Response;
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