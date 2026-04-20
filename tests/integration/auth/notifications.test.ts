/**
 * @file notifications.test.ts
 * @description Tests for connection notification settings
 * @endpoints GET /auth/notification-settings, PUT /auth/notification-settings,
 *   POST /auth/notification-settings/test
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

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

describe('Notification Settings', () => {
  afterAll(async () => {
    // Save captured examples for documentation
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/notifications-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }

    // Reset to defaults for other tests
    await makeRequest(replaceVariables(authConfig('PUT', '/auth/notification-settings', {
      notifyOnConnect: true,
      discordWebhookUrl: '',
      notifyEmail: '',
    }), testVariables)).catch(() => {});
  });

  test('GET /auth/notification-settings returns defaults', async () => {
    const captured = await captureExample(
      authConfig('GET', '/auth/notification-settings'),
      testVariables,
      '/auth/notification-settings - Get defaults'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('notifyOnConnect');
    // Default should be true
    expect(captured.response.data.notifyOnConnect).toBe(true);
    expect(captured.response.data).toHaveProperty('smtpAvailable');

    capturedExamples.push(captured);
  });

  test('PUT /auth/notification-settings updates webhook URL', async () => {
    const captured = await captureExample(
      authConfig('PUT', '/auth/notification-settings', {
        notifyOnConnect: true,
        discordWebhookUrl: 'https://discord.com/api/webhooks/test/test',
        notifyEmail: '',
      }),
      testVariables,
      '/auth/notification-settings - Update webhook'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('discordWebhookUrl', 'https://discord.com/api/webhooks/test/test');
    expect(captured.response.data).toHaveProperty('notifyOnConnect', true);

    capturedExamples.push(captured);
  });

  test('PUT /auth/notification-settings can disable notifications', async () => {
    const captured = await captureExample(
      authConfig('PUT', '/auth/notification-settings', {
        notifyOnConnect: false,
      }),
      testVariables,
      '/auth/notification-settings - Disable notifications'
    );

    expect(captured.response.status).toBe(200);
    expect(captured.response.data).toHaveProperty('notifyOnConnect', false);

    capturedExamples.push(captured);
  });

  test('GET /auth/notification-settings reflects updated values', async () => {
    const response = await makeRequest(replaceVariables(authConfig('GET', '/auth/notification-settings'), testVariables));
    expect(response.status).toBe(200);
    // Should reflect the most recent PUT
    expect(response.data).toHaveProperty('notifyOnConnect', false);
  });
});
