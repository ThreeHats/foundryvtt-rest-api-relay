/**
 * Active Effects API routes (system-agnostic).
 */
import express, { Router } from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';

export const effectsRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage, express.json()];

/**
 * Get all active effects on an actor or token.
 *
 * Returns the collection of ActiveEffect documents currently applied
 * to the specified actor or token.
 *
 * @route GET /effects
 * @returns {object} Array of active effects
 */
effectsRouter.get("/effects", ...commonMiddleware, createApiRoute({
    type: 'get-effects',
    requiredParams: [
        { name: 'uuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor or token to query
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Add an active effect to an actor or token.
 *
 * Adds a status condition (by statusId) or a custom ActiveEffect
 * (via effectData) to the specified actor or token.
 *
 * @route POST /effects
 * @returns {object} Result of the add operation
 */
effectsRouter.post("/effects", ...commonMiddleware, createApiRoute({
    type: 'add-effect',
    requiredParams: [
        { name: 'uuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor or token to add the effect to
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'statusId', from: ['body', 'query'], type: 'string' }, // Standard status condition ID (e.g., "poisoned", "blinded", "prone")
        { name: 'effectData', from: ['body', 'query'], type: 'object' }, // Custom ActiveEffect data object (name, icon, duration, changes)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Remove an active effect from an actor or token.
 *
 * Removes an effect by its document ID (effectId) or by status condition
 * identifier (statusId).
 *
 * @route DELETE /effects
 * @returns {object} Result of the remove operation
 */
effectsRouter.delete("/effects", ...commonMiddleware, createApiRoute({
    type: 'remove-effect',
    requiredParams: [
        { name: 'uuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor or token to remove the effect from
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'effectId', from: ['body', 'query'], type: 'string' }, // The ActiveEffect document ID to remove
        { name: 'statusId', from: ['body', 'query'], type: 'string' }, // Standard status condition ID to remove (e.g., "poisoned")
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));
