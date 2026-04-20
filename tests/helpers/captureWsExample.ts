import * as fs from 'fs';
import * as path from 'path';

interface CapturedWsExample {
  messageType: string;
  description: string;
  request: Record<string, any>;
  response: Record<string, any>;
  codeExamples: {
    javascript: string;
    python: string;
    typescript: string;
  };
}

/**
 * Capture a WebSocket request/response pair for documentation
 */
export function captureWsExample(
  messageType: string,
  description: string,
  request: Record<string, any>,
  response: Record<string, any>,
  wsUrl: string
): CapturedWsExample {
  return {
    messageType,
    description,
    request,
    response,
    codeExamples: {
      javascript: generateJsWsExample(messageType, request, wsUrl),
      python: generatePythonWsExample(messageType, request, wsUrl),
      typescript: generateTsWsExample(messageType, request, wsUrl),
    },
  };
}

function generateJsWsExample(type: string, request: Record<string, any>, wsUrl: string): string {
  const msg = JSON.stringify({ ...request, type, requestId: 'unique-id' }, null, 2);
  return `const ws = new WebSocket('${wsUrl}?clientId=YOUR_CLIENT_ID');

ws.onopen = () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify(${msg.replace(/\n/g, '\n    ')}));
  }
  if (data.type === '${type}-result') {
    console.log(data);
  }
};`;
}

function generatePythonWsExample(type: string, request: Record<string, any>, wsUrl: string): string {
  const msg = JSON.stringify({ ...request, type, requestId: 'unique-id' });
  // Convert JSON booleans for Python display
  const pyMsg = msg.replace(/\btrue\b/g, 'True').replace(/\bfalse\b/g, 'False').replace(/\bnull\b/g, 'None');
  return `import asyncio
import websockets
import json

async def main():
    uri = '${wsUrl}?clientId=YOUR_CLIENT_ID'
    async with websockets.connect(uri) as ws:
        # Send auth message first — token must not be in the URL
        await ws.send(json.dumps({'type': 'auth', 'token': 'YOUR_API_KEY'}))
        connected = json.loads(await ws.recv())
        if connected.get('type') != 'connected':
            raise Exception('Auth failed')
        await ws.send(json.dumps(${pyMsg}))
        response = await ws.recv()
        data = json.loads(response)
        print(data)

asyncio.run(main())`;
}

function generateTsWsExample(type: string, request: Record<string, any>, wsUrl: string): string {
  const msg = JSON.stringify({ ...request, type, requestId: 'unique-id' }, null, 2);
  return `import WebSocket from 'ws';

const ws = new WebSocket('${wsUrl}?clientId=YOUR_CLIENT_ID');

ws.on('open', () => {
  // Send auth message first — token must not be in the URL
  ws.send(JSON.stringify({ type: 'auth', token: 'YOUR_API_KEY' }));
});

ws.on('message', (raw: string) => {
  const data = JSON.parse(raw);
  if (data.type === 'connected') {
    // Now send your request
    ws.send(JSON.stringify(${msg.replace(/\n/g, '\n    ')}));
  }
  if (data.type === '${type}-result') {
    console.log(data);
  }
}):`;
}

/**
 * Save captured WS examples to a file (deduplicates by messageType)
 */
export function saveWsExamples(examples: CapturedWsExample[], outputPath: string): void {
  const dir = path.dirname(outputPath);
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
  // Deduplicate by messageType
  const seen = new Map<string, CapturedWsExample>();
  for (const example of examples) {
    if (!seen.has(example.messageType)) {
      seen.set(example.messageType, example);
    }
  }
  fs.writeFileSync(outputPath, JSON.stringify(Array.from(seen.values()), null, 2));
}
