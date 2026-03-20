import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';

export const sceneRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Get scene(s)
 *
 * Retrieves one or more scenes by ID, name, active status, viewed status, or all.
 *
 * @route GET /scene
 * @returns {object} Scene data
 */
sceneRouter.get("/scene", ...commonMiddleware, createApiRoute({
    type: 'get-scene',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'sceneId', from: 'query', type: 'string' }, // ID of a specific scene to retrieve
        { name: 'name', from: 'query', type: 'string' }, // Name of the scene to retrieve
        { name: 'active', from: 'query', type: 'boolean' }, // Set to true to get the currently active scene
        { name: 'viewed', from: 'query', type: 'boolean' }, // Set to true to get the currently viewed scene
        { name: 'all', from: 'query', type: 'boolean' }, // Set to true to get all scenes
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: (params) => {
        if (!params.sceneId && !params.name && !params.active && !params.viewed && !params.all) {
            return { error: "At least one of 'sceneId', 'name', 'active', 'viewed', or 'all' is required" };
        }
        return null;
    }
}));

/**
 * Create a new scene
 *
 * @route POST /scene
 * @returns {object} Created scene data
 */
sceneRouter.post("/scene", ...commonMiddleware, express.json(), createApiRoute({
    type: 'create-scene',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'data', from: 'body', type: 'object' } // Scene data object (name, width, height, grid, etc.)
    ],
    optionalParams: [
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    timeout: 30000
}));

/**
 * Update an existing scene
 *
 * @route PUT /scene
 * @returns {object} Updated scene data
 */
sceneRouter.put("/scene", ...commonMiddleware, express.json(), createApiRoute({
    type: 'update-scene',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'data', from: 'body', type: 'object' } // Object containing the scene fields to update
    ],
    optionalParams: [
        { name: 'sceneId', from: ['query', 'body'], type: 'string' }, // ID of the scene to update
        { name: 'name', from: ['query', 'body'], type: 'string' }, // Name of the scene to update (alternative to sceneId)
        { name: 'active', from: ['query', 'body'], type: 'boolean' }, // Set to true to target the active scene
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: (params) => {
        if (!params.sceneId && !params.name && !params.active) {
            return { error: "At least one of 'sceneId', 'name', or 'active' is required to identify the scene" };
        }
        return null;
    },
    timeout: 30000
}));

/**
 * Delete a scene
 *
 * @route DELETE /scene
 * @returns {object} Deletion result
 */
sceneRouter.delete("/scene", ...commonMiddleware, createApiRoute({
    type: 'delete-scene',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'sceneId', from: 'query', type: 'string' }, // ID of the scene to delete
        { name: 'name', from: 'query', type: 'string' }, // Name of the scene to delete (alternative to sceneId)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: (params) => {
        if (!params.sceneId && !params.name) {
            return { error: "Either 'sceneId' or 'name' is required" };
        }
        return null;
    }
}));

/**
 * Switch the active scene
 *
 * @route POST /switch-scene
 * @returns {object} Result of the scene switch
 */
sceneRouter.post("/switch-scene", ...commonMiddleware, express.json(), createApiRoute({
    type: 'switch-scene',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'sceneId', from: ['body', 'query'], type: 'string' }, // ID of the scene to activate
        { name: 'name', from: ['body', 'query'], type: 'string' }, // Name of the scene to activate (alternative to sceneId)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: (params) => {
        if (!params.sceneId && !params.name) {
            return { error: "Either 'sceneId' or 'name' is required" };
        }
        return null;
    },
    timeout: 30000
}));
