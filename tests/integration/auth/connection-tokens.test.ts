/**
 * @file connection-tokens.test.ts
 * @description Tests for connection token pairing, auth-via-first-message, and audit logs
 * @endpoints POST /auth/connection-tokens, POST /auth/pair, GET /auth/connection-tokens,
 *   DELETE /auth/connection-tokens/:id, GET /auth/connection-logs
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];


/** Helper to make authenticated requests using template variables */
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

describe('Pairing Flow', () => {
  let pairingCode = '';
  let connectionToken = '';
  let pairedClientId = '';

  afterAll(async () => {
    // Save captured examples for documentation
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/connection-tokens-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }

    if (pairedClientId) {
      // Delete all connection tokens created for the test clientId
      const tokenListRes = await makeRequest(replaceVariables(authConfig('GET', '/auth/connection-tokens'), testVariables)).catch(() => null);
      if (tokenListRes?.status === 200 && Array.isArray(tokenListRes.data?.tokens)) {
        const testTokens = tokenListRes.data.tokens.filter((t: any) => t.clientId === pairedClientId);
        for (const t of testTokens) {
          await makeRequest(replaceVariables(authConfig('DELETE', `/auth/connection-tokens/${t.id}`), testVariables)).catch(() => {});
        }
        if (testTokens.length > 0) {
          console.log(`  ✓ Deleted ${testTokens.length} test connection token(s) for ${pairedClientId}`);
        }
      }

      // Delete the KnownClient row created by pairing with worldId='test-world-integration'
      const clientListRes = await makeRequest(replaceVariables(authConfig('GET', '/auth/known-clients'), testVariables)).catch(() => null);
      if (clientListRes?.status === 200 && Array.isArray(clientListRes.data?.clients)) {
        const testClient = clientListRes.data.clients.find((c: any) => c.clientId === pairedClientId);
        if (testClient?.id) {
          await makeRequest(replaceVariables(authConfig('DELETE', `/auth/known-clients/${testClient.id}`), testVariables)).catch(() => {});
          console.log(`  ✓ Deleted test KnownClient: ${pairedClientId} (worldId=test-world-integration)`);
        }
      }
    }
  });

  test('POST /auth/connection-tokens generates a pairing code', async () => {
    const captured = await captureExample(
      authConfig('POST', '/auth/connection-tokens'),
      testVariables,
      '/auth/connection-tokens - Generate pairing code'
    );

    expect(captured.response.status).toBe(201);
    expect(captured.response.data).toHaveProperty('code');
    expect(captured.response.data.code).toHaveLength(6);
    expect(captured.response.data).toHaveProperty('expiresIn', 600);
    pairingCode = captured.response.data.code;

    capturedExamples.push(captured);
  });

  test('POST /auth/pair with valid code returns token, clientId, and relayUrl', async () => {
    if (!pairingCode) return;

    const requestConfig: ApiRequestConfig = {
      url: { raw: '{{baseUrl}}/auth/pair', host: ['{{baseUrl}}'], path: ['auth', 'pair'] },
      method: 'POST',
      body: { mode: 'raw', raw: JSON.stringify({ code: pairingCode, worldId: 'test-world-integration', worldTitle: 'Test World' }) },
    };

    const captured = await captureExample(requestConfig, testVariables, '/auth/pair - Complete pairing');

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('token');
    expect(captured.response.data).toHaveProperty('clientId');
    expect(captured.response.data).toHaveProperty('relayUrl');
    expect(captured.response.data.token).toHaveLength(64); // 32 bytes hex
    expect(captured.response.data.clientId).toMatch(/^fvtt_[a-z0-9]{16}$/);

    connectionToken = captured.response.data.token;
    pairedClientId = captured.response.data.clientId;

    capturedExamples.push(captured);
  });

  test('returned clientId is opaque format, not guessable', async () => {
    expect(pairedClientId).toMatch(/^fvtt_/);
    expect(pairedClientId).not.toContain('foundry-'); // Not the old format
  });

  test('POST /auth/pair with invalid code returns 404', async () => {
    const config: ApiRequestConfig = {
      url: { raw: '{{baseUrl}}/auth/pair', host: ['{{baseUrl}}'], path: ['auth', 'pair'] },
      method: 'POST',
      body: { mode: 'raw', raw: JSON.stringify({ code: 'XXXXXX' }) },
    };

    const response = await makeRequest(replaceVariables(config, testVariables));
    expect(response.status).toBe(404);
    expect(response.data).toHaveProperty('error');
  });

  test('POST /auth/pair with used code returns 410', async () => {
    if (!pairingCode) return;

    const config: ApiRequestConfig = {
      url: { raw: '{{baseUrl}}/auth/pair', host: ['{{baseUrl}}'], path: ['auth', 'pair'] },
      method: 'POST',
      body: { mode: 'raw', raw: JSON.stringify({ code: pairingCode, worldId: 'test-world-integration', worldTitle: 'Test World' }) },
    };

    const response = await makeRequest(replaceVariables(config, testVariables));
    expect(response.status).toBe(410);
    expect(response.data).toHaveProperty('error');
  });
});

describe('Connection Token Management', () => {
  test('GET /auth/connection-tokens lists tokens (no raw values)', async () => {
    const captured = await captureExample(
      authConfig('GET', '/auth/connection-tokens'),
      testVariables,
      '/auth/connection-tokens - List tokens'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('tokens');
    expect(Array.isArray(captured.response.data.tokens)).toBe(true);

    // Tokens should never contain raw token values
    for (const token of captured.response.data.tokens) {
      expect(token).not.toHaveProperty('tokenHash');
      expect(token).toHaveProperty('id');
      expect(token).toHaveProperty('name');
    }

    capturedExamples.push(captured);
  });

  test('connection token cannot be used as X-API-Key header (401)', async () => {
    // Generate a new pairing and get a token
    const codeResponse = await makeRequest(replaceVariables(authConfig('POST', '/auth/connection-tokens'), testVariables));
    if (codeResponse.status !== 201) return;

    const pairResponse = await makeRequest(replaceVariables({
      url: { raw: '{{baseUrl}}/auth/pair', host: ['{{baseUrl}}'], path: ['auth', 'pair'] },
      method: 'POST',
      body: { mode: 'raw', raw: JSON.stringify({ code: codeResponse.data.code, worldId: 'test-world-integration', worldTitle: 'Test World' }) },
    }, testVariables));
    if (pairResponse.status !== 200) return;

    const connToken = pairResponse.data.token;

    // Try using connection token as HTTP API key — should fail
    const apiResponse = await makeRequest(replaceVariables({
      url: { raw: '{{baseUrl}}/world-info', host: ['{{baseUrl}}'], path: ['world-info'] },
      method: 'GET',
      header: [{ key: 'x-api-key', value: connToken }],
    }, testVariables));

    expect(apiResponse.status).toBe(401);
  });
});

describe('Connection Audit Log', () => {
  test('GET /auth/connection-logs returns paginated logs', async () => {
    const captured = await captureExample(
      authConfig('GET', '/auth/connection-logs'),
      testVariables,
      '/auth/connection-logs - List connection logs'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('logs');
    expect(captured.response.data).toHaveProperty('limit');
    expect(captured.response.data).toHaveProperty('offset');

    capturedExamples.push(captured);
  });
});
