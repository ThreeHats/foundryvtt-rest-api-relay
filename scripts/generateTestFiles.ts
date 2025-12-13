/**
 * Generate Jest test files from route file doc strings
 * This script parses route files and generates integration test files for endpoints that don't have tests yet
 */

import * as fs from 'fs';
import * as path from 'path';

interface RouteParam {
  name: string;
  from: string[];
  type: string;
  description: string;
}

interface RouteInfo {
  method: string;
  path: string;
  description: string;
  requiredParams: RouteParam[];
  optionalParams: RouteParam[];
}

/**
 * Extract route information from TypeScript route files
 * This uses similar parsing logic to generateApiDocs.js
 */
function extractRoutes(filePath: string): RouteInfo[] {
  const content = fs.readFileSync(filePath, 'utf8');
  const lines = content.split('\n');
  const routes: RouteInfo[] = [];

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    
    // Find router method calls
    if ((line.includes('Router') || line.includes('router')) && 
        (line.includes('.get(') || line.includes('.post(') || 
         line.includes('.put(') || line.includes('.delete('))) {
      
      const methodMatch = line.match(/(\w+[Rr]outer)\.(get|post|put|delete)\(\s*["']([^"']+)["']/);
      if (!methodMatch) continue;

      const [, , method, routePath] = methodMatch;

      // Look backwards for JSDoc comment
      let j = i - 1;
      let jsdocContent = '';
      let foundJSDocEnd = false;
      let foundJSDocStart = false;

      while (j >= 0 && j >= i - 50) {
        const currentLine = lines[j].trim();
        
        if (foundJSDocEnd && !foundJSDocStart && 
            !currentLine.startsWith('//') && !currentLine.startsWith('/*') && 
            currentLine !== '' && !currentLine.includes('*/')) {
          break;
        }

        if (currentLine.includes('*/')) {
          foundJSDocEnd = true;
          const endPos = j;
          
          while (j >= 0 && j >= i - 50) {
            if (lines[j].includes('/**')) {
              foundJSDocStart = true;
              jsdocContent = lines.slice(j, endPos + 1).join('\n');
              break;
            }
            j--;
          }
          
          if (foundJSDocStart) break;
        }
        
        j--;
      }

      // Parse JSDoc
      const jsdoc = parseJSDoc(jsdocContent);

      // Look for createApiRoute configuration
      let k = i;
      let routeConfigLines: string[] = [];
      let bracketCount = 0;
      let foundCreateApiRoute = false;

      // Check if this uses createApiRoute
      const lookAhead = Math.min(i + 10, lines.length);
      let hasCreateApiRoute = false;
      
      for (let la = i; la < lookAhead; la++) {
        if (lines[la].includes('createApiRoute')) {
          hasCreateApiRoute = true;
          break;
        }
      }

      if (hasCreateApiRoute) {
        while (k < lines.length) {
          const currentLine = lines[k];
          
          if (currentLine.includes('createApiRoute')) {
            foundCreateApiRoute = true;
          }
          
          if (foundCreateApiRoute) {
            routeConfigLines.push(currentLine);
            
            for (const char of currentLine) {
              if (char === '{') bracketCount++;
              if (char === '}') bracketCount--;
            }
            
            if (bracketCount === 0 && currentLine.includes('})')) {
              break;
            }
          }
          
          k++;
        }

        // Parse createApiRoute parameters
        const { requiredParams, optionalParams } = parseCreateApiRouteParams(routeConfigLines);
        
        routes.push({
          method: method.toUpperCase(),
          path: routePath,
          description: jsdoc.description,
          requiredParams,
          optionalParams
        });
      } else {
        // Traditional endpoint - use JSDoc params
        const requiredParams: RouteParam[] = [];
        const optionalParams: RouteParam[] = [];

        jsdoc.params.forEach(param => {
          const paramInfo: RouteParam = {
            name: param.name,
            from: param.from,
            type: param.type,
            description: param.description
          };

          if (param.optional) {
            optionalParams.push(paramInfo);
          } else {
            requiredParams.push(paramInfo);
          }
        });

        routes.push({
          method: method.toUpperCase(),
          path: routePath,
          description: jsdoc.description,
          requiredParams,
          optionalParams
        });
      }
    }
  }

  return routes;
}

/**
 * Parse JSDoc comments
 */
