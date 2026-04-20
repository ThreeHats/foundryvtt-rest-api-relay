/**
 * Endpoint Metadata Helper
 *
 * Reads the generated api-docs.json to discover all API endpoints and their
 * parameters. This is the single source of truth — when new endpoints are
 * added, just re-run `go run ./cmd/docgen` and the tests auto-discover them.
 *
 * Used by auth-requirements.test.ts and permission-filtering.test.ts to
 * dynamically test all endpoints without hand-maintained endpoint lists.
 */

import * as fs from 'fs';
import * as path from 'path';

export interface EndpointParam {
  name: string;
  type: string;       // "string", "number", "boolean", "array", "object"
  location: string;   // "query", "body", "query/body", "params"
  description: string;
  required: boolean;
}

export interface EndpointInfo {
  path: string;
  method: string;
  description: string;
  tag: string;
  requiredParams: EndpointParam[];
  optionalParams: EndpointParam[];
  hasClientId: boolean;
  hasUserId: boolean;
  isSSE: boolean;
}

/**
 * Auth routes are mounted separately from the API auth middleware.
 * They handle their own auth (session-based) or are public endpoints.
 * They should NOT be tested for x-api-key enforcement.
 */
const AUTH_ROUTE_PREFIX = '/auth/';

/**
 * Endpoints where userId is accepted but invalid userId does NOT produce
 * a "User not found" error (e.g., the endpoint ignores it or the module
 * handles it differently).
 */
const SKIP_INVALID_USERID_ENDPOINTS = [
  '/effects/list',   // Returns system effects regardless of userId
  '/upload',         // Custom handler validates file data before userId reaches module
];

/**
 * Extra params needed by endpoints with custom ValidateParams that
 * api-docs.json lists as optional but the handler requires at least one.
 * Keys are "METHOD /path".
 */
const CUSTOM_REQUIRED_PARAMS: Record<string, Record<string, any>> = {
  'GET /scene': { all: 'true' },
  'PUT /scene': { sceneId: 'fakeScene123' },
  'DELETE /scene': { name: 'fakeScene' },
  'POST /switch-scene': { name: 'fakeScene' },
  'POST /select': { uuids: ['Actor.fakeId12345678'] },
  'POST /execute-js': { script: 'return 1' },
  'GET /user': { id: 'fakeUser123' },
  'PUT /user': { id: 'fakeUser123' },
  'DELETE /user': { id: 'fakeUser123' },
  'GET /get-folder': { name: 'test' },
  'POST /create-folder': { name: 'test', folderType: 'Actor' },
  'DELETE /delete-folder': { folderId: 'fakeFolderId' },
  'POST /upload': { path: 'test', filename: 'test.txt' },
  'GET /sheet': { uuid: 'Actor.fakeId12345678' },
  'GET /search': { query: 'test' },
};

let cachedEndpoints: EndpointInfo[] | null = null;

/**
 * Load all endpoint metadata from api-docs.json.
 * Results are cached for the process lifetime.
 */
export function loadEndpoints(): EndpointInfo[] {
  if (cachedEndpoints) return cachedEndpoints;

  const docsPath = path.join(__dirname, '../../public/api-docs.json');
  if (!fs.existsSync(docsPath)) {
    console.warn(`  ⚠ api-docs.json not found at ${docsPath}. Run: cd go-relay && go run ./cmd/docgen`);
    return [];
  }

  const raw = JSON.parse(fs.readFileSync(docsPath, 'utf8'));
  const endpoints: EndpointInfo[] = [];

  for (const ep of raw.endpoints || []) {
    const requiredParams: EndpointParam[] = (ep.requiredParameters || []).map((p: any) => ({
      name: p.name,
      type: p.type || 'string',
      location: p.location || 'body',
      description: p.description || '',
      required: true,
    }));

    const optionalParams: EndpointParam[] = (ep.optionalParameters || []).map((p: any) => ({
      name: p.name,
      type: p.type || 'string',
      location: p.location || 'query',
      description: p.description || '',
      required: false,
    }));

    const allParams = [...requiredParams, ...optionalParams];
    const hasClientId = allParams.some(p => p.name === 'clientId');
    const hasUserId = allParams.some(p => p.name === 'userId');
    const isSSE = (ep.path || '').includes('/subscribe');

    endpoints.push({
      path: ep.path,
      method: (ep.method || 'GET').toUpperCase(),
      description: ep.description || '',
      tag: ep.tag || '',
      requiredParams,
      optionalParams,
      hasClientId,
      hasUserId,
      isSSE,
    });
  }

  cachedEndpoints = endpoints;
  console.log(`  Loaded ${endpoints.length} endpoints from api-docs.json`);
  return endpoints;
}

/**
 * Provide a sensible dummy value for a parameter based on its name and type.
 * Used to satisfy required param validation so userId-related tests can reach
 * the module layer.
 */
export function dummyParamValue(name: string, type: string): any {
  // Named overrides
  switch (name) {
    case 'clientId': return '{{clientId}}';
    case 'uuid': return 'Actor.fakeId12345678';
    case 'actorUuid': return 'Actor.fakeId12345678';
    case 'sourceActorUuid': return 'Actor.fakeId12345678';
    case 'targetActorUuid': return 'Actor.fakeId12345679';
    case 'itemUuid': return 'Item.fakeId12345678';
    case 'query': return 'test';
    case 'formula': return '1d20';
    case 'entityType': return 'Actor';
    case 'content': return 'test';
    case 'name': return 'test';
    case 'folderId': return 'fakeFolder123';
    case 'folderType': return 'Actor';
    case 'attribute': return 'name';
    case 'sceneId': return 'fakeScene123';
    case 'script': return 'return 1';
    case 'skill': return 'acr';
    case 'ability': return 'str';
    case 'damage': return 10;
    case 'amount': return 1;
    case 'data': return { name: 'test' };
    case 'details': return ['resources'];
    case 'currency': return { gp: 1 };
    case 'statusId': return 'poisoned';
    case 'documentId': return 'fakeDoc123';
    case 'hooks': return 'updateActor';
  }

  // Type-based fallbacks
  switch (type) {
    case 'number': return 1;
    case 'boolean': return true;
    case 'array': return ['test'];
    case 'object': return { name: 'test' };
    default: return 'test-value';
  }
}

/**
 * Whether this endpoint is an auth route (has its own auth, not API key based).
 */
export function isAuthRoute(ep: EndpointInfo): boolean {
  return ep.path.startsWith(AUTH_ROUTE_PREFIX);
}

/**
 * Whether this endpoint should be skipped for invalid userId testing.
 */
export function shouldSkipInvalidUserId(ep: EndpointInfo): boolean {
  return SKIP_INVALID_USERID_ENDPOINTS.includes(ep.path);
}

/**
 * Get extra params needed by endpoints with custom validation.
 * Returns a map of param name → value to merge into the request.
 */
export function getCustomRequiredParams(method: string, path: string): Record<string, any> {
  return CUSTOM_REQUIRED_PARAMS[`${method} ${path}`] || {};
}

/**
 * Resolve a Chi-style path parameter like {documentType} to a test value.
 */
export function resolvePathParam(name: string): string {
  switch (name) {
    case 'documentType': return 'Token';
    case 'messageId': return 'fakeMsg123';
    case 'uuid': return 'Actor.fakeId12345678';
    default: return 'dummy';
  }
}
