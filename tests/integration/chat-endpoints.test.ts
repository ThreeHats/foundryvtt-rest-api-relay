/**
 * @file chat-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description Chat Endpoint Tests
 * @endpoints GET /chat, POST /chat, DELETE /chat/:messageId, DELETE /chat, GET /chat/subscribe
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

/**
 * Generate custom code examples for the SSE /chat/subscribe endpoint.
 * SSE is a streaming protocol that requires special handling in each language.
 */
function generateSSECodeExamples() {
  return {
    javascript: `const { EventSource } = require('eventsource'); // npm install eventsource

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = \`\${baseUrl}/chat/subscribe?clientId=your-client-id\`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

function formatMessage(prefix, message) {
  const speaker = message.author?.name || message.speaker?.alias || '?';
  console.log(\`[\${prefix}] \${speaker}: \${message.content}\`);
  if (message.flavor) console.log(\`  Flavor: \${message.flavor}\`);
  if (message.isRoll && message.rolls?.length > 0) {
    for (const roll of message.rolls) {
      const dice = roll.dice?.map(d =>
        \`\${d.results.map(r => \`\${r.result}\${r.active ? '' : '(dropped)'}\`).join(', ')} (d\${d.faces})\`
      ).join(' + ') || '';
      console.log(\`  Roll: \${roll.formula} = \${roll.total}\${roll.isCritical ? ' CRITICAL!' : ''}\${roll.isFumble ? ' FUMBLE!' : ''}\`);
      if (dice) console.log(\`  Dice: \${dice}\`);
    }
  }
}

eventSource.addEventListener('connected', (event) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('chat-create', (event) => {
  const message = JSON.parse(event.data);
  formatMessage('new', message);
});

eventSource.addEventListener('chat-update', (event) => {
  const message = JSON.parse(event.data);
  formatMessage('updated', message);
});

eventSource.addEventListener('chat-delete', (event) => {
  const data = JSON.parse(event.data);
  console.log('Message deleted:', JSON.stringify(data));
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
};

// To disconnect later:
// eventSource.close();`,

    curl: `# Connect to the SSE stream (streams events until interrupted with Ctrl+C)
curl -N 'http://localhost:3010/chat/subscribe?clientId=your-client-id' \\
  -H "x-api-key: your-api-key-here" \\
  -H "Accept: text/event-stream"

# Example output:
# event: connected
# data: {"clientId":"your-client-id"}
#
# event: chat-create
# data: {"id":"abc123","content":"Hello!","author":{"id":"xyz","name":"GM"},"isRoll":false,...}
#
# event: chat-create (dice roll)
# data: {"id":"def456","content":"16","flavor":"Attack Roll","isRoll":true,"rolls":[{"formula":"1d20+5","total":16,"isCritical":false,"isFumble":false,"dice":[{"faces":20,"results":[{"result":11,"active":true}]}]}],...}`,

    python: `import sseclient  # pip install sseclient-py
import requests
import json

base_url = 'http://localhost:3010'
url = f'{base_url}/chat/subscribe'
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
    elif event.event in ('chat-create', 'chat-update'):
        prefix = 'new' if event.event == 'chat-create' else 'updated'
        speaker = (data.get('author') or {}).get('name') or (data.get('speaker') or {}).get('alias') or '?'
        print(f'[{prefix}] {speaker}: {data.get("content", "")}')
        if data.get('flavor'):
            print(f'  Flavor: {data["flavor"]}')
        if data.get('isRoll') and data.get('rolls'):
            for roll in data['rolls']:
                dice_parts = []
                for d in roll.get('dice', []):
                    results = ', '.join(
                        f'{r["result"]}{"" if r.get("active", True) else "(dropped)"}'
                        for r in d.get('results', [])
                    )
                    dice_parts.append(f'{results} (d{d["faces"]})')
                crit = ' CRITICAL!' if roll.get('isCritical') else ''
                fumble = ' FUMBLE!' if roll.get('isFumble') else ''
                print(f'  Roll: {roll["formula"]} = {roll["total"]}{crit}{fumble}')
                if dice_parts:
                    print(f'  Dice: {" + ".join(dice_parts)}')
    elif event.event == 'chat-delete':
        print(f'Message deleted: {json.dumps(data)}')`,

    typescript: `// npm install eventsource
import { EventSource } from 'eventsource';

const baseUrl = 'http://localhost:3010';
const apiKey = 'your-api-key-here';
const url = \`\${baseUrl}/chat/subscribe?clientId=your-client-id\`;

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

interface ChatMessage {
  id: string;
  content: string;
  type: string;
  author: { id: string; name: string };
  speaker: any;
  timestamp: number;
  flavor: string;
  isRoll: boolean;
  rolls: {
    formula: string;
    total: number;
    isCritical: boolean;
    isFumble: boolean;
    dice: { faces: number; results: { result: number; active: boolean }[] }[];
  }[];
  whisper: string[];
  flags: Record<string, any>;
}

function formatMessage(prefix: string, message: ChatMessage) {
  const speaker = message.author?.name || message.speaker?.alias || '?';
  console.log(\`[\${prefix}] \${speaker}: \${message.content}\`);
  if (message.flavor) console.log(\`  Flavor: \${message.flavor}\`);
  if (message.isRoll && message.rolls?.length > 0) {
    for (const roll of message.rolls) {
      const dice = roll.dice?.map(d =>
        \`\${d.results.map(r => \`\${r.result}\${r.active ? '' : '(dropped)'}\`).join(', ')} (d\${d.faces})\`
      ).join(' + ') || '';
      console.log(\`  Roll: \${roll.formula} = \${roll.total}\${roll.isCritical ? ' CRITICAL!' : ''}\${roll.isFumble ? ' FUMBLE!' : ''}\`);
      if (dice) console.log(\`  Dice: \${dice}\`);
    }
  }
}

eventSource.addEventListener('connected', (event: MessageEvent) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
});

eventSource.addEventListener('chat-create', (event: MessageEvent) => {
  const message: ChatMessage = JSON.parse(event.data);
  formatMessage('new', message);
});

eventSource.addEventListener('chat-update', (event: MessageEvent) => {
  const message: ChatMessage = JSON.parse(event.data);
  formatMessage('updated', message);
});

eventSource.addEventListener('chat-delete', (event: MessageEvent) => {
  const data = JSON.parse(event.data);
  console.log('Message deleted:', JSON.stringify(data));
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
};

// To disconnect: eventSource.close();`,

    emojicode: `Just don't 😂`
  };
}

