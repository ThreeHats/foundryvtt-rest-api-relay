/**
 * @file auth-cleanup.test.ts
 * @description Cleanup: delete throwaway auth test user and test account deletion flow
 * @endpoints DELETE /auth/account, POST /auth/login
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { captureExample, appendExamples } from '../../helpers/captureExample';
import { getGlobalVariable } from '../../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Read throwaway user credentials from global variables (set by auth-endpoints.test.ts)
const throwawayEmail = getGlobalVariable('auth', 'throwawayEmail', '');
const throwawayPassword = getGlobalVariable('auth', 'throwawayPassword', '');
const throwawayApiKey = getGlobalVariable('auth', 'throwawayApiKey', '');

// Read main test account credentials (set by register-account.test.ts in ephemeral mode)
const mainAccountRegistered = getGlobalVariable('account', 'registered', false);
const mainAccountEmail = getGlobalVariable('account', 'email', '');
const mainAccountPassword = getGlobalVariable('account', 'password', '');
const mainAccountApiKey = getGlobalVariable('account', 'apiKey', '');

const hasThrowawayUser = throwawayEmail !== '' && throwawayApiKey !== '';
const describeOrSkip = hasThrowawayUser ? describe : describe.skip;

describeOrSkip('Auth Cleanup — Account Deletion', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/auth-examples.json');
      appendExamples(capturedExamples, outputPath);
      console.log(`\nAppended ${capturedExamples.length} auth-cleanup examples to ${outputPath}`);
    }
  });

  test('DELETE /auth/account - should reject wrong email confirmation', async () => {
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: `${testVariables.baseUrl}/auth/account`,
        host: [testVariables.baseUrl],
        path: ['auth', 'account'],
      },
      method: 'DELETE',
      header: [
        { key: 'x-api-key', value: throwawayApiKey }
      ],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          confirmEmail: 'wrong@example.com',
          password: throwawayPassword
        })
      }
    };

    const captured = await captureExample(requestConfig, {}, '/auth/account (wrong email)');
    expect(captured.response.status).toBe(400);
  });

  test('DELETE /auth/account - should reject wrong password', async () => {
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: `${testVariables.baseUrl}/auth/account`,
        host: [testVariables.baseUrl],
        path: ['auth', 'account'],
      },
      method: 'DELETE',
      header: [
        { key: 'x-api-key', value: throwawayApiKey }
      ],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          confirmEmail: throwawayEmail,
          password: 'WrongPassword1'
        })
      }
    };

    const captured = await captureExample(requestConfig, {}, '/auth/account (wrong password)');
    expect(captured.response.status).toBe(401);
  });

  test('DELETE /auth/account - delete account successfully', async () => {
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: `${testVariables.baseUrl}/auth/account`,
        host: [testVariables.baseUrl],
        path: ['auth', 'account'],
      },
      method: 'DELETE',
      header: [
        { key: 'x-api-key', value: throwawayApiKey }
      ],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          confirmEmail: throwawayEmail,
          password: throwawayPassword
        })
      }
    };

    const captured = await captureExample(requestConfig, {}, '/auth/account');
    capturedExamples.push(captured);

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('success', true);
    expect(captured.response.data).toHaveProperty('message');
  });

  test('POST /auth/login - verify deleted user cannot login', async () => {
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: `${testVariables.baseUrl}/auth/login`,
        host: [testVariables.baseUrl],
        path: ['auth', 'login'],
      },
      method: 'POST',
      header: [],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          email: throwawayEmail,
          password: throwawayPassword
        })
      }
    };

    const captured = await captureExample(requestConfig, {}, '/auth/login (deleted user)');
    expect(captured.response.status).toBe(401);
  });
});

// Delete the main test account registered by register-account.test.ts (ephemeral mode only).
// Skipped when using a pre-provisioned TEST_API_KEY account.
const describeMainCleanup = mainAccountRegistered ? describe : describe.skip;

describeMainCleanup('Auth Cleanup — Main Test Account Deletion', () => {
  test('DELETE /auth/account - delete ephemeral test account', async () => {
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: `${testVariables.baseUrl}/auth/account`,
        host: [testVariables.baseUrl],
        path: ['auth', 'account'],
      },
      method: 'DELETE',
      header: [{ key: 'x-api-key', value: mainAccountApiKey }],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          confirmEmail: mainAccountEmail,
          password: mainAccountPassword
        })
      }
    };

    const captured = await captureExample(requestConfig, {}, '/auth/account (main test account)');
    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('success', true);
    console.log(`  ✓ Deleted main test account: ${mainAccountEmail}`);
  });

  test('POST /auth/login - verify main test account cannot login after deletion', async () => {
    const requestConfig: ApiRequestConfig = {
      url: {
        raw: `${testVariables.baseUrl}/auth/login`,
        host: [testVariables.baseUrl],
        path: ['auth', 'login'],
      },
      method: 'POST',
      header: [],
      body: {
        mode: 'raw',
        raw: JSON.stringify({
          email: mainAccountEmail,
          password: mainAccountPassword
        })
      }
    };

    const captured = await captureExample(requestConfig, {}, '/auth/login (deleted main account)');
    expect(captured.response.status).toBe(401);
  });
});
