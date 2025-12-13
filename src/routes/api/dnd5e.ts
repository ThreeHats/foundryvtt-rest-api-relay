/**
 * TODO: add POST /ability-check, POST /saving-throw, POST /skill-check
 * 
 */
import express, { Router } from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { createApiRoute } from '../route-helpers';

export const dnd5eRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage, express.json()];

/**
 * Get detailed information for a specific D&D 5e actor.
 * 
 * Retrieves comprehensive details about an actor including stats, inventory,
 * spells, features, and other character information based on the requested details array.
 * 
 * @route GET /dnd5e/get-actor-details
 * @returns {object} Actor details object containing requested information
 */
dnd5eRouter.get("/get-actor-details", ...commonMiddleware, createApiRoute({
    type: 'get-actor-details',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor
        { name: 'details', from: ['body', 'query'], type: 'array' } // Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"])
    ]
}));

/**
 * Modify the charges for a specific item owned by an actor.
 * 
 * Increases or decreases the charges/uses of an item in an actor's inventory.
 * Useful for consumable items like potions, scrolls, or charged magic items.
 * 
 * @route POST /dnd5e/modify-item-charges
 * @returns {object} Result of the charge modification operation
 */
dnd5eRouter.post("/modify-item-charges", ...commonMiddleware, createApiRoute({
    type: 'modify-item-charges',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor who owns the item
        { name: 'amount', from: ['body', 'query'], type: 'number' } // The amount to modify charges by (positive or negative)
    ],
    optionalParams: [
        { name: 'itemUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific item (optional if itemName provided)
        { name: 'itemName', from: ['body', 'query'], type: 'string' } // The name of the item if UUID not provided (optional if itemUuid provided)
    ]
}));

/**
 * Use a general ability for an actor.
 * 
 * Triggers the use of any ability, feature, spell, or item for an actor.
 * This is a generic endpoint that can handle various types of abilities.
 * 
 * @route POST /dnd5e/use-ability
 * @returns {object} Result of the ability use operation
 */
dnd5eRouter.post("/use-ability", ...commonMiddleware, createApiRoute({
    type: 'use-ability',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor using the ability
    ],
    optionalParams: [
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific ability (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the ability if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the ability (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' } // The name of the target if UUID not provided (optional)
    ]
}));

/**
 * Use a class or racial feature for an actor.
 * 
 * Activates class features (like Action Surge, Rage) or racial features 
 * (like Dragonborn Breath Weapon) for a character.
 * 
 * @route POST /dnd5e/use-feature
 * @returns {object} Result of the feature use operation
 */
dnd5eRouter.post("/use-feature", ...commonMiddleware, createApiRoute({
    type: 'use-feature',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor using the feature
    ],
    optionalParams: [
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific feature (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the feature if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the feature (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' } // The name of the target if UUID not provided (optional)
    ]
}));

/**
 * Cast a spell for an actor.
 * 
 * Casts a spell from the actor's spell list, consuming spell slots as appropriate.
 * Handles cantrips, leveled spells, and spell-like abilities.
 * 
 * @route POST /dnd5e/use-spell
 * @returns {object} Result of the spell casting operation
 */
dnd5eRouter.post("/use-spell", ...commonMiddleware, createApiRoute({
    type: 'use-spell',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor casting the spell
    ],
    optionalParams: [
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific spell (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the spell if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the spell (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' } // The name of the target if UUID not provided (optional)
    ]
}));

/**
 * Use an item for an actor.
 * 
 * Activates an item from the actor's inventory, such as drinking a potion,
 * using a magic item, or activating equipment with special properties.
 * 
 * @route POST /dnd5e/use-item
 * @returns {object} Result of the item use operation
 */
dnd5eRouter.post("/use-item", ...commonMiddleware, createApiRoute({
    type: 'use-item',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor using the item
    ],
    optionalParams: [
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific item (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the item if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the item (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' } // The name of the target if UUID not provided (optional)
    ]
}));

/**
 * Modify the experience points for a specific actor.
 * 
 * Adds or removes experience points from an actor.
 * 
 * @route POST /dnd5e/modify-experience
 * @returns {object} Result of the experience modification operation
 */
dnd5eRouter.post("/modify-experience", ...commonMiddleware, createApiRoute({
    type: 'modify-experience',
    requiredParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'amount', from: ['body', 'query'], type: 'number' } // The amount of experience to add (can be negative)
    ],
    optionalParams: [
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor to modify
        { name: 'selected', from: ['body', 'query'], type: 'boolean' } // Modify the selected token's actor
    ]
}));
