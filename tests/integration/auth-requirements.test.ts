/**
 * @file auth-requirements.test.ts
 * @generated Partially auto-generated - scans route files dynamically
 * @description Authentication Requirements Test Suite
 * @endpoints Dynamically tests all endpoints for auth enforcement
 */

import { describe, test, expect } from '@jest/globals';
import { ApiRequestConfig, makeRequest, replaceVariables } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { forEachVersion } from '../helpers/multiVersion';
import * as fs from 'fs';
import * as path from 'path';

interface EndpointTest {
  path: string;
  method: string;
  hasClientId: boolean;
  hasApiKey: boolean;
  config: ApiRequestConfig;
}

/**
 * Extract endpoint information from route files
 */
function extractEndpoints(): EndpointTest[] {
  const endpoints: EndpointTest[] = [];
  const apiDir = path.join(__dirname, '../../src/routes/api');
  const routeFiles = fs.readdirSync(apiDir).filter(file => file.endsWith('.ts'));

  // Read api.ts to find router mounting paths
  const apiTsPath = path.join(__dirname, '../../src/routes/api.ts');
  const apiTsContent = fs.readFileSync(apiTsPath, 'utf8');
  const routerMounts = new Map<string, string>();
  
  // Match patterns like: app.use('/dnd5e', dnd5eRouter);
  const mountRegex = /app\.use\(['"]([^'"]+)['"],\s*(\w+Router)\)/g;
  let mountMatch;
  while ((mountMatch = mountRegex.exec(apiTsContent)) !== null) {
    const [, mountPath, routerName] = mountMatch;
    // Convert router name to file name (e.g., dnd5eRouter -> dnd5e.ts)
    const fileName = routerName.replace('Router', '') + '.ts';
    // Only set prefix if it's not '/' (root mount has no prefix)
    if (mountPath !== '/') {
      routerMounts.set(fileName, mountPath);
    }
  }

  routeFiles.forEach(file => {
    const content = fs.readFileSync(path.join(apiDir, file), 'utf8');
    const lines = content.split('\n');

    // Get router prefix from the mounting map (empty string if mounted at root)
    const routerPrefix = routerMounts.get(file) || '';

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];

      // Find router method calls
      if ((line.includes('Router') || line.includes('router')) &&
          (line.includes('.get(') || line.includes('.post(') ||
           line.includes('.put(') || line.includes('.delete('))) {

        const methodMatch = line.match(/(\w+[Rr]outer)\.(get|post|put|delete)\(\s*["']([^"']+)["']/);
        if (!methodMatch) continue;

        const [, , method, routePath] = methodMatch;
        const fullPath = routerPrefix + routePath;

        // Check if this is a deprecated endpoint by looking at nearby lines
        let isDeprecated = false;
        for (let k = Math.max(0, i - 10); k < Math.min(i + 5, lines.length); k++) {
          if (lines[k].includes('deprecated')) {
            isDeprecated = true;
            break;
          }
        }
        if (isDeprecated) continue; // Skip deprecated endpoints

        // Look for createApiRoute configuration
        let hasClientId = false;
        let hasApiKey = true; // Assume all routes have API key requirement
        const queryParams: Array<{ key: string; value: string }> = [];
        const bodyParams: any = {};

        // Scan ahead for requiredParams and optionalParams
        for (let j = i; j < Math.min(i + 100, lines.length); j++) {
          const scanLine = lines[j];

          // Look for clientId in params
          if (scanLine.includes("name: 'clientId'") || scanLine.includes('name: "clientId"')) {
            hasClientId = true;
            queryParams.push({ key: 'clientId', value: '{{clientId}}' });
          }

          // Look for other common query params
          if (scanLine.includes("from: 'query'") || scanLine.includes('from: "query"')) {
            // Try to extract param name from previous lines
            for (let k = j - 5; k < j; k++) {
              if (k >= 0) {
                const paramMatch = lines[k].match(/name:\s*['"]([^'"]+)['"]/);
                if (paramMatch && paramMatch[1] !== 'clientId') {
                  // Skip adding other params for now, we just need minimal config
                }
              }
            }
          }

          // Stop scanning after we find the closing of createApiRoute
          if (scanLine.includes('}));') && j > i + 10) {
            break;
          }
        }

        // Replace path parameters with dummy values for testing
        const testPath = fullPath.replace(/:(\w+)/g, (match, paramName) => {
          // Provide sensible dummy values based on parameter name
          if (paramName === 'uuid') return 'Actor.abc123';
          if (paramName === 'path') return 'test-file.json';
          if (paramName === 'id') return 'test123';
          return 'dummy-value';
        });

        // Build a minimal request config for testing
        const config: ApiRequestConfig = {
          url: {
            raw: `{{baseUrl}}${testPath}`,
            host: ['{{baseUrl}}'],
            path: [testPath.replace(/^\//, '')],
            query: queryParams
          },
          method: method.toUpperCase() as any,
          header: [
            {
              key: 'x-api-key',
              value: '{{apiKey}}',
              type: 'text'
            }
          ]
        };

        // Add minimal body for POST/PUT/DELETE with required fields
        if (['POST', 'PUT', 'DELETE'].includes(method.toUpperCase())) {
          // Provide minimal valid body to pass validation
          const bodyData: any = {};
          
          // Scan for common required body params to avoid validation errors
          for (let k = i; k < Math.min(i + 100, lines.length); k++) {
            const scanLine = lines[k];
            if (scanLine.includes("from: 'body'") || scanLine.includes('from: "body"')) {
              // Try to find param name in surrounding lines
              for (let m = k - 5; m < k; m++) {
                if (m >= 0) {
                  const nameMatch = lines[m].match(/name:\s*['"]([^'"]+)['"]/);
                  if (nameMatch) {
                    const paramName = nameMatch[1];
                    // Add dummy values for common params
                    if (paramName === 'uuid') bodyData[paramName] = 'Actor.abc123';
                    else if (paramName === 'data') bodyData[paramName] = {};
                    else if (paramName === 'type') bodyData[paramName] = 'Actor';
                    else if (paramName === 'name') bodyData[paramName] = 'Test';
                    else bodyData[paramName] = 'test-value';
                  }
                }
              }
            }
          }
          
          config.body = {
            mode: 'raw',
            raw: JSON.stringify(bodyData)
          };
        }

        endpoints.push({
          path: fullPath,
          method: method.toUpperCase(),
          hasClientId,
          hasApiKey,
          config
        });
      }
    }
  });

  return endpoints;
}

describe('Authentication Requirements', () => {
  const endpoints = extractEndpoints();

  forEachVersion((version, getClientId) => {
    describe(`API Key Requirements (v${version})`, () => {
      endpoints.filter(e => e.hasApiKey).forEach(endpoint => {
        test(`${endpoint.method} ${endpoint.path} - should reject missing API key`, async () => {
          setVariable('clientId', getClientId());

          // Test without API key
          const configNoKey = JSON.parse(JSON.stringify(endpoint.config));
          configNoKey.header = configNoKey.header?.filter((h: any) => h.key !== 'x-api-key') || [];
          const replaced = replaceVariables(configNoKey, testVariables);

          const response = await makeRequest(replaced);
          expect(response.status).toBe(401);
          expect(response.data).toHaveProperty('error');
        });

        test(`${endpoint.method} ${endpoint.path} - should reject invalid API key`, async () => {
          setVariable('clientId', getClientId());

          // Test with invalid API key
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
      endpoints.filter(e => e.hasClientId).forEach(endpoint => {
        test(`${endpoint.method} ${endpoint.path} - should reject missing clientId`, async () => {
          setVariable('clientId', getClientId());

          // Test without clientId
          const configNoClient = JSON.parse(JSON.stringify(endpoint.config));
          if (configNoClient.url.query) {
            configNoClient.url.query = configNoClient.url.query.filter((q: any) => q.key !== 'clientId');
          }
          const replaced = replaceVariables(configNoClient, testVariables);

          const response = await makeRequest(replaced);
          // Missing clientId should return 400
          expect(response.status).toBe(400);
          expect(response.data).toHaveProperty('error');
        });

        test(`${endpoint.method} ${endpoint.path} - should reject invalid clientId`, async () => {
          setVariable('clientId', getClientId());

          // Test with invalid clientId
          const configInvalidClient = JSON.parse(JSON.stringify(endpoint.config));
          const clientIdParam = configInvalidClient.url.query?.find((q: any) => q.key === 'clientId');
          if (clientIdParam) {
            clientIdParam.value = 'invalid-client-id-12345';
          }
          const replaced = replaceVariables(configInvalidClient, testVariables);

          const response = await makeRequest(replaced);
          // Invalid clientId should return 404
          expect(response.status).toBe(404);
          expect(response.data).toHaveProperty('error');
        });
      });
    });
  });
});
