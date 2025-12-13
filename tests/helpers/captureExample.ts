/**
 * Test Response Capturer
 * This script runs tests and captures actual API responses for documentation
 */

import { makeRequest, replaceVariables, ApiRequestConfig } from './apiRequest';
import * as fs from 'fs';
import * as path from 'path';

interface CapturedExample {
  endpoint: string;
  method: string;
  description: string;
  request: {
    url: string;
    method: string;
    headers: Record<string, string>;
    body?: any;
  };
  response: {
    status: number;
    data: any;
  };
  codeExamples: {
    javascript: string;
    curl: string;
    python: string;
    typescript: string;
    emojicode: string;
  };
}

interface CaptureOptions {
  /** Custom code examples to use instead of auto-generated ones */
  customCodeExamples?: Partial<CapturedExample['codeExamples']>;
}

/**
 * Capture an API request and response for documentation
 */
export async function captureExample(
  config: ApiRequestConfig,
  variables: Record<string, string>,
  description: string,
  options?: CaptureOptions
): Promise<CapturedExample> {
  const resolvedConfig = replaceVariables(config, variables);
  const response = await makeRequest(resolvedConfig);

  // Extract request details
  const request = {
    url: resolvedConfig.url.raw,
    method: resolvedConfig.method,
    headers: resolvedConfig.header?.reduce((acc, h) => {
      acc[h.key] = h.value;
      return acc;
    }, {} as Record<string, string>) || {},
    body: resolvedConfig.body?.raw ? (() => {
      try {
        return JSON.parse(resolvedConfig.body.raw);
      } catch {
        return resolvedConfig.body.raw;
      }
    })() : undefined
  };

  // Generate code examples (raw values - sanitization happens in updateDocsWithExamples)
  const codeExamples = {
    javascript: options?.customCodeExamples?.javascript || generateJavaScript(resolvedConfig),
    curl: options?.customCodeExamples?.curl || generateCurl(resolvedConfig),
    python: options?.customCodeExamples?.python || generatePython(resolvedConfig),
    typescript: options?.customCodeExamples?.typescript || generateTypeScript(resolvedConfig),
    emojicode: options?.customCodeExamples?.emojicode || generateEmojicode(resolvedConfig)
  };

  const endpoint = extractEndpoint(config.url.path || config.url.raw);

  return {
    endpoint,
    method: config.method,
    description,
    request,
    response: {
      status: response.status,
      data: response.data
    },
    codeExamples
  };
}

function extractEndpoint(pathInfo: any): string {
  if (typeof pathInfo === 'string') {
    const match = pathInfo.match(/\/[^?]*/);
    return match ? match[0] : pathInfo;
  }
  if (Array.isArray(pathInfo)) {
    return '/' + pathInfo.join('/');
  }
  return '';
}

function generateJavaScript(config: any): string {
  const method = config.method.toLowerCase();
  
  // Extract base URL and path
  const urlMatch = config.url.raw.match(/^(https?:\/\/[^\/]+)(.*)$/);
  const baseUrl = urlMatch ? urlMatch[1] : config.url.raw;
  const path = urlMatch ? urlMatch[2].split('?')[0] : '';
  
  // Build URL construction code
  let urlCode = `const baseUrl = '${baseUrl}';\nconst path = '${path}';`;
  
  // Add query params if present
  if (config.url.query && config.url.query.length > 0) {
    const queryParamsCode = config.url.query
      .filter((q: any) => !q.disabled)
      .map((q: any) => `  ${q.key}: '${q.value}'`)
      .join(',\n');
    urlCode += `\nconst params = {\n${queryParamsCode}\n};\nconst queryString = new URLSearchParams(params).toString();\nconst url = \`\${baseUrl}\${path}?\${queryString}\`;`;
  } else {
    urlCode += `\nconst url = \`\${baseUrl}\${path}\`;`;
  }
  
  // Build headers, always include Content-Type if there's a body
  const headers: string[] = [];
  if (config.header && config.header.length > 0) {
    headers.push(...config.header.map((h: any) => `    '${h.key}': '${h.value}'`));
  }
  
  let bodyStr = '';
  if (config.body?.raw) {
    // Add Content-Type if not already present
    const hasContentType = config.header?.some((h: any) => h.key.toLowerCase() === 'content-type');
    if (!hasContentType) {
      headers.push(`    'Content-Type': 'application/json'`);
    }
    
    try {
      const parsed = JSON.parse(config.body.raw);
      bodyStr = `,\n  body: JSON.stringify(${JSON.stringify(parsed, null, 2).replace(/\n/g, '\n    ')})`;
    } catch {
      bodyStr = `,\n  body: ${JSON.stringify(config.body.raw)}`;
    }
  }

  const headersStr = headers.length > 0 ? `,\n  headers: {\n${headers.join(',\n')}\n  }` : '';

  return `${urlCode}\n\nconst response = await fetch(url, {\n  method: '${method.toUpperCase()}'${headersStr}${bodyStr}\n});\nconst data = await response.json();\nconsole.log(data);`;
}

