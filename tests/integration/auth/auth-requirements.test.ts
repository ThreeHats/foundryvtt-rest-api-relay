/**
 * @file auth-requirements.test.ts
 * @description Authentication Requirements Test Suite
 *
 * Uses the generated api-docs.json to dynamically discover all endpoints and
 * verify each correctly rejects missing or invalid API keys. When new endpoints
 * are added, re-run `cd go-relay && go run ./cmd/docgen` and these tests
 * automatically cover them.
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../../helpers/apiRequest';
import { testVariables, setVariable } from '../../helpers/testVariables';
import { forEachVersion } from '../../helpers/multiVersion';
import {
  loadEndpoints, EndpointInfo, dummyParamValue, resolvePathParam, isAuthRoute
} from '../../helpers/endpointMetadata';

interface EndpointTest {
  path: string;
  method: string;
  hasClientId: boolean;
  config: ApiRequestConfig;
}

function buildEndpointTests(): EndpointTest[] {
  const endpoints = loadEndpoints();
  const tests: EndpointTest[] = [];

  for (const ep of endpoints) {
    // Skip SSE/subscribe endpoints (long-lived connections)
    if (ep.isSSE) continue;

    // Skip auth routes — they have their own auth (session-based) or are public
    if (isAuthRoute(ep)) continue;

    // Resolve path parameters like {documentType}
    const testPath = ep.path.replace(/\{(\w+)\}/g, (_, name) => resolvePathParam(name));

    const queryParams: Array<{ key: string; value: string }> = [];
    const bodyData: Record<string, any> = {};

    // Add clientId if the endpoint uses it
    if (ep.hasClientId) {
      queryParams.push({ key: 'clientId', value: '{{clientId}}' });
    }

    // Satisfy required params with dummy values
    for (const p of ep.requiredParams) {
      if (p.name === 'clientId') continue; // already handled
      const val = dummyParamValue(p.name, p.type);
      if (p.location === 'query' || ep.method === 'GET') {
        queryParams.push({ key: p.name, value: typeof val === 'object' ? JSON.stringify(val) : String(val) });
      } else {
        bodyData[p.name] = val;
      }
    }

    const config: ApiRequestConfig = {
      url: {
        raw: `{{baseUrl}}${testPath}`,
        host: ['{{baseUrl}}'],
        path: [testPath.replace(/^\//, '')],
        query: queryParams,
      },
      method: ep.method as any,
      header: [
        { key: 'x-api-key', value: '{{apiKey}}', type: 'text' },
      ],
    };

    if (['POST', 'PUT', 'DELETE'].includes(ep.method)) {
      config.body = { mode: 'raw', raw: JSON.stringify(bodyData) };
    }

    tests.push({
      path: ep.path,
      method: ep.method,
      hasClientId: ep.hasClientId,
      config,
    });
  }

  return tests;
}

describe('Authentication Requirements', () => {
  const endpointTests = buildEndpointTests();

  forEachVersion((version, getClientId) => {
    describe(`API Key Requirements (v${version})`, () => {
      endpointTests.forEach(endpoint => {
        test(`${endpoint.method} ${endpoint.path} - should reject missing API key`, async () => {
          setVariable('clientId', getClientId());

          const configNoKey = JSON.parse(JSON.stringify(endpoint.config));
          configNoKey.header = configNoKey.header?.filter((h: any) => h.key !== 'x-api-key') || [];
          const replaced = replaceVariables(configNoKey, testVariables);

          const response = await makeRequest(replaced);
          expect(response.status).toBe(401);
          expect(response.data).toHaveProperty('error');
        });

        test(`${endpoint.method} ${endpoint.path} - should reject invalid API key`, async () => {
          setVariable('clientId', getClientId());

          const configInvalidKey = JSON.parse(JSON.stringify(endpoint.config));
          const apiKeyHeader = configInvalidKey.header?.find((h: any) => h.key === 'x-api-key');
          if (apiKeyHeader) {
            apiKeyHeader.value = 'invalid-api-key-12345';
          }
          const replaced = replaceVariables(configInvalidKey, testVariables);

          const response = await makeRequest(replaced);
          expect(response.status).toBe(401);
          expect(response.data).toHaveProperty('error');
        });
      });
    });

    describe(`Client ID Requirements (v${version})`, () => {
      let connectedClientCount = 0;

      beforeAll(async () => {
        try {
          const baseUrl = testVariables.baseUrl;
          const apiKey = testVariables.apiKey;
          const resp = await axios.get(`${baseUrl}/clients`, {
            headers: { 'x-api-key': apiKey },
          });
          connectedClientCount = resp.data?.total ?? 0;
        } catch {
          connectedClientCount = 0;
        }
      });

      // Only test clientId requirements for endpoints that actually use clientId
      endpointTests.filter(e => e.hasClientId).forEach(endpoint => {
        test(`${endpoint.method} ${endpoint.path} - missing clientId behavior`, async () => {
          setVariable('clientId', getClientId());

          const configNoClient = JSON.parse(JSON.stringify(endpoint.config));
          if (configNoClient.url.query) {
            configNoClient.url.query = configNoClient.url.query.filter((q: any) => q.key !== 'clientId');
          }
          const replaced = replaceVariables(configNoClient, testVariables);

          const response = await makeRequest(replaced);

          if (connectedClientCount === 1) {
            expect([200, 400, 408]).toContain(response.status);
          } else if (connectedClientCount > 1) {
            expect([400, 404]).toContain(response.status);
          } else {
            expect([400, 404]).toContain(response.status);
          }
        });

        test(`${endpoint.method} ${endpoint.path} - should reject invalid clientId`, async () => {
          setVariable('clientId', getClientId());

          const configInvalidClient = JSON.parse(JSON.stringify(endpoint.config));
          const clientIdParam = configInvalidClient.url.query?.find((q: any) => q.key === 'clientId');
          if (clientIdParam) {
            clientIdParam.value = 'invalid-client-id-12345';
          }
          const replaced = replaceVariables(configInvalidClient, testVariables);

          const response = await makeRequest(replaced);
          expect(response.status).toBe(404);
          expect(response.data).toHaveProperty('error');
        });
      });
    });
  });
});
