/**
 * Test Entity Helper
 * 
 * Provides utilities for creating, tracking, and cleaning up test entities.
 * Entities created through this helper are automatically registered for cleanup.
 */

import { ApiRequestConfig, makeRequest, replaceVariables } from './apiRequest';
import { testVariables, setVariable } from './testVariables';
import { setGlobalVariable, getGlobalVariable } from './globalVariables';
import { captureExample } from './captureExample';
import { getClientId } from './multiVersion';

// Key used to store the cleanup list in global variables
const CLEANUP_KEY = '_createdEntities';

// Good example macro
const GOOD_EXAMPLE_MACRO = `// Example macro that uses parameters
function myMacro(args) {
  const targetName = args.targetName || "Target";
  const damage = args.damage || 0;
  const effect = args.effect || "none";
  
  // Use the parameters
  console.log(\`Attacking \${targetName} for \${damage} \${effect} damage\`);
  
  // Return a value (can be any data type)
  return {
    success: true,
    damageDealt: damage,
    target: targetName
  };
}

// Don't forget to return the result of your function
return myMacro(args);`;

/**
 * Specification for creating a test entity
 */
export interface EntitySpec {
  /** Semantic key for this entity (e.g., 'primary', 'expendable') */
  key: string;
  /** Entity type (e.g., 'Actor', 'Item', 'JournalEntry') */
  entityType: string;
  /** Optional data override for the entity */
  data?: Record<string, any>;
  /** Whether to capture this entity creation for documentation (default: false) */
  captureForDocs?: boolean;
}

/**
 * Options for createTestEntities
 */
export interface CreateEntitiesOptions {
  /** Array to push captured examples to (for documentation) */
  capturedExamples?: any[];
}

/**
 * Default data for different entity types
 */
function getDefaultEntityData(entityType: string): Record<string, any> {
  switch (entityType) {
    case 'Actor':
      return {
        name: `test-${entityType.toLowerCase()}`,
        type: 'base'
      };
    case 'Item':
      return {
        name: `test-${entityType.toLowerCase()}`,
        type: 'base'
      };
    case 'JournalEntry':
      return {
        name: `test-${entityType.toLowerCase()}`
      };
    case 'Scene':
      return {
        name: `test-${entityType.toLowerCase()}`
      };
    case 'Macro':
      return {
        name: `test-${entityType.toLowerCase()}`,
        type: 'script',
        command: GOOD_EXAMPLE_MACRO
      };
    case 'RollTable':
      return {
        name: `test-${entityType.toLowerCase()}`
      };
    case 'Playlist':
      return {
        name: `test-${entityType.toLowerCase()}`
      };
    default:
      return {
        name: `test-${entityType.toLowerCase()}`
      };
  }
}

/**
 * Build the global variable key for an entity
 */
function buildEntityKey(entityType: string, key: string): string {
  return `${entityType.toLowerCase()}_${key}_uuid`;
}

/**
 * Register an entity UUID for cleanup
 */
export function registerForCleanup(version: string, uuid: string): void {
  const existing = getGlobalVariable(version, CLEANUP_KEY, []) as string[];
  if (!existing.includes(uuid)) {
    existing.push(uuid);
    setGlobalVariable(version, CLEANUP_KEY, existing);
  }
}

/**
 * Get all registered entity UUIDs for cleanup
 */
export function getCleanupList(version: string): string[] {
  return getGlobalVariable(version, CLEANUP_KEY, []) as string[];
}

/**
 * Clear the cleanup list (after cleanup is done)
 */
export function clearCleanupList(version: string): void {
  setGlobalVariable(version, CLEANUP_KEY, []);
}

/**
 * Get a named entity UUID
 * @param version - Foundry version
 * @param entityType - Entity type (e.g., 'Actor', 'Item')
 * @param key - Semantic key (e.g., 'primary', 'expendable')
 */
export function getEntityUuid(version: string, entityType: string, key: string): string | undefined {
  const varKey = buildEntityKey(entityType, key);
  return getGlobalVariable(version, varKey);
}

/**
 * Create multiple test entities and register them for cleanup
 * 
 * @example
 * await createTestEntities('13', [
 *   { key: 'primary', entityType: 'Actor', captureForDocs: true },
 *   { key: 'secondary', entityType: 'Actor' },
 *   { key: 'expendable', entityType: 'Actor' },
 * ], { capturedExamples });
 * 
 * // Later, retrieve with:
 * const uuid = getEntityUuid('13', 'Actor', 'primary');
 */
export async function createTestEntities(
  version: string,
  specs: EntitySpec[],
  options?: CreateEntitiesOptions
): Promise<Map<string, string>> {
  const results = new Map<string, string>();
  const clientId = getClientId(version);
  
  if (!clientId) {
    throw new Error(`No clientId found for version ${version}. Make sure session tests have run.`);
  }
  
  // Set clientId in testVariables for request replacement
  setVariable('clientId', clientId);
  
  for (const spec of specs) {
    const entityData = {
      ...getDefaultEntityData(spec.entityType),
      ...spec.data
    };
    
    // Build request config
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: '{{baseUrl}}/create',
        host: ['{{baseUrl}}'],
        path: ['create'],
        query: [
          { key: 'clientId', value: '{{clientId}}' }
        ]
      },
      method: 'POST',
      header: [
        { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
        { key: 'Content-Type', value: 'application/json', type: 'text' }
      ],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          entityType: spec.entityType,
          data: entityData
        }, null, 2)
      }
    };
    
    let response;
    
    if (spec.captureForDocs && options?.capturedExamples) {
      // Use captureExample to get documentation
      const captured = await captureExample(
        requestConfig,
        testVariables,
        '/create'
      );
      options.capturedExamples.push(captured);
      response = captured.response;
    } else {
      // Just make the request without capturing
      const resolvedConfig = replaceVariables(requestConfig, testVariables);
      const apiResponse = await makeRequest(resolvedConfig);
      response = {
        status: apiResponse.status,
        data: apiResponse.data
      };
    }
    
    if (response.status === 200 && response.data?.uuid) {
      const uuid = response.data.uuid;
      const varKey = buildEntityKey(spec.entityType, spec.key);
      
      // Store as named global variable
      setGlobalVariable(version, varKey, uuid);
      
      // Register for cleanup
      registerForCleanup(version, uuid);
      
      // Add to results map
      results.set(`${spec.entityType}:${spec.key}`, uuid);
      
      console.log(`  ✓ Created ${spec.entityType} '${spec.key}': ${uuid}`);
    } else {
      console.warn(`  ✗ Failed to create ${spec.entityType} '${spec.key}': ${response.status} - ${JSON.stringify(response.data)}`);
    }
  }
  
  return results;
}

/**
 * Create a single test entity (convenience wrapper)
 */
export async function createTestEntity(
  version: string,
  entityType: string,
  key: string,
  options?: {
    data?: Record<string, any>;
    captureForDocs?: boolean;
    capturedExamples?: any[];
  }
): Promise<string | undefined> {
  const results = await createTestEntities(version, [{
    key,
    entityType,
    data: options?.data,
    captureForDocs: options?.captureForDocs
  }], {
    capturedExamples: options?.capturedExamples
  });
  
  return results.get(`${entityType}:${key}`);
}
