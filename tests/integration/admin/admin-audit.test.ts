/**
 * @file admin-audit.test.ts
 * @description Admin audit log query endpoint.
 * @endpoints GET /admin/api/audit-logs
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Audit Log', () => {
  let session: AdminSession;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/audit-logs returns paginated entries', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/audit-logs', query: { limit: 20 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('entries');
    expect(response.data).toHaveProperty('total');
    expect(Array.isArray(response.data.entries)).toBe(true);
  });

  test('GET /admin/api/audit-logs supports pagination params', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/audit-logs', query: { limit: 1, offset: 0 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data.entries.length).toBeLessThanOrEqual(1);
    expect(response.data).toHaveProperty('limit', 1);
    expect(response.data).toHaveProperty('offset', 0);
  });

  test('login from this test session creates an admin.login audit entry', async () => {
    // We just logged in via beforeAll. Filter for admin.login.
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/audit-logs', query: { action: 'admin.login', limit: 5 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data.entries.length).toBeGreaterThan(0);
    const entry = response.data.entries[0];
    expect(entry.action).toBe('admin.login');
    // Sensitive content checks
    const entryStr = JSON.stringify(entry);
    expect(entryStr).not.toContain(testVariables.adminPassword);
  });

  test('filter by adminUserId works', async () => {
    const me = await makeAdminRequest({ method: 'GET', path: '/admin/auth/me' }, session);
    expect(me.status).toBe(200);
    const myId = me.data.id;
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/audit-logs', query: { adminUserId: myId, limit: 10 } },
      session
    );
    expect(response.status).toBe(200);
    if (response.data.entries.length > 0) {
      expect(response.data.entries[0].adminUserId).toBe(myId);
    }
  });

  test('audit log entries do not leak PII', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/audit-logs', query: { limit: 10 } },
      session
    );
    expect(response.status).toBe(200);
    const bodyStr = JSON.stringify(response.data);
    expect(bodyStr).not.toContain(testVariables.adminPassword);
    // Should not contain full API keys
    if (testVariables.apiKey.length >= 16) {
      expect(bodyStr).not.toContain(testVariables.apiKey);
    }
  });

  test('non-admin cannot access /admin/api/audit-logs', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/audit-logs`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
