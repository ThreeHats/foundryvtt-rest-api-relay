/**
 * System-Specific Test Setup
 *
 * Provides per-system configuration for test entities. Each system module
 * defines how to set up actors with proper items, spells, features, etc.
 * for system-specific endpoint testing.
 *
 * To add support for a new system:
 * 1. Add a new SystemTestConfig entry to SYSTEM_CONFIGS
 * 2. Implement the setup functions for that system's needs
 */

import { ApiRequestConfig, makeRequest, replaceVariables } from './apiRequest';
import { testVariables, setVariable } from './testVariables';
import { getGlobalVariable, setGlobalVariable } from './globalVariables';
import { getClientId } from './multiVersion';
import { fetchCompendiumEntityData } from './compendiumData';
import { registerForCleanup } from './testEntities';

export interface SystemTestConfig {
  /** System ID (e.g., 'dnd5e', 'pf2e') */
  id: string;

  /** Compendium pack filters for fetching specific entity types */
  compendiumPacks: {
    /** Pack filter for heroes/characters (for Actor creation) */
    actors?: string;
    /** Pack filter for spells */
    spells?: string;
    /** Pack filter for items with charges/uses */
    consumables?: string;
  };

  /**
   * Give the test actor a spell from the system's compendium.
   * Returns the spell name if successful, null otherwise.
   */
  giveSpell: (version: string, actorUuid: string) => Promise<string | null>;

  /**
   * Give the test actor a consumable item with charges.
   * Returns { name, maxCharges } if successful, null otherwise.
   */
  giveConsumable: (version: string, actorUuid: string) => Promise<{ name: string; maxCharges: number } | null>;
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

/**
 * Give a compendium item to an actor by fetching it from the compendium,
 * creating it as a world item, then using /give to add it to the actor.
 */
async function giveCompendiumItemToActor(
  version: string,
  actorUuid: string,
  entityType: string,
  options: { packFilter?: string; entryFilter?: string }
): Promise<{ uuid: string; name: string; data: Record<string, any> } | null> {
  const clientId = getClientId(version);
  if (!clientId) return null;

  setVariable('clientId', clientId);

  // Fetch the item data from compendium
  const itemData = await fetchCompendiumEntityData(version, entityType, options);
  if (!itemData) return null;

  // Create the item as a world Item
  const createConfig: ApiRequestConfig = {
    url: {
      raw: '{{baseUrl}}/create',
      host: ['{{baseUrl}}'],
      path: ['create'],
      query: [{ key: 'clientId', value: '{{clientId}}' }]
    },
    method: 'POST',
    header: [
      { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
      { key: 'Content-Type', value: 'application/json', type: 'text' }
    ],
    body: {
      mode: 'raw',
      raw: JSON.stringify({ entityType: 'Item', data: itemData })
    }
  };

  const createResolved = replaceVariables(createConfig, testVariables);
  const createResponse = await makeRequest(createResolved);

  if (createResponse.status !== 200 || !createResponse.data?.uuid) {
    console.warn(`  ⚠ Failed to create world item from compendium: ${createResponse.status}`);
    return null;
  }

  const itemUuid = createResponse.data.uuid;
  registerForCleanup(version, itemUuid);

  // Give the item to the actor
  const giveConfig: ApiRequestConfig = {
    url: {
      raw: '{{baseUrl}}/give',
      host: ['{{baseUrl}}'],
      path: ['give'],
      query: [{ key: 'clientId', value: '{{clientId}}' }]
    },
    method: 'POST',
    header: [
      { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
      { key: 'Content-Type', value: 'application/json', type: 'text' }
    ],
    body: {
      mode: 'raw',
      raw: JSON.stringify({
        toUuid: actorUuid,
        itemUuid: itemUuid,
        quantity: 1
      })
    }
  };

  const giveResolved = replaceVariables(giveConfig, testVariables);
  const giveResponse = await makeRequest(giveResolved);

  if (giveResponse.status !== 200 || !giveResponse.data?.success) {
    console.warn(`  ⚠ Failed to give item to actor: ${giveResponse.status} - ${JSON.stringify(giveResponse.data)}`);
    return null;
  }

  console.log(`  ✓ Gave "${itemData.name}" to actor`);
  return { uuid: itemUuid, name: itemData.name, data: itemData };
}

// ──────────────────────────────────────────────
// System configurations
// ──────────────────────────────────────────────

const dnd5eConfig: SystemTestConfig = {
  id: 'dnd5e',
  compendiumPacks: {
    actors: 'heroes',
    spells: 'spells',
    consumables: 'items',
  },

  giveSpell: async (version, actorUuid) => {
    const result = await giveCompendiumItemToActor(version, actorUuid, 'Item', {
      packFilter: 'spells',
    });
    if (!result) return null;

    // Store spell name for test assertions
    setGlobalVariable(version, 'dnd5e_test_spell_name', result.name);
    return result.name;
  },

  giveConsumable: async (version, actorUuid) => {
    // Try to find a consumable with uses/charges in the items compendium
    const result = await giveCompendiumItemToActor(version, actorUuid, 'Item', {
      packFilter: 'items',
      entryFilter: 'potion',
    });
    if (!result) return null;

    const maxCharges = parseInt(result.data.system?.uses?.max) || 1;
    setGlobalVariable(version, 'dnd5e_test_consumable_name', result.name);
    return { name: result.name, maxCharges };
  },
};

// ──────────────────────────────────────────────
// Registry
// ──────────────────────────────────────────────

const SYSTEM_CONFIGS: Record<string, SystemTestConfig> = {
  dnd5e: dnd5eConfig,
};

/**
 * Get the test setup config for a system, or null if not supported.
 */
export function getSystemConfig(systemId: string): SystemTestConfig | null {
  return SYSTEM_CONFIGS[systemId] || null;
}

/**
 * Run system-specific setup for a version's actor.
 * Gives the actor spells, consumables, etc. based on the system config.
 * Returns a summary of what was set up.
 */
export async function setupSystemTestData(
  version: string,
  actorUuid: string
): Promise<{ spellName: string | null; consumableName: string | null }> {
  const systemId = getGlobalVariable(version, 'systemId');
  const config = systemId ? getSystemConfig(systemId) : null;

  if (!config) {
    return { spellName: null, consumableName: null };
  }

  console.log(`  Setting up ${config.id}-specific test data...`);

  const spellName = await config.giveSpell(version, actorUuid);
  const consumableName = (await config.giveConsumable(version, actorUuid))?.name ?? null;

  return { spellName, consumableName };
}
