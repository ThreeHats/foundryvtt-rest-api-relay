/**
 * Global Variables Helper
 * Manages global test data that needs to be shared across ALL test files
 * Uses a file-based approach since Jest isolates the global object per test file
 * This is useful for storing UUIDs and other identifiers created during tests
 */

import * as fs from 'fs';
import * as path from 'path';

interface GlobalVariables {
  [version: string]: {
    [key: string]: any;
  };
}

// File path for persisting global variables
const GLOBAL_VARS_FILE = path.join(__dirname, '..', '.global-vars.json');

/**
 * Read the global variables from file
 */
function readStore(): GlobalVariables {
  try {
    if (fs.existsSync(GLOBAL_VARS_FILE)) {
      const content = fs.readFileSync(GLOBAL_VARS_FILE, 'utf8');
      return JSON.parse(content);
    }
  } catch (e) {
    // File doesn't exist or is invalid - return empty store
  }
  return {};
}

/**
 * Write the global variables to file
 */
function writeStore(store: GlobalVariables): void {
  try {
    fs.writeFileSync(GLOBAL_VARS_FILE, JSON.stringify(store, null, 2));
  } catch (e) {
    console.error('Failed to write global variables file:', e);
    throw e;
  }
}

/**
 * Create or update a global variable for a specific Foundry version
 * @param version - Foundry version (e.g., '13', '12')
 * @param key - Variable name
 * @param value - Value to store
 */
export function setGlobalVariable(version: string, key: string, value: any): void {
  const store = readStore();
  if (!store[version]) {
    store[version] = {};
  }
  store[version][key] = value;
  writeStore(store);
}

/**
 * Get a global variable for a specific Foundry version
 * @param version - Foundry version (e.g., '13', '12')
 * @param key - Variable name
 * @param defaultValue - Default value if variable doesn't exist
 * @returns The stored value or default value
 */
export function getGlobalVariable(version: string, key: string, defaultValue?: any): any {
  const store = readStore();
  if (!store[version]) {
    return defaultValue;
  }
  const value = store[version][key];
  return value !== undefined ? value : defaultValue;
}

/**
 * Get all global variables for a specific version
 * @param version - Foundry version (e.g., '13', '12')
 * @returns Object containing all variables for that version
 */
export function getVersionVariables(version: string): Record<string, any> {
  const store = readStore();
  return store[version] || {};
}

/**
 * Clear all global variables for a specific version
 * @param version - Foundry version (e.g., '13', '12')
 */
export function clearVersionVariables(version: string): void {
  const store = readStore();
  if (store[version]) {
    delete store[version];
    writeStore(store);
  }
}

/**
 * Clear all global variables for all versions
 */
export function clearAllVariables(): void {
  writeStore({});
}

/**
 * Check if a global variable exists for a specific version
 * @param version - Foundry version (e.g., '13', '12')
 * @param key - Variable name
 * @returns True if the variable exists
 */
export function hasGlobalVariable(version: string, key: string): boolean {
  const store = readStore();
  return store[version] && store[version][key] !== undefined;
}

/**
 * Delete the global variables file (for cleanup)
 */
export function deleteGlobalVarsFile(): void {
  try {
    if (fs.existsSync(GLOBAL_VARS_FILE)) {
      fs.unlinkSync(GLOBAL_VARS_FILE);
    }
  } catch (e) {
    // Ignore errors
  }
}
