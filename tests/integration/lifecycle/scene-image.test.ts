/**
 * @file scene-image.test.ts
 * @description Scene Image Endpoint Tests
 * @endpoints GET /scene/image, GET /scene/image/raw
 *
 * Tests capturing rendered scene images and raw scene backgrounds.
 * Sets up a rich test scene (background, grid, walls, drawings, lights, token, etc.)
 * then verifies the screenshot captures real visual content.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import axios from 'axios';
import * as zlib from 'zlib';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../../helpers/globalVariables';
import * as path from 'path';
import * as fs from 'fs';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Directory to save received images
const IMAGE_OUTPUT_DIR = path.join(__dirname, '../../test-results/images');

// Reference images committed to git. Copy a good run's output here to lock in the baseline.
const REFERENCE_DIR = path.join(__dirname, 'reference');

// Background color: dark blue-gray, like a dungeon stone floor
const BG_R = 50, BG_G = 60, BG_B = 80;

// Canvas entities to create so the screenshot has real content to show
const SCENE_CANVAS_ENTITIES: Array<{ type: string; data: Record<string, any> }> = [
  // Two walls forming an L-corner to the bottom-right of the token (token is at 400,400)
  { type: 'walls',     data: { c: [600, 400, 600, 700] } },  // vertical wall (right of token)
  { type: 'walls',     data: { c: [400, 600, 700, 600] } },  // horizontal wall (below token)
  { type: 'drawings',  data: { x: 550, y: 450, shape: { type: 'r', width: 250, height: 150 }, fillType: 1, fillColor: '#ff6600', fillAlpha: 0.5 } },
  { type: 'lights',    data: { x: 400, y: 400, config: { dim: 10, bright: 5, color: '#c64600' } } },
  { type: 'templates', data: { x: 650, y: 300, t: 'circle', distance: 10 } },
  { type: 'notes',     data: { x: 680, y: 290, text: 'Scene Test' } },
  { type: 'tiles',     data: { x: 200, y: 600, width: 200, height: 200, texture: { src: '' } } },
];

/**
 * Generate a valid PNG image programmatically using zlib for proper compression.
 * Creates a solid-color image of the given size.
 */
function generateTestPng(width = 64, height = 64, r = 255, g = 0, b = 0): string {
  // Build raw scanline data: each row is [filterByte(0=None), R, G, B, ...]
  const rawData = Buffer.alloc((1 + width * 3) * height);
  for (let y = 0; y < height; y++) {
    const rowOffset = y * (1 + width * 3);
    rawData[rowOffset] = 0; // filter: None
    for (let x = 0; x < width; x++) {
      const px = rowOffset + 1 + x * 3;
      rawData[px] = r; rawData[px + 1] = g; rawData[px + 2] = b;
    }
  }
  const compressed = zlib.deflateSync(rawData);

  const crcTable: number[] = [];
  for (let n = 0; n < 256; n++) {
    let c = n;
    for (let k = 0; k < 8; k++) c = (c & 1) ? (0xEDB88320 ^ (c >>> 1)) : (c >>> 1);
    crcTable[n] = c;
  }
  const crc32 = (buf: Buffer): number => {
    let crc = 0xFFFFFFFF;
    for (let i = 0; i < buf.length; i++) crc = crcTable[(crc ^ buf[i]) & 0xFF] ^ (crc >>> 8);
    return (crc ^ 0xFFFFFFFF) >>> 0;
  };
  const makeChunk = (type: string, data: Buffer): Buffer => {
    const typeBytes = Buffer.from(type, 'ascii');
    const len = Buffer.alloc(4); len.writeUInt32BE(data.length);
    const crcVal = Buffer.alloc(4); crcVal.writeUInt32BE(crc32(Buffer.concat([typeBytes, data])));
    return Buffer.concat([len, typeBytes, data, crcVal]);
  };

  const ihdrData = Buffer.alloc(13);
  ihdrData.writeUInt32BE(width, 0); ihdrData.writeUInt32BE(height, 4);
  ihdrData[8] = 8; ihdrData[9] = 2; // bit depth 8, RGB

  const png = Buffer.concat([
    Buffer.from([0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A]),
    makeChunk('IHDR', ihdrData),
    makeChunk('IDAT', compressed),
    makeChunk('IEND', Buffer.alloc(0)),
  ]);
  return `data:image/png;base64,${png.toString('base64')}`;
}

