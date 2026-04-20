/**
 * @file admin-sessions.test.ts
 * @description Admin headless ChromeDP session management.
 * @endpoints GET /admin/api/headless-sessions, DELETE /admin/api/headless-sessions/{id}
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Headless Sessions', () => {
  let session: AdminSession;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/headless-sessions returns list with session details', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/headless-sessions' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('total');
    expect(typeof response.data.total).toBe('number');
    expect(response.data).toHaveProperty('sessions');
    expect(Array.isArray(response.data.sessions)).toBe(true);

    if (response.data.sessions.length > 0) {
      const s = response.data.sessions[0];
      expect(s).toHaveProperty('sessionId');
      expect(s).toHaveProperty('clientId');
      expect(s).toHaveProperty('foundryUrl');
      expect(s).toHaveProperty('startedAt');
    }
  });

  test('GET /admin/api/headless-sessions does not leak credentials', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/headless-sessions' },
      session
    );
    expect(response.status).toBe(200);

    // PII safety: never return Foundry username/password
    if (response.data.sessions.length > 0) {
      const s = response.data.sessions[0];
      expect(s).not.toHaveProperty('username');
      expect(s).not.toHaveProperty('password');
      expect(s).not.toHaveProperty('foundryPassword');
      expect(s).not.toHaveProperty('foundryUsername');
    }
  });

  test('DELETE /admin/api/headless-sessions/{id} returns 404 for non-existent session', async () => {
    const response = await makeAdminRequest(
      { method: 'DELETE', path: '/admin/api/headless-sessions/non-existent-session-id' },
      session
    );
    expect(response.status).toBe(404);
    expect(response.data).toHaveProperty('error');
  });

  test('non-admin cannot access /admin/api/headless-sessions', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/headless-sessions`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