function generateCurl(config: any): string {
  let url = config.url.raw;
  
  // Build query string if not already in URL
  if (!url.includes('?') && config.url.query && config.url.query.length > 0) {
    const queryParams = config.url.query
      .filter((q: any) => !q.disabled)
      .map((q: any) => `${q.key}=${encodeURIComponent(q.value)}`)
      .join('&');
    url = `${url}?${queryParams}`;
  }
  
  const method = config.method;
  
  let cmd = `curl -X ${method} '${url}'`;
  
  if (config.header && config.header.length > 0) {
    for (const header of config.header) {
      cmd += ` \\\n  -H "${header.key}: ${header.value}"`;
    }
  }

  if (config.body?.raw) {
    // Add Content-Type if not already present
    const hasContentType = config.header?.some((h: any) => h.key.toLowerCase() === 'content-type');
    if (!hasContentType) {
      cmd += ` \\\n  -H "Content-Type: application/json"`;
    }
    
    try {
      const parsed = JSON.parse(config.body.raw);
      cmd += ` \\\n  -d '${JSON.stringify(parsed)}'`;
    } catch {
      cmd += ` \\\n  -d '${config.body.raw}'`;
    }
  }

  return cmd;
}

function generatePython(config: any): string {
  const method = config.method.toLowerCase();
  
  // Extract base URL and path
  const urlMatch = config.url.raw.match(/^(https?:\/\/[^\/]+)(.*)$/);
  const baseUrl = urlMatch ? urlMatch[1] : config.url.raw;
  const path = urlMatch ? urlMatch[2].split('?')[0] : '';
  
  // Build URL construction code
  let urlCode = `base_url = '${baseUrl}'\npath = '${path}'`;
  
  // Add query params if present
  let paramsStr = '';
  if (config.url.query && config.url.query.length > 0) {
    const queryParamsCode = config.url.query
      .filter((q: any) => !q.disabled)
      .map((q: any) => `    '${q.key}': '${q.value}'`)
      .join(',\n');
    urlCode += `\nparams = {\n${queryParamsCode}\n}`;
    paramsStr = ',\n    params=params';
  }
  
  urlCode += `\nurl = f'{base_url}{path}'`;
  
  let headersStr = '';
  if (config.header && config.header.length > 0) {
    const headers = config.header
      .map((h: any) => `        '${h.key}': '${h.value}'`)
      .join(',\n');
    headersStr = `,\n    headers={\n${headers}\n    }`;
  }

  let bodyStr = '';
  if (config.body?.raw) {
    try {
      const parsed = JSON.parse(config.body.raw);
      // Convert JSON to Python syntax (true -> True, false -> False, null -> None)
      const pythonJson = JSON.stringify(parsed, null, 2)
        .replace(/\btrue\b/g, 'True')
        .replace(/\bfalse\b/g, 'False')
        .replace(/\bnull\b/g, 'None');
      // Indent the JSON body to align with the parameter context (8 spaces base + 2 for content)
      const indentedJson = pythonJson.split('\n').map((line, i) => {
        if (i === 0) return line; // First line stays inline with json=
        return '    ' + line; // Indent subsequent lines to align with opening brace
      }).join('\n');
      bodyStr = `,\n    json=${indentedJson}`;
    } catch {
      bodyStr = `,\n    data='${config.body.raw}'`;
    }
  }

  return `import requests\n\n${urlCode}\n\nresponse = requests.${method}(\n    url${paramsStr}${headersStr}${bodyStr}\n)\ndata = response.json()\nprint(data)`
}