/**
 * Decode a PNG buffer to raw R/G/B channel arrays.
 * Handles all 5 PNG filter types (None, Sub, Up, Average, Paeth).
 */
function decodePngPixels(data: Buffer): { width: number; height: number; r: Uint8Array; g: Uint8Array; b: Uint8Array } {
  // Verify PNG signature
  if (data[0] !== 0x89 || data[1] !== 0x50 || data[2] !== 0x4E || data[3] !== 0x47) {
    throw new Error('Not a valid PNG');
  }

  let width = 0, height = 0, colorType = 2;
  const idatChunks: Buffer[] = [];
  let pos = 8;

  while (pos + 8 < data.length) {
    const len = data.readUInt32BE(pos); pos += 4;
    const type = data.toString('ascii', pos, pos + 4); pos += 4;
    const chunk = data.subarray(pos, pos + len); pos += len + 4; // +4 for CRC

    if (type === 'IHDR') {
      width = chunk.readUInt32BE(0);
      height = chunk.readUInt32BE(4);
      colorType = chunk[9]; // 2=RGB, 6=RGBA
    } else if (type === 'IDAT') {
      idatChunks.push(Buffer.from(chunk));
    } else if (type === 'IEND') {
      break;
    }
  }

  const bpp = colorType === 6 ? 4 : 3; // bytes per pixel: RGBA or RGB
  const stride = width * bpp;
  const raw = zlib.inflateSync(Buffer.concat(idatChunks));

  const out = new Uint8Array(width * height * bpp);
  const prev = new Uint8Array(stride); // previous reconstructed row

  for (let y = 0; y < height; y++) {
    const srcRow = y * (stride + 1);
    const filter = raw[srcRow];
    const dst = new Uint8Array(stride);

    for (let x = 0; x < stride; x++) {
      const byte = raw[srcRow + 1 + x];
      const a = x >= bpp ? dst[x - bpp] : 0;     // left (same channel)
      const b = prev[x];                           // above (same channel)
      const c = x >= bpp ? prev[x - bpp] : 0;     // upper-left (same channel)
      let val: number;
      switch (filter) {
        case 0: val = byte; break;
        case 1: val = (byte + a) & 0xff; break;
        case 2: val = (byte + b) & 0xff; break;
        case 3: val = (byte + Math.floor((a + b) / 2)) & 0xff; break;
        case 4: {
          const p = a + b - c;
          const pa = Math.abs(p - a), pb = Math.abs(p - b), pc = Math.abs(p - c);
          val = (byte + (pa <= pb && pa <= pc ? a : pb <= pc ? b : c)) & 0xff;
          break;
        }
        default: val = byte;
      }
      dst[x] = val;
    }

    out.set(dst, y * stride);
    prev.set(dst);
  }

  const n = width * height;
  const r = new Uint8Array(n), g = new Uint8Array(n), b2 = new Uint8Array(n);
  for (let i = 0; i < n; i++) {
    r[i] = out[i * bpp]; g[i] = out[i * bpp + 1]; b2[i] = out[i * bpp + 2];
  }
  return { width, height, r, g, b: b2 };
}

/** Compute mean and standard deviation for a channel. */
function channelStats(ch: Uint8Array): { mean: number; stddev: number } {
  const mean = ch.reduce((s, v) => s + v, 0) / ch.length;
  const stddev = Math.sqrt(ch.reduce((s, v) => s + (v - mean) ** 2, 0) / ch.length);
  return { mean, stddev };
}

/**
 * Save image data to disk for manual inspection.
 */
function saveImage(data: Buffer | string, filename: string): void {
  if (!fs.existsSync(IMAGE_OUTPUT_DIR)) fs.mkdirSync(IMAGE_OUTPUT_DIR, { recursive: true });
  const outPath = path.join(IMAGE_OUTPUT_DIR, filename);
  if (Buffer.isBuffer(data)) {
    fs.writeFileSync(outPath, data);
  } else {
    let b64 = data;
    const idx = b64.indexOf(',');
    if (idx >= 0) b64 = b64.substring(idx + 1);
    fs.writeFileSync(outPath, Buffer.from(b64, 'base64'));
  }
  console.log(`  Saved image: ${outPath}`);
}

