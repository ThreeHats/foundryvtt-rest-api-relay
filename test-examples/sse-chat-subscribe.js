// Browser-compatible JavaScript SSE example
// Run with: node test-examples/sse-chat-subscribe.js <clientId>
//
// Note: In a browser, the native EventSource API does not support custom headers.
// This Node.js example uses the 'eventsource' package which supports custom fetch.
// In a browser, you would need to pass auth via query params or cookies instead.

const { EventSource } = require('eventsource');
require('dotenv').config({ path: '.env.test' });

const baseUrl = process.env.TEST_BASE_URL || 'http://localhost:3010';
const apiKey = process.env.TEST_API_KEY;
const clientId = process.argv[2] || process.env.TEST_CLIENT_ID_V13 || process.env.TEST_CLIENT_ID_V12;

if (!clientId) {
  console.error('Usage: node test-examples/sse-chat-subscribe.js <clientId>');
  console.error('  Or set TEST_CLIENT_ID_V13 in .env.test');
  process.exit(1);
}

const url = `${baseUrl}/chat/subscribe?clientId=${clientId}`;
console.log(`Connecting to ${url}...`);

// eventsource v4 uses a custom fetch function to inject headers
const eventSource = new EventSource(url, {
  fetch: (input, init) => fetch(input, {
    ...init,
    headers: { ...init?.headers, 'x-api-key': apiKey }
  })
});

function formatMessage(prefix, message) {
  const speaker = message.author?.name || message.speaker?.alias || '?';
  console.log(`[${prefix}] ${speaker}: ${message.content}`);
  if (message.flavor) console.log(`  Flavor: ${message.flavor}`);
  if (message.isRoll && message.rolls?.length > 0) {
    for (const roll of message.rolls) {
      const dice = roll.dice?.map(d =>
        `${d.results.map(r => `${r.result}${r.active ? '' : '(dropped)'}`).join(', ')} (d${d.faces})`
      ).join(' + ') || '';
      console.log(`  Roll: ${roll.formula} = ${roll.total}${roll.isCritical ? ' CRITICAL!' : ''}${roll.isFumble ? ' FUMBLE!' : ''}`);
      if (dice) console.log(`  Dice: ${dice}`);
    }
  }
}

eventSource.addEventListener('connected', (event) => {
  const data = JSON.parse(event.data);
  console.log('Connected:', data.clientId);
  console.log('Listening for chat events... (Ctrl+C to stop)\n');
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
  eventSource.close();
  process.exit(1);
};

process.on('SIGINT', () => {
  console.log('\nDisconnecting...');
  eventSource.close();
  process.exit(0);
});
