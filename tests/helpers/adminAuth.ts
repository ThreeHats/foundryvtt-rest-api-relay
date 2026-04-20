/**
 * Admin Auth Helper for Tests
 *
 * Admin endpoints use cookie-based JWT auth + CSRF token, not the x-api-key header.
 * This helper logs in and provides a wrapper for admin-authenticated requests.
 */

import axios, { AxiosRequestConfig, Method } from 'axios';
import { testVariables } from './testVariables';

export interface AdminSession {
  cookie: string; // raw Cookie header value (e.g. "__Host-admin_token=...")
  csrfToken: string;
}

export interface AdminRequestConfig {
  method: Method;
  path: string; // e.g. "/admin/api/users"
  query?: Record<string, string | number>;
  body?: any;
}

export interface AdminResponse<T = any> {
  data: T;
  status: number;
  statusText: string;
  headers: any;
}

// Module-level session cache. Shared across all test files in the same Jest process.
// This prevents repeated login calls from exhausting the admin rate limit.
let _cachedSession: AdminSession | null = null;

/**
 * Log in as the configured admin and return a session (cookie + CSRF).
 * By default returns a cached session if one exists. Pass { fresh: true } to
 * force a new login (required for logout / session-invalidation tests).
 * Throws if login fails — tests should call this in a beforeAll/beforeEach.
 */
export async function adminLogin(
  email: string = testVariables.adminEmail,
  password: string = testVariables.adminPassword,
  options: { fresh?: boolean } = {}
): Promise<AdminSession> {
  if (!options.fresh && _cachedSession) {
    return _cachedSession;
  }
  if (!email || !password) {
    throw new Error('TEST_ADMIN_EMAIL and TEST_ADMIN_PASSWORD must be set in .env.test');
  }
  const response = await axios.post(
    `${testVariables.baseUrl}/admin/auth/login`,
    { email, password },
    { validateStatus: () => true, timeout: 30000 }
  );
  if (response.status !== 200) {
    throw new Error(`Admin login failed: ${response.status} ${JSON.stringify(response.data)}`);
  }
  const setCookie = response.headers['set-cookie'];
  if (!setCookie || !Array.isArray(setCookie) || setCookie.length === 0) {
    throw new Error('Admin login response missing Set-Cookie header');
  }
  const cookieValue = setCookie
    .map((c: string) => c.split(';')[0]) // strip attributes
    .find((c: string) => c.startsWith('admin_token=') || c.startsWith('__Host-admin_token='));
  if (!cookieValue) {
    throw new Error('Admin login response missing admin_token cookie');
  }
  const csrfToken = response.headers['x-csrf-token'] || response.data?.csrfToken;
  if (!csrfToken) {
    throw new Error('Admin login response missing CSRF token');
  }
  const session = { cookie: cookieValue, csrfToken };
  if (!options.fresh) {
    _cachedSession = session;
  }
  return session;
}

/**
 * Make an admin-authenticated HTTP request.
 * Attaches the session cookie and CSRF header automatically.
 */
export async function makeAdminRequest<T = any>(
  config: AdminRequestConfig,
  session: AdminSession
): Promise<AdminResponse<T>> {
  let url = `${testVariables.baseUrl}${config.path}`;
  if (config.query) {
    const params = new URLSearchParams();
    for (const [k, v] of Object.entries(config.query)) {
      params.append(k, String(v));
    }
    url += `?${params.toString()}`;
  }

  const headers: Record<string, string> = {
    Cookie: session.cookie,
    'X-CSRF-Token': session.csrfToken,
  };
  if (config.body !== undefined) {
    headers['Content-Type'] = 'application/json';
  }

  const axiosConfig: AxiosRequestConfig = {
    url,
    method: config.method,
    headers,
    data: config.body,
    timeout: 60000,
    validateStatus: () => true,
    maxRedirects: 0,
  };

  const response = await axios.request(axiosConfig);
  return {
    data: response.data,
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
  };
}

/**
 * Log out — invalidates the cookie via /admin/auth/logout.
 * Clears the module-level session cache if the logged-out session was the cached one.
 */
export async function adminLogout(session: AdminSession): Promise<void> {
  if (_cachedSession === session) {
    _cachedSession = null;
  }
  await axios.post(
    `${testVariables.baseUrl}/admin/auth/logout`,
    {},
    {
      headers: {
        Cookie: session.cookie,
        'X-CSRF-Token': session.csrfToken,
      },
      validateStatus: () => true,
      timeout: 30000,
    }
  );
}
