/**
 * Helper for running tests across multiple Foundry versions
 */

import { getGlobalVariable } from './globalVariables';

// Check if using existing sessions from env
const useExistingSession = process.env.USE_EXISTING_SESSION === 'true';

/**
 * Get all configured Foundry versions from environment
 */
export function getConfiguredVersions(): string[] {
  const versions = (process.env.TEST_FOUNDRY_VERSIONS || '13')
    .split(',')
    .map(v => v.trim())
    .filter(v => v.length > 0);
  return versions;
}

/**
 * Run a test suite for each configured Foundry version
 * 
 * @example
 * forEachVersion((version, getClientId) => {
 *   describe(`/roll (v${version})`, () => {
 *     test('POST /roll', async () => {
 *       const clientId = getClientId();
 *       // ... test code using clientId
 *     });
 *   });
 * });
 */
export function forEachVersion(
  testFn: (version: string, getClientId: () => string) => void
): void {
  const versions = getConfiguredVersions();
  
  versions.forEach(version => {
    const getClientId = (): string => {
      // If using existing session, get clientId from env
      if (useExistingSession) {
        const envClientId = process.env[`TEST_CLIENT_ID_V${version}`];
        if (envClientId) {
          return envClientId;
        }
        console.warn(`⚠️ USE_EXISTING_SESSION=true but TEST_CLIENT_ID_V${version} not set`);
        return '';
      }
      
      // Otherwise get clientId from globalVariables (set by session tests)
      const clientId = getGlobalVariable(version, 'clientId');
      if (!clientId) {
        console.warn(`⚠️ No clientId found for v${version} - session tests may not have run yet`);
      }
      return clientId || '';
    };
    
    testFn(version, getClientId);
  });
}

/**
 * Get client ID for a specific version or default
 */
export function getClientId(version?: string): string {
  const targetVersion = version || getConfiguredVersions()[0];
  
  // If using existing session, get clientId from env
  if (useExistingSession) {
    return process.env[`TEST_CLIENT_ID_V${targetVersion}`] || '';
  }
  
  return getGlobalVariable(targetVersion, 'clientId') || '';
}

// TODO: scrutinize the following methods, and implement system specific tests
/**
 * Get system ID for a specific version
 */
export function getSystemId(version?: string): string {
  const targetVersion = version || getConfiguredVersions()[0];
  
  // TODO: will need to add this to the .env.test.example file
  // If using existing session, get systemId from env
  if (useExistingSession) {
    return process.env[`TEST_SYSTEM_ID_V${targetVersion}`] || '';
  }
  
  return getGlobalVariable(targetVersion, 'systemId') || '';
}

/**
 * Check if a specific version is running a particular system
 */
export function isSystem(systemId: string, version?: string): boolean {
  return getSystemId(version) === systemId;
}

/**
 * Check if any configured version is running a particular system
 */
export function hasSystemVersion(systemId: string): boolean {
  const versions = getConfiguredVersions();
  return versions.some(v => getSystemId(v) === systemId);
}

/**
 * Get only the versions running a specific system
 */
export function getVersionsWithSystem(systemId: string): string[] {
  const versions = getConfiguredVersions();
  return versions.filter(v => getSystemId(v) === systemId);
}

/**
 * Run a test suite for each configured Foundry version that has a specific system
 */
export function forEachVersionWithSystem(
  systemId: string,
  testFn: (version: string, getClientId: () => string) => void
): void {
  const versions = getVersionsWithSystem(systemId);
  
  if (versions.length === 0) {
    // No versions with this system - create a skipped describe block
    describe.skip(`(No clients with ${systemId} system)`, () => {
      test('skipped', () => {});
    });
    return;
  }
  
  versions.forEach(version => {
    const getClientIdFn = (): string => {
      if (useExistingSession) {
        const envClientId = process.env[`TEST_CLIENT_ID_V${version}`];
        if (envClientId) {
          return envClientId;
        }
        console.warn(`⚠️ USE_EXISTING_SESSION=true but TEST_CLIENT_ID_V${version} not set`);
        return '';
      }
      
      const clientId = getGlobalVariable(version, 'clientId');
      if (!clientId) {
        console.warn(`⚠️ No clientId found for v${version} - session tests may not have run yet`);
      }
      return clientId || '';
    };
    
    testFn(version, getClientIdFn);
  });
}
