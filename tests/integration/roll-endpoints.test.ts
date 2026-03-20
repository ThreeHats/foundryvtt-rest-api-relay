/**
 * @file roll-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Dice Rolling Endpoint Tests
 * @endpoints POST /roll, GET /rolls, GET /lastroll, GET /rolls/subscribe
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

/**
 * Generate custom code examples for the SSE /rolls/subscribe endpoint.
 * SSE is a streaming protocol that requires special handling in each language.
 */
function generateRollSSECodeExamples() {
  return {
    javascript: `const { EventSource } = require('eventsource'); // npm install eventsource

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = \`\${baseUrl}/rolls/subscribe?clientId=your-client-id\`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

eventSource.addEventListener('connected', (event) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('roll', (event) => {
  const roll = JSON.parse(event.data);
  const dice = roll.dice?.map(d =>
    \`\${d.results.map(r => \`\${r.result}\${r.active ? '' : '(dropped)'}\`).join(', ')} (d\${d.faces})\`
  ).join(' + ') || '';
  console.log(\`[\${roll.user?.name}] \${roll.formula} = \${roll.rollTotal}\${roll.isCritical ? ' CRITICAL!' : ''}\${roll.isFumble ? ' FUMBLE!' : ''}\`);
  if (roll.flavor) console.log(\`  Flavor: \${roll.flavor}\`);
  if (dice) console.log(\`  Dice: \${dice}\`);
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
};

// To disconnect later:
// eventSource.close();`,

    curl: `# Connect to the roll SSE stream (streams events until interrupted with Ctrl+C)
curl -N 'http://localhost:3010/rolls/subscribe?clientId=your-client-id' \\
  -H "x-api-key: your-api-key-here" \\
  -H "Accept: text/event-stream"

# Example output:
# event: connected
# data: {"clientId":"your-client-id"}
#
# event: roll
# data: {"id":"abc123","user":{"id":"xyz","name":"GM"},"formula":"1d20+5","rollTotal":18,"isCritical":false,"isFumble":false,"dice":[{"faces":20,"results":[{"result":13,"active":true}]}],"flavor":"Attack Roll","timestamp":1234567890}`,

    python: `import sseclient  # pip install sseclient-py
import requests
import json

base_url = 'http://localhost:3010'
url = f'{base_url}/rolls/subscribe'
params = {'clientId': 'your-client-id'}
headers = {
    'x-api-key': 'your-api-key-here',
    'Accept': 'text/event-stream'
}

# Connect to the SSE stream
response = requests.get(url, params=params, headers=headers, stream=True)
client = sseclient.SSEClient(response)

for event in client.events():
    data = json.loads(event.data)

    if event.event == 'connected':
        print(f'Connected: {data["clientId"]}')
    elif event.event == 'roll':
        user = (data.get('user') or {}).get('name', '?')
        crit = ' CRITICAL!' if data.get('isCritical') else ''
        fumble = ' FUMBLE!' if data.get('isFumble') else ''
        print(f'[{user}] {data["formula"]} = {data["rollTotal"]}{crit}{fumble}')
        if data.get('flavor'):
            print(f'  Flavor: {data["flavor"]}')
        for d in data.get('dice', []):
            results = ', '.join(
                f'{r["result"]}{"" if r.get("active", True) else "(dropped)"}'
                for r in d.get('results', [])
            )
            print(f'  Dice: {results} (d{d["faces"]})')`,

    typescript: `// npm install eventsource
import { EventSource } from 'eventsource';

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = \`\${baseUrl}/rolls/subscribe?clientId=your-client-id\`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

interface RollEvent {
  id: string;
  messageId: string;
  user: { id: string; name: string };
  speaker: any;
  flavor: string;
  rollTotal: number;
  formula: string;
  isCritical: boolean;
  isFumble: boolean;
  dice: { faces: number; results: { result: number; active: boolean }[] }[];
  timestamp: number;
}

eventSource.addEventListener('connected', (event: MessageEvent) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('roll', (event: MessageEvent) => {
  const roll: RollEvent = JSON.parse(event.data);
  const dice = roll.dice?.map(d =>
    \`\${d.results.map(r => \`\${r.result}\${r.active ? '' : '(dropped)'}\`).join(', ')} (d\${d.faces})\`
  ).join(' + ') || '';
  console.log(\`[\${roll.user?.name}] \${roll.formula} = \${roll.rollTotal}\${roll.isCritical ? ' CRITICAL!' : ''}\${roll.isFumble ? ' FUMBLE!' : ''}\`);
  if (roll.flavor) console.log(\`  Flavor: \${roll.flavor}\`);
  if (dice) console.log(\`  Dice: \${dice}\`);
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
};

// To disconnect: eventSource.close();`,

    emojicode: `Just don't 😂`
  };
}

