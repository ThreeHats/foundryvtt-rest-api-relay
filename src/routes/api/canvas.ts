import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';
import { safeResponse } from '../shared';

export const canvasRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

const VALID_DOCUMENT_TYPES = ['tokens', 'tiles', 'drawings', 'lights', 'sounds', 'notes', 'templates', 'walls'] as const;

const DOCUMENT_TYPE_TO_CLASS: Record<string, string> = {
    tokens: 'Token',
    tiles: 'Tile',
    drawings: 'Drawing',
    lights: 'AmbientLight',
    sounds: 'AmbientSound',
    notes: 'Note',
    templates: 'MeasuredTemplate',
    walls: 'Wall'
};

function validateDocumentType(params: Record<string, any>) {
    const dt = params.documentType;
    if (!VALID_DOCUMENT_TYPES.includes(dt)) {
        return {
            error: `Invalid documentType '${dt}'. Must be one of: ${VALID_DOCUMENT_TYPES.join(', ')}`
        };
    }
    return null;
}

/**
 * Get canvas embedded documents
 *
 * @route GET /canvas/:documentType
 * @returns {object} Array of embedded documents
 */
canvasRouter.get("/canvas/:documentType", ...commonMiddleware, createApiRoute({
    type: 'get-canvas-documents',
    requiredParams: [
        { name: 'documentType', from: 'params', type: 'string' } // Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)
    ],
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'sceneId', from: 'query', type: 'string' }, // Scene ID to query (defaults to the active scene)
        { name: 'documentId', from: 'query', type: 'string' }, // Specific document ID to retrieve
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: validateDocumentType,
    buildPayload: (params) => ({
        ...params,
        className: DOCUMENT_TYPE_TO_CLASS[params.documentType]
    })
}));

/**
 * Create canvas embedded document(s)
 *
 * @route POST /canvas/:documentType
 * @returns {object} Created document(s)
 */
canvasRouter.post("/canvas/:documentType", ...commonMiddleware, express.json(), createApiRoute({
    type: 'create-canvas-document',
    requiredParams: [
        { name: 'documentType', from: 'params', type: 'string' }, // Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)
        { name: 'data', from: 'body' } // Document data object or array of objects to create
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'sceneId', from: ['query', 'body'], type: 'string' }, // Scene ID to create in (defaults to the active scene)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: validateDocumentType,
    buildPayload: (params) => ({
        ...params,
        className: DOCUMENT_TYPE_TO_CLASS[params.documentType]
    }),
    timeout: 30000
}));

/**
 * Update a canvas embedded document
 *
 * @route PUT /canvas/:documentType
 * @returns {object} Updated document
 */
canvasRouter.put("/canvas/:documentType", ...commonMiddleware, express.json(), createApiRoute({
    type: 'update-canvas-document',
    requiredParams: [
        { name: 'documentType', from: 'params', type: 'string' }, // Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)
        { name: 'documentId', from: ['body', 'query'], type: 'string' }, // ID of the document to update
        { name: 'data', from: 'body', type: 'object' } // Object containing the fields to update
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'sceneId', from: ['query', 'body'], type: 'string' }, // Scene ID containing the document (defaults to the active scene)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: validateDocumentType,
    buildPayload: (params) => ({
        ...params,
        className: DOCUMENT_TYPE_TO_CLASS[params.documentType]
    }),
    timeout: 30000
}));

/**
 * Delete a canvas embedded document
 *
 * @route DELETE /canvas/:documentType
 * @returns {object} Deletion result
 */
canvasRouter.delete("/canvas/:documentType", ...commonMiddleware, createApiRoute({
    type: 'delete-canvas-document',
    requiredParams: [
        { name: 'documentType', from: 'params', type: 'string' }, // Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)
        { name: 'documentId', from: 'query', type: 'string' } // ID of the document to delete
    ],
    optionalParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'sceneId', from: 'query', type: 'string' }, // Scene ID containing the document (defaults to the active scene)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ],
    validateParams: validateDocumentType,
    buildPayload: (params) => ({
        ...params,
        className: DOCUMENT_TYPE_TO_CLASS[params.documentType]
    })
}));
