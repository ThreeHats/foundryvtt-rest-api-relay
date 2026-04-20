/**
 * @file admin-clients.test.ts
 * @description Admin connected clients endpoints — list, detail, disconnect, broadcast.
 * @endpoints GET /admin/api/clients, GET /admin/api/clients/{id},
 *   POST /admin/api/clients/{id}/disconnect, POST /admin/api/clients/broadcast
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Connected Clients', () => {
  let session: AdminSession;
  let firstClientId: string | null = null;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/clients returns list with total count', async () => {
    const response = await makeAdminRequest({ method: 'GET', path: '/admin/api/clients' }, session);
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('total');
    expect(response.data).toHaveProperty('clients');
    expect(Array.isArray(response.data.clients)).toBe(true);
    expect(typeof response.data.total).toBe('number');

    if (response.data.clients.length > 0) {
      const client = response.data.clients[0];
      firstClientId = client.clientId || client.id;
      expect(client).toHaveProperty('clientId');
    }
  });

  test('GET /admin/api/clients/{id} returns client detail when client exists', async () => {
    if (!firstClientId) {
      console.log('  Skipping: no connected clients');
      return;
    }

    const response = await makeAdminRequest(
      { method: 'GET', path: `/admin/api/clients/${firstClientId}` },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('clientId', firstClientId);
  });

  test('GET /admin/api/clients/{id} returns 404 for unknown id', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/clients/nonexistent-fake-id' },
      session
    );
    expect(response.status).toBe(404);
    expect(response.data).toHaveProperty('error');
  });

  test('non-admin cannot access /admin/api/clients (no cookie)', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/clients`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
