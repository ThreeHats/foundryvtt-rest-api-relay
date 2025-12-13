/**
 * @file fileSystem-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description File System Endpoint Tests
 * @endpoints POST /upload, GET /file-system, GET /download
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { setGlobalVariable, getGlobalVariable } from '../helpers/globalVariables'
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('FileSystem', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/fileSystem-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/upload (v${version})`, () => {
      test('POST /upload', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());

        // Create a simple text file as base64
        const testContent = 'Hello from REST API test!';
        const base64Content = Buffer.from(testContent).toString('base64');
        const fileData = `data:text/plain;base64,${base64Content}`;
        setGlobalVariable(version, 'testFileData', fileData);
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/upload',
            host: ['{{baseUrl}}'],
            path: ['upload'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'The ID of the Foundry client to connect to'
              },
              {
                key: 'path',
                value: 'rest-api-tests'
              },
              {
                key: 'source',
                value: 'data'
              },
              {
                key: 'filename',
                value: 'test-file.txt'
              },
              {
                key: 'mimeType',
                value: 'text/plain'
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
              fileData: fileData,
              mimeType: 'text/plain',
              overwrite: true
            }, null, 2)
        }
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/upload'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('path', 'rest-api-tests/test-file.txt');
      });
    });

    describe(`/file-system (v${version})`, () => {
      test('GET /file-system', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/file-system',
            host: ['{{baseUrl}}'],
            path: ['file-system'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'The ID of the Foundry client to connect to'
              },
              {
                key: 'source',
                value: 'data',
                description: 'The source directory to use (data, systems, modules, etc.)'
              },
              {
                key: 'recursive',
                value: 'false',
                description: 'Whether to recursively list all subdirectories'
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
          '/file-system'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('path');
        expect(captured.response.data).toHaveProperty('source', 'data');
        expect(captured.response.data).toHaveProperty('results');
        expect(captured.response.data.results).toBeInstanceOf(Array);
        expect(captured.response.data.results.length).toBeGreaterThan(0);
        expect(captured.response.data.results[0]).toHaveProperty('name');
        expect(captured.response.data.results[0]).toHaveProperty('path');
        expect(captured.response.data.results[0]).toHaveProperty('type');
        expect(captured.response.data.results).toContainEqual(
          expect.objectContaining({
            name: 'rest-api-tests',
            path: 'rest-api-tests',
            type: 'directory'
          })
        );
      });
    });

    describe(`/download (v${version})`, () => {
      test('GET /download', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/download',
            host: ['{{baseUrl}}'],
            path: ['download'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
                description: 'The ID of the Foundry client to connect to'
              },
              {
                key: 'path',
                value: 'rest-api-tests/test-file.txt',
                description: 'The full path to the file to download'
              },
              {
                key: 'source',
                value: 'data',
                description: 'The source directory to use (data, systems, modules, etc.)'
              },
              {
                key: 'format',
                value: 'base64',
                description: 'The format to return the file in (binary, base64)'
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
          '/download'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('success', true);
        expect(captured.response.data).toHaveProperty('path', 'rest-api-tests/test-file.txt');
        expect(captured.response.data).toHaveProperty('mimeType');
        expect(captured.response.data.mimeType).toContain('text/plain');
        expect(captured.response.data).toHaveProperty('fileData');
        // Extract and compare only the base64 content, not the full data URL
        const expectedBase64 = getGlobalVariable(version, 'testFileData').split(',')[1];
        const receivedBase64 = captured.response.data.fileData.split(',')[1];
        expect(receivedBase64).toBe(expectedBase64);
      });
    });

  });

});