function parseJSDoc(jsdocContent: string) {
  const lines = jsdocContent.split('\n')
    .map(line => line.replace(/^\s*\/\*\*\s*/, ''))
    .map(line => line.replace(/^\s*\*\s?/, ''))
    .map(line => line.replace(/\s*\*\/\s*$/, ''))
    .map(line => line.trim());

  const result = {
    description: '',
    params: [] as Array<{
      name: string;
      type: string;
      description: string;
      optional: boolean;
      from: string[];
    }>
  };

  let descriptionLines: string[] = [];
  
  for (const line of lines) {
    if (line.startsWith('@param')) {
      const paramMatch = line.match(/@param\s+(?:\{([^}]+)\})?\s*(\[?[\w\-\[\]\.]+\]?)(?:\s*-\s*)?\s*(?:\[([^\]]+)\])?\s*(.*)/);
      
      if (paramMatch) {
        const [, type, name, source, description] = paramMatch;
        
        const isOptionalByName = name.startsWith('[') && name.endsWith(']');
        const cleanName = isOptionalByName ? name.slice(1, -1) : name;
        
        let from = ['query'];
        let isOptionalBySource = false;
        
        if (source) {
          isOptionalBySource = source.includes('?');
          from = source.split(',')
            .map(s => s.trim())
            .filter(s => s !== '?');
        }
        
        const isOptional = isOptionalByName || isOptionalBySource;
        
        result.params.push({
          name: cleanName,
          type: type || 'string',
          description: description.trim(),
          optional: isOptional,
          from
        });
      }
    } else if (!line.startsWith('@')) {
      descriptionLines.push(line);
    }
  }

  result.description = descriptionLines.join(' ').trim();
  return result;
}

/**
 * Parse createApiRoute parameters
 */