function generateTypeScript(config: any): string {
  const method = config.method.toLowerCase();
  
  // Extract base URL and path
  const urlMatch = config.url.raw.match(/^(https?:\/\/[^\/]+)(.*)$/);
  const baseUrl = urlMatch ? urlMatch[1] : config.url.raw;
  const path = urlMatch ? urlMatch[2].split('?')[0] : '';
  
  // Build URL construction code
  let urlCode = `const baseUrl = '${baseUrl}';\nconst path = '${path}';`;
  
  // Add query params if present
  let paramsStr = '';
  if (config.url.query && config.url.query.length > 0) {
    const queryParamsCode = config.url.query
      .filter((q: any) => !q.disabled)
      .map((q: any) => `  ${q.key}: '${q.value}'`)
      .join(',\n');
    urlCode += `\nconst params = {\n${queryParamsCode}\n};\nconst queryString = new URLSearchParams(params).toString();\nconst url = \`\${baseUrl}\${path}?\${queryString}\`;`;
  } else {
    urlCode += `\nconst url = \`\${baseUrl}\${path}\`;`;
  }
  
  // Build headers, always include Content-Type if there's a body
  const headers: string[] = [];
  if (config.header && config.header.length > 0) {
    headers.push(...config.header.map((h: any) => `    '${h.key}': '${h.value}'`));
  }
  
  let bodyStr = '';
  if (config.body?.raw) {
    // Add Content-Type if not already present
    const hasContentType = config.header?.some((h: any) => h.key.toLowerCase() === 'content-type');
    if (!hasContentType) {
      headers.push(`    'Content-Type': 'application/json'`);
    }
    
    try {
      const parsed = JSON.parse(config.body.raw);
      bodyStr = `,\n  data: ${JSON.stringify(parsed, null, 2).replace(/\n/g, '\n    ')}`;
    } catch {
      bodyStr = `,\n  data: ${JSON.stringify(config.body.raw)}`;
    }
  }

  const headersStr = headers.length > 0 ? `,\n    headers: {\n  ${headers.join(',\n  ')}\n    }` : '';

  return `import axios from 'axios';

(async () => {
  ${urlCode.replace(/\n/g, '\n  ')}

  const response = await axios({
    method: '${method}'${headersStr},
    url${bodyStr.replace(/\n/g, '\n  ')}
  });
  const data = response.data;
  console.log(data);
})();`;
}

/**
 * Generate Emojicode example using raw sockets for HTTP requests
 * Emojicode is an esoteric language that uses emoji as syntax
 * Requires: emojicodec compiler, sockets package
 * Compile with: emojicodec example.🍇 -o example
 */
