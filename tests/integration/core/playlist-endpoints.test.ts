/**
 * @file playlist-endpoints.test.ts
 * @description Playlist Control Endpoint Tests
 * @endpoints GET /playlists, POST /playlist/play, POST /playlist/stop,
 *            POST /playlist/next, POST /playlist/volume, POST /play-sound,
 *            POST /stop-sound
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../../helpers/globalVariables';
import * as path from 'path';

const capturedExamples: any[] = [];

const TEST_PLAYLIST_NAME = 'test-rest-api-playlist';

describe('Playlist Control', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/playlist-examples.json');
      saveExamples(capturedExamples, outputPath);
    }
  });

  forEachVersion((version, getClientId) => {

    // ═══════════════════════════════════════════
    // Setup: Create a test playlist with a sound
    // ═══════════════════════════════════════════

    describe(`Setup: Create test playlist (v${version})`, () => {
      test('POST /create - Create test playlist with embedded sound', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/create',
            host: ['{{baseUrl}}'],
            path: ['create'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              entityType: 'Playlist',
              data: {
                name: TEST_PLAYLIST_NAME,
                sounds: [
                  {
                    name: 'test-sound',
                    path: 'sounds/dice.wav',
                    volume: 0.5,
                    repeat: false
                  }
                ]
              }
            })
          }
        };

        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolvedConfig);

        expect(response.status).toBe(200);
        // POST /create returns { entity: {...}, uuid: '...', type: '...', requestId: '...' }
        expect(response.data.entity).toHaveProperty('name', TEST_PLAYLIST_NAME);
        expect(response.data).toHaveProperty('uuid');

        setGlobalVariable(version, 'testPlaylistName', TEST_PLAYLIST_NAME);
        setGlobalVariable(version, 'testPlaylistUuid', response.data.uuid);
        console.log(`  ✓ Created test playlist: ${response.data.uuid}`);
      }, 30000);
    });

    // ═══════════════════════════════════════════
    // GET /playlists
    // ═══════════════════════════════════════════

    describe(`GET /playlists (v${version})`, () => {
      test('GET /playlists - list all playlists', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/playlists',
            host: ['{{baseUrl}}'],
            path: ['playlists'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const captured = await captureExample(requestConfig, testVariables, '/playlists');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        const data = captured.response.data.data;
        expect(data).toHaveProperty('playlists');
        expect(data.playlists).toBeInstanceOf(Array);

        // Verify our test playlist appears in the list
        const testPlaylist = data.playlists.find((p: any) => p.name === TEST_PLAYLIST_NAME);
        expect(testPlaylist).toBeTruthy();
        console.log(`  Found ${data.playlists.length} playlists (including test playlist)`);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /playlist/play
    // ═══════════════════════════════════════════

    describe(`POST /playlist/play (v${version})`, () => {
      test('POST /playlist/play - play a playlist by name', async () => {
        setVariable('clientId', getClientId());

        const playlistName = getGlobalVariable(version, 'testPlaylistName') as string;
        expect(playlistName).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/playlist/play',
            host: ['{{baseUrl}}'],
            path: ['playlist', 'play'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ playlistName })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/playlist/play');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('playlist');
        expect(captured.response.data.data.playlist).toHaveProperty('name', playlistName);
        expect(captured.response.data.data.playlist).toHaveProperty('playing', true);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /playlist/next
    // ═══════════════════════════════════════════

    describe(`POST /playlist/next (v${version})`, () => {
      test('POST /playlist/next - skip to next track', async () => {
        setVariable('clientId', getClientId());

        const playlistName = getGlobalVariable(version, 'testPlaylistName') as string;
        expect(playlistName).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/playlist/next',
            host: ['{{baseUrl}}'],
            path: ['playlist', 'next'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ playlistName })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/playlist/next');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('playlist');
        expect(captured.response.data.data.playlist).toHaveProperty('name', playlistName);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /playlist/volume
    // ═══════════════════════════════════════════

    describe(`POST /playlist/volume (v${version})`, () => {
      test('POST /playlist/volume - set playlist volume', async () => {
        setVariable('clientId', getClientId());

        const playlistName = getGlobalVariable(version, 'testPlaylistName') as string;
        expect(playlistName).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/playlist/volume',
            host: ['{{baseUrl}}'],
            path: ['playlist', 'volume'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ playlistName, volume: 0.75 })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/playlist/volume');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('playlist');
        expect(captured.response.data.data.playlist).toHaveProperty('name', playlistName);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /playlist/stop
    // ═══════════════════════════════════════════

    describe(`POST /playlist/stop (v${version})`, () => {
      test('POST /playlist/stop - stop a playlist by name', async () => {
        setVariable('clientId', getClientId());

        const playlistName = getGlobalVariable(version, 'testPlaylistName') as string;
        expect(playlistName).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/playlist/stop',
            host: ['{{baseUrl}}'],
            path: ['playlist', 'stop'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ playlistName })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/playlist/stop');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('playlist');
        expect(captured.response.data.data.playlist).toHaveProperty('name', playlistName);
        expect(captured.response.data.data.playlist).toHaveProperty('playing', false);
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /play-sound
    // ═══════════════════════════════════════════

    describe(`POST /play-sound (v${version})`, () => {
      test('POST /play-sound - play a one-shot sound effect', async () => {
        setVariable('clientId', getClientId());

        // sounds/dice.wav ships with every Foundry installation
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/play-sound',
            host: ['{{baseUrl}}'],
            path: ['play-sound'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              src: 'sounds/dice.wav',
              volume: 0.3,
              loop: false
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/play-sound');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('src', 'sounds/dice.wav');
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // POST /stop-sound
    // ═══════════════════════════════════════════

    describe(`POST /stop-sound (v${version})`, () => {
      test('POST /stop-sound - stop a playing sound by src', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/stop-sound',
            host: ['{{baseUrl}}'],
            path: ['stop-sound'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              src: 'sounds/dice.wav'
            })
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/stop-sound');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('stopped');
        expect(captured.response.data.data).toHaveProperty('src', 'sounds/dice.wav');
      }, 15000);

      test('POST /stop-sound - stop all playing sounds', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/stop-sound',
            host: ['{{baseUrl}}'],
            path: ['stop-sound'],
            query: [{ key: 'clientId', value: '{{clientId}}' }]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({})
          }
        };

        const captured = await captureExample(requestConfig, testVariables, '/stop-sound - all');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data.data).toHaveProperty('stopped');
        expect(captured.response.data.data.src).toBeNull();
      }, 15000);
    });

    // ═══════════════════════════════════════════
    // Cleanup: Delete test playlist
    // ═══════════════════════════════════════════

    describe(`Cleanup: Delete test playlist (v${version})`, () => {
      test('DELETE /delete - Delete test playlist', async () => {
        setVariable('clientId', getClientId());

        const playlistUuid = getGlobalVariable(version, 'testPlaylistUuid') as string;
        expect(playlistUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/delete',
            host: ['{{baseUrl}}'],
            path: ['delete'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: playlistUuid }
            ]
          },
          method: 'DELETE',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }]
        };

        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const response = await makeRequest(resolvedConfig);

        expect(response.status).toBe(200);
        expect(response.data.success).toBe(true);
        console.log(`  ✓ Deleted test playlist: ${playlistUuid}`);
      }, 30000);
    });
  });
});
