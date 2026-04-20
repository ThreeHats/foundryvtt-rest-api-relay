/**
 * Test environment variables
 * Centralized place for all test configuration
 */

import * as fs from 'fs';
import * as path from 'path';

// Read credentials written by register-account.test.ts (ephemeral account mode).
// Returns empty strings if the file doesn't exist yet (first test file to load this module).
function readRegisteredAccount(): { apiKey: string; email: string; password: string } {
  try {
    const file = path.join(__dirname, '../.global-vars.json');
    if (fs.existsSync(file)) {
      const store = JSON.parse(fs.readFileSync(file, 'utf8'));
      return {
        apiKey: store?.account?.apiKey || '',
        email: store?.account?.email || '',
        password: store?.account?.password || '',
      };
    }
  } catch {
    // File doesn't exist or is unreadable — fall through to env vars / defaults
  }
  return { apiKey: '', email: '', password: '' };
}

const registeredAccount = readRegisteredAccount();

// TEST_API_KEY, TEST_USER_EMAIL, and TEST_USER_PASSWORD are one semantic unit —
// they're all credentials for the same pre-provisioned account. If TEST_API_KEY
// is not set (ephemeral mode), ignore the email/password env vars too and use
// the registered account's credentials instead. Mixing them would point the
// change-password and login tests at a different account than the API key.
const preProvisioned = !!process.env.TEST_API_KEY;

export const testVariables = {
  baseUrl: process.env.TEST_BASE_URL || 'http://localhost:3010',
  apiKey: process.env.TEST_API_KEY || registeredAccount.apiKey || 'test-api-key',

  // These will be set during test execution
  clientId: '', // Set from global session
  foundryUrl: 'http://localhost:30013',
  worldName: 'testing',
  username: 'Gamemaster',

  // Test user credentials (for login, change-password, and account-deletion tests).
  // In pre-provisioned mode: must come from TEST_USER_EMAIL / TEST_USER_PASSWORD.
  // In ephemeral mode: come from the account registered by register-account.test.ts.
  userEmail: preProvisioned ? (process.env.TEST_USER_EMAIL || '') : (registeredAccount.email || ''),
  userPassword: preProvisioned ? (process.env.TEST_USER_PASSWORD || '') : (registeredAccount.password || ''),

  // Admin credentials (for admin dashboard tests).
  // Falls back to TEST_USER_EMAIL/PASSWORD if dedicated admin vars unset.
  adminEmail: process.env.TEST_ADMIN_EMAIL || process.env.TEST_USER_EMAIL || '',
  adminPassword: process.env.TEST_ADMIN_PASSWORD || process.env.TEST_USER_PASSWORD || '',

  // Test data UUIDs (will be populated as entities are created)
  testActorUuid: '',
  testItemUuid: '',
  testSceneUuid: '',
  testMacroUuid: ''
};

/**
 * Update a variable value
 */
export function setVariable(key: keyof typeof testVariables, value: string): void {
  testVariables[key] = value;
}

/**
 * Get all variables as a record
 */
export function getVariables(): Record<string, string> {
  return { ...testVariables };
}
