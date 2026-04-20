/**
 * @file account-security.test.ts
 * @description Tests for email verification
 * @endpoints POST /auth/register, GET /auth/verify, POST /auth/resend-verification
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Whether email verification is configured depends on server config
// Tests detect this dynamically

describe('Email Verification', () => {
  const testEmail = `verify-test-${Date.now()}@example.com`;
  const testPassword = 'TestPassword123!';
  let registeredApiKey = '';
  let registrationDisabled = false;

  afterAll(async () => {
    // Save captured examples for documentation
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/account-security-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }

    // Cleanup: delete test user
    if (!registrationDisabled && registeredApiKey) {
      const config: ApiRequestConfig = {
        url: { raw: '{{baseUrl}}/auth/account', host: ['{{baseUrl}}'], path: ['auth', 'account'] },
        method: 'DELETE',
        header: [{ key: 'x-api-key', value: registeredApiKey }],
        body: { mode: 'raw', raw: JSON.stringify({ confirmEmail: testEmail, password: testPassword }) },
      };
      await makeRequest(replaceVariables(config, testVariables)).catch(() => {});
    }
  });

  test('POST /auth/register returns user with emailVerified field', async () => {
    const requestConfig: ApiRequestConfig = {
      url: { raw: '{{baseUrl}}/auth/register', host: ['{{baseUrl}}'], path: ['auth', 'register'] },
      method: 'POST',
      body: { mode: 'raw', raw: JSON.stringify({ email: testEmail, password: testPassword }) },
    };

    const captured = await captureExample(requestConfig, testVariables, '/auth/register - with email verification');

    if (captured.response.status === 403) {
      registrationDisabled = true;
      console.log('Registration disabled — skipping email verification tests');
      return;
    }

    expect(captured.response.status).toBe(201);
    expect(captured.response.data).toHaveProperty('emailVerified');
    registeredApiKey = captured.response.data.apiKey;

    if (captured.response.status < 400) {
      capturedExamples.push(captured);
    }
  });

  test('GET /auth/verify with invalid token returns 404', async () => {
    if (registrationDisabled) return;

    const requestConfig: ApiRequestConfig = {
      url: {
        raw: '{{baseUrl}}/auth/verify',
        host: ['{{baseUrl}}'],
        path: ['auth', 'verify'],
        query: [{ key: 'token', value: 'invalid-token-abc123' }],
      },
      method: 'GET',
    };

    const captured = await captureExample(requestConfig, testVariables, '/auth/verify - invalid token');

    // Should be 404 or 400 for invalid token
    expect([400, 404]).toContain(captured.response.status);
    expect(captured.response.data).toHaveProperty('error');
  });
});
