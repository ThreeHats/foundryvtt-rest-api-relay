/**
 * @file structure-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings
 * @description World Structure Endpoint Tests
 * @endpoints POST /create-folder, GET /structure, GET /contents/:path, GET /get-folder, DELETE /delete-folder
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

describe('Structure', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/structure-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(`\nSaved ${capturedExamples.length} examples to ${outputPath}`);
  });

  forEachVersion((version, getClientId) => {
    describe(`/create-folder (v${version})`, () => {
      test('POST /create-folder', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/create-folder',
            host: ['{{baseUrl}}'],
            path: ['create-folder'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'name',
                value: 'test-folder',
              },
              {
                key: 'folderType',
                value: 'Scene',
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
          ]
      };

        // Capture this example for documentation (also makes the request)
        const captured = await captureExample(
          requestConfig,
          testVariables,
          '/create-folder'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('id');
        expect(captured.response.data.data).toHaveProperty('uuid');
        expect(captured.response.data.data).toHaveProperty('name');
        expect(captured.response.data.data).toHaveProperty('type');
        expect(captured.response.data.data).toHaveProperty('parentFolder');

        // Set global variable for use in other tests
        setGlobalVariable(version, 'folderName', captured.response.data.data.name);
        setGlobalVariable(version, 'folderType', captured.response.data.data.type);
        setGlobalVariable(version, 'folderId', captured.response.data.data.id);
        setGlobalVariable(version, 'folderUuid', captured.response.data.data.uuid);
      });
    });

    describe(`/structure (v${version})`, () => {
      test('GET /structure', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/structure',
            host: ['{{baseUrl}}'],
            path: ['structure'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'includeEntityData',
                value: 'true',
              },
              {
                key: 'recursive',
                value: 'true',
              },
              {
                key: 'types',
                value: 'Scene',
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
          '/structure'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('folders');
        expect(captured.response.data.data.folders).toBeInstanceOf(Object);
        expect(captured.response.data.data.folders).toHaveProperty(getGlobalVariable(version, 'folderName'));
        expect(captured.response.data.data.folders[getGlobalVariable(version, 'folderName')]).toHaveProperty('id');
        expect(captured.response.data.data.folders[getGlobalVariable(version, 'folderName')]).toHaveProperty('uuid');
        expect(captured.response.data.data.folders[getGlobalVariable(version, 'folderName')]).toHaveProperty('type');
        expect(captured.response.data.data).toHaveProperty('entities');
      });
    });

    describe(`/contents (v${version})`, () => {
      test('GET /contents/:path', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/contents/:path',
            host: ['{{baseUrl}}'],
            path: ['contents'],
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
          '/contents/:path'
        );
        // Deprecated
        // capturedExamples.push(captured);

        // Assertions
        // Deprecated. Should return 400 with error message and example
        expect(captured.response.status).toBe(400);
        expect(captured.response.data.error).toBeTruthy();
        expect(captured.response.data.message).toBeTruthy();
        expect(captured.response.data.example).toBeTruthy();
      });
    });

    describe(`/get-folder (v${version})`, () => {
      test('GET /get-folder', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/get-folder',
            host: ['{{baseUrl}}'],
            path: ['get-folder'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'name',
                value: getGlobalVariable(version, 'folderName'),
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
          '/get-folder'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('id');
        expect(captured.response.data.data).toHaveProperty('uuid');
        expect(captured.response.data.data).toHaveProperty('name');
        expect(captured.response.data.data).toHaveProperty('type');
        expect(captured.response.data.data).toHaveProperty('parentFolder');
        expect(captured.response.data.data).toHaveProperty('contents');
        expect(captured.response.data.data.contents).toBeInstanceOf(Array);
      });
    });

    describe(`/delete-folder (v${version})`, () => {
      test('DELETE /delete-folder', async () => {
        // Set clientId for this version
        setVariable('clientId', getClientId());
        
        // Request configuration
        const requestConfig: ApiRequestConfig = {
          url: {
            raw: '{{baseUrl}}/delete-folder',
            host: ['{{baseUrl}}'],
            path: ['delete-folder'],
            query: [
              {
                key: 'clientId',
                value: '{{clientId}}',
              },
              {
                key: 'folderId',
                value: getGlobalVariable(version, 'folderId'),
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
          '/delete-folder'
        );
        capturedExamples.push(captured);

        // Assertions
        expect(captured.response.status).toBe(200);
        expect(captured.response.data).toHaveProperty('data');
        expect(captured.response.data.data).toHaveProperty('deleted');
        expect(captured.response.data.data.deleted).toBe(true);
        expect(captured.response.data.data).toHaveProperty('folderId');
        expect(captured.response.data.data.folderId).toBe(getGlobalVariable(version, 'folderId'));
        expect(captured.response.data.data).toHaveProperty('entitiesDeleted');
        expect(captured.response.data.data).toHaveProperty('foldersDeleted');
      });
    });

  });

});
