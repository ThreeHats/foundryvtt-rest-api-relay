/**
 * @file admin-health.test.ts
 * @description Admin extended system health endpoint.
 * @endpoints GET /admin/api/system/health
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin System Health', () => {
  let session: AdminSession;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/system/health returns status and version', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/system/health' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('status');
    expect(['ok', 'degraded']).toContain(response.data.status);
    expect(response.data).toHaveProperty('version');
    expect(typeof response.data.version).toBe('string');
    expect(response.data).toHaveProperty('timestamp');
  });

  test('GET /admin/api/system/health returns runtime metrics', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/system/health' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('goroutines');
    expect(typeof response.data.goroutines).toBe('number');
    expect(response.data.goroutines).toBeGreaterThan(0);
    expect(response.data).toHaveProperty('memoryAllocBytes');
    expect(typeof response.data.memoryAllocBytes).toBe('number');
    expect(response.data).toHaveProperty('memoryHeapBytes');
    expect(response.data).toHaveProperty('memorySysBytes');
    expect(response.data).toHaveProperty('gcCount');
  });

  test('GET /admin/api/system/health returns connection counts', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/system/health' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('wsConnectionCount');
    expect(typeof response.data.wsConnectionCount).toBe('number');
    expect(response.data).toHaveProperty('pendingRequestCount');
    expect(typeof response.data.pendingRequestCount).toBe('number');
    expect(response.data).toHaveProperty('headlessSessionCount');
    expect(response.data).toHaveProperty('redisStatus');
    expect(['connected', 'degraded', 'disabled']).toContain(response.data.redisStatus);
  });

  test('GET /admin/api/system/health does not leak PII', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/system/health' },
      session
    );
    expect(response.status).toBe(200);
    const bodyStr = JSON.stringify(response.data);
    expect(bodyStr).not.toContain(testVariables.adminEmail);
    expect(bodyStr).not.toContain(testVariables.adminPassword);
    // Should not contain API keys
    expect(bodyStr).not.toContain(testVariables.apiKey);
  });

  test('non-admin cannot access /admin/api/system/health', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/system/health`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});