function generateEmojicode(config: any): string {
  const method = config.method.toUpperCase();
  
  // Extract host and port from URL
  const urlMatch = config.url.raw.match(/^https?:\/\/([^:\/]+)(?::(\d+))?(.*)$/);
  const host = urlMatch ? urlMatch[1] : 'localhost';
  const port = urlMatch && urlMatch[2] ? urlMatch[2] : '80';
  const pathOnly = urlMatch ? urlMatch[3].split('?')[0] : '/';
  
  // Build query params section if present
  let queryParamsCode = '';
  let queryStringBuild = '';
  const queryParams = config.url.query?.filter((q: any) => !q.disabled) || [];
  
  if (queryParams.length > 0) {
    // Build readable param declarations
    queryParamsCode = '\n  💭 Query parameters\n';
    const paramParts: string[] = [];
    for (const q of queryParams) {
      const varName = q.key.replace(/[^a-zA-Z0-9]/g, '_');
      queryParamsCode += `  🔤${q.key}=${q.value}🔤 ➡️ ${varName}\n`;
      paramParts.push(`🧲${varName}🧲`);
    }
    // Build the query string by concatenating
    queryParamsCode += `  🔤?${paramParts.join('&')}🔤 ➡️ queryString\n`;
    queryStringBuild = '🧲queryString🧲';
  }
  
  // Build headers
  const headerLines: string[] = [];
  if (config.header && config.header.length > 0) {
    for (const h of config.header) {
      headerLines.push(`${h.key}: ${h.value}`);
    }
  }
  
  // Build body and calculate Content-Length
  let body = '';
  let bodyCode = '';
  if (config.body?.raw) {
    try {
      const parsed = JSON.parse(config.body.raw);
      body = JSON.stringify(parsed);
    } catch {
      body = config.body.raw;
    }
    
    // Add Content-Type if not present
    const hasContentType = config.header?.some((h: any) => h.key.toLowerCase() === 'content-type');
    if (!hasContentType) {
      headerLines.push('Content-Type: application/json');
    }
    headerLines.push(`Content-Length: ${Buffer.byteLength(body, 'utf8')}`);
    
    bodyCode = `\n  💭 Request body\n  🔤${body}🔤 ➡️ body\n`;
  }
  
  // Build the HTTP request parts
  const crlf = '❌r❌n';
  let headersStr = '';
  for (const header of headerLines) {
    headersStr += `${header}${crlf}`;
  }
  
  // Construct the request - if we have query params, use string interpolation
  let requestConstruction: string;
  if (queryParams.length > 0) {
    requestConstruction = `  💭 Build HTTP request\n  🔤${method} ${pathOnly}🧲queryString🧲 HTTP/1.1${crlf}Host: ${host}:${port}${crlf}${headersStr}${crlf}${body ? '🧲body🧲' : ''}🔤 ➡️ request`;
  } else {
    requestConstruction = `  💭 Build HTTP request\n  🔤${method} ${pathOnly} HTTP/1.1${crlf}Host: ${host}:${port}${crlf}${headersStr}${crlf}${body}🔤 ➡️ request`;
  }
  
  return `📦 sockets 🏠

💭 Emojicode HTTP Client
💭 Compile: emojicodec example.🍇 -o example
💭 Run: ./example

🏁 🍇
  💭 Connection settings
  🔤${host}🔤 ➡️ host
  ${port} ➡️ port
  🔤${pathOnly}🔤 ➡️ path
${queryParamsCode}${bodyCode}
${requestConstruction}

  💭 Connect and send
  🍺 🆕📞 host port❗ ➡️ socket
  🍺 💬 socket 📇 request❗❗
  
  💭 Read and print response
  🍺 👂 socket 4096❗ ➡️ data
  😀 🍺 🔡 data❗❗
  
  💭 Close socket
  🚪 socket❗
🍉`;
}

/**
 * Deduplicate examples by endpoint+method, keeping only the first successful one
 */
function deduplicateExamples(examples: CapturedExample[]): CapturedExample[] {
  const seen = new Map<string, CapturedExample>();
  
  for (const example of examples) {
    const key = `${example.method} ${example.endpoint}`;
    
    // Keep first successful example (status 2xx)
    if (!seen.has(key)) {
      if (example.response.status >= 200 && example.response.status < 300) {
        seen.set(key, example);
      }
    }
  }
  
  // If we didn't get a successful one, fall back to first example
  for (const example of examples) {
    const key = `${example.method} ${example.endpoint}`;
    if (!seen.has(key)) {
      seen.set(key, example);
    }
  }
  
  return Array.from(seen.values());
}

/**
 * Save captured examples to a file (deduplicates by endpoint)
 */
export function saveExamples(examples: CapturedExample[], outputPath: string): void {
  const dir = path.dirname(outputPath);
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
  const deduplicated = deduplicateExamples(examples);
  fs.writeFileSync(outputPath, JSON.stringify(deduplicated, null, 2));
}

/**
 * Append captured examples to an existing file (or create if doesn't exist)
 * Deduplicates by endpoint, preferring existing examples
 */
export function appendExamples(examples: CapturedExample[], outputPath: string): void {
  const existing = loadExamples(outputPath);
  const combined = [...existing, ...examples];
  saveExamples(combined, outputPath);
}

/**
 * Load captured examples from a file
 */
export function loadExamples(inputPath: string): CapturedExample[] {
  if (!fs.existsSync(inputPath)) {
    return [];
  }
  try {
    const content = fs.readFileSync(inputPath, 'utf8');
    return JSON.parse(content);
  } catch (error) {
    console.error(`Failed to load examples from ${inputPath}:`, error);
    return [];
  }
}
