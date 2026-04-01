/**
 * @file auth-endpoints.test.ts
 * @description Auth endpoint tests — registration, login, user data, password management, Stripe
 * @endpoints POST /auth/register, POST /auth/login, GET /auth/user-data, GET /auth/export-data,
 *   POST /auth/change-password, POST /auth/regenerate-key, POST /auth/forgot-password,
 *   POST /auth/reset-password, GET /auth/validate-reset-token/:token,
 *   GET /api/subscriptions/status, POST /api/subscriptions/create-checkout-session
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { setGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Throwaway user for register/delete tests
const throwawayEmail = `auth-test-${Date.now()}@example.com`;
const throwawayPassword = 'TestPassword1';
let throwawayApiKey = '';

// Existing test user credentials
const userEmail = testVariables.userEmail;
const userPassword = testVariables.userPassword;
const masterApiKey = testVariables.apiKey;

// Whether we have test user credentials configured
const hasUserCredentials = userEmail !== '' && userPassword !== '';

// Whether Stripe test keys are configured (check via a status request)
let stripeEnabled = false;

describe('Auth Endpoints', () => {
  afterAll(() => {
    const outputPath = path.join(__dirname, '../../docs/examples/auth-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} auth examples to ${outputPath}`);
  });

  // ==================== Group 1: Registration ====================

  describe('/auth/register', () => {
    test('POST /auth/register - register a new user', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/register`,
          host: [testVariables.baseUrl],
          path: ['auth', 'register'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: throwawayEmail,
            password: throwawayPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/register');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(201);
      expect(captured.response.data).toHaveProperty('id');
      expect(captured.response.data).toHaveProperty('email', throwawayEmail);
      expect(captured.response.data).toHaveProperty('apiKey');
      expect(captured.response.data).toHaveProperty('subscriptionStatus', 'free');

      // Store for cleanup
      throwawayApiKey = captured.response.data.apiKey;
      setGlobalVariable('auth', 'throwawayEmail', throwawayEmail);
      setGlobalVariable('auth', 'throwawayPassword', throwawayPassword);
      setGlobalVariable('auth', 'throwawayApiKey', throwawayApiKey);
    });

    test('POST /auth/register - should reject missing fields', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/register`,
          host: [testVariables.baseUrl],
          path: ['auth', 'register'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ email: 'incomplete@example.com' })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/register (missing fields)');
      expect(captured.response.status).toBe(400);
      expect(captured.response.data).toHaveProperty('error');
    });

    test('POST /auth/register - should reject weak password', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/register`,
          host: [testVariables.baseUrl],
          path: ['auth', 'register'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ email: 'weak@example.com', password: 'short' })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/register (weak password)');
      expect(captured.response.status).toBe(400);
      expect(captured.response.data).toHaveProperty('error');
    });

    test('POST /auth/register - should reject duplicate email', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/register`,
          host: [testVariables.baseUrl],
          path: ['auth', 'register'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: throwawayEmail,
            password: throwawayPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/register (duplicate)');
      expect(captured.response.status).toBe(409);
      expect(captured.response.data).toHaveProperty('error');
    });
  });

  // ==================== Group 2: Login ====================

  describe('/auth/login', () => {
    const describeOrSkip = hasUserCredentials ? describe : describe.skip;

    describeOrSkip('with test user credentials', () => {
      test('POST /auth/login - login with valid credentials', async () => {
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: `${testVariables.baseUrl}/auth/login`,
            host: [testVariables.baseUrl],
            path: ['auth', 'login'],
          },
          method: 'POST',
          header: [],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              email: userEmail,
              password: userPassword
            })
          }
        };

        const captured = await captureExample(requestConfig, {}, '/auth/login');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('id');
        expect(captured.response.data).toHaveProperty('email', userEmail);
        expect(captured.response.data).toHaveProperty('apiKey');
        expect(captured.response.data).toHaveProperty('requestsThisMonth');
      });
    });

    test('POST /auth/login - should reject wrong password', async () => {
      // Use throwaway email since we know it exists
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/login`,
          host: [testVariables.baseUrl],
          path: ['auth', 'login'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: throwawayEmail,
            password: 'WrongPassword1'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/login (wrong password)');
      expect(captured.response.status).toBe(401);
      expect(captured.response.data).toHaveProperty('error');
    });

    test('POST /auth/login - should reject missing fields', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/login`,
          host: [testVariables.baseUrl],
          path: ['auth', 'login'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ email: throwawayEmail })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/login (missing fields)');
      expect(captured.response.status).toBe(400);
      expect(captured.response.data).toHaveProperty('error');
    });
  });

  // ==================== Group 3: User Data & Export ====================

  describe('/auth/user-data', () => {
    test('GET /auth/user-data - with valid API key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/user-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'user-data'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/user-data');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('id');
      expect(captured.response.data).toHaveProperty('email');
      expect(captured.response.data).toHaveProperty('apiKey');
      expect(captured.response.data).toHaveProperty('subscriptionStatus');
      expect(captured.response.data).toHaveProperty('limits');
      expect(captured.response.data.limits).toHaveProperty('dailyLimit');
      expect(captured.response.data.limits).toHaveProperty('monthlyLimit');
    });

    test('GET /auth/user-data - should reject invalid API key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/user-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'user-data'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: 'invalid-api-key-that-does-not-exist' }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/user-data (invalid key)');
      expect(captured.response.status).toBe(404);
    });

    test('GET /auth/user-data - should reject missing API key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/user-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'user-data'],
        },
        method: 'GET',
        header: []
      };

      const captured = await captureExample(requestConfig, {}, '/auth/user-data (no key)');
      expect(captured.response.status).toBe(401);
    });
  });

  describe('/auth/export-data', () => {
    test('GET /auth/export-data - with valid API key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/export-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'export-data'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/auth/export-data');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('exportDate');
      expect(captured.response.data).toHaveProperty('user');
      expect(captured.response.data).toHaveProperty('subscription');
      expect(captured.response.data).toHaveProperty('usage');
      expect(captured.response.data).toHaveProperty('apiAccess');
      expect(captured.response.data).toHaveProperty('scopedKeys');
    });

    test('GET /auth/export-data - should reject missing API key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/export-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'export-data'],
        },
        method: 'GET',
        header: []
      };

      const captured = await captureExample(requestConfig, {}, '/auth/export-data (no key)');
      expect(captured.response.status).toBe(401);
    });
  });

  // ==================== Group 4: Change Password ====================

  const describeChangePassword = hasUserCredentials ? describe : describe.skip;

  describeChangePassword('/auth/change-password', () => {
    const newPassword = 'ChangedPassword2';

    test('POST /auth/change-password - change password', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/change-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'change-password'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            currentPassword: userPassword,
            newPassword: newPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/change-password');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('message');
    });

    test('POST /auth/login - verify new password works', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/login`,
          host: [testVariables.baseUrl],
          path: ['auth', 'login'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: userEmail,
            password: newPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/login (new password)');
      expect(captured.response.status).toBe(200);
    });

    test('POST /auth/change-password - restore original password', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/change-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'change-password'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            currentPassword: newPassword,
            newPassword: userPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/change-password (restore)');
      expect(captured.response.status).toBe(200);
    });

    test('POST /auth/change-password - should reject wrong current password', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/change-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'change-password'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            currentPassword: 'WrongCurrent1',
            newPassword: 'NewPassword1'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/change-password (wrong current)');
      expect(captured.response.status).toBe(401);
    });

    test('POST /auth/change-password - should reject weak new password', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/change-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'change-password'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            currentPassword: userPassword,
            newPassword: 'weak'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/change-password (weak password)');
      expect(captured.response.status).toBe(400);
    });
  });

  // ==================== Group 5: Regenerate API Key (throwaway user) ====================

  describe('/auth/regenerate-key', () => {
    test('POST /auth/regenerate-key - regenerate API key', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/regenerate-key`,
          host: [testVariables.baseUrl],
          path: ['auth', 'regenerate-key'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: throwawayEmail,
            password: throwawayPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/regenerate-key');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('apiKey');
      expect(captured.response.data.apiKey).not.toBe(throwawayApiKey);

      // Update stored key for subsequent tests
      const oldKey = throwawayApiKey;
      throwawayApiKey = captured.response.data.apiKey;
      setGlobalVariable('auth', 'throwawayApiKey', throwawayApiKey);

      // Verify old key no longer works
      const verifyConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/user-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'user-data'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: oldKey }
        ]
      };
      const verifyOld = await captureExample(verifyConfig, {}, '/auth/user-data (old key)');
      expect(verifyOld.response.status).toBe(404);

      // Verify new key works
      const verifyNewConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/user-data`,
          host: [testVariables.baseUrl],
          path: ['auth', 'user-data'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: throwawayApiKey }
        ]
      };
      const verifyNew = await captureExample(verifyNewConfig, {}, '/auth/user-data (new key)');
      expect(verifyNew.response.status).toBe(200);
    });

    test('POST /auth/regenerate-key - should reject wrong password', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/regenerate-key`,
          host: [testVariables.baseUrl],
          path: ['auth', 'regenerate-key'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: throwawayEmail,
            password: 'WrongPassword1'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/regenerate-key (wrong password)');
      expect(captured.response.status).toBe(401);
    });

    test('POST /auth/regenerate-key - should reject missing fields', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/regenerate-key`,
          host: [testVariables.baseUrl],
          path: ['auth', 'regenerate-key'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ email: throwawayEmail })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/regenerate-key (missing fields)');
      expect(captured.response.status).toBe(400);
    });
  });

  // ==================== Group 6: Password Reset ====================

  // Track reset token across tests in this group
  let resetToken = '';
  let resetTokenAvailable = false;

  describe('/auth/forgot-password', () => {
    test('POST /auth/forgot-password - request reset for valid email', async () => {
      // Use throwaway email to avoid sending real emails
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/forgot-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'forgot-password'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ email: throwawayEmail })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/forgot-password');

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('message');

      // If RETURN_RESET_TOKEN=true is set on the server, we get the raw token back
      if (captured.response.data.token) {
        resetToken = captured.response.data.token;
        resetTokenAvailable = true;
        console.log('Reset token available — full password reset flow will be tested');
      } else {
        console.log('RETURN_RESET_TOKEN not enabled on server — only error paths will be tested');
      }

      // Strip test-only token field before capturing for docs
      const docCapture = { ...captured, response: { ...captured.response, data: { ...captured.response.data } } };
      delete docCapture.response.data.token;
      capturedExamples.push(docCapture);
    });

    test('POST /auth/forgot-password - non-existent email returns same response (anti-enumeration)', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/forgot-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'forgot-password'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({ email: 'nonexistent@example.com' })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/forgot-password (nonexistent)');
      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('message');
      // Non-existent email should NOT return a token even when RETURN_RESET_TOKEN is set
      // (because no user was found, no token is generated)
      expect(captured.response.data.token).toBeUndefined();
    });

    test('POST /auth/forgot-password - should reject missing email', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/forgot-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'forgot-password'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({})
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/forgot-password (missing email)');
      expect(captured.response.status).toBe(400);
    });
  });

  describe('/auth/validate-reset-token', () => {
    test('GET /auth/validate-reset-token/:token - invalid token returns valid:false', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/validate-reset-token/invalidtoken123456`,
          host: [testVariables.baseUrl],
          path: ['auth', 'validate-reset-token', 'invalidtoken123456'],
        },
        method: 'GET',
        header: []
      };

      const captured = await captureExample(requestConfig, {}, '/auth/validate-reset-token/:token (invalid)');
      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('valid', false);
    });

    test('GET /auth/validate-reset-token/:token - valid token returns valid:true', async () => {
      if (!resetTokenAvailable) {
        console.log('Skipping — RETURN_RESET_TOKEN not enabled');
        return;
      }

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/validate-reset-token/${resetToken}`,
          host: [testVariables.baseUrl],
          path: ['auth', 'validate-reset-token', resetToken],
        },
        method: 'GET',
        header: []
      };

      const captured = await captureExample(requestConfig, {}, '/auth/validate-reset-token/:token');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('valid', true);
    });
  });

  describe('/auth/reset-password', () => {
    test('POST /auth/reset-password - should reject invalid token', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/reset-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'reset-password'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            token: 'invalidtoken123456',
            password: 'NewValidPassword1'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/reset-password (invalid token)');
      expect(captured.response.status).toBe(400);
      expect(captured.response.data).toHaveProperty('error');
    });

    test('POST /auth/reset-password - reset password with valid token', async () => {
      if (!resetTokenAvailable) {
        console.log('Skipping — RETURN_RESET_TOKEN not enabled');
        return;
      }

      const newResetPassword = 'ResetPassword1';
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/reset-password`,
          host: [testVariables.baseUrl],
          path: ['auth', 'reset-password'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            token: resetToken,
            password: newResetPassword
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/reset-password');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('message');

      // Update throwaway password for cleanup tests
      setGlobalVariable('auth', 'throwawayPassword', newResetPassword);
    });

    test('POST /auth/login - verify login works with reset password', async () => {
      if (!resetTokenAvailable) {
        console.log('Skipping — RETURN_RESET_TOKEN not enabled');
        return;
      }

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/login`,
          host: [testVariables.baseUrl],
          path: ['auth', 'login'],
        },
        method: 'POST',
        header: [],
        body: {
          mode: 'raw',
          raw: JSON.stringify({
            email: throwawayEmail,
            password: 'ResetPassword1'
          })
        }
      };

      const captured = await captureExample(requestConfig, {}, '/auth/login (after reset)');
      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('apiKey');
    });

    test('GET /auth/validate-reset-token/:token - used token returns valid:false', async () => {
      if (!resetTokenAvailable) {
        console.log('Skipping — RETURN_RESET_TOKEN not enabled');
        return;
      }

      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/auth/validate-reset-token/${resetToken}`,
          host: [testVariables.baseUrl],
          path: ['auth', 'validate-reset-token', resetToken],
        },
        method: 'GET',
        header: []
      };

      const captured = await captureExample(requestConfig, {}, '/auth/validate-reset-token/:token (used)');
      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('valid', false);
    });
  });

  // ==================== Group 6: Stripe Subscription (conditional) ====================

  describe('Stripe Subscriptions', () => {
    // Check if Stripe is enabled by hitting the status endpoint
    test('GET /api/subscriptions/status', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/api/subscriptions/status`,
          host: [testVariables.baseUrl],
          path: ['api', 'subscriptions', 'status'],
        },
        method: 'GET',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/api/subscriptions/status');
      capturedExamples.push(captured);

      expect(captured.response.status).toBe(200);
      expect(captured.response.data).toHaveProperty('subscriptionStatus');

      // Check if Stripe is actually enabled (not returning disabled message)
      stripeEnabled = !captured.response.data.message?.includes('disabled');
    });

    test('POST /api/subscriptions/create-checkout-session', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/api/subscriptions/create-checkout-session`,
          host: [testVariables.baseUrl],
          path: ['api', 'subscriptions', 'create-checkout-session'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/api/subscriptions/create-checkout-session');

      if (stripeEnabled) {
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('url');
      } else {
        // Stripe disabled — expect 503 or disabled response
        expect([200, 503]).toContain(captured.response.status);
      }
    });

    test('POST /api/subscriptions/create-portal-session', async () => {
      const requestConfig: ApiRequestConfig = {
        url: {
          raw: `${testVariables.baseUrl}/api/subscriptions/create-portal-session`,
          host: [testVariables.baseUrl],
          path: ['api', 'subscriptions', 'create-portal-session'],
        },
        method: 'POST',
        header: [
          { key: 'x-api-key', value: masterApiKey }
        ]
      };

      const captured = await captureExample(requestConfig, {}, '/api/subscriptions/create-portal-session');

      if (stripeEnabled) {
        // May return URL or 503 if portal URL not configured
        expect([200, 503]).toContain(captured.response.status);
      } else {
        expect([200, 503]).toContain(captured.response.status);
      }
    });
  });
});
