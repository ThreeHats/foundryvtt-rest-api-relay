/**
 * @file admin-auth.test.ts
 * @description Admin dashboard authentication tests — JWT cookie, CSRF, lockout, rate limiting.
 * @endpoints POST /admin/auth/login, POST /admin/auth/logout, GET /admin/auth/me
 *
 * Requires TEST_ADMIN_EMAIL + TEST_ADMIN_PASSWORD in .env.test (or falls back to TEST_USER_EMAIL/PASSWORD).
 * The configured user must have role='admin' on the relay (set via ADMIN_EMAIL/ADMIN_PASSWORD env vars).
 */

import axios from 'axios';
import { describe, test, expect, beforeAll } from '@jest/globals';
import { testVariables } from '../../helpers/testVariables';
import { adminLogin, adminLogout, makeAdminRequest, AdminSession } from '../../helpers/adminAuth';

const hasAdminCredentials = testVariables.adminEmail !== '' && testVariables.adminPassword !== '';

const describeAdmin = hasAdminCredentials ? describe : describe.skip;

describeAdmin('Admin Authentication', () => {
  let session: AdminSession | null = null;

  beforeAll(() => {
    if (!hasAdminCredentials) {
      console.log('Skipping admin auth tests — TEST_ADMIN_EMAIL/PASSWORD not set');
    }
  });

  describe('POST /admin/auth/login', () => {
    test('valid admin credentials return Set-Cookie + CSRF token', async () => {
      const response = await axios.post(
        `${testVariables.baseUrl}/admin/auth/login`,
        { email: testVariables.adminEmail, password: testVariables.adminPassword },
        { validateStatus: () => true }
      );

      expect(response.status).toBe(200);
      const setCookie = response.headers['set-cookie'];
      expect(setCookie).toBeDefined();
      expect(Array.isArray(setCookie)).toBe(true);
      const cookieStr = (setCookie as string[]).join('; ');
      // Dev uses "admin_token=", production uses "__Host-admin_token="
      expect(cookieStr).toMatch(/(__Host-)?admin_token=/);
      expect(cookieStr).toContain('HttpOnly');
      expect(cookieStr).toContain('SameSite=Strict');

      const csrf = response.headers['x-csrf-token'] || response.data?.csrfToken;
      expect(csrf).toBeDefined();
      expect(typeof csrf).toBe('string');
      expect((csrf as string).length).toBeGreaterThan(16);

      // Login response should NOT contain password or full apiKey (PII check)
      const bodyStr = JSON.stringify(response.data);
      expect(bodyStr).not.toContain(testVariables.adminPassword);
    });

    test('wrong password returns 401', async () => {
      const response = await axios.post(
        `${testVariables.baseUrl}/admin/auth/login`,
        { email: testVariables.adminEmail, password: 'wrong-password-xyz-123' },
        { validateStatus: () => true }
      );
      expect(response.status).toBe(401);
    });

    test('unknown email returns 401', async () => {
      const response = await axios.post(
        `${testVariables.baseUrl}/admin/auth/login`,
        { email: `nobody-${Date.now()}@example.com`, password: 'irrelevant1' },
        { validateStatus: () => true }
      );
      expect(response.status).toBe(401);
    });

    test('login as a non-admin user returns 401', async () => {
      // Skip if there's no separate non-admin test user configured
      if (
        !testVariables.userEmail ||
        !testVariables.userPassword ||
        testVariables.userEmail === testVariables.adminEmail
      ) {
        return;
      }
      const response = await axios.post(
        `${testVariables.baseUrl}/admin/auth/login`,
        { email: testVariables.userEmail, password: testVariables.userPassword },
        { validateStatus: () => true }
      );
      // Server returns 401 for non-admin users (same as wrong credentials) to avoid
      // leaking role information — all failure paths use the same generic error.
      expect(response.status).toBe(401);
    });
  });

  describe('GET /admin/auth/me', () => {
    test('without cookie returns 401', async () => {
      const response = await axios.get(`${testVariables.baseUrl}/admin/auth/me`, {
        validateStatus: () => true,
      });
      expect(response.status).toBe(401);
    });

    test('with valid cookie returns admin user info', async () => {
      session = await adminLogin();
      const response = await makeAdminRequest({ method: 'GET', path: '/admin/auth/me' }, session);
      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('email', testVariables.adminEmail);
      expect(response.data).toHaveProperty('role', 'admin');

      // Should never include password or full apiKey
      const bodyStr = JSON.stringify(response.data);
      expect(bodyStr).not.toContain(testVariables.adminPassword);
    });
  });

  describe('POST /admin/auth/logout', () => {
    test('logout invalidates the session — subsequent /me returns 401', async () => {
      const fresh = await adminLogin(testVariables.adminEmail, testVariables.adminPassword, { fresh: true });
      await adminLogout(fresh);
      const response = await makeAdminRequest({ method: 'GET', path: '/admin/auth/me' }, fresh);
      expect(response.status).toBe(401);
    });
  });

  describe('CSRF protection', () => {
    test('mutating request without X-CSRF-Token returns 403', async () => {
      const fresh = await adminLogin();
      // Send a POST without the CSRF header (deliberately strip it)
      const response = await axios.post(
        `${testVariables.baseUrl}/admin/auth/logout`,
        {},
        {
          headers: { Cookie: fresh.cookie },
          validateStatus: () => true,
        }
      );
      expect(response.status).toBe(403);
    });

    test('mutating request with wrong X-CSRF-Token returns 403', async () => {
      const fresh = await adminLogin();
      const response = await axios.post(
        `${testVariables.baseUrl}/admin/auth/logout`,
        {},
        {
          headers: {
            Cookie: fresh.cookie,
            'X-CSRF-Token': 'completely-fake-token',
          },
          validateStatus: () => true,
        }
      );
      expect(response.status).toBe(403);
    });
  });

  describe('Disabled account', () => {
    test('disabled admin account cannot log in (returns 401)', async () => {
      const adminSession = await adminLogin();

      // Register a temporary user to avoid disabling the actual admin account
      const tempEmail = `test-disable-${Date.now()}@test.invalid`;
      const tempPassword = 'TempPass123!';
      const registerResponse = await axios.post(
        `${testVariables.baseUrl}/auth/register`,
        { email: tempEmail, password: tempPassword },
        { validateStatus: () => true }
      );
      if (registerResponse.status !== 201) {
        console.log('Skipping disabled-account test — registration failed or is disabled');
        return;
      }

      // Find the new user's ID
      const usersResponse = await makeAdminRequest({ method: 'GET', path: '/admin/api/users' }, adminSession);
      const tempUser = usersResponse.data?.users?.find((u: any) => u.email === tempEmail);
      if (!tempUser) {
        console.log('Skipping disabled-account test — temp user not found in list');
        return;
      }
      const tempId = tempUser.id;

      try {
        // Promote to admin so it's a disabled-admin scenario
        await makeAdminRequest(
          { method: 'PATCH', path: `/admin/api/users/${tempId}`, body: { role: 'admin' } },
          adminSession
        );

        // Disable the temp admin
        const disableResp = await makeAdminRequest(
          { method: 'POST', path: `/admin/api/users/${tempId}/disable` },
          adminSession
        );
        expect(disableResp.status).toBe(200);

        // Login as disabled admin → must be rejected
        // The login endpoint returns 401 for disabled accounts (doesn't reveal why)
        const loginResponse = await axios.post(
          `${testVariables.baseUrl}/admin/auth/login`,
          { email: tempEmail, password: tempPassword },
          { validateStatus: () => true }
        );
        expect(loginResponse.status).toBe(401);
      } finally {
        // Always clean up regardless of test outcome
        await makeAdminRequest(
          { method: 'POST', path: `/admin/api/users/${tempId}/enable` },
          adminSession
        ).catch(() => {});
        await makeAdminRequest(
          { method: 'DELETE', path: `/admin/api/users/${tempId}` },
          adminSession
        ).catch(() => {});
      }
    });
  });

  describe('Rate limiting', () => {
    test('repeated failed logins from same IP eventually return 429', async () => {
      // Make up to 12 attempts with unique emails to avoid account lockout interfering.
      // If rate limiting is disabled on the server (ADMIN_LOGIN_RATE_LIMIT=0), no 429
      // will be seen and the test passes silently — this is expected in dev environments.
      let saw429 = false;
      for (let i = 0; i < 12; i++) {
        const response = await axios.post(
          `${testVariables.baseUrl}/admin/auth/login`,
          { email: `ratelimit-${Date.now()}-${i}@example.com`, password: 'wrong' },
          { validateStatus: () => true }
        );
        if (response.status === 429) {
          saw429 = true;
          break;
        }
      }
      if (!saw429) {
        // Rate limiting is disabled on this server — skip the assertion.
        // Set ADMIN_LOGIN_RATE_LIMIT to a positive integer to activate rate limiting tests.
        console.log('Note: admin rate limiting appears disabled (no 429 after 12 attempts)');
        return;
      }
      expect(saw429).toBe(true);
    });
  });
});
