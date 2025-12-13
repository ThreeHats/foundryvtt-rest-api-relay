import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';
import { log } from '../../utils/logger';
import { validateScript } from './utility';

export const entityRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Get entity details
 *
 * This endpoint retrieves the details of a specific entity.
 * 
 * @route GET /get
 * @returns {object} Entity details object containing requested information
 */
entityRouter.get("/get", ...commonMiddleware, createApiRoute({
    type: 'entity',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'uuid', from: 'query', type: 'string' }, // UUID of the entity to retrieve (optional if selected=true)
        { name: 'selected', from: 'query', type: 'boolean' }, // Whether to get the selected entity
        { name: 'actor', from: 'query', type: 'boolean' } // Return the actor of specified entity
    ]
}));

/**
 * Create a new entity
 *
 * This endpoint creates a new entity in the Foundry world.
 * 
 * @route POST /create
 * @returns {object} Result of the entity creation operation
 */
entityRouter.post("/create", ...commonMiddleware, express.json(), createApiRoute({
    type: 'create',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'entityType', from: 'body', type: 'string' }, // Document type of entity to create (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist, ext.)
        { name: 'data', from: 'body', type: 'object' } // Data for the new entity
    ],
    optionalParams: [
        { name: 'folder', from: 'body', type: 'string' } // Optional folder UUID to place the new entity in
    ],
    validateParams: (params) => {
        if (params.entityType === "Macro") {
            if (!validateScript(params.data.command)) {
                log.warn(`Request for ${params.clientId} contains forbidden patterns in script`);
                return {
                    error: "Script contains forbidden patterns",
                    suggestion: "Ensure the script does not access localStorage, sessionStorage, or eval()"
                };
            }
        }
        return null;
    }
}));

/**
 * Update an existing entity
 *
 * This endpoint updates an existing entity in the Foundry world.
 *
 * @route PUT /update
 * @returns {object} Result of the entity update operation
 */
entityRouter.put("/update", ...commonMiddleware, express.json(), createApiRoute({
    type: 'update',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
        { name: 'data', from: 'body', type: 'object' } // Data to update the entity with
    ],
    optionalParams: [
        { name: 'uuid', from: 'query', type: 'string' }, // UUID of the entity to update (optional if selected=true)
        { name: 'selected', from: 'query', type: 'boolean' }, // Whether to update the selected entity
        { name: 'actor', from: 'query', type: 'boolean' } // Update the actor of selected entity when selected=true
    ]
}));

/**
 * Delete an entity
 * 
 * This endpoint deletes an entity from the Foundry world.
 * 
 * @route DELETE /delete
 * @returns {object} Result of the entity deletion operation
 */
entityRouter.delete("/delete", ...commonMiddleware, createApiRoute({
    type: 'delete',
    requiredParams: [
        { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
    ],
    optionalParams: [
        { name: 'uuid', from: 'query', type: 'string' }, // UUID of the entity to delete (optional if selected=true)
        { name: 'selected', from: 'query', type: 'boolean' } // Whether to delete the selected entity
    ]
}));

/**
 * Give an item to an entity
 *
 * This endpoint gives an item to a specified entity.
 * Optionally, removes the item from the giver.
 * 
 * @route POST /give
 * @returns {object} Result of the item giving operation
 */
entityRouter.post("/give", ...commonMiddleware, express.json(), createApiRoute({
type: 'give',
requiredParams: [
    { name: 'clientId', from: ['body', 'query'], type: 'string' } // Client ID for the Foundry world
],
optionalParams: [
    { name: 'fromUuid', from: 'body', type: 'string' }, // UUID of the entity giving the item
    { name: 'toUuid', from: 'body', type: 'string' }, // UUID of the entity receiving the item
    { name: 'selected', from: 'body', type: 'boolean' }, // Whether to give to the selected token's actor
    { name: 'itemUuid', from: 'body', type: 'string' }, // UUID of the item to give (optional if itemName provided)
    { name: 'itemName', from: 'body', type: 'string' }, // Name of the item to give (search with Quick Insert if UUID not provided)
    { name: 'quantity', from: 'body', type: 'number' } // Quantity of the item to give (negative values decrease quantity to 0)
]
}));

/**
 * Remove an item from an entity
 *
 * This endpoint removes an item from a specified entity.
 * @route POST /remove
 * @returns {object} Result of the item removal operation
 */
entityRouter.post("/remove", ...commonMiddleware, express.json(), createApiRoute({
type: 'remove',
requiredParams: [
    { name: 'clientId', from: ['body', 'query'], type: 'string' } // Client ID for the Foundry world
],
optionalParams: [
    { name: 'actorUuid', from: 'body', type: 'string' }, // UUID of the actor to remove the item from (optional if selected=true)
    { name: 'selected', from: 'body', type: 'boolean' }, // Whether to remove from the selected token's actor
    { name: 'itemUuid', from: 'body', type: 'string' }, // UUID of the item to remove
    { name: 'itemName', from: 'body', type: 'string' }, // Name of the item to remove (search with Quick Insert if UUID not provided)
    { name: 'quantity', from: 'body', type: 'number' } // Quantity of the item to remove
]
}));

/**
 * Decrease an attribute
 *
 * This endpoint decreases an attribute of a specified entity.
 * 
 * @route POST /decrease
 * @returns {object} Result of the attribute decrease operation
 */
entityRouter.post("/decrease", ...commonMiddleware, express.json(), createApiRoute({
  type: 'decrease',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
    { name: 'attribute', from: 'body', type: 'string' }, // The attribute data path to decrease (e.g., "system.attributes.hp.value")
    { name: 'amount', from: 'body', type: 'number' } // The amount to decrease the attribute by
  ],
  optionalParams: [
    { name: 'uuid', from: 'query', type: 'string' }, // UUID of the entity to decrease the attribute for (optional if selected=true)
    { name: 'selected', from: 'query', type: 'boolean' } // Whether to decrease the attribute for the selected entity
  ]
}));

/**
 * Increase an attribute
 *
 * This endpoint increases an attribute of a specified entity.
 * 
 * @route POST /increase
 * @returns {object} Result of the attribute increase operation
 */
entityRouter.post("/increase", ...commonMiddleware, express.json(), createApiRoute({
  type: 'increase',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' }, // Client ID for the Foundry world
    { name: 'attribute', from: 'body', type: 'string' }, // The attribute data path to increase (e.g., "system.attributes.hp.value")
    { name: 'amount', from: 'body', type: 'number' } // The amount to increase the attribute by
  ],
  optionalParams: [
    { name: 'uuid', from: 'query', type: 'string' }, // UUID of the entity to increase the attribute for (optional if selected=true)
    { name: 'selected', from: 'query', type: 'boolean' } // Whether to increase the attribute for the selected entity
  ]
}));

/**
 * Kill an entity
 *
 * Marks an entity as killed in the combat tracker,
 * gives it the "dead" status,
 * and sets its health to 0 in the Foundry world.
 * 
 * @route POST /kill
 * @returns {object} Result of the entity kill operation
 */
entityRouter.post("/kill", ...commonMiddleware, express.json(), createApiRoute({
  type: 'kill',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // Client ID for the Foundry world
  ],
  optionalParams: [
    { name: 'uuid', from: 'query', type: 'string' }, // UUID of the entity to kill (optional if selected=true)
    { name: 'selected', from: 'query', type: 'boolean' } // Whether to kill the selected entity
  ]
}));