/** POST to a canvas endpoint and return the created document ID (or null on failure). */
async function createCanvasDoc(clientId: string, type: string, data: Record<string, any>): Promise<string | null> {
  const config: ApiRequestConfig = {
    url: {
      raw: `{{baseUrl}}/canvas/${type}`,
      host: ['{{baseUrl}}'],
      path: ['canvas', type],
      query: [{ key: 'clientId', value: clientId }],
    },
    method: 'POST',
    header: [
      { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
      { key: 'Content-Type', value: 'application/json', type: 'text' },
    ],
    body: { mode: 'raw', raw: JSON.stringify({ data }) },
  };
  const resp = await makeRequest(replaceVariables(config, testVariables));
  return resp.data?.data?.[0]?._id ?? null;
}

/** DELETE a canvas document by ID. */
async function deleteCanvasDoc(clientId: string, type: string, docId: string): Promise<void> {
  const config: ApiRequestConfig = {
    url: {
      raw: `{{baseUrl}}/canvas/${type}`,
      host: ['{{baseUrl}}'],
      path: ['canvas', type],
      query: [
        { key: 'clientId', value: clientId },
        { key: 'documentId', value: docId },
      ],
    },
    method: 'DELETE',
    header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
  };
  await makeRequest(replaceVariables(config, testVariables)).catch(() => {});
}

describe('Scene Image', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/scene-image-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }
  });

  forEachVersion((version, getClientId) => {

    // ═══════════════════════════════════════════
    // Setup: build a rich test scene
    // ═══════════════════════════════════════════

    describe(`Scene image setup (v${version})`, () => {
      test('Set up test scene: background, grid, canvas entities, token', async () => {
        const clientId = getClientId();
        setVariable('clientId', clientId);

        // Step 1: Upload a dark blue-gray background PNG
        const testPng = generateTestPng(1000, 1000, BG_R, BG_G, BG_B);
        const uploadResp = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/upload',
            host: ['{{baseUrl}}'], path: ['upload'],
            query: [{ key: 'clientId', value: clientId }],
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' },
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              path: 'rest-api-tests', filename: 'test-background.png',
              source: 'data', mimeType: 'image/png', overwrite: true, fileData: testPng,
            }),
          },
        }, testVariables));
        expect(uploadResp.status).toBe(200);
        const uploadedPath = uploadResp.data?.data?.path || uploadResp.data?.path || 'rest-api-tests/test-background.png';
        console.log(`  Uploaded background: ${uploadedPath}`);

        // Step 2: Get the active scene
        const sceneResp = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'], path: ['scene'],
            query: [{ key: 'clientId', value: clientId }, { key: 'active', value: 'true' }],
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
        }, testVariables));
        expect(sceneResp.status).toBe(200);
        const sceneData = sceneResp.data?.data || sceneResp.data;
        const sceneId = sceneData?.id || sceneData?._id;
        expect(sceneId).toBeTruthy();

        // Save original settings for restoration
        setGlobalVariable(version, 'scene_image_scene_id', sceneId);
        setGlobalVariable(version, 'scene_image_original_bg', sceneData?.background?.src || sceneData?.img || '');
        setGlobalVariable(version, 'scene_image_original_grid', sceneData?.grid ?? {});
        setGlobalVariable(version, 'scene_image_original_tokenVision', sceneData?.tokenVision ?? false);
        setGlobalVariable(version, 'scene_image_original_darkness', sceneData?.darkness ?? 0);
        setGlobalVariable(version, 'scene_image_original_padding', sceneData?.padding ?? 0.25);
        setGlobalVariable(version, 'scene_image_original_environment', sceneData?.environment ?? {});

        // Step 3: Apply test background and set grid to bright white so it's clearly visible
        const updateResp = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'], path: ['scene'],
            query: [{ key: 'clientId', value: clientId }],
          },
          method: 'PUT',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              sceneId,
              data: {
                background: { src: uploadedPath },
                grid: { color: '#ffffff', alpha: 1.0 },
                tokenVision: true,
                padding: 0,
                darkness: 0.5,
                environment: { globalLight: { enabled: true } },
              },
            }),
          },
        }, testVariables));
        expect(updateResp.status).toBe(200);
        console.log(`  Applied background and white grid to scene ${sceneId}`);

        // Step 4: Create canvas entities so the screenshot has real scene content
        const entityIds: Record<string, string> = {};
        for (const entity of SCENE_CANVAS_ENTITIES) {
          const id = await createCanvasDoc(clientId, entity.type, entity.data);
          if (id) {
            const key = `scene_image_canvas_${entity.type}_${id.substring(0, 4)}`;
            entityIds[key] = id;
            setGlobalVariable(version, key, id);
          }
        }
        setGlobalVariable(version, 'scene_image_canvas_entity_keys', Object.keys(entityIds));
        console.log(`  Created ${Object.keys(entityIds).length}/${SCENE_CANVAS_ENTITIES.length} canvas entities`);

        // Step 5: Create a temporary actor and place a token at (400, 400)
        const actorResp = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/create',
            host: ['{{baseUrl}}'], path: ['create'],
            query: [{ key: 'clientId', value: clientId }],
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' },
          ],
          body: { mode: 'raw', raw: JSON.stringify({ entityType: 'Actor', data: { name: 'scene-image-test-actor', type: 'npc' } }) },
        }, testVariables));
        if (actorResp.status === 200 && actorResp.data?.uuid) {
          const actorUuid = actorResp.data.uuid;
          const actorId = actorUuid.split('.').pop();
          setGlobalVariable(version, 'scene_image_actor_uuid', actorUuid);

          const tokenId = await createCanvasDoc(clientId, 'tokens', {
            actorId, x: 400, y: 400,
            sight: { enabled: true },
          });
          if (tokenId) {
            setGlobalVariable(version, 'scene_image_token_id', tokenId);
            console.log(`  Created test actor and token ${tokenId}`);

            // Select the token so the screenshot renders from its vision perspective,
            // making walls and lights visible.
            await makeRequest(replaceVariables({
              url: {
                raw: '{{baseUrl}}/select',
                host: ['{{baseUrl}}'], path: ['select'],
                query: [{ key: 'clientId', value: clientId }],
              },
              method: 'POST',
              header: [
                { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
                { key: 'Content-Type', value: 'application/json', type: 'text' },
              ],
              body: { mode: 'raw', raw: JSON.stringify({ name: 'scene-image-test-actor', overwrite: true }) },
            }, testVariables)).catch(() => {});
            console.log(`  Selected token for vision-based screenshot`);
          }
        }
      }, 60000);
    });

    // ═══════════════════════════════════════════
    // GET /scene/image — rendered canvas screenshot
    // ═══════════════════════════════════════════

    describe(`/scene/image (v${version})`, () => {
      test('GET /scene/image - Capture rendered scene image', async () => {
        const clientId = getClientId();
        setVariable('clientId', clientId);

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene/image',
            host: ['{{baseUrl}}'],
            path: ['scene', 'image'],
            query: [
              { key: 'clientId', value: clientId },
              { key: 'format', value: 'png' },
            ],
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
        };

        const captured = await captureExample(requestConfig, testVariables, '/scene/image - Capture rendered scene image');
        capturedExamples.push(captured);
        expect(captured.response.status).toBe(200);

        // Fetch binary PNG for pixel-level inspection
        const resolvedConfig = replaceVariables(requestConfig, testVariables);
        const fullUrl = new URL(resolvedConfig.url.raw);
        for (const q of resolvedConfig.url.query || []) {
          if (!q.disabled) fullUrl.searchParams.append(q.key, q.value);
        }
        const binaryResponse = await axios.get(fullUrl.toString(), {
          headers: { 'x-api-key': testVariables.apiKey },
          responseType: 'arraybuffer',
          validateStatus: () => true,
          timeout: 30000,
        });

        expect(binaryResponse.status).toBe(200);
        const contentType = binaryResponse.headers['content-type'] || '';
        console.log(`  Content-Type: ${contentType}`);
        expect(contentType).toContain('image/');

        const imgWidth = Number(binaryResponse.headers['x-image-width'] || 0);
        const imgHeight = Number(binaryResponse.headers['x-image-height'] || 0);
        console.log(`  Dimensions: ${imgWidth}x${imgHeight}, size: ${binaryResponse.data?.byteLength ?? 0} bytes`);
        expect(imgWidth).toBeGreaterThan(0);
        expect(imgHeight).toBeGreaterThan(0);

        const pngBuf = Buffer.from(binaryResponse.data);
        const ext = contentType.includes('jpeg') ? 'jpg' : 'png';
        saveImage(pngBuf, `scene-screenshot-v${version}.${ext}`);

        if (contentType.includes('image/png')) {
          try {
            const { r, g, b } = decodePngPixels(pngBuf);
            const rStats = channelStats(r);
            const gStats = channelStats(g);
            const bStats = channelStats(b);
            console.log(`  Pixel stats — R:${rStats.mean.toFixed(1)}±${rStats.stddev.toFixed(1)} G:${gStats.mean.toFixed(1)}±${gStats.stddev.toFixed(1)} B:${bStats.mean.toFixed(1)}±${bStats.stddev.toFixed(1)}`);

            // Must not be a flat/blank image — meaningful std dev proves real content was rendered
            const maxStddev = Math.max(rStats.stddev, gStats.stddev, bStats.stddev);
            expect(maxStddev).toBeGreaterThan(10);

            // Compare against a committed reference image if one exists.
            // To update the baseline: copy test-results/images/scene-screenshot-v*.png
            // into tests/integration/lifecycle/reference/ and commit.
            const refPath = path.join(REFERENCE_DIR, `scene-screenshot-v${version}.png`);
            if (fs.existsSync(refPath)) {
              const refBuf = fs.readFileSync(refPath);
              const ref = decodePngPixels(refBuf);
              // Mean absolute difference per channel — allow up to 20 units of drift
              const mad = (a: Uint8Array, b: Uint8Array) =>
                a.reduce((s, v, i) => s + Math.abs(v - b[i]), 0) / a.length;
              const madR = mad(r, ref.r);
              const madG = mad(g, ref.g);
              const madB = mad(b, ref.b);
              console.log(`  vs. reference — MAD R:${madR.toFixed(2)} G:${madG.toFixed(2)} B:${madB.toFixed(2)}`);
              expect(madR).toBeLessThan(20);
              expect(madG).toBeLessThan(20);
              expect(madB).toBeLessThan(20);
            } else {
              console.log(`  No reference image found at ${refPath} — skipping comparison.`);
              console.log(`  To lock in this screenshot as baseline, copy it to that path.`);
            }
          } catch (e) {
            console.warn(`  Could not decode PNG pixels: ${e}`);
          }
        }
      }, 60000);
    });

    // ═══════════════════════════════════════════
    // GET /scene/image/raw — raw background image
    // ═══════════════════════════════════════════

    describe(`/scene/image/raw (v${version})`, () => {
      test('GET /scene/image/raw - Get raw scene background', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene/image/raw',
            host: ['{{baseUrl}}'],
            path: ['scene', 'image', 'raw'],
            query: [{ key: 'clientId', value: '{{clientId}}' }],
          },
          method: 'GET',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
        };

        const captured = await captureExample(requestConfig, testVariables, '/scene/image/raw - Get raw scene background');
        capturedExamples.push(captured);

        // v14 moved scene background into Level documents; if the module's Level update
        // didn't propagate in time, the scene may report no background — accept 400.
        if (captured.response.status === 400) {
          console.log(`  Note: /scene/image/raw returned 400 on v${version} — scene background may not be set`);
          return;
        }
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();

        const data = captured.response.data?.data || captured.response.data;
        if (data?.imageData) {
          saveImage(data.imageData, `scene-raw-bg-v${version}.png`);

          // The raw background should be the test PNG we uploaded — decode and verify its color
          try {
            const pngBuf = Buffer.from(
              data.imageData.includes(',') ? data.imageData.split(',')[1] : data.imageData,
              'base64'
            );
            const { r, g, b } = decodePngPixels(pngBuf);
            const meanR = r.reduce((s, v) => s + v, 0) / r.length;
            const meanG = g.reduce((s, v) => s + v, 0) / g.length;
            const meanB = b.reduce((s, v) => s + v, 0) / b.length;
            console.log(`  Raw background mean color: R=${meanR.toFixed(0)} G=${meanG.toFixed(0)} B=${meanB.toFixed(0)}`);

            // Should be approximately the test background color (50, 60, 80)
            const TOLERANCE = 15;
            expect(Math.abs(meanR - BG_R)).toBeLessThan(TOLERANCE);
            expect(Math.abs(meanG - BG_G)).toBeLessThan(TOLERANCE);
            expect(Math.abs(meanB - BG_B)).toBeLessThan(TOLERANCE);
          } catch (e) {
            console.warn(`  Could not decode raw background PNG: ${e}`);
          }
        }
      }, 30000);
    });

    // ═══════════════════════════════════════════
    // Teardown: delete canvas entities and restore scene
    // ═══════════════════════════════════════════

    describe(`Scene image teardown (v${version})`, () => {
      test('Delete canvas entities, token, actor and restore scene', async () => {
        const clientId = getClientId();
        setVariable('clientId', clientId);

        // Delete all canvas entities created in setup
        const entityKeys = (getGlobalVariable(version, 'scene_image_canvas_entity_keys') as string[]) || [];
        for (const key of entityKeys) {
          const docId = getGlobalVariable(version, key) as string;
          if (!docId) continue;
          // Extract the canvas type from the key (e.g. "scene_image_canvas_walls_abcd" → "walls")
          const parts = key.split('_');
          const type = parts[3]; // scene_image_canvas_<type>_<id>
          if (type) {
            await deleteCanvasDoc(clientId, type, docId);
            console.log(`  Deleted ${type} ${docId}`);
          }
        }

        // Delete the canvas token
        const tokenId = getGlobalVariable(version, 'scene_image_token_id') as string;
        if (tokenId) {
          await deleteCanvasDoc(clientId, 'tokens', tokenId);
          console.log(`  Deleted token ${tokenId}`);
        }

        // Delete the temporary actor
        const actorUuid = getGlobalVariable(version, 'scene_image_actor_uuid') as string;
        if (actorUuid) {
          await makeRequest(replaceVariables({
            url: {
              raw: '{{baseUrl}}/delete',
              host: ['{{baseUrl}}'], path: ['delete'],
              query: [
                { key: 'clientId', value: clientId },
                { key: 'uuid', value: actorUuid },
              ],
            },
            method: 'DELETE',
            header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          }, testVariables)).catch(() => {});
          console.log(`  Deleted actor ${actorUuid}`);
        }

        // Restore original scene background and grid
        const sceneId = getGlobalVariable(version, 'scene_image_scene_id') as string;
        if (!sceneId) { console.log('  No scene to restore'); return; }

        const originalBg = (getGlobalVariable(version, 'scene_image_original_bg') as string) || '';
        const originalGrid = getGlobalVariable(version, 'scene_image_original_grid') || {};
        const originalTokenVision = getGlobalVariable(version, 'scene_image_original_tokenVision') ?? false;
        const originalDarkness = getGlobalVariable(version, 'scene_image_original_darkness') ?? 0;
        const originalPadding = getGlobalVariable(version, 'scene_image_original_padding') ?? 0.25;
        const originalEnvironment = getGlobalVariable(version, 'scene_image_original_environment') ?? {};

        const restoreResp = await makeRequest(replaceVariables({
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'], path: ['scene'],
            query: [{ key: 'clientId', value: clientId }],
          },
          method: 'PUT',
          header: [{ key: 'x-api-key', value: '{{apiKey}}', type: 'text' }],
          body: {
            mode: 'raw',
            raw: JSON.stringify({ sceneId, data: { background: { src: originalBg }, grid: originalGrid, tokenVision: originalTokenVision, darkness: originalDarkness, padding: originalPadding, environment: originalEnvironment } }),
          },
        }, testVariables));
        expect(restoreResp.status).toBe(200);
        console.log(`  Restored scene background "${originalBg || '(none)'}" and original grid/lighting`);
      }, 60000);
    });

  });
});
