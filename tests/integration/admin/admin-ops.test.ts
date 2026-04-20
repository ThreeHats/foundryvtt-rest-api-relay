/**
 * @file admin-ops.test.ts
 * @description Admin operational tools — feature flags, maintenance mode.
 * @endpoints GET /admin/api/ops/feature-flags, POST /admin/api/ops/feature-flags
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Operational Tools', () => {
  let session: AdminSession;

  beforeAll(async () => {
    session = await adminLogin();
  });

  afterAll(async () => {
    // Defensive: ensure maintenance_mode is OFF in case a test left it on
    if (session) {
      await makeAdminRequest(
        {
          method: 'POST',
          path: '/admin/api/ops/feature-flags',
          body: { flag: 'maintenance_mode', value: false },
        },
        session
      ).catch(() => {});
    }
  });

  test('GET /admin/api/ops/feature-flags returns current flags', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/ops/feature-flags' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('disable_registration');
    expect(response.data).toHaveProperty('maintenance_mode');
  });

  test('POST /admin/api/ops/feature-flags toggles maintenance_mode', async () => {
    // Enable maintenance mode
    let response = await makeAdminRequest(
      {
        method: 'POST',
        path: '/admin/api/ops/feature-flags',
        body: { flag: 'maintenance_mode', value: true },
      },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data.maintenance_mode).toBe(true);

    // /api/health should still be accessible (exempt)
    const health = await axios.get(`${testVariables.baseUrl}/api/health`, {
      validateStatus: () => true,
    });
    expect(health.status).toBe(200);

    // /clients should be blocked (503 maintenance)
    const clients = await axios.get(`${testVariables.baseUrl}/clients`, {
      headers: { 'x-api-key': testVariables.apiKey },
      validateStatus: () => true,
    });
    expect(clients.status).toBe(503);

    // Disable maintenance mode
    response = await makeAdminRequest(
      {
        method: 'POST',
        path: '/admin/api/ops/feature-flags',
        body: { flag: 'maintenance_mode', value: false },
      },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data.maintenance_mode).toBe(false);
  });

  test('POST /admin/api/ops/feature-flags rejects unknown flags', async () => {
    const response = await makeAdminRequest(
      {
        method: 'POST',
        path: '/admin/api/ops/feature-flags',
        body: { flag: 'totally_made_up_flag', value: true },
      },
      session
    );
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });

  test('POST /admin/api/ops/feature-flags without flag name returns 400', async () => {
    const response = await makeAdminRequest(
      {
        method: 'POST',
        path: '/admin/api/ops/feature-flags',
        body: { value: true }, // flag intentionally missing
      },
      session
    );
    expect(response.status).toBe(400);
    expect(response.data).toHaveProperty('error');
  });

  test('non-admin cannot access /admin/api/ops/feature-flags', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/ops/feature-flags`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
