/**
 * @file admin-keys.test.ts
 * @description Admin API key management — list across all users, detail, enable/disable, revoke.
 * @endpoints GET /admin/api/keys, GET /admin/api/keys/{id}, PATCH /admin/api/keys/{id},
 *   DELETE /admin/api/keys/{id}
 *
 * PII safety: responses must contain keyPrefix only, NEVER full key value.
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin API Key Management', () => {
  let session: AdminSession;
  let firstKeyId: number | null = null;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/keys returns paginated list with PII-safe fields', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/keys', query: { limit: 10 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('keys');
    expect(response.data).toHaveProperty('total');
    expect(Array.isArray(response.data.keys)).toBe(true);

    if (response.data.keys.length > 0) {
      const k = response.data.keys[0];
      firstKeyId = k.id;

      expect(k).toHaveProperty('id');
      expect(k).toHaveProperty('userId');
      expect(k).toHaveProperty('keyPrefix');
      expect(k).toHaveProperty('enabled');
      expect(k).toHaveProperty('scopes');

      // PII checks: must NOT contain the full key
      expect(k).not.toHaveProperty('key');
      expect((k.keyPrefix as string).length).toBeLessThan(20);
      // Must NOT contain encrypted password fields
      expect(k).not.toHaveProperty('encryptedFoundryPassword');
      expect(k).not.toHaveProperty('passwordIv');
      expect(k).not.toHaveProperty('passwordAuthTag');
    }
  });

  test('GET /admin/api/keys supports pagination', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/keys', query: { limit: 1, offset: 0 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data.keys.length).toBeLessThanOrEqual(1);
    expect(response.data).toHaveProperty('limit', 1);
    expect(response.data).toHaveProperty('offset', 0);
  });

  test('GET /admin/api/keys/{id} returns key detail', async () => {
    if (!firstKeyId) return;

    const response = await makeAdminRequest(
      { method: 'GET', path: `/admin/api/keys/${firstKeyId}` },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('id', firstKeyId);
    expect(response.data).toHaveProperty('keyPrefix');
    expect(response.data).toHaveProperty('enabled');

    // PII: no full key
    expect(response.data).not.toHaveProperty('key');
  });

  test('GET /admin/api/keys/{id} returns 404 for unknown id', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/keys/999999' },
      session
    );
    expect(response.status).toBe(404);
    expect(response.data).toHaveProperty('error');
  });

  test('PATCH /admin/api/keys/{id} can update scopes', async () => {
    if (!firstKeyId) return;

    // First get current scopes
    const before = await makeAdminRequest(
      { method: 'GET', path: `/admin/api/keys/${firstKeyId}` },
      session
    );
    const originalScopes = before.data.scopes || [];

    // Update scopes
    const response = await makeAdminRequest(
      {
        method: 'PATCH',
        path: `/admin/api/keys/${firstKeyId}`,
        body: { scopes: ['entity:read', 'roll:read'] },
      },
      session
    );
    expect(response.status).toBe(200);

    // Restore original scopes if update succeeded
    if (response.status === 200) {
      await makeAdminRequest(
        {
          method: 'PATCH',
          path: `/admin/api/keys/${firstKeyId}`,
          body: { scopes: originalScopes },
        },
        session
      );
    }
  });

  test('non-admin cannot access /admin/api/keys', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/keys`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });

  test('non-admin cannot access /admin/api/keys with x-api-key (wrong auth method)', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/keys`, {
      headers: { 'x-api-key': testVariables.apiKey },
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
