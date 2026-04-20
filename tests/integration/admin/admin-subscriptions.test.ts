/**
 * @file admin-subscriptions.test.ts
 * @description Admin Stripe subscription overview.
 * @endpoints GET /admin/api/subscriptions
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Subscriptions', () => {
  let session: AdminSession;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/subscriptions returns overview with statusCounts', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/subscriptions' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('statusCounts');
    expect(typeof response.data.statusCounts).toBe('object');
    expect(response.data).toHaveProperty('subscribers');
    expect(Array.isArray(response.data.subscribers)).toBe(true);
  });

  test('GET /admin/api/subscriptions supports pagination', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/subscriptions', query: { limit: 1, offset: 0 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('total');
    expect(response.data).toHaveProperty('limit');
    expect(response.data).toHaveProperty('offset');
    expect(response.data.subscribers.length).toBeLessThanOrEqual(1);
  });

  test('GET /admin/api/subscriptions does not expose PII', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/subscriptions' },
      session
    );
    expect(response.status).toBe(200);

    // PII safety: never expose Stripe customer IDs
    const bodyStr = JSON.stringify(response.data);
    expect(bodyStr).not.toMatch(/cus_[A-Za-z0-9]+/);

    // Subscribers should not contain passwords or full API keys
    for (const sub of response.data.subscribers) {
      expect(sub).not.toHaveProperty('password');
      expect(sub).not.toHaveProperty('apiKey');
    }
  });

  test('non-admin cannot access /admin/api/subscriptions', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/subscriptions`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
