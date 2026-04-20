/**
 * @file admin-metrics.test.ts
 * @description Admin metrics — both in-app rolling counters and Prometheus /metrics.
 * @endpoints GET /admin/api/metrics/overview, GET /admin/api/metrics/by-endpoint,
 *   GET /admin/api/metrics/top-consumers, GET /metrics
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';
const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin In-App Metrics', () => {
  let session: AdminSession;

  beforeAll(async () => {
    session = await adminLogin();
  });

  test('GET /admin/api/metrics/overview returns request counters', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/metrics/overview' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('requestsPerMinute');
    expect(response.data).toHaveProperty('requestsPerHour');
    expect(response.data).toHaveProperty('requestsPerDay');
    expect(response.data).toHaveProperty('errorsTotal');
    // All should be numbers
    expect(typeof response.data.requestsPerMinute).toBe('number');
    expect(typeof response.data.requestsPerDay).toBe('number');
  });

  test('GET /admin/api/metrics/by-endpoint returns endpoint breakdown', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/metrics/by-endpoint' },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('endpoints');
    expect(Array.isArray(response.data.endpoints)).toBe(true);

    if (response.data.endpoints.length > 0) {
      const ep = response.data.endpoints[0];
      expect(ep).toHaveProperty('path');
      expect(ep).toHaveProperty('count');
    }
  });

  test('GET /admin/api/metrics/top-consumers returns user IDs only (no emails)', async () => {
    const response = await makeAdminRequest(
      { method: 'GET', path: '/admin/api/metrics/top-consumers', query: { limit: 5 } },
      session
    );
    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('users');
    expect(Array.isArray(response.data.users)).toBe(true);
    if (response.data.users.length > 0) {
      const u = response.data.users[0];
      expect(u).toHaveProperty('userId');
      expect(u).toHaveProperty('count');
      // PII: should not contain email
      expect(u).not.toHaveProperty('email');
    }
  });

  test('non-admin cannot access /admin/api/metrics/overview', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/admin/api/metrics/overview`, {
      validateStatus: () => true,
    });
    expect(response.status).toBe(401);
  });
});

describe('Prometheus /metrics endpoint', () => {
  test('GET /metrics returns text/plain Prometheus format', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/metrics`, {
      validateStatus: () => true,
    });
    // May return 403 if IP allowlist blocks us — that's fine in some envs
    if (response.status === 403) {
      console.log('Skipping /metrics test: blocked by IP allowlist (expected in some envs)');
      return;
    }
    expect(response.status).toBe(200);
    expect(typeof response.data).toBe('string');
    expect(response.data).toContain('http_requests_total');
    expect(response.data).toContain('ws_connections_active');
  });

  test('GET /metrics does not expose PII in labels', async () => {
    const response = await axios.get(`${testVariables.baseUrl}/metrics`, {
      validateStatus: () => true,
    });
    if (response.status === 403) return;

    expect(response.status).toBe(200);
    // PII safety: never expose emails or API keys in metric labels
    expect(response.data).not.toContain('@');
    expect(response.data).not.toContain(testVariables.adminEmail);
    if (testVariables.apiKey.length >= 16) {
      expect(response.data).not.toContain(testVariables.apiKey);
    }
  });
});
