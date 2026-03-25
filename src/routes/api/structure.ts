import { Router } from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';
export const structureRouter = Router();
const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Get the structure of the Foundry world
 * 
 * Retrieves the folder and compendium structure for the specified Foundry world.
 * 
 * @route GET /structure
 * @returns {object} The folder and compendium structure
 */
structureRouter.get("/structure", ...commonMiddleware, createApiRoute({
    type: 'structure',
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'includeEntityData', from: 'query', type: 'boolean' }, // Whether to include full entity data or just UUIDs and names
        { name: 'path', from: 'query', type: 'string' }, // Path to read structure from (null = root)
        { name: 'recursive', from: 'query', type: 'boolean' }, // Whether to read down the folder tree
        { name: 'recursiveDepth', from: 'query', type: 'number' }, // Depth to recurse into folders (default 5)
        { name: 'types', from: 'query', type: 'string' }, // Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    buildPayload: (params) => {
        // Handle types parameter - can be comma-separated string or JSON array
        if (params.types && typeof params.types === 'string') {
            try {
                // Try to parse as JSON first
                params.types = JSON.parse(params.types);
            } catch {
                // If not JSON, split by comma
                params.types = params.types.split(',').map((t: string) => t.trim());
            }
        }
        return params;
    }
}));

/**
 * This route is deprecated - use /structure with the path query parameter instead
 * 
 * @route GET /contents/:path
 * @returns {object} Error message directing to use /structure endpoint
 */
structureRouter.get("/contents/:path", (req, res) => {
    res.status(400).json({
        error: "This endpoint is deprecated",
        message: "Please use GET /structure with the 'path' query parameter instead",
        example: `/structure?clientId=${req.query.clientId}&path=${req.params.path}`
    });
});

/**
 * Get a specific folder by name
 * 
 * @route GET /get-folder
 * @returns {object} The folder information and its contents
 */
structureRouter.get("/get-folder", ...commonMiddleware, createApiRoute({
    type: 'get-folder',
    requiredParams: [
        { name: 'name', from: ['body', 'query'], type: 'string' } // Name of the folder to retrieve
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Create a new folder
 * 
 * @route POST /create-folder
 * @returns {object} The created folder information
 */
structureRouter.post("/create-folder", ...commonMiddleware, createApiRoute({
    type: 'create-folder',
    requiredParams: [
        { name: 'name', from: ['body', 'query'], type: 'string' }, // Name of the new folder
        { name: 'folderType', from: ['body', 'query'], type: 'string' } // Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist)
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'parentFolderId', from: ['body', 'query'], type: 'string' }, // ID of the parent folder (optional for root level)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Delete a folder
 * 
 * @route DELETE /delete-folder
 * @returns {object} Confirmation of deletion
 */
structureRouter.delete("/delete-folder", ...commonMiddleware, createApiRoute({
    type: 'delete-folder',
    requiredParams: [
        { name: 'folderId', from: ['body', 'query'], type: 'string' } // ID of the folder to delete
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'deleteAll', from: ['body', 'query'], type: 'boolean' }, // Whether to delete all entities in the folder
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));