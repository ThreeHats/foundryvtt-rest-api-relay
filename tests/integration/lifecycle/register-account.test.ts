/**
 * @file register-account.test.ts
 * @description Phase 0: Register a fresh test account for this test run.
 *
 * If TEST_API_KEY is set in the environment, the pre-provisioned account is used
 * and no registration occurs. Otherwise, a fresh account is registered and its
 * credentials are written to .global-vars.json so subsequent test files pick them
 * up via testVariables. auth-cleanup.test.ts deletes the account at the end.
 */

import { describe, test, expect } from '@jest/globals';
import { makeRequest } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { setGlobalVariable } from '../../helpers/globalVariables';

describe('Account Setup', () => {
  test('Register test account (or use pre-provisioned)', async () => {
    if (process.env.TEST_API_KEY) {
      // Pre-provisioned mode — nothing to register or clean up.
      setGlobalVariable('account', 'registered', false);
      console.log('  Using pre-provisioned TEST_API_KEY — skipping registration');
      return;
    }

    const email = `relay-test-${Date.now()}@example.com`;
    const password = 'TestPassword123!';

    const response = await makeRequest({
      url: {
        raw: `${testVariables.baseUrl}/auth/register`,
        host: [testVariables.baseUrl],
        path: ['auth', 'register'],
      },
      method: 'POST',
      header: [],
      body: { mode: 'raw', raw: JSON.stringify({ email, password }) },
    });

    if (response.status === 403) {
      throw new Error(
        'Registration is disabled on this server. ' +
        'Set TEST_API_KEY in .env.test to use a pre-provisioned account.'
      );
    }

    expect(response.status).toBe(201);
    expect(response.data).toHaveProperty('apiKey');

    const apiKey = response.data.apiKey as string;

    // Write credentials to global vars — testVariables.ts reads these as fallback
    // so every subsequent test file automatically uses this account.
    setGlobalVariable('account', 'apiKey', apiKey);
    setGlobalVariable('account', 'email', email);
    setGlobalVariable('account', 'password', password);
    setGlobalVariable('account', 'registered', true);

    console.log(`  ✓ Registered test account: ${email}`);
  });
});
