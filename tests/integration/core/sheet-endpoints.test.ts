/**
 * @file sheet-endpoints.test.ts
 * @description Actor/Item Sheet Screenshot Endpoint Tests
 * @endpoints GET /sheet
 *
 * Captures actor sheet screenshots as PNG/JPEG images.
 * Works on both Foundry v12 and v13+.
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import axios from 'axios';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { captureExample, saveExamples } from '../../helpers/captureExample';
import { forEachVersion } from '../../helpers/multiVersion';
import { getEntityUuid } from '../../helpers/testEntities';
import * as path from 'path';
import * as fs from 'fs';

// Store captured examples for documentation
const capturedExamples: any[] = [];

// Directory to save received images
const IMAGE_OUTPUT_DIR = path.join(__dirname, '../../test-results/images');

function saveImage(data: Buffer | string, filename: string): void {
  if (!fs.existsSync(IMAGE_OUTPUT_DIR)) {
    fs.mkdirSync(IMAGE_OUTPUT_DIR, { recursive: true });
  }
  const outPath = path.join(IMAGE_OUTPUT_DIR, filename);
  if (Buffer.isBuffer(data)) {
    fs.writeFileSync(outPath, data);
  } else if (typeof data === 'string') {
    let b64 = data;
    const commaIdx = b64.indexOf(',');
    if (commaIdx >= 0) b64 = b64.substring(commaIdx + 1);
    fs.writeFileSync(outPath, Buffer.from(b64, 'base64'));
  }
  console.log(`  Saved image: ${outPath}`);
}

describe('Sheet', () => {
  afterAll(() => {
    if (capturedExamples.length > 0) {
      const outputPath = path.join(__dirname, '../../../docs/examples/sheet-examples.json');
      saveExamples(capturedExamples, outputPath);
      console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
    }
  });

  forEachVersion((version, getClientId) => {
    describe(`/sheet (v${version})`, () => {
      test('GET /sheet - Capture actor sheet as PNG', async () => {
        setVariable('clientId', getClientId());

        const actorUuid = getEntityUuid(version, 'Actor', 'primary');
        expect(actorUuid).toBeTruthy();

        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/sheet',
            host: ['{{baseUrl}}'],
            path: ['sheet'],
            query: [
              { key: 'clientId', value: '{{clientId}}' },
              { key: 'uuid', value: actorUuid! },
              { key: 'format', value: 'png' }
            ]
          },
          method: 'GET',
          header: [
            { key: 'x-api-key', value: '{{apiKey}}', type: 'text' }
          ]
        };

        const captured = await captureExample(requestConfig, testVariables, '/sheet - Actor sheet screenshot');
        capturedExamples.push(captured);

        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toBeTruthy();

        // Make a binary request to save the image
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

        if (binaryResponse.status === 200 && contentType.includes('image/')) {
          const ext = contentType.includes('jpeg') ? 'jpg' : 'png';
          saveImage(Buffer.from(binaryResponse.data), `sheet-actor-v${version}.${ext}`);
        }
      }, 30000);
    });
  });
});
