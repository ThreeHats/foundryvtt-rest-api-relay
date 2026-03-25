import { Router } from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';

export const searchRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Search entities
 * 
 * This endpoint allows searching for entities in the Foundry world based on a query string.
 * Requires Quick Insert module to be installed and enabled.
 * 
 * @route GET /search
 * @returns {object} Search results containing matching entities
 */
searchRouter.get("/search", ...commonMiddleware, createApiRoute({
    type: 'search',
    requiredParams: [
        { name: 'query', from: 'query', type: 'string' } // Search query string
    ],
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'filter', from: 'query', type: 'string' }, // Filter to apply (simple: filter="Actor", property-based: filter="key:value,key2:value2")
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));
