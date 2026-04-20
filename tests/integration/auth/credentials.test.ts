/**
 * @file credentials.test.ts
 * @description Tests for credential vault and known clients registry
 * @endpoints POST /auth/credentials, GET /auth/credentials, PATCH /auth/credentials/:id,
 *   DELETE /auth/credentials/:id, GET /auth/known-clients, DELETE /auth/known-clients/:id
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

const credsToCleanup: string[] = [];

function authConfig(method: string, urlPath: string, body?: any): ApiRequestConfig {
  return {
    url: {
      raw: `{{baseUrl}}${urlPath}`,
      host: ['{{baseUrl}}'],
      path: urlPath.split('/').filter(Boolean),
    },
    method: method as any,
    header: [{ key: 'x-api-key', value: '{{apiKey}}' }],
    body: body ? { mode: 'raw', raw: JSON.stringify(body) } : undefined,
  };
}

describe('Credential Vault', () => {
  let credentialId = '';

  afterAll(async () => {
    // Save captured examples for documentation
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/credentials-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }

    for (const id of credsToCleanup) {
      await makeRequest(authConfig('DELETE', `/auth/credentials/${id}`)).catch(() => {});
    }
  });

  test('POST /auth/credentials creates credential set', async () => {
    const captured = await captureExample(
      authConfig('POST', '/auth/credentials', {
        name: 'Test Foundry Server',
        foundryUrl: 'http://localhost:30000',
        foundryUsername: 'testgm',
        foundryPassword: 'testpassword123',
      }),
      testVariables,
      '/auth/credentials - Create credential set'
    );

    // May fail if CREDENTIALS_ENCRYPTION_KEY not set
    if (captured.response.status === 400 && captured.response.data?.error?.includes('CREDENTIALS_ENCRYPTION_KEY')) {
      console.log('Credential encryption not configured — skipping credential vault tests');
      return;
    }

    expect(captured.response.status).toBe(201);
    expect(captured.response.data).toHaveProperty('id');
    expect(captured.response.data).toHaveProperty('name', 'Test Foundry Server');
    expect(captured.response.data).toHaveProperty('foundryUrl', 'http://localhost:30000');
    expect(captured.response.data).toHaveProperty('foundryUsername', 'testgm');

    credentialId = String(captured.response.data.id);
    credsToCleanup.push(credentialId);

    capturedExamples.push(captured);
  });

  test('GET /auth/credentials lists credentials without passwords', async () => {
    const captured = await captureExample(
      authConfig('GET', '/auth/credentials'),
      testVariables,
      '/auth/credentials - List credentials'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('credentials');
    expect(Array.isArray(captured.response.data.credentials)).toBe(true);

    // Verify no password data is returned
    for (const cred of captured.response.data.credentials) {
      expect(cred).not.toHaveProperty('encryptedFoundryPassword');
      expect(cred).not.toHaveProperty('passwordIv');
      expect(cred).not.toHaveProperty('passwordAuthTag');
      expect(cred).not.toHaveProperty('foundryPassword');
    }

    capturedExamples.push(captured);
  });

  test('PATCH /auth/credentials/:id updates name/url/username', async () => {
    if (!credentialId) return;

    const captured = await captureExample(
      authConfig('PATCH', `/auth/credentials/${credentialId}`, {
        name: 'Updated Server Name',
        foundryUrl: 'http://localhost:31000',
        foundryPassword: 'testpassword123',
      }),
      testVariables,
      '/auth/credentials/:id - Update credential'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('name', 'Updated Server Name');
    expect(captured.response.data).toHaveProperty('foundryUrl', 'http://localhost:31000');

    capturedExamples.push(captured);
  });

  test('DELETE /auth/credentials/:id removes credential', async () => {
    if (!credentialId) return;

    const captured = await captureExample(
      authConfig('DELETE', `/auth/credentials/${credentialId}`),
      testVariables,
      '/auth/credentials/:id - Delete credential'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('success', true);

    // Remove from cleanup since we already deleted
    const idx = credsToCleanup.indexOf(credentialId);
    if (idx >= 0) credsToCleanup.splice(idx, 1);

    capturedExamples.push(captured);
  });
});

describe('Known Clients', () => {
  test('GET /auth/known-clients returns list with online/offline status', async () => {
    const captured = await captureExample(
      authConfig('GET', '/auth/known-clients'),
      testVariables,
      '/auth/known-clients - List known clients'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('clients');
    expect(Array.isArray(captured.response.data.clients)).toBe(true);

    // If any clients exist, they should have the expected fields
    if (captured.response.data.clients.length > 0) {
      const client = captured.response.data.clients[0];
      expect(client).toHaveProperty('clientId');
      expect(client).toHaveProperty('isOnline');
    }

    capturedExamples.push(captured);
  });
});
