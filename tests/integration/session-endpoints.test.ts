/**
 * @file session-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Session Authentication and Management Tests
 * @endpoints POST /session-handshake, POST /start-session, GET /session
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';
import crypto from 'crypto';

// Skip session tests when using existing sessions
const useExistingSession = process.env.USE_EXISTING_SESSION === 'true';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Store handshake data per version for use across tests
const handshakeData: Record<string, { token: string; publicKey: string; nonce: string }> = {};

const describeOrSkip = useExistingSession ? describe.skip : describe;

/**
 * Generate code examples that show the complete session flow:
 * 1. Call /session-handshake to get token, publicKey, nonce
 * 2. Encrypt password with RSA-OAEP using publicKey and nonce
 * 3. Call /start-session with handshakeToken and encryptedPassword
 */
function generateSessionCodeExamples() {
  return {
    javascript: `// Step 1: Get handshake credentials
const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key';

const handshakeResponse = await fetch(\`\${baseUrl}/session-handshake\`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'x-foundry-url': 'http://localhost:30000',
    'x-world-name': 'my-world',
    'x-username': 'Gamemaster'
  }
});
const { token, publicKey, nonce } = await handshakeResponse.json();

// Step 2: Encrypt password using Web Crypto API (RSA-OAEP with SHA-1)
const password = 'your-password';
const payload = JSON.stringify({ password, nonce });

// Import the public key
const pemContents = publicKey
  .replace('-----BEGIN PUBLIC KEY-----', '')
  .replace('-----END PUBLIC KEY-----', '')
  .replace(/\\n/g, '');
const binaryKey = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0));

const cryptoKey = await crypto.subtle.importKey(
  'spki',
  binaryKey,
  { name: 'RSA-OAEP', hash: 'SHA-1' },  // Must use SHA-1 to match server
  false,
  ['encrypt']
);

// Encrypt the payload
const encrypted = await crypto.subtle.encrypt(
  { name: 'RSA-OAEP' },
  cryptoKey,
  new TextEncoder().encode(payload)
);
const encryptedPassword = btoa(String.fromCharCode(...new Uint8Array(encrypted)));

// Step 3: Start the session
const sessionResponse = await fetch(\`\${baseUrl}/start-session\`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ handshakeToken: token, encryptedPassword })
});
const { sessionId, clientId } = await sessionResponse.json();
console.log('Session started:', { sessionId, clientId });`,

    typescript: `import crypto from 'crypto';

// Step 1: Get handshake credentials
const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key';

const handshakeResponse = await fetch(\`\${baseUrl}/session-handshake\`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'x-foundry-url': 'http://localhost:30000',
    'x-world-name': 'my-world',
    'x-username': 'Gamemaster'
  }
});
const { token, publicKey, nonce } = await handshakeResponse.json();

// Step 2: Encrypt password using Node.js crypto (RSA-OAEP)
const password = 'your-password';
const payload = JSON.stringify({ password, nonce });
const encryptedPassword = crypto.publicEncrypt(
  { key: publicKey, padding: crypto.constants.RSA_PKCS1_OAEP_PADDING },
  Buffer.from(payload, 'utf8')
).toString('base64');

// Step 3: Start the session
const sessionResponse = await fetch(\`\${baseUrl}/start-session\`, {
  method: 'POST',
  headers: {
    'x-api-key': apiKey,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ handshakeToken: token, encryptedPassword })
});
const { sessionId, clientId } = await sessionResponse.json();
console.log('Session started:', { sessionId, clientId });`,

    python: `import requests
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import padding
import base64
import json

base_url = 'http://localhost:3010'
api_key = 'your-api-key'

# Step 1: Get handshake credentials
handshake_response = requests.post(
    f'{base_url}/session-handshake',
    headers={
        'x-api-key': api_key,
        'x-foundry-url': 'http://localhost:30000',
        'x-world-name': 'my-world',
        'x-username': 'Gamemaster'
    }
)
handshake_data = handshake_response.json()
token = handshake_data['token']
public_key_pem = handshake_data['publicKey']
nonce = handshake_data['nonce']

# Step 2: Encrypt password using RSA-OAEP with SHA-1 (must match server)
password = 'your-password'
payload = json.dumps({'password': password, 'nonce': nonce})

public_key = serialization.load_pem_public_key(public_key_pem.encode())
encrypted = public_key.encrypt(
    payload.encode(),
    padding.OAEP(
        mgf=padding.MGF1(algorithm=hashes.SHA1()),
        algorithm=hashes.SHA1(),
        label=None
    )
)
encrypted_password = base64.b64encode(encrypted).decode()

# Step 3: Start the session
session_response = requests.post(
    f'{base_url}/start-session',
    headers={'x-api-key': api_key},
    json={'handshakeToken': token, 'encryptedPassword': encrypted_password}
)
session_data = session_response.json()
print('Session started:', session_data)`,

    curl: `# Session creation requires encryption - use the JavaScript, Python, or TypeScript examples.
# 
# The flow is:
# 1. POST /session-handshake to get token, publicKey, nonce
# 2. Encrypt JSON payload { password, nonce } with RSA-OAEP using publicKey
# 3. POST /start-session with { handshakeToken, encryptedPassword }
#
# Example handshake request:
curl -X POST 'http://localhost:3010/session-handshake' \\
  -H "x-api-key: your-api-key" \\
  -H "x-foundry-url: http://localhost:30000" \\
  -H "x-world-name: my-world" \\
  -H "x-username: Gamemaster"`
  };
}

