/**
 * @file admin-alerts.test.ts
 * @description Admin alert config (unified URL + checkbox matrix) and recent events.
 * @endpoints GET /admin/api/alerts/config, PUT /admin/api/alerts/config,
 *   POST /admin/api/alerts/test, GET /admin/api/alerts/recent
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Alert Config', () => {
  let session: AdminSession;
  // Original config restored in afterAll
  let originalConfig: { discordWebhookUrl: string; emailDestination: string; subscriptions: any[] } | null = null;

  beforeAll(async () => {
    session = await adminLogin();
    // Save current config so we can restore it
    const r = await makeAdminRequest({ method: 'GET', path: '/admin/api/alerts/config' }, session);
    if (r.status === 200) originalConfig = r.data;
  });

  afterAll(async () => {
    if (originalConfig) {
      await makeAdminRequest(
        { method: 'PUT', path: '/admin/api/alerts/config', body: originalConfig },
        session
      ).catch(() => {});
    }
  });

  test('GET /admin/api/alerts/config returns expected shape', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/alerts/config' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('discordWebhookUrl');
    expect(response.data).toHaveProperty('emailDestination');
    expect(response.data).toHaveProperty('subscriptions');
    expect(Array.isArray(response.data.subscriptions)).toBe(true);
  });

  test('PUT /admin/api/alerts/config saves and round-trips', async () => {
    const body = {
      discordWebhookUrl: 'https://discord.com/api/webhooks/test/abc123',
      emailDestination: 'alerts-test@example.com',
      subscriptions: [
        { alertType: 'admin_login', channel: 'discord' },
        { alertType: 'failed_auth_spike', channel: 'email' },
      ],
    };

    const putResponse = await makeAdminRequest(
      { method: 'PUT', path: '/admin/api/alerts/config', body },
      session
    );
    expect(putResponse.status).toBe(200);
    expect(putResponse.data.discordWebhookUrl).toBe(body.discordWebhookUrl);
    expect(putResponse.data.emailDestination).toBe(body.emailDestination);

    // Verify round-trip via GET
    const getResponse = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/alerts/config' },
      session
    );
    expect(getResponse.status).toBe(200);
    expect(getResponse.data.discordWebhookUrl).toBe(body.discordWebhookUrl);
    expect(getResponse.data.emailDestination).toBe(body.emailDestination);

    const subs = getResponse.data.subscriptions as any[];
    expect(subs.some((s: any) => s.alertType === 'admin_login' && s.channel === 'discord')).toBe(true);
    expect(subs.some((s: any) => s.alertType === 'failed_auth_spike' && s.channel === 'email')).toBe(true);
  });

  test('PUT /admin/api/alerts/config with invalid alertType returns 400', async () => {
    const response = await makeAdminRequest(
      {
        method: 'PUT',
        path: '/admin/api/alerts/config',
        body: {
          discordWebhookUrl: '',
          emailDestination: '',
          subscriptions: [{ alertType: 'totally_fake_alert', channel: 'discord' }],
        },
      },
      session
    );
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });

  test('PUT /admin/api/alerts/config with invalid channel returns 400', async () => {
    const response = await makeAdminRequest(
      {
        method: 'PUT',
        path: '/admin/api/alerts/config',
        body: {
          discordWebhookUrl: '',
          emailDestination: '',
          subscriptions: [{ alertType: 'admin_login', channel: 'sms' }],
        },
      },
      session
    );
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });

  test('POST /admin/api/alerts/test with invalid channel returns 400', async () => {
    const response = await makeAdminRequest(
      {
        method: 'POST',
        path: '/admin/api/alerts/test',
        body: { channel: 'sms' },
      },
      session
    );
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });

  test('GET /admin/api/alerts/recent returns events array (PII safe)', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/alerts/recent' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('events');
    expect(Array.isArray(response.data.events)).toBe(true);

    if (response.data.events.length > 0) {
      const bodyStr = JSON.stringify(response.data);
      expect(bodyStr).not.toContain(testVariables.adminEmail);
    }
  });

  test('non-admin cannot access /admin/api/alerts/config', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/alerts/config`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
