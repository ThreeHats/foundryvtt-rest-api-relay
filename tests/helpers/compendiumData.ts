/**
 * Compendium Data Helper
 *
 * Fetches entity data from the active system's compendiums to create
 * properly initialized test entities (with system-specific attributes like
 * initiative formulas, HP, etc.) instead of bare `type: 'base'` entities.
 */

import { ApiRequestConfig, makeRequest, replaceVariables } from './apiRequest';
import { testVariables, setVariable } from './testVariables';
import { getGlobalVariable } from './globalVariables';
import { getClientId } from './multiVersion';

// Cache: cacheKey → fetched data
const cache = new Map<string, Record<string, any> | null>();

export interface CompendiumFetchOptions {
  /** Override the system ID (default: read from global variables) */
  systemId?: string;
  /** Filter to a specific pack by ID substring (e.g., 'monsters', 'heroes', 'spells') */
  packFilter?: string;
  /** Pick a specific entry by name (case-insensitive substring match) */
  entryFilter?: string;
}

/**
 * Fetch entity data from a system compendium pack.
 * Returns a clean entity data object suitable for POST /create, or null on failure.
 */
export async function fetchCompendiumEntityData(
  version: string,
  entityType: string,
  options?: CompendiumFetchOptions
): Promise<Record<string, any> | null> {
  const cacheKey = `${version}:${entityType}:${options?.packFilter || ''}:${options?.entryFilter || ''}`;
  if (cache.has(cacheKey)) {
    return cache.get(cacheKey) ?? null;
  }

  try {
    const systemId = options?.systemId || getGlobalVariable(version, 'systemId');
    if (!systemId) {
      console.warn(`  ⚠ No systemId for v${version}, skipping compendium lookup`);
      cache.set(cacheKey, null);
      return null;
    }

    const clientId = getClientId(version);
    if (!clientId) {
      console.warn(`  ⚠ No clientId for v${version}, skipping compendium lookup`);
      cache.set(cacheKey, null);
      return null;
    }

    // Set clientId for variable replacement
    setVariable('clientId', clientId);

    // Step 1: GET /structure to discover compendium packs
    const structureConfig: ApiRequestConfig = {
      url: {
        raw: '{{baseUrl}}/structure',
        host: ['{{baseUrl}}'],
        path: ['structure'],
        query: [
          { key: 'clientId', value: '{{clientId}}' },
          { key: 'types', value: entityType }
        ]
      },
      method: 'GET',
      header: [
        { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
      ]
    };

    const resolvedStructure = replaceVariables(structureConfig, testVariables);
    const structureResponse = await makeRequest(resolvedStructure);

    // Structure response nests data inside a `data` field: response.data.data.compendiumPacks
    const structureData = structureResponse.data?.data || structureResponse.data;
    if (structureResponse.status !== 200 || !structureData?.compendiumPacks) {
      console.warn(`  ⚠ /structure call failed or no compendiumPacks for ${entityType} (status ${structureResponse.status})`);
      cache.set(cacheKey, null);
      return null;
    }

    const packs = structureData.compendiumPacks;

    // Step 2: Find a system pack with entries
    let targetEntry: { uuid: string; name: string } | null = null;

    for (const packKey of Object.keys(packs)) {
      const pack = packs[packKey];
      // Filter to system-specific packs (e.g., dnd5e.monsters)
      if (!pack.id?.startsWith(`${systemId}.`)) continue;
      if (!pack.entities || pack.entities.length === 0) continue;

      // Apply pack filter if specified
      if (options?.packFilter && !pack.id.includes(options.packFilter)) continue;

      // Pick entry — apply entry filter or take first
      let entry = pack.entities[0];
      if (options?.entryFilter) {
        const filter = options.entryFilter.toLowerCase();
        const match = pack.entities.find((e: any) => e.name?.toLowerCase().includes(filter));
        if (!match) continue; // Try next pack
        entry = match;
      }

      targetEntry = {
        uuid: `Compendium.${pack.id}.${entry.id}`,
        name: entry.name
      };
      break;
    }

    if (!targetEntry) {
      const filterDesc = options?.packFilter ? ` (pack: ${options.packFilter})` : '';
      console.warn(`  ⚠ No ${systemId} compendium entries found for ${entityType}${filterDesc}`);
      cache.set(cacheKey, null);
      return null;
    }

    // Step 3: GET /get to fetch full document data
    const getConfig: ApiRequestConfig = {
      url: {
        raw: '{{baseUrl}}/get',
        host: ['{{baseUrl}}'],
        path: ['get'],
        query: [
          { key: 'clientId', value: '{{clientId}}' },
          { key: 'uuid', value: targetEntry.uuid }
        ]
      },
      method: 'GET',
      header: [
        { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
      ]
    };

    const resolvedGet = replaceVariables(getConfig, testVariables);
    const getResponse = await makeRequest(resolvedGet);

    if (getResponse.status !== 200 || !getResponse.data?.data) {
      console.warn(`  ⚠ /get failed for compendium entry ${targetEntry.uuid} (status ${getResponse.status})`);
      cache.set(cacheKey, null);
      return null;
    }

    // Step 4: Clean up the data for creation
    const entityData = { ...getResponse.data.data };
    delete entityData._id;
    delete entityData._stats;
    delete entityData.folder;
    delete entityData.sort;
    delete entityData.ownership;
    delete entityData.flags;

    // Prefix name with test- for identification
    entityData.name = `test-${(entityData.name || entityType).toLowerCase()}`;

    console.log(`  ✓ Fetched compendium data for ${entityType}: type="${entityData.type}" from ${targetEntry.uuid}`);

    cache.set(cacheKey, entityData);
    return entityData;
  } catch (error) {
    console.warn(`  ⚠ Compendium lookup failed for ${entityType}:`, error);
    cache.set(cacheKey, null);
    return null;
  }
}
