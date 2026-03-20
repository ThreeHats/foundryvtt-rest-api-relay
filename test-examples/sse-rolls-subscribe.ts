// npm install eventsource
import { EventSource } from 'eventsource';
import { config } from 'dotenv';
config({ path: '.env.test' });

const baseUrl = process.env.TEST_BASE_URL || 'http://localhost:3010';
const apiKey = process.env.TEST_API_KEY!;
const clientId = process.argv[2] || process.env.TEST_CLIENT_ID_V13 || process.env.TEST_CLIENT_ID_V12;

if (!clientId) {
  console.error('Usage: tsx test-examples/sse-rolls-subscribe.ts <clientId>');
  console.error('  Or set TEST_CLIENT_ID_V13 in .env.test');
  process.exit(1);
}

const url = `${baseUrl}/rolls/subscribe?clientId=${clientId}`;
console.log(`Connecting to ${url}...`);

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
  console.log('Listening for roll events... (Ctrl+C to stop)\n');
});

eventSource.addEventListener('roll', (event: MessageEvent) => {
  const roll: RollEvent = JSON.parse(event.data);
  const dice = roll.dice?.map(d =>
    `${d.results.map(r => `${r.result}${r.active ? '' : '(dropped)'}`).join(', ')} (d${d.faces})`
  ).join(' + ') || '';
  console.log(`[${roll.user?.name}] ${roll.formula} = ${roll.rollTotal}${roll.isCritical ? ' CRITICAL!' : ''}${roll.isFumble ? ' FUMBLE!' : ''}`);
  if (roll.flavor) console.log(`  Flavor: ${roll.flavor}`);
  if (dice) console.log(`  Dice: ${dice}`);
});

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
  eventSource.close();
  process.exit(1);
};

// Graceful shutdown
process.on('SIGINT', () => {
  console.log('\nDisconnecting...');
  eventSource.close();
  process.exit(0);
});
