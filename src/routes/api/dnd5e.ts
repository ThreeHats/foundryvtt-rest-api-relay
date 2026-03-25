/**
 * D&D 5e system-specific API routes.
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
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor
        { name: 'details', from: ['body', 'query'], type: 'array' } // Array of detail types to retrieve (e.g., ["resources", "items", "spells", "features"])
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
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
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor who owns the item
        { name: 'amount', from: ['body', 'query'], type: 'number' } // The amount to modify charges by (positive or negative)
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'itemUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific item (optional if itemName provided)
        { name: 'itemName', from: ['body', 'query'], type: 'string' }, // The name of the item if UUID not provided (optional if itemUuid provided)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
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
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor using the ability
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific ability (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the ability if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the ability (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' }, // The name of the target if UUID not provided (optional)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
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
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor using the feature
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific feature (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the feature if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the feature (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' }, // The name of the target if UUID not provided (optional)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
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
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor casting the spell
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific spell (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the spell if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the spell (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' }, // The name of the target if UUID not provided (optional)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
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
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' } // UUID of the actor using the item
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'abilityUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the specific item (optional if abilityName provided)
        { name: 'abilityName', from: ['body', 'query'], type: 'string' }, // The name of the item if UUID not provided (optional if abilityUuid provided)
        { name: 'targetUuid', from: ['body', 'query'], type: 'string' }, // The UUID of the target for the item (optional)
        { name: 'targetName', from: ['body', 'query'], type: 'string' }, // The name of the target if UUID not provided (optional)
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Perform a short rest for an actor.
 *
 * Triggers the D&D 5e short rest workflow including hit dice recovery,
 * class feature resets, and HP recovery.
 *
 * @route POST /dnd5e/short-rest
 * @returns {object} Result of the short rest operation
 */
dnd5eRouter.post("/short-rest", ...commonMiddleware, createApiRoute({
    type: 'short-rest',
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor to rest (optional if selected is true)
        { name: 'selected', from: ['body', 'query'], type: 'boolean' }, // Use the currently selected token's actor
        { name: 'autoHD', from: ['body', 'query'], type: 'boolean' }, // Automatically spend hit dice during short rest
        { name: 'autoHDThreshold', from: ['body', 'query'], type: 'number' }, // HP threshold below which to auto-spend hit dice (0-1 as fraction of max HP)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Perform a long rest for an actor.
 *
 * Triggers the D&D 5e long rest workflow including full HP recovery,
 * spell slot restoration, hit dice recovery, and feature resets.
 *
 * @route POST /dnd5e/long-rest
 * @returns {object} Result of the long rest operation
 */
dnd5eRouter.post("/long-rest", ...commonMiddleware, createApiRoute({
    type: 'long-rest',
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor to rest (optional if selected is true)
        { name: 'selected', from: ['body', 'query'], type: 'boolean' }, // Use the currently selected token's actor
        { name: 'newDay', from: ['body', 'query'], type: 'boolean' }, // Whether the long rest marks a new day (default: true)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Roll a skill check for an actor.
 *
 * Rolls a D&D 5e skill check with all applicable modifiers including
 * proficiency, expertise, Jack of All Trades, and conditional bonuses.
 *
 * @route POST /dnd5e/skill-check
 * @returns {object} Result of the skill check roll
 */
dnd5eRouter.post("/skill-check", ...commonMiddleware, createApiRoute({
    type: 'skill-check',
    requiredParams: [
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor rolling the check
        { name: 'skill', from: ['body', 'query'], type: 'string' }, // Skill abbreviation (e.g., "acr", "ath", "ste", "prc")
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'advantage', from: ['body', 'query'], type: 'boolean' }, // Roll with advantage
        { name: 'disadvantage', from: ['body', 'query'], type: 'boolean' }, // Roll with disadvantage
        { name: 'bonus', from: ['body', 'query'], type: 'string' }, // Extra bonus formula to add (e.g., "1d4", "+2")
        { name: 'createChatMessage', from: ['body', 'query'], type: 'boolean' }, // Whether to post the roll to chat (default: true)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Roll an ability saving throw for an actor.
 *
 * Rolls a D&D 5e ability saving throw with all applicable modifiers.
 *
 * @route POST /dnd5e/ability-save
 * @returns {object} Result of the saving throw roll
 */
dnd5eRouter.post("/ability-save", ...commonMiddleware, createApiRoute({
    type: 'ability-save',
    requiredParams: [
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor rolling the save
        { name: 'ability', from: ['body', 'query'], type: 'string' }, // Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha")
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'advantage', from: ['body', 'query'], type: 'boolean' }, // Roll with advantage
        { name: 'disadvantage', from: ['body', 'query'], type: 'boolean' }, // Roll with disadvantage
        { name: 'bonus', from: ['body', 'query'], type: 'string' }, // Extra bonus formula to add (e.g., "1d4", "+2")
        { name: 'createChatMessage', from: ['body', 'query'], type: 'boolean' }, // Whether to post the roll to chat (default: true)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Roll an ability check for an actor.
 *
 * Rolls a D&D 5e ability check (raw ability test, not a skill check)
 * with all applicable modifiers.
 *
 * @route POST /dnd5e/ability-check
 * @returns {object} Result of the ability check roll
 */
dnd5eRouter.post("/ability-check", ...commonMiddleware, createApiRoute({
    type: 'ability-check',
    requiredParams: [
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor rolling the check
        { name: 'ability', from: ['body', 'query'], type: 'string' }, // Ability abbreviation (e.g., "str", "dex", "con", "int", "wis", "cha")
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'advantage', from: ['body', 'query'], type: 'boolean' }, // Roll with advantage
        { name: 'disadvantage', from: ['body', 'query'], type: 'boolean' }, // Roll with disadvantage
        { name: 'bonus', from: ['body', 'query'], type: 'string' }, // Extra bonus formula to add (e.g., "1d4", "+2")
        { name: 'createChatMessage', from: ['body', 'query'], type: 'boolean' }, // Whether to post the roll to chat (default: true)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));

/**
 * Roll a death saving throw for an actor.
 *
 * Rolls a D&D 5e death saving throw, handling DC 10 CON save,
 * three successes/failures tracking, nat 20 healing, and nat 1 double failure.
 *
 * @route POST /dnd5e/death-save
 * @returns {object} Result of the death saving throw
 */
dnd5eRouter.post("/death-save", ...commonMiddleware, createApiRoute({
    type: 'death-save',
    requiredParams: [
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor rolling the death save
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'advantage', from: ['body', 'query'], type: 'boolean' }, // Roll with advantage
        { name: 'createChatMessage', from: ['body', 'query'], type: 'boolean' }, // Whether to post the roll to chat (default: true)
        { name: 'userId', from: ['query', 'body'], type: 'string' }, // Foundry user ID or username to scope permissions (omit for GM-level access)
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
        { name: 'amount', from: ['body', 'query'], type: 'number' } // The amount of experience to add (can be negative)
    ],
    optionalParams: [
        { name: 'clientId', from: ['body', 'query'], type: 'string' }, // Client ID for the Foundry world
        { name: 'actorUuid', from: ['body', 'query'], type: 'string' }, // UUID of the actor to modify
        { name: 'selected', from: ['body', 'query'], type: 'boolean' }, // Modify the selected token's actor
        { name: 'userId', from: ['query', 'body'], type: 'string' } // Foundry user ID or username to scope permissions (omit for GM-level access)
    ]
}));
