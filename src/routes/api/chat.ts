import { Router, Request, Response } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';
import { addSSEConnection, removeSSEConnection, SSEConnection } from '../shared';
import { ClientManager } from '../../core/ClientManager';
import { log } from '../../utils/logger';

export const chatRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage, express.json()];

/**
 * Get chat messages
 *
 * Retrieves chat messages from the Foundry world with optional pagination and filtering.
 *
 * @route GET /chat
 * @returns {object} Paginated list of chat messages
 */
chatRouter.get("/chat", ...commonMiddleware, createApiRoute({
    type: 'chat-messages',
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'limit', from: 'query', type: 'number' }, // Maximum number of messages to return (default: 10)
        { name: 'offset', from: 'query', type: 'number' }, // Number of messages to skip for pagination
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
        { name: 'chatType', from: 'query', type: 'number' }, // Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field.
        { name: 'speaker', from: 'query', type: 'string' } // Filter messages by speaker name or actor ID
    ]
}));

/**
 * Send a chat message
 *
 * Creates a new chat message in the Foundry world.
 *
 * @route POST /chat
 * @returns {object} The created chat message
 */
chatRouter.post("/chat", ...commonMiddleware, createApiRoute({
    type: 'chat-send',
    requiredParams: [
        { name: 'content', from: 'body', type: 'string' } // The message content (supports HTML)
    ],
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'whisper', from: 'body', type: 'array' }, // Array of user IDs to whisper the message to
        { name: 'speaker', from: 'body', type: 'string' }, // Actor ID to use as the message speaker
        { name: 'alias', from: 'body', type: 'string' }, // Display name alias for the speaker
        { name: 'chatType', from: 'body', type: 'number' }, // Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field.
        { name: 'flavor', from: 'body', type: 'string' }, // Flavor text displayed above the message content
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Delete a specific chat message
 *
 * Deletes a chat message by its ID. Only the message author or a GM can delete messages.
 *
 * @route DELETE /chat/:messageId
 * @returns {object} Success confirmation
 */
chatRouter.delete("/chat/:messageId", ...commonMiddleware, createApiRoute({
    type: 'chat-delete',
    requiredParams: [
        { name: 'messageId', from: 'params', type: 'string' } // ID of the chat message to delete
    ],
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Clear all chat messages
 *
 * Flushes all chat message history. Only GMs can perform this action.
 *
 * @route DELETE /chat
 * @returns {object} Success confirmation
 */
chatRouter.delete("/chat", ...commonMiddleware, createApiRoute({
    type: 'chat-flush',
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Subscribe to real-time chat events via Server-Sent Events (SSE)
 *
 * Opens a persistent SSE connection that streams chat events (create, update, delete)
 * as they occur in the Foundry world.
 *
 * @route GET /chat/subscribe
 * @param {string} clientId - [query] Client ID for the Foundry world
 * @param {string} speaker - [query,?] Filter events by speaker name or actor ID
 * @param {number} type - [query,?] Filter events by chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll)
 * @param {boolean} whisperOnly - [query,?] Only receive whispered messages
 * @param {string} userId - [query,?] Foundry user ID or username to scope permissions (omit for GM-level access)
 * @returns {stream} SSE event stream
 */
chatRouter.get("/chat/subscribe", requestForwarderMiddleware, authMiddleware, trackApiUsage, async (req: Request, res: Response) => {
    const clientId = req.query.clientId as string;
    if (!clientId) {
        res.status(400).json({ error: "'clientId' is required" });
        return;
    }

    const client = await ClientManager.getClient(clientId);
    if (!client) {
        res.status(404).json({ error: "Invalid client ID" });
        return;
    }

    // Set SSE headers
    res.setHeader('Content-Type', 'text/event-stream');
    res.setHeader('Cache-Control', 'no-cache');
    res.setHeader('Connection', 'keep-alive');
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.flushHeaders();

    const filters: SSEConnection['filters'] = {};
    if (req.query.speaker) filters.speaker = req.query.speaker as string;
    if (req.query.type) filters.type = parseInt(req.query.type as string, 10);
    if (req.query.whisperOnly) filters.whisperOnly = req.query.whisperOnly === 'true';
    if (req.query.userId) filters.userId = req.query.userId as string;

    const connection: SSEConnection = { res, clientId, filters };

    const added = addSSEConnection(clientId, connection);
    if (!added) {
        res.write(`event: error\ndata: ${JSON.stringify({ error: "Too many SSE connections for this client" })}\n\n`);
        res.end();
        return;
    }

    log.info(`SSE connection opened for client ${clientId}`);

    // Send initial connection event
    res.write(`event: connected\ndata: ${JSON.stringify({ clientId })}\n\n`);

    // Keepalive every 15 seconds
    const keepaliveInterval = setInterval(() => {
        res.write(`: keepalive\n\n`);
    }, 15000);

    // Cleanup on close
    req.on('close', () => {
        clearInterval(keepaliveInterval);
        removeSSEConnection(clientId, connection);
        log.info(`SSE connection closed for client ${clientId}`);
    });
});
