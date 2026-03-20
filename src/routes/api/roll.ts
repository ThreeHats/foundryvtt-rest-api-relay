import { Router, Request, Response } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';
import { addRollSSEConnection, removeRollSSEConnection, RollSSEConnection } from '../shared';
import { ClientManager } from '../../core/ClientManager';
import { log } from '../../utils/logger';

export const rollRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage, express.json()];

/**
 * Get recent rolls
 * 
 * Retrieves a list of up to 20 recent rolls made in the Foundry world.
 * 
 * @route GET /rolls
 * @returns {object} An array of recent rolls with details
 */
rollRouter.get("/rolls", ...commonMiddleware, createApiRoute({
    type: 'rolls',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'limit', from: 'query', type: 'number' }, // Optional limit on the number of rolls to return (default is 20)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Get the last roll
 * 
 * Retrieves the most recent roll made in the Foundry world.
 * 
 * @route GET /lastroll
 * @returns {object} The most recent roll with details
 */
rollRouter.get("/lastroll", ...commonMiddleware, createApiRoute({
    type: 'last-roll',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Make a roll
 * 
 * Executes a roll with the specified formula
 * 
 * @route POST /roll
 * @returns {object} Result of the roll operation
 */
rollRouter.post("/roll", ...commonMiddleware, createApiRoute({
    type: 'roll',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'formula', from: 'body', type: 'string' } // The roll formula to evaluate (e.g., "1d20 + 5")
    ],
    optionalParams: [
        { name: 'flavor', from: 'body', type: 'string' }, // Optional flavor text for the roll
        { name: 'createChatMessage', from: 'body', type: 'boolean' }, // Whether to create a chat message for the roll
        { name: 'speaker', from: 'body', type: 'string' }, // The speaker for the roll
        { name: 'whisper', from: 'body', type: 'array' }, // Users to whisper the roll result to
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Subscribe to real-time roll events via Server-Sent Events (SSE)
 *
 * Opens a persistent SSE connection that streams roll events as dice rolls
 * occur in the Foundry world. Each event includes the full roll details
 * including formula, total, individual dice results, and critical/fumble status.
 *
 * @route GET /rolls/subscribe
 * @param {string} clientId - [query] Client ID for the Foundry world
 * @param {string} userId - [query,?] Foundry user ID or username to scope permissions (omit for GM-level access)
 * @returns {stream} SSE event stream
 */
rollRouter.get("/rolls/subscribe", requestForwarderMiddleware, authMiddleware, trackApiUsage, async (req: Request, res: Response) => {
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

    const filters: RollSSEConnection['filters'] = {};
    if (req.query.userId) filters.userId = req.query.userId as string;

    const connection: RollSSEConnection = { res, clientId, filters };

    const added = addRollSSEConnection(clientId, connection);
    if (!added) {
        res.write(`event: error\ndata: ${JSON.stringify({ error: "Too many SSE connections for this client" })}\n\n`);
        res.end();
        return;
    }

    log.info(`Roll SSE connection opened for client ${clientId}`);

    // Send initial connection event
    res.write(`event: connected\ndata: ${JSON.stringify({ clientId })}\n\n`);

    // Keepalive every 15 seconds
    const keepaliveInterval = setInterval(() => {
        res.write(`: keepalive\n\n`);
    }, 15000);

    // Cleanup on close
    req.on('close', () => {
        clearInterval(keepaliveInterval);
        removeRollSSEConnection(clientId, connection);
        log.info(`Roll SSE connection closed for client ${clientId}`);
    });
});