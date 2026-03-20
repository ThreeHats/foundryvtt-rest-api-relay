import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';
export const encounterRouter = Router();
const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Get all active encounters
 * 
 * Retrieves a list of all currently active encounters in the Foundry world.
 * 
 * @route GET /encounters
 * @returns {object} An array of active encounters with details
 */
encounterRouter.get("/encounters", ...commonMiddleware, createApiRoute({
  type: 'encounters',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Start a new encounter
 * 
 * Initiates a new encounter in the Foundry world.
 * 
 * @route POST /start-encounter
 * @returns {object} Details of the started encounter
 */
encounterRouter.post("/start-encounter", ...commonMiddleware, express.json(), createApiRoute({
  type: 'start-encounter',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'tokens', from: 'body', type: 'array' }, // Array of token UUIDs to include in the encounter
    { name: 'startWithSelected', from: 'body', type: 'boolean' }, // Whether to start with selected tokens
    { name: 'startWithPlayers', from: 'body', type: 'boolean' }, // Whether to start with players
    { name: 'rollNPC', from: 'body', type: 'boolean' }, // Whether to roll for NPCs
    { name: 'rollAll', from: 'body', type: 'boolean' }, // Whether to roll for all tokens
    { name: 'name', from: 'body', type: 'string' }, // The name of the encounter (unused)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Advance to the next turn in the encounter
 * 
 * Moves the encounter to the next turn.
 * 
 * @route POST /next-turn
 * @returns {object} Details of the next turn
 */
encounterRouter.post("/next-turn", ...commonMiddleware, express.json(), createApiRoute({
  type: 'next-turn',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to advance (optional, defaults to current encounter)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Advance to the next round in the encounter
 * 
 * Moves the encounter to the next round.
 * 
 * @route POST /next-round
 * @returns {object} Details of the next round
 */
encounterRouter.post("/next-round", ...commonMiddleware, express.json(), createApiRoute({
  type: 'next-round',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to advance (optional, defaults to current encounter)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Advance to the last turn in the encounter
 * 
 * Moves the encounter to the last turn.
 * 
 * @route POST /last-turn
 * @returns {object} Details of the last turn
 */
encounterRouter.post("/last-turn", ...commonMiddleware, express.json(), createApiRoute({
  type: 'last-turn',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to advance (optional, defaults to current encounter)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Advance to the last round in the encounter
 * 
 * Moves the encounter to the last round.
 * 
 * @route POST /last-round
 * @returns {object} Details of the last round
 */
encounterRouter.post("/last-round", ...commonMiddleware, express.json(), createApiRoute({
  type: 'last-round',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to advance (optional, defaults to current encounter)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * End an encounter
 * 
 * Ends the current encounter in the Foundry world.
 * 
 * @route POST /end-encounter
 * @returns {object} Details of the ended encounter
 */
encounterRouter.post("/end-encounter", ...commonMiddleware, express.json(), createApiRoute({
  type: 'end-encounter',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to end (optional, defaults to current encounter)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Add tokens to an encounter
 * 
 * Adds selected tokens or specified UUIDs to the current encounter.
 * 
 * @route POST /add-to-encounter
 * @returns {object} Details of the updated encounter
 */
encounterRouter.post("/add-to-encounter", ...commonMiddleware, express.json(), createApiRoute({
  type: 'add-to-encounter',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to add tokens to (optional, defaults to current encounter)
    { name: 'selected', from: 'body', type: 'boolean' }, // Whether to add selected tokens (optional, defaults to false)
    { name: 'uuids', from: 'body', type: 'array' }, // The UUIDs of the tokens to add (optional, defaults to empty array)
    { name: 'rollInitiative', from: 'body', type: 'boolean' }, // Whether to roll initiative for the added tokens (optional, defaults to false)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));

/**
 * Remove tokens from an encounter
 * 
 * Removes selected tokens or specified UUIDs from the current encounter.
 * 
 * @route POST /remove-from-encounter
 * @returns {object} Details of the updated encounter
 */
encounterRouter.post("/remove-from-encounter", ...commonMiddleware, express.json(), createApiRoute({
  type: 'remove-from-encounter',
  requiredParams: [
    { name: 'clientId', from: 'query', type: 'string' } // The ID of the Foundry client to connect to 
  ],
  optionalParams: [
    { name: 'encounter', from: ['query', 'body'], type: 'string' }, // The ID of the encounter to remove tokens from (optional, defaults to current encounter)
    { name: 'selected', from: 'body', type: 'boolean' }, // Whether to remove selected tokens (optional, defaults to false)
    { name: 'uuids', from: 'body', type: 'array' }, // The UUIDs of the tokens to remove (optional, defaults to empty array)
    { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
  ]
}));