describeOrSkip('Session', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/session-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/session-handshake (v${version})`, () => {
      test('POST /session-handshake', async () => {
        // Get version-specific Foundry connection details from env
        const foundryUrl = process.env[`FOUNDRY_V${version}_URL`] || '';
        const worldName = process.env[`FOUNDRY_V${version}_WORLD`] || process.env.TEST_DEFAULT_WORLD || 'test-world';
        const username = process.env.FOUNDRY_USERNAME || 'testuser';
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/session-handshake',
            host: ['{{baseUrl}}'],
            path: ['session-handshake'],
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'x-foundry-url', value: foundryUrl, type: 'text' },
            { key: 'x-world-name', value: worldName, type: 'text' },
            { key: 'x-username', value: username, type: 'text' }
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/session-handshake'
        );
        // Skip capturing this example for documentation since we cover it in the start-session example
        // capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data.token).toBeTruthy();
        expect(captured.response.data.publicKey).toBeTruthy();
        expect(captured.response.data.nonce).toBeTruthy();

        // Store handshake data for use in start-session test
        handshakeData[version] = {
          token: captured.response.data.token,
          publicKey: captured.response.data.publicKey,
          nonce: captured.response.data.nonce
        };
      });
    });

    describe(`/start-session (v${version})`, () => {
      test('POST /start-session', async () => {
        // Ensure we have handshake data from previous test
        const hs = handshakeData[version];
        expect(hs).toBeTruthy();
        
        // Encrypt password with nonce using RSA-OAEP
        const password = process.env.FOUNDRY_PASSWORD || 'password';
        const payload = JSON.stringify({ password, nonce: hs.nonce });
        const encryptedPassword = crypto.publicEncrypt(
          { key: hs.publicKey, padding: crypto.constants.RSA_PKCS1_OAEP_PADDING },
          Buffer.from(payload, 'utf8')
        ).toString('base64');
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/start-session',
            host: ['{{baseUrl}}'],
            path: ['start-session'],
          },
          method: 'POST',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              handshakeToken: hs.token,
              encryptedPassword: encryptedPassword
            }, null, 2)
        }
      };

        // Custom code examples that show the full handshake + encryption flow
        const customCodeExamples = generateSessionCodeExamples();

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/start-session',
          { customCodeExamples }
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data.sessionId).toBeTruthy();

        // Store session data in globalVariables for other test files
        if (captured.response.data.clientId) {
          const clientId = captured.response.data.clientId;
          const sessionId = captured.response.data.sessionId;
          
          setGlobalVariable(version, 'clientId', clientId);
          setGlobalVariable(version, 'sessionId', sessionId);
          
          console.log(`Stored clientId=${clientId} and sessionId=${sessionId} for v${version}`);
          return;
        }
        
        // If we get here without a clientId, fail the test
        throw new Error(`Timed out waiting for clientId for v${version} session`);
      }, 90000); // Extended timeout for headless startup
    });

    describe(`/session (v${version})`, () => {
      test('GET /session', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/session',
            host: ['{{baseUrl}}'],
            path: ['session'],
          },
          method: 'GET',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/session'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data.activeSessions).toBeTruthy();
        expect(captured.response.data.activeSessions.length).toBeGreaterThan(0);
        
        // Store client details (systemId, etc.) for use by other tests
        const sessions = captured.response.data.activeSessions;
        if (sessions && sessions.length > 0) {
          const session = sessions[0];
          if (session.systemId) {
            setGlobalVariable(version, 'systemId', session.systemId);
            console.log(`✅ Stored systemId=${session.systemId} for v${version}`);
          }
          if (session.worldId) {
            setGlobalVariable(version, 'worldId', session.worldId);
          }
          if (session.foundryVersion) {
            setGlobalVariable(version, 'foundryVersion', session.foundryVersion);
          }
        }
      });
    });

  });

});
