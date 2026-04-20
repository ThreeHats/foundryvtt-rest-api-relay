/**
 * @file scene-image.test.ts
 * @description Scene Image Endpoint Tests
 * @endpoints GET /scene/image, GET /scene/image/raw
 *
 * Tests capturing rendered scene images and raw scene backgrounds.
 * Uploads a test background image to the active scene before testing /scene/image/raw.
 * Saves received images to test-results/images/ for manual inspection.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import axios from 'axios';
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

/**
 * Generate a valid PNG image programmatically using zlib for proper compression.
 * Creates a solid-color image of the given size.
 */
function generateTestPng(width = 64, height = 64, r = 255, g = 0, b = 0): string {
  const zlib = require('zlib');

  // Build raw scanline data: each row is [filterByte, R, G, B, R, G, B, ...]
  const rawData = Buffer.alloc((1 + width * 3) * height);
  for (let y = 0; y < height; y++) {
    const rowOffset = y * (1 + width * 3);
    rawData[rowOffset] = 0; // filter: None
    for (let x = 0; x < width; x++) {
      const px = rowOffset + 1 + x * 3;
      rawData[px] = r;
      rawData[px + 1] = g;
      rawData[px + 2] = b;
    }
  }

  const compressed = zlib.deflateSync(rawData);

  // CRC32 implementation for PNG chunks
  const crcTable: number[] = [];
  for (let n = 0; n < 256; n++) {
    let c = n;
    for (let k = 0; k < 8; k++) c = (c & 1) ? (0xEDB88320 ^ (c >>> 1)) : (c >>> 1);
    crcTable[n] = c;
  }
  function crc32(buf: Buffer): number {
    let crc = 0xFFFFFFFF;
    for (let i = 0; i < buf.length; i++) crc = crcTable[(crc ^ buf[i]) & 0xFF] ^ (crc >>> 8);
    return (crc ^ 0xFFFFFFFF) >>> 0;
  }

  function makeChunk(type: string, data: Buffer): Buffer {
    const typeBytes = Buffer.from(type, 'ascii');
    const len = Buffer.alloc(4);
    len.writeUInt32BE(data.length);
    const crcInput = Buffer.concat([typeBytes, data]);
    const crcVal = Buffer.alloc(4);
    crcVal.writeUInt32BE(crc32(crcInput));
    return Buffer.concat([len, typeBytes, data, crcVal]);
  }

  // IHDR: width, height, bit depth 8, color type 2 (RGB), compression 0, filter 0, interlace 0
  const ihdrData = Buffer.alloc(13);
  ihdrData.writeUInt32BE(width, 0);
  ihdrData.writeUInt32BE(height, 4);
  ihdrData[8] = 8;  // bit depth
  ihdrData[9] = 2;  // color type RGB
  ihdrData[10] = 0; // compression
  ihdrData[11] = 0; // filter
  ihdrData[12] = 0; // interlace

  const signature = Buffer.from([0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A]);
  const ihdr = makeChunk('IHDR', ihdrData);
  const idat = makeChunk('IDAT', compressed);
  const iend = makeChunk('IEND', Buffer.alloc(0));

  const png = Buffer.concat([signature, ihdr, idat, iend]);
  return `data:image/png;base64,${png.toString('base64')}`;
}

/**
 * Save image data to disk for inspection.
 * @param data - Raw binary buffer, base64 string, or data URL string
 * @param filename - Output filename (e.g., "scene-screenshot-v13.png")
 */
