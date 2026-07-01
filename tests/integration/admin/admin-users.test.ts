/**
 * @file admin-users.test.ts
 * @description Admin user management endpoints — list (search/filter/sort), detail,
 *   disable/enable, subscription/verified edits, rotate-key, password reset, delete.
 * @endpoints GET /admin/api/users (search, role/status/verified/rotation/subscription
 *   filters, sortable columns), GET /admin/api/users/{id}, PATCH /admin/api/users/{id}
 *   (role, emailVerified, subscriptionStatus, maxHeadlessSessions),
 *   POST /admin/api/users/{id}/disable, POST /admin/api/users/{id}/enable,
 *   POST /admin/api/users/{id}/rotate-key, POST /admin/api/users/{id}/send-password-reset,
 *   DELETE /admin/api/users/{id}
 *
 * Includes PII safety checks: list responses must NOT contain password, full apiKey,
 * stripeCustomerId, or verificationTokenHash.
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials && process.env.TEST_SKIP_ADMIN !== 'true' ? describe : describe.skip;

describeAdmin('Admin User Management', () => {
  let session: AdminSession;
  let throwawayUserId: number | null = null;
  let throwawayEmail = '';
  let throwawayApiKey = '';
  let throwawaySessionToken = '';

  beforeAll(async () => {
    session = await adminLogin();

    // Create a throwaway test user via the public registration endpoint
    // (so we have a non-admin user to disable/enable/delete in tests).
    const email = `admin-test-${Date.now()}@example.com`;
    const response = await axios.post(
      `${testVariables.baseUrl}/auth/register`,
      { email, password: 'TestPassword1' },
      { validateStatus: () => true }
    );
    if (response.status === 201) {
      throwawayUserId = response.data.id;
      throwawayEmail = email;
      throwawayApiKey = response.data.apiKey;
      throwawaySessionToken = response.data.sessionToken;
    } else {
      console.log(`Could not create throwaway user (status ${response.status}). Some tests will be skipped.`);
    }
  });

  test('GET /admin/api/users returns paginated list with PII-safe fields', async () => {
    const response = await makeAdminRequest({ method: 'GET', path: '/admin/api/users', query: { limit: 10 } }, session);
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('users');
    expect(response.data).toHaveProperty('total');
    expect(Array.isArray(response.data.users)).toBe(true);

    if (response.data.users.length > 0) {
      const u = response.data.users[0];
      expect(u).toHaveProperty('id');
      expect(u).toHaveProperty('email');
      expect(u).toHaveProperty('role');
      expect(u).toHaveProperty('disabled');
      expect(u).toHaveProperty('subscriptionStatus');
      expect(u).toHaveProperty('requestsToday');

      // PII checks: must NOT contain password, apiKey, stripeCustomerId, verificationTokenHash
      expect(u).not.toHaveProperty('password');
      expect(u).not.toHaveProperty('apiKey');
      expect(u).not.toHaveProperty('stripeCustomerId');
      expect(u).not.toHaveProperty('verificationTokenHash');
    }
  });

  test('GET /admin/api/users supports pagination', async () => {
    const r1 = await makeAdminRequest({ method: 'GET', path: '/admin/api/users', query: { limit: 1, offset: 0 } }, session);
    expect(r1.status).toBe(200);
    expect(r1.data.users.length).toBeLessThanOrEqual(1);
    expect(r1.data.limit).toBe(1);
    expect(r1.data.offset).toBe(0);
  });

  test('GET /admin/api/users/{id} returns user detail', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest({ method: 'GET', path: `/admin/api/users/${throwawayUserId}` }, session);
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('id', throwawayUserId);
    expect(response.data).not.toHaveProperty('password');
    expect(response.data).not.toHaveProperty('apiKey');
  });

  test('POST /admin/api/users/{id}/disable disables user', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'POST', path: `/admin/api/users/${throwawayUserId}/disable` },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('message', 'User disabled');
  });

  test('disabled user session is rejected by user-data endpoint', async () => {
    if (!throwawaySessionToken) return;
    const response = await axios.get(`${testVariables.baseUrl}/auth/user-data`, {
      headers: { 'Authorization': `Bearer ${throwawaySessionToken}` },
      validateStatus: () => true,
    });
    expect(response.status).toBe(403);
  });

  test('POST /admin/api/users/{id}/enable re-enables user', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'POST', path: `/admin/api/users/${throwawayUserId}/enable` },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('message', 'User enabled');
  });

  test('POST /admin/api/users/{id}/rotate-key returns key prefix only (PII safe)', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'POST', path: `/admin/api/users/${throwawayUserId}/rotate-key` },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('keyPrefix');
    // Should be truncated form, not full 64-char hex key
    expect((response.data.keyPrefix as string).length).toBeLessThan(20);
  });

  // --- Search / filter / sort (GET /admin/api/users query params) ---

  test('GET /admin/api/users?search=<email> finds the user by email', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/users', query: { search: throwawayEmail } },
      session
    );
    expect(response.status).toBe(200);
    const match = response.data.users.find((u: any) => u.id === throwawayUserId);
    expect(match).toBeDefined();
    expect(match.email).toBe(throwawayEmail);
    // PII stays stripped under filtered queries too.
    expect(match).not.toHaveProperty('password');
    expect(match).not.toHaveProperty('apiKey');
  });

  test('GET /admin/api/users?search=<id> matches the numeric id', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/users', query: { search: String(throwawayUserId) } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data.users.some((u: any) => u.id === throwawayUserId)).toBe(true);
  });

  test('GET /admin/api/users?role=user returns only non-admin users', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/users', query: { role: 'user', limit: 100 } },
      session
    );
    expect(response.status).toBe(200);
    for (const u of response.data.users) {
      expect(u.role).toBe('user');
    }
  });

  test('GET /admin/api/users?sort=id&dir=desc returns descending ids', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/users', query: { sort: 'id', dir: 'desc', limit: 100 } },
      session
    );
    expect(response.status).toBe(200);
    const ids = response.data.users.map((u: any) => u.id);
    const sorted = [...ids].sort((a, b) => b - a);
    expect(ids).toEqual(sorted);
  });

  test('GET /admin/api/users with an injection-y sort falls back safely (no 500)', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/users', query: { sort: 'id); DROP TABLE Users;--' } },
      session
    );
    expect(response.status).toBe(200);
    expect(Array.isArray(response.data.users)).toBe(true);
  });

  // --- Inline edits via PATCH ---

  test('PATCH /admin/api/users/{id} updates subscriptionStatus', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'PATCH', path: `/admin/api/users/${throwawayUserId}`, body: { subscriptionStatus: 'active' } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('subscriptionStatus', 'active');
  });

  test('PATCH /admin/api/users/{id} rejects an invalid subscriptionStatus', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'PATCH', path: `/admin/api/users/${throwawayUserId}`, body: { subscriptionStatus: 'not-a-real-status' } },
      session
    );
    expect(response.status).toBe(400);
  });

  test('PATCH /admin/api/users/{id} can set emailVerified', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'PATCH', path: `/admin/api/users/${throwawayUserId}`, body: { emailVerified: true } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('emailVerified', true);
  });

  // --- Admin-triggered password reset ---

  test('POST /admin/api/users/{id}/send-password-reset succeeds', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'POST', path: `/admin/api/users/${throwawayUserId}/send-password-reset` },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('message');
  });

  test('POST /admin/api/users/{id}/send-password-reset 404s for an unknown user', async () => {
    const response = await makeAdminRequest(
      { method: 'POST', path: '/admin/api/users/999999999/send-password-reset' },
      session
    );
    expect(response.status).toBe(404);
  });

  test('admin cannot disable themselves', async () => {
    // First fetch /me to get our admin ID
    const me = await makeAdminRequest({ method: 'GET', path: '/admin/auth/me' }, session);
    expect(me.status).toBe(200);
    const myId = me.data.id;
    const response = await makeAdminRequest(
      { method: 'POST', path: `/admin/api/users/${myId}/disable` },
      session
    );
    expect(response.status).toBe(400);
  });

  test('non-admin cannot access /admin/api/users (no cookie)', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/users`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });

  test('DELETE /admin/api/users/{id} removes the user', async () => {
    if (!throwawayUserId) return;
    const response = await makeAdminRequest(
      { method: 'DELETE', path: `/admin/api/users/${throwawayUserId}` },
      session
    );
    expect(response.status).toBe(200);
    // Confirm gone
    const check = await makeAdminRequest(
      { method: 'GET', path: `/admin/api/users/${throwawayUserId}` },
      session
    );
    expect(check.status).toBe(404);
    throwawayUserId = null; // already deleted
  });
});
