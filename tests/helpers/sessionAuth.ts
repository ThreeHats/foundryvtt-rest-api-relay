/**
 * Session-token helper for suites that hit session-only /auth routes
 * (connection tokens, credentials, notifications, key requests, scoped keys).
 *
 * Ephemeral mode: register-account.test.ts stores a sessionToken in
 * .global-vars.json and testVariables picks it up automatically.
 *
 * Pre-provisioned mode (TEST_API_KEY set): there is no registered account, so
 * we log in once with TEST_USER_EMAIL / TEST_USER_PASSWORD and persist the
 * token to the global-vars store so every later test file reuses it instead
 * of logging in again (which would pressure the auth rate limit).
 */
import { testVariables, setVariable } from './testVariables';
import { makeRequest } from './apiRequest';
import { setGlobalVariable } from './globalVariables';

/**
 * Ensure testVariables.sessionToken is populated, logging in if necessary.
 * Call from a beforeAll in any suite that sends `Bearer {{sessionToken}}`.
 * Throws with an actionable message when no credentials are available.
 */
export async function ensureSessionToken(): Promise<string> {
  if (testVariables.sessionToken) {
    return testVariables.sessionToken;
  }

  if (!testVariables.userEmail || !testVariables.userPassword) {
    throw new Error(
      'No session token available: set TEST_USER_EMAIL and TEST_USER_PASSWORD in .env.test ' +
      '(pre-provisioned mode) so session-only /auth tests can log in.'
    );
  }

  const response = await makeRequest({
    url: {
      raw: `${testVariables.baseUrl}/auth/login`,
      host: [testVariables.baseUrl],
      path: ['auth', 'login'],
    },
    method: 'POST',
    header: [],
    body: {
      mode: 'raw',
      raw: JSON.stringify({ email: testVariables.userEmail, password: testVariables.userPassword }),
    },
  });

  if (response.status !== 200 || !response.data?.sessionToken) {
    throw new Error(
      `Login for session-only tests failed: ${response.status} ` +
      `${JSON.stringify(response.data ?? {}).slice(0, 200)}`
    );
  }

  const token = response.data.sessionToken as string;
  // Current file sees it via {{sessionToken}} substitution immediately…
  setVariable('sessionToken', token);
  // …and later test files pick it up at load through readRegisteredAccount().
  setGlobalVariable('account', 'sessionToken', token);
  return token;
}