function saveImage(data: Buffer | string, filename: string): void {
  if (!fs.existsSync(IMAGE_OUTPUT_DIR)) {
    fs.mkdirSync(IMAGE_OUTPUT_DIR, { recursive: true });
  }
  const outPath = path.join(IMAGE_OUTPUT_DIR, filename);

  if (Buffer.isBuffer(data)) {
    fs.writeFileSync(outPath, data);
  } else if (typeof data === 'string') {
    // Strip data URL prefix if present
    let b64 = data;
    const commaIdx = b64.indexOf(',');
    if (commaIdx >= 0) {
      b64 = b64.substring(commaIdx + 1);
    }
    fs.writeFileSync(outPath, Buffer.from(b64, 'base64'));
  }
  console.log(`  Saved image: ${outPath}`);
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
    // Setup: upload a test background image to the active scene
    // ═══════════════════════════════════════════

    describe(`Scene image setup (v${version})`, () => {
      test('Upload test background image and set it on the active scene', async () => {
        setVariable('clientId', getClientId());

        // Step 1: Upload a test PNG to Foundry's data directory
        // Use a large enough image so the canvas has meaningful content to render
        const testPng = generateTestPng(400, 400, 180, 40, 40);

        const uploadConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/upload',
            host: ['{{baseUrl}}'],
            path: ['upload'],
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
          },
          method: 'POST',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
            { key: 'Content-Type', value: 'application/json', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              path: 'rest-api-tests',
              filename: 'test-background.png',
              source: 'data',
              mimeType: 'image/png',
              overwrite: true,
              fileData: testPng
            })
          }
        };

        const uploadResolved = replaceVariables(uploadConfig, testVariables);
        const uploadResponse = await makeRequest(uploadResolved);
        expect(uploadResponse.status).toBe(200);

        // Extract the uploaded file path from the response
        const uploadedPath = uploadResponse.data?.data?.path
          || uploadResponse.data?.path
          || 'rest-api-tests/test-background.png';
        console.log(`  Uploaded test background: ${uploadedPath}`);

        // Step 2: Get the active scene to record its current background
        const sceneConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'active', value: 'true' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const sceneResolved = replaceVariables(sceneConfig, testVariables);
        const sceneResponse = await makeRequest(sceneResolved);
        expect(sceneResponse.status).toBe(200);

        const sceneData = sceneResponse.data?.data || sceneResponse.data;
        const sceneId = sceneData?.id || sceneData?._id;
        expect(sceneId).toBeTruthy();

        // Save original background so we can restore it later
        const originalBg = sceneData?.background?.src || sceneData?.img || '';
        setGlobalVariable(version, 'scene_image_original_bg', originalBg);
        setGlobalVariable(version, 'scene_image_scene_id', sceneId);

        // Step 3: Set the test image as the scene background
        const updateConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
          },
          method: 'PUT',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              sceneId,
              data: {
                background: { src: uploadedPath }
              }
            })
          }
        };

        const updateResolved = replaceVariables(updateConfig, testVariables);
        const updateResponse = await makeRequest(updateResolved);
        expect(updateResponse.status).toBe(200);
        console.log(`  Set background on scene ${sceneId}`);
      }, 30000);
    });

    // ═══════════════════════════════════════════
    // GET /scene/image — rendered canvas screenshot
    // ═══════════════════════════════════════════

    describe(`/scene/image (v${version})`, () => {
      test('GET /scene/image - Capture rendered scene image', async () => {
        setVariable('clientId', getClientId());

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene/image',
            host: ['{{baseUrl}}'],
            path: ['scene', 'image'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'format', value: 'jpeg' },
              { key: 'quality', value: '0.5' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const captured = await captureExample(requestConfig, testVariables, '/scene/image - Capture rendered scene image');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();

        // Make a binary request to get the actual image data + headers
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

        const contentType = binaryResponse.headers['content-type'] || '';
        console.log(`  Response Content-Type: ${contentType}`);

        const isImageResponse = contentType.includes('image/') || contentType.includes('application/json');
        expect(isImageResponse).toBe(true);

        // Log image dimensions from response headers
        const imgWidth = binaryResponse.headers['x-image-width'];
        const imgHeight = binaryResponse.headers['x-image-height'];
        console.log(`  Image dimensions: ${imgWidth}x${imgHeight}, size: ${binaryResponse.data?.byteLength || 0} bytes`);

        // Save the screenshot image to disk
        if (binaryResponse.status === 200 && binaryResponse.data && contentType.includes('image/')) {
          const ext = contentType.includes('jpeg') ? 'jpg' : 'png';
          saveImage(Buffer.from(binaryResponse.data), `scene-screenshot-v${version}.${ext}`);
        } else if (binaryResponse.status === 200 && binaryResponse.data) {
          // Response might be JSON with base64 imageData (fallback when binary decode fails)
          try {
            const jsonStr = Buffer.from(binaryResponse.data).toString('utf8');
            const json = JSON.parse(jsonStr);
            const imageData = json.imageData || json.data?.imageData;
            if (imageData) {
              saveImage(imageData, `scene-screenshot-v${version}.png`);
            } else {
              console.log(`  Screenshot response is JSON without imageData: ${jsonStr.substring(0, 200)}`);
            }
          } catch {
            console.log(`  Could not parse screenshot response as JSON`);
          }
        }
      }, 30000);
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
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const captured = await captureExample(requestConfig, testVariables, '/scene/image/raw - Get raw scene background');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();

        // Save the raw background image
        const data = captured.response.data?.data || captured.response.data;
        if (data?.imageData) {
          saveImage(data.imageData, `scene-raw-bg-v${version}.png`);
        }
      }, 30000);
    });

    // ═══════════════════════════════════════════
    // Teardown: restore original background
    // ═══════════════════════════════════════════

    describe(`Scene image teardown (v${version})`, () => {
      test('Restore original scene background', async () => {
        setVariable('clientId', getClientId());

        const sceneId = getGlobalVariable(version, 'scene_image_scene_id');
        const originalBg = getGlobalVariable(version, 'scene_image_original_bg') || '';

        if (!sceneId) {
          console.log('  No scene to restore');
          return;
        }

        const restoreConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/scene',
            host: ['{{baseUrl}}'],
            path: ['scene'],
            query: [
              { key: 'clientId', value: '{{clientId}}' }
            ]
          },
          method: 'PUT',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ],
          body: {
            mode: 'raw',
            raw: JSON.stringify({
              sceneId,
              data: {
                background: { src: originalBg }
              }
            })
          }
        };

        const resolved = replaceVariables(restoreConfig, testVariables);
        const response = await makeRequest(resolved);
        expect(response.status).toBe(200);
        console.log(`  Restored original background: "${originalBg || '(none)'}"`);
      }, 15000);
    });

  });
});