describe('Chat', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/chat-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`POST /chat - Send message (v${version})`, () => {
      test('POST /chat', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
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
              content: 'Hello from the REST API test suite!',
              flavor: 'Test Message'
            }, null, 2)
          }
        };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/chat - Send message'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('id');
        expect(captured.response.data.data).toHaveProperty('content', 'Hello from the REST API test suite!');
        expect(captured.response.data.data).toHaveProperty('flavor', 'Test Message');
        expect(captured.response.data.data).toHaveProperty('timestamp');
        expect(captured.response.data.data).toHaveProperty('speaker');
        expect(captured.response.data.data).toHaveProperty('whisper');
        expect(captured.response.data.data).toHaveProperty('author');

        // Save the message ID for subsequent tests
        setGlobalVariable(version, 'chatMessageId', captured.response.data.data.id);
      });
    });

    describe(`POST /chat - Send whispered message (v${version})`, () => {
      test('POST /chat with whisper', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
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
              content: 'This is a whispered test message',
              whisper: [],
              alias: 'API Test Bot'
            }, null, 2)
          }
        };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/chat - Send whispered message'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('id');
        expect(captured.response.data.data).toHaveProperty('content', 'This is a whispered test message');

        // Save for cleanup
        setGlobalVariable(version, 'chatWhisperMessageId', captured.response.data.data.id);
      });
    });

    describe(`GET /chat - Get messages (v${version})`, () => {
      test('GET /chat', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'limit',
                value: '10',
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
          '/chat - Get messages'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('messages');
        expect(captured.response.data.data).toHaveProperty('total');
        expect(captured.response.data.data).toHaveProperty('offset', 0);
        expect(captured.response.data.data).toHaveProperty('limit', 10);
        expect(Array.isArray(captured.response.data.data.messages)).toBe(true);
        expect(captured.response.data.data.messages.length).toBeGreaterThan(0);
        expect(captured.response.data.data.messages.length).toBeLessThanOrEqual(10);

        // Verify message structure
        const firstMessage = captured.response.data.data.messages[0];
        expect(firstMessage).toHaveProperty('id');
        expect(firstMessage).toHaveProperty('content');
        expect(firstMessage).toHaveProperty('timestamp');
        expect(firstMessage).toHaveProperty('speaker');
        expect(firstMessage).toHaveProperty('whisper');
        expect(firstMessage).toHaveProperty('author');
        expect(firstMessage).toHaveProperty('type');
      });
    });

    describe(`GET /chat - Get messages with pagination (v${version})`, () => {
      test('GET /chat with offset', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'limit',
                value: '5',
              },
              {
                key: 'offset',
                value: '0',
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
          '/chat - Get messages with pagination'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data.data).toHaveProperty('messages');
        expect(captured.response.data.data).toHaveProperty('offset', 0);
        expect(captured.response.data.data).toHaveProperty('limit', 5);
        expect(captured.response.data.data.messages.length).toBeLessThanOrEqual(5);
      });
    });

    describe(`DELETE /chat/:messageId - Delete message (v${version})`, () => {
      test('DELETE /chat/:messageId', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        const messageId = getGlobalVariable(version, 'chatMessageId');
        expect(messageId).toBeTruthy();

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: `{{baseUrl}}/chat/${messageId}`,
            host: ['{{baseUrl}}'],
            path: ['chat', messageId],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'DELETE',
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
          '/chat/:messageId - Delete message'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('messageId', messageId);
      });
    });

    describe(`DELETE /chat/:messageId - Delete whispered message (v${version})`, () => {
      test('DELETE /chat/:messageId (whispered)', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        const messageId = getGlobalVariable(version, 'chatWhisperMessageId');
        expect(messageId).toBeTruthy();

        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: `{{baseUrl}}/chat/${messageId}`,
            host: ['{{baseUrl}}'],
            path: ['chat', messageId],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'DELETE',
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/chat/:messageId - Delete whispered message'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
      });
    });

    describe(`DELETE /chat - Flush all messages (v${version})`, () => {
      test('DELETE /chat', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // First send a message so there's something to flush
        const setupConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
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
              content: 'Message to be flushed'
            }, null, 2)
          }
        };

        const setupCaptured = await captureExample(
          setupConfig,
          testVariables,
          '/chat - Setup message for flush'
        );
        expect(setupCaptured.response.status).toBe(200);

        // Now flush all chat messages
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              }
            ]
          },
          method: 'DELETE',
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
          '/chat - Flush all messages'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('message');

        // Verify flush worked by fetching messages
        const verifyConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/chat',
            host: ['{{baseUrl}}'],
            path: ['chat'],
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

        const verifyCaptured = await captureExample(
          verifyConfig,
          testVariables,
          '/chat - Verify flush'
        );
        expect(verifyCaptured.response.status).toBe(200);
        expect(verifyCaptured.response.data.data.total).toBe(0);
      });
    });

    describe(`GET /chat/subscribe - SSE connection (v${version})`, () => {
      test('GET /chat/subscribe receives chat events', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        const axios = (await import('axios')).default;
        const baseUrl = testVariables.baseUrl;
        const clientId = getClientId();
        const apiKey = testVariables.apiKey;

        // Step 1: Open SSE connection
        const sseResponse = await axios({
          url: `${baseUrl}/chat/subscribe?clientId=${clientId}`,
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

        // Step 2: Wait for connected event, then send a chat message, then wait for the chat event
        const chatEvent = await new Promise<{ connected: string; message: string }>((resolve, reject) => {
          let buffer = '';
          let connected = '';
          let messageTriggered = false;

          const timer = setTimeout(() => {
            sseResponse.data.destroy();
            reject(new Error(`Timed out waiting for chat SSE event. Buffer: ${buffer}`));
          }, 10000);

          sseResponse.data.on('data', async (chunk: Buffer) => {
            buffer += chunk.toString();

            // Once we see the connected event, send a chat message
            if (!messageTriggered && buffer.includes('event: connected')) {
              messageTriggered = true;
              connected = buffer;

              // Step 3: Send a chat message via the API
              try {
                await axios({
                  url: `${baseUrl}/chat?clientId=${clientId}`,
                  method: 'POST',
                  headers: {
                    'x-api-key': apiKey,
                    'Content-Type': 'application/json'
                  },
                  data: { content: 'SSE test message' }
                });
              } catch (err) {
                clearTimeout(timer);
                sseResponse.data.destroy();
                reject(new Error(`Failed to send chat message: ${err}`));
              }
            }

            // Check if we've received the chat-create event
            if (messageTriggered && buffer.includes('event: chat-create')) {
              clearTimeout(timer);
              sseResponse.data.destroy();
              resolve({ connected, message: buffer });
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
        expect(chatEvent.connected).toContain('event: connected');
        expect(chatEvent.connected).toContain(clientId);
        expect(chatEvent.message).toContain('event: chat-create');
        expect(chatEvent.message).toContain('SSE test message');
        expect(chatEvent.message).toContain('"content"');

        // Clean up: delete the test message we created
        try {
          // Extract the message ID from the SSE event data
          const eventDataMatch = chatEvent.message.match(/event: chat-create\ndata: (.+)\n/);
          if (eventDataMatch) {
            const eventData = JSON.parse(eventDataMatch[1]);
            if (eventData.id) {
              await axios({
                url: `${baseUrl}/chat/${eventData.id}?clientId=${clientId}`,
                method: 'DELETE',
                headers: { 'x-api-key': apiKey }
              });
            }
          }
        } catch { /* best effort cleanup */ }

        // Push a manually-crafted captured example with custom SSE code examples
        capturedExamples.push({
          endpoint: '/chat/subscribe',
          method: 'GET',
          description: '/chat/subscribe',
          request: {
            url: `${baseUrl}/chat/subscribe?clientId=${clientId}`,
            method: 'GET',
            headers: { 'x-api-key': apiKey }
          },
          response: {
            status: 200,
            data: { event: 'connected', data: { clientId } }
          },
          codeExamples: generateSSECodeExamples()
        });
      }, 20000);
    });

  });

});