describe('Roll', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/roll-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/roll (v${version})`, () => {
      test('POST /roll', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/roll',
            host: ['{{baseUrl}}'],
            path: ['roll'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
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
              formula: '2d20kh',
              flavor: 'Test Roll',
              createChatMessage: true
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/roll'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('roll');
        expect(captured.response.data.data.roll).toHaveProperty('formula');
        expect(captured.response.data.data.roll).toHaveProperty('total');
        expect(captured.response.data.data.roll).toHaveProperty('isCritical');
        expect(captured.response.data.data.roll).toHaveProperty('isFumble');
        expect(captured.response.data.data.roll).toHaveProperty('dice');

        // Save variables for use in subsequent tests
        setGlobalVariable(version, 'lastRollTotal', captured.response.data.data.roll.total);
      });
    });

    describe(`/rolls (v${version})`, () => {
      test('GET /rolls', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/rolls',
            host: ['{{baseUrl}}'],
            path: ['rolls'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'limit',
                value: '20',
              }
            ]
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
          '/rolls'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(Array.isArray(captured.response.data.data)).toBe(true);
        expect(captured.response.data.data.length).toBeLessThanOrEqual(20);
        expect(captured.response.data.data.length).toBeGreaterThan(0);
        expect(captured.response.data.data[0]).toHaveProperty('formula');
        expect(captured.response.data.data[0]).toHaveProperty('rollTotal', getGlobalVariable(version, 'lastRollTotal'));
        expect(captured.response.data.data[0]).toHaveProperty('isCritical');
        expect(captured.response.data.data[0]).toHaveProperty('isFumble');
        expect(captured.response.data.data[0]).toHaveProperty('dice');
      });
    });

    describe(`/lastroll (v${version})`, () => {
      test('GET /lastroll', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/lastroll',
            host: ['{{baseUrl}}'],
            path: ['lastroll'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
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
          '/lastroll'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('formula');
        expect(captured.response.data.data).toHaveProperty('rollTotal', getGlobalVariable(version, 'lastRollTotal'));
        expect(captured.response.data.data).toHaveProperty('isCritical');
        expect(captured.response.data.data).toHaveProperty('isFumble');
        expect(captured.response.data.data).toHaveProperty('dice');
      });
    });

    describe(`GET /rolls/subscribe - SSE connection (v${version})`, () => {
      test('GET /rolls/subscribe receives roll events', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        const axios = (await import('axios')).default;
        const baseUrl = testVariables.baseUrl;
        const clientId = getClientId();
        const apiKey = testVariables.apiKey;

        // Step 1: Open SSE connection
        const sseResponse = await axios({
          url: `${baseUrl}/rolls/subscribe?clientId=${clientId}`,
          method: 'GET',
          headers: { 'x-api-key': apiKey },
          responseType: 'stream',
          timeout: 15000,
          validateStatus: () => true
        });

        // Verify SSE headers
        expect(sseResponse.status).toBe(200);
        expect(sseResponse.headers['content-type']).toBe('text/event-stream');
        expect(sseResponse.headers['cache-control']).toBe('no-cache');
        expect(sseResponse.headers['connection']).toBe('keep-alive');

        // Step 2: Wait for connected event, then trigger a roll, then wait for the roll event
        const rollEvent = await new Promise<{ connected: string; roll: string }>((resolve, reject) => {
          let buffer = '';
          let connected = '';
          let rollTriggered = false;

          const timer = setTimeout(() => {
            sseResponse.data.destroy();
            reject(new Error(`Timed out waiting for roll SSE event. Buffer: ${buffer}`));
          }, 10000);

          sseResponse.data.on('data', async (chunk: Buffer) => {
            buffer += chunk.toString();

            // Once we see the connected event, trigger a roll
            if (!rollTriggered && buffer.includes('event: connected')) {
              rollTriggered = true;
              connected = buffer;

              // Step 3: Make a roll via the API
              try {
                await axios({
                  url: `${baseUrl}/roll?clientId=${clientId}`,
                  method: 'POST',
                  headers: {
                    'x-api-key': apiKey,
                    'Content-Type': 'application/json'
                  },
                  data: { formula: '1d20', flavor: 'SSE Test Roll', createChatMessage: true }
                });
              } catch (err) {
                clearTimeout(timer);
                sseResponse.data.destroy();
                reject(new Error(`Failed to trigger roll: ${err}`));
              }
            }

            // Check if we've received the roll event
            if (rollTriggered && buffer.includes('event: roll')) {
              clearTimeout(timer);
              sseResponse.data.destroy();
              resolve({ connected, roll: buffer });
            }
          });

          sseResponse.data.on('error', (err: Error) => {
            if (!err.message.includes('aborted') && !err.message.includes('destroyed')) {
              clearTimeout(timer);
              reject(err);
            }
          });
        });

        // Step 4: Verify the events
        expect(rollEvent.connected).toContain('event: connected');
        expect(rollEvent.connected).toContain(clientId);
        expect(rollEvent.roll).toContain('event: roll');
        expect(rollEvent.roll).toContain('formula');
        expect(rollEvent.roll).toContain('rollTotal');
        expect(rollEvent.roll).toContain('1d20');

        // Push a manually-crafted captured example with custom SSE code examples
        capturedExamples.push({
          endpoint: '/rolls/subscribe',
          method: 'GET',
          description: '/rolls/subscribe',
          request: {
            url: `${baseUrl}/rolls/subscribe?clientId=${clientId}`,
            method: 'GET',
            headers: { 'x-api-key': apiKey }
          },
          response: {
            status: 200,
            data: { event: 'connected', data: { clientId } }
          },
          codeExamples: generateRollSSECodeExamples()
        });
      }, 20000);
    });

  });

});