function parseCreateApiRouteParams(routeConfigLines: string[]) {
  const requiredParams: RouteParam[] = [];
  const optionalParams: RouteParam[] = [];
  
  let inRequiredParams = false;
  let inOptionalParams = false;
  let requiredBracketCount = 0;
  let optionalBracketCount = 0;

  for (const line of routeConfigLines) {
    // Check for requiredParams section
    if (line.includes('requiredParams:')) {
      inRequiredParams = true;
      inOptionalParams = false;
    }
    
    // Check for optionalParams section
    if (line.includes('optionalParams:')) {
      inOptionalParams = true;
      inRequiredParams = false;
    }

    if (inRequiredParams) {
      for (const char of line) {
        if (char === '[') requiredBracketCount++;
        if (char === ']') requiredBracketCount--;
      }

      if (line.includes('name:') && line.includes('from:') && line.includes('type:')) {
        const nameMatch = line.match(/name:\s*['"]([^'"]+)['"]/);
        const typeMatch = line.match(/type:\s*['"]([^'"]+)['"]/);
        const commentMatch = line.match(/\/\/\s*(.+)$/);
        
        // Handle both string and array formats for 'from'
        // from: 'query' or from: ['body', 'query']
        let fromSources: string[] = [];
        const fromStringMatch = line.match(/from:\s*['"]([^'"]+)['"]/);
        const fromArrayMatch = line.match(/from:\s*\[([^\]]+)\]/);
        
        if (fromStringMatch) {
          fromSources = [fromStringMatch[1]];
        } else if (fromArrayMatch) {
          // Parse array of quoted strings
          fromSources = fromArrayMatch[1]
            .split(',')
            .map(s => s.trim())
            .map(s => s.replace(/['"]/g, ''));
        }

        if (nameMatch && fromSources.length > 0 && typeMatch) {
          requiredParams.push({
            name: nameMatch[1],
            from: fromSources,
            type: typeMatch[1],
            description: commentMatch ? commentMatch[1].trim() : ''
          });
        }
      }

      if (requiredBracketCount === 0 && line.includes(']')) {
        inRequiredParams = false;
      }
    }

    if (inOptionalParams) {
      for (const char of line) {
        if (char === '[') optionalBracketCount++;
        if (char === ']') optionalBracketCount--;
      }

      if (line.includes('name:') && line.includes('from:') && line.includes('type:')) {
        const nameMatch = line.match(/name:\s*['"]([^'"]+)['"]/);
        const typeMatch = line.match(/type:\s*['"]([^'"]+)['"]/);
        const commentMatch = line.match(/\/\/\s*(.+)$/);
        
        // Handle both string and array formats for 'from'
        // from: 'query' or from: ['body', 'query']
        let fromSources: string[] = [];
        const fromStringMatch = line.match(/from:\s*['"]([^'"]+)['"]/);
        const fromArrayMatch = line.match(/from:\s*\[([^\]]+)\]/);
        
        if (fromStringMatch) {
          fromSources = [fromStringMatch[1]];
        } else if (fromArrayMatch) {
          // Parse array of quoted strings
          fromSources = fromArrayMatch[1]
            .split(',')
            .map(s => s.trim())
            .map(s => s.replace(/['"]/g, ''));
        }

        if (nameMatch && fromSources.length > 0 && typeMatch) {
          optionalParams.push({
            name: nameMatch[1],
            from: fromSources,
            type: typeMatch[1],
            description: commentMatch ? commentMatch[1].trim() : ''
          });
        }
      }

      if (optionalBracketCount === 0 && line.includes(']')) {
        inOptionalParams = false;
      }
    }
  }

  return { requiredParams, optionalParams };
}

/**
 * Generate example value based on parameter type and name
 */
function generateExampleValue(param: RouteParam): any {
  // Check parameter name for hints
  const name = param.name.toLowerCase();
  
  if (name === 'uuid') {
    return 'Actor.abc123def456';
  }
  
  if (name === 'clientid') {
    return '{{clientId}}';
  }

  // Fall back to type-based defaults
  switch (param.type) {
    case 'string':
      return 'example-value';
    case 'number':
      return '123';
    case 'boolean':
      return 'true';
    case 'array':
      return '[]';
    case 'object':
      return '{}';
    default:
      return 'value';
  }
}

/**
 * Generate test file content for a module
 */
function generateTestFile(moduleName: string, routes: RouteInfo[]): string {
  const moduleTitle = moduleName.charAt(0).toUpperCase() + moduleName.slice(1);
  
  // Build endpoints list from routes
  const endpointsList = routes.map(r => `${r.method} ${r.path}`).join(', ');
  
  let content = `/**
 * @file ${moduleName}-endpoints.test.ts
 * @generated Partially auto-generated from route docstrings (incomplete)
 * @description ${moduleTitle} Endpoint Tests
 * @endpoints ${endpointsList}
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables, setVariable } from '../helpers/testVariables';
import { captureExample, saveExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import * as path from 'path';

// Store captured examples for documentation
const capturedExamples: any[] = [];

describe('${moduleTitle}', () => {
  afterAll(() => {
    // Save captured examples for documentation
    const outputPath = path.join(__dirname, '../../docs/examples/${moduleName}-examples.json');
    saveExamples(capturedExamples, outputPath);
    console.log(\`\\nSaved \${capturedExamples.length} examples to \${outputPath}\`);
  });

`;

  // Group routes by path for better organization
  const routesByPath = new Map<string, RouteInfo[]>();
  routes.forEach(route => {
    const basePath = route.path.split('/:')[0]; // Group similar paths
    if (!routesByPath.has(basePath)) {
      routesByPath.set(basePath, []);
    }
    routesByPath.get(basePath)!.push(route);
  });

  // Generate tests for each route wrapped in forEachVersion
  content += `  forEachVersion((version, getClientId) => {\n`;
  
  routesByPath.forEach((routeGroup, basePath) => {
    content += `    describe(\`${basePath} (v\${version})\`, () => {\n`;
    
    routeGroup.forEach(route => {
      const testName = route.path;
      const method = route.method;
      
      content += `      test('${method} ${testName}', async () => {\n`;
      content += `        // Set clientId for this version\n`;
      content += `        setVariable('clientId', getClientId());\n`;
      content += `        \n`;
      content += `        // Request configuration\n`;
      content += `        // TODO: Replace placeholder values with actual test data\n`;
      content += `        const requestConfig: ApiRequestConfig = {\n`;

      // Build URL object
      content += `          url: {\n`;
      
      // Build query parameters
      const queryParams: RouteParam[] = [];
      const pathParams: RouteParam[] = [];
      const bodyParams: RouteParam[] = [];

      [...route.requiredParams, ...route.optionalParams].forEach(param => {
        if (param.from.includes('query')) {
          queryParams.push(param);
        }
        if (param.from.includes('params')) {
          pathParams.push(param);
        }
        if (param.from.includes('body')) {
          bodyParams.push(param);
        }
      });

      // Build path with params replaced
      let pathWithParams = route.path;
      pathParams.forEach(param => {
        const exampleValue = generateExampleValue(param);
        pathWithParams = pathWithParams.replace(`:${param.name}`, exampleValue);
      });

      content += `            raw: '{{baseUrl}}${pathWithParams}',\n`;
      content += `            host: ['{{baseUrl}}'],\n`;
      content += `            path: ['${route.path.replace(/^\//, '').split('/')[0]}'],\n`;

      // Add query parameters detail
      if (queryParams.length > 0) {
        content += `            query: [\n`;
        queryParams.forEach(param => {
          content += `              {\n`;
          content += `                key: '${param.name}',\n`;
          content += `                value: '${generateExampleValue(param)}',\n`;
          if (param.description) {
            content += `                description: '${param.description}'\n`;
          }
          content += `              }${param === queryParams[queryParams.length - 1] ? '' : ','}\n`;
        });
        content += `            ]\n`;
      }

      content += `          },\n`;
      content += `          method: '${method}',\n`;

      // Add headers
      content += `          header: [\n`;
      content += `            {\n`;
      content += `              key: 'x-api-key',\n`;
      content += `              value: '{{apiKey}}',\n`;
      content += `              type: 'text'\n`;
      content += `            }\n`;
      content += `          ]`;

      // Add body if POST/PUT/DELETE with body params
      if ((method === 'POST' || method === 'PUT' || method === 'DELETE') && bodyParams.length > 0) {
        content += `,\n`;
        content += `        body: {\n`;
        content += `          mode: 'raw',\n`;
        content += `          raw: JSON.stringify({\n`;
        
        bodyParams.forEach((param, index) => {
          const value = generateExampleValue(param);
          const isLast = index === bodyParams.length - 1;
          
          // Format the value based on type
          let formattedValue: string;
          if (param.type === 'boolean') {
            formattedValue = value; // Already 'true' or 'false' as string
          } else if (param.type === 'number') {
            formattedValue = value; // Already a number string like '123'
          } else if (param.type === 'array') {
            formattedValue = value; // Already '[]' format
          } else if (param.type === 'object') {
            formattedValue = value; // Already '{}' format
          } else {
            // String value - wrap in single quotes
            formattedValue = `'${value}'`;
          }
          
          content += `              ${param.name}: ${formattedValue}${isLast ? '' : ','}\n`;
        });
        
        content += `            }, null, 2)\n`;
        content += `        }\n`;
      } else {
        content += `\n`;
      }

      content += `      };\n\n`;

      // Add example capture (this also makes the request)
      content += `        // Capture this example for documentation (also makes the request)\n`;
      content += `        const captured = await captureExample(\n`;
      content += `          requestConfig,\n`;
      content += `          testVariables,\n`;
      content += `          '${route.path}'\n`;
      content += `        );\n`;
      content += `        capturedExamples.push(captured);\n\n`;

      // Add assertions
      content += `        // Assertions\n`;
      content += `        // TODO: Add test assertions\n`;
      content += `        expect(captured.response.status).toBe(200);\n`;
      content += `      });\n`;
    });
    
    content += `    });\n\n`;
  });

  content += `  });\n\n`;
  content += `});\n`;

  return content;
}

/**
 * Main execution
 */
function main() {
  const apiDir = path.join(__dirname, '../src/routes/api');
  const testDir = path.join(__dirname, '../tests/integration');

  // Get all route files
  const routeFiles = fs.readdirSync(apiDir).filter(file => file.endsWith('.ts'));

  console.log('Scanning for route files that need test generation...\n');

  let generatedCount = 0;

  routeFiles.forEach(file => {
    const moduleName = path.basename(file, '.ts');
    const testFileName = `${moduleName}-endpoints.test.ts`;
    const testFilePath = path.join(testDir, testFileName);

    // Skip if test file already exists
    if (fs.existsSync(testFilePath)) {
      console.log(`✓ ${testFileName} already exists, skipping...`);
      return;
    }

    const routeFilePath = path.join(apiDir, file);
    console.log(`Generating tests for ${moduleName}...`);

    try {
      const routes = extractRoutes(routeFilePath);

      if (routes.length === 0) {
        console.log(`No routes found in ${file}, skipping...`);
        return;
      }

      const testContent = generateTestFile(moduleName, routes);
      fs.writeFileSync(testFilePath, testContent);
      
      console.log(`Generated ${testFileName} with ${routes.length} test(s)`);
      generatedCount++;
    } catch (error) {
      console.error(`Error processing ${file}:`, error instanceof Error ? error.message : error);
    }
  });

  console.log(`\nTest generation complete! Generated ${generatedCount} new test file(s).`);
  
  if (generatedCount === 0) {
    console.log('All route files already have corresponding test files.');
  }
}

main();
