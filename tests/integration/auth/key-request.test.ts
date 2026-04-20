/**
 * @file key-request.test.ts
 * @description Tests for the OAuth-like app-initiated key request flow
 * @endpoints POST /auth/key-request, GET /auth/key-request/:code/status,
 *   GET /auth/key-request/:code, POST /auth/key-request/:code/approve,
 *   POST /auth/key-request/:code/deny, POST /auth/key-request/exchange
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

const keysToCleanup: string[] = [];

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

function noAuthConfig(method: string, urlPath: string, body?: any): ApiRequestConfig {
  return {
    url: {
      raw: `{{baseUrl}}${urlPath}`,
      host: ['{{baseUrl}}'],
      path: urlPath.split('/').filter(Boolean),
    },
    method: method as any,
    body: body ? { mode: 'raw', raw: JSON.stringify(body) } : undefined,
  };
}

describe('Device Flow', () => {
  let requestCode = '';

  afterAll(async () => {
    // Save captured examples for documentation
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/key-request-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }

    for (const id of keysToCleanup) {
      await makeRequest(replaceVariables(authConfig('DELETE', `/auth/api-keys/${id}`), testVariables)).catch(() => {});
    }
  });

  test('POST /auth/key-request creates pending request with code', async () => {
    const captured = await captureExample(
      noAuthConfig('POST', '/auth/key-request', {
        appName: 'Test Discord Bot',
        appDescription: 'A test integration',
        scopes: ['entity:read', 'roll:read', 'chat:read'],
      }),
      testVariables,
      '/auth/key-request - Create key request'
    );

    expect(captured.response.status).toBe(201);
    expect(captured.response.data).toHaveProperty('code');
    expect(captured.response.data.code).toHaveLength(6);
    expect(captured.response.data).toHaveProperty('approvalUrl');
    expect(captured.response.data).toHaveProperty('expiresIn');

    requestCode = captured.response.data.code;

    capturedExamples.push(captured);
  });

  test('GET /auth/key-request/:code/status returns "pending" before approval', async () => {
    if (!requestCode) return;

    const captured = await captureExample(
      noAuthConfig('GET', `/auth/key-request/${requestCode}/status`),
      testVariables,
      '/auth/key-request/:code/status - Poll status (pending)'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('status', 'pending');

    capturedExamples.push(captured);
  });

  test('GET /auth/key-request/:code returns details (authenticated)', async () => {
    if (!requestCode) return;

    const captured = await captureExample(
      authConfig('GET', `/auth/key-request/${requestCode}`),
      testVariables,
      '/auth/key-request/:code - Get request details'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('appName', 'Test Discord Bot');
    expect(captured.response.data).toHaveProperty('requestedScopes');
    expect(captured.response.data.requestedScopes).toContain('entity:read');

    capturedExamples.push(captured);
  });

  test('POST /auth/key-request/:code/approve creates scoped key', async () => {
    if (!requestCode) return;

    const captured = await captureExample(
      authConfig('POST', `/auth/key-request/${requestCode}/approve`, {
        scopes: ['entity:read', 'roll:read'], // Remove chat:read
      }),
      testVariables,
      '/auth/key-request/:code/approve - Approve request'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('success', true);
    expect(captured.response.data).toHaveProperty('keyId');
    keysToCleanup.push(String(captured.response.data.keyId));

    capturedExamples.push(captured);
  });

  test('GET /auth/key-request/:code/status returns "approved" with apiKey', async () => {
    if (!requestCode) return;

    const captured = await captureExample(
      noAuthConfig('GET', `/auth/key-request/${requestCode}/status`),
      testVariables,
      '/auth/key-request/:code/status - Poll status (approved)'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('status', 'approved');
    expect(captured.response.data).toHaveProperty('apiKey');
    expect(captured.response.data).toHaveProperty('scopes');
    expect(captured.response.data.scopes).toContain('entity:read');
    expect(captured.response.data.scopes).toContain('roll:read');
    expect(captured.response.data.scopes).not.toContain('chat:read'); // Was removed during approval

    capturedExamples.push(captured);
  });

  test('polling after key retrieval returns "exchanged" (no key)', async () => {
    if (!requestCode) return;

    const response = await makeRequest(replaceVariables(noAuthConfig('GET', `/auth/key-request/${requestCode}/status`), testVariables));
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('status', 'exchanged');
    expect(response.data).not.toHaveProperty('apiKey');
  });
});

describe('Denial Flow', () => {
  test('POST /auth/key-request/:code/deny marks as denied', async () => {
    // Create a request
    const createResponse = await makeRequest(replaceVariables(noAuthConfig('POST', '/auth/key-request', {
      appName: 'Denied App',
      scopes: ['entity:read'],
    }), testVariables));
    if (createResponse.status !== 201) return;

    const code = createResponse.data.code;

    // Deny it
    const captured = await captureExample(
      authConfig('POST', `/auth/key-request/${code}/deny`),
      testVariables,
      '/auth/key-request/:code/deny - Deny request'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('success', true);
    capturedExamples.push(captured);

    // Poll — should show denied
    const statusResponse = await makeRequest(replaceVariables(noAuthConfig('GET', `/auth/key-request/${code}/status`), testVariables));
    expect(statusResponse.status).toBe(200);
    expect(statusResponse.data).toHaveProperty('status', 'denied');
  });
});

describe('Approval Controls', () => {
  test('approver cannot grant scopes not in original request', async () => {
    const createResponse = await makeRequest(replaceVariables(noAuthConfig('POST', '/auth/key-request', {
      appName: 'Limited App',
      scopes: ['entity:read'],
    }), testVariables));
    if (createResponse.status !== 201) return;

    const code = createResponse.data.code;

    // Try to approve with extra scopes
    const response = await makeRequest(replaceVariables(authConfig('POST', `/auth/key-request/${code}/approve`, {
      scopes: ['entity:read', 'execute-js'], // execute-js was NOT requested
    }), testVariables));

    expect(response.status).toBe(400);
    expect(response.data.error).toContain('execute-js');
  });
});

describe('Expired Request', () => {
  test('POST /auth/key-request with missing fields returns 400', async () => {
    const response = await makeRequest(replaceVariables(noAuthConfig('POST', '/auth/key-request', {
      appName: '', // empty
      scopes: [],  // empty
    }), testVariables));
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });

  test('POST /auth/key-request with invalid scopes returns 400', async () => {
    const response = await makeRequest(replaceVariables(noAuthConfig('POST', '/auth/key-request', {
      appName: 'Bad Scopes App',
      scopes: ['not-a-real-scope'],
    }), testVariables));
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });
});
