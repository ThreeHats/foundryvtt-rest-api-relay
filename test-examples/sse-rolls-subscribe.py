#!/usr/bin/env python3
"""
SSE Rolls Subscribe example
pip install sseclient-py requests python-dotenv

Usage: python test-examples/sse-rolls-subscribe.py <clientId>
  Or set TEST_CLIENT_ID_V13 in .env.test
"""

import sys
import os
import json

try:
    import sseclient
    import requests
    from dotenv import load_dotenv
except ImportError as e:
    print(f"Missing dependency: {e}")
    print("Install with: pip install sseclient-py requests python-dotenv")
    sys.exit(1)

load_dotenv('.env.test')

base_url = os.environ.get('TEST_BASE_URL', 'http://localhost:3010')
api_key = os.environ.get('TEST_API_KEY', '')
client_id = sys.argv[1] if len(sys.argv) > 1 else os.environ.get('TEST_CLIENT_ID_V13') or os.environ.get('TEST_CLIENT_ID_V12')

if not client_id:
    print('Usage: python test-examples/sse-rolls-subscribe.py <clientId>')
    print('  Or set TEST_CLIENT_ID_V13 in .env.test')
    sys.exit(1)

url = f'{base_url}/rolls/subscribe'
params = {'clientId': client_id}
headers = {
    'x-api-key': api_key,
    'Accept': 'text/event-stream'
}

print(f'Connecting to {url}?clientId={client_id}...')

response = requests.get(url, params=params, headers=headers, stream=True)

if response.status_code != 200:
    print(f'Error: HTTP {response.status_code}')
    print(response.text)
    sys.exit(1)

client = sseclient.SSEClient(response)

try:
    for event in client.events():
        data = json.loads(event.data)

        if event.event == 'connected':
            print(f'Connected: {data["clientId"]}')
            print('Listening for roll events... (Ctrl+C to stop)\n')
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
                print(f'  Dice: {results} (d{d["faces"]})')
        else:
            print(f'Unknown event: {event.event} -> {data}')
except KeyboardInterrupt:
    print('\nDisconnecting...')
