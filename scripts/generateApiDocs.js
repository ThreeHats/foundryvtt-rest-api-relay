const fs = require('fs');
const path = require('path');

/**
 * Extract API documentation from TypeScript route files
 * This script parses the route files to extract JSDoc comments and createApiRoute information
 */

function extractApiInfo(filePath) {
  const content = fs.readFileSync(filePath, 'utf8');
  
  // Read the file directly as lines to better parse structure
  const lines = content.split('\n');
  const routes = [];
  
    // Find router method calls - both traditional endpoints and createApiRoute style
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if ((line.includes('Router') || line.includes('router')) && 
        (line.includes('.get(') || line.includes('.post(') || 
         line.includes('.put(') || line.includes('.delete('))) {
      
      // Extract method and path - more flexible pattern to match various formats
      const methodMatch = line.match(/(\w*[Rr]outer)\.(get|post|put|delete)\(\s*["']([^"']+)["']/);
      if (methodMatch) {
        const [, routerName, method, path] = methodMatch;        // Look backwards for JSDoc comment
        let j = i - 1;
        let foundJSDocEnd = false;
        let foundJSDocStart = false;
        let jsdocContent = '';
        
        // First look for a JSDoc block above this line
        while (j >= 0 && j >= i - 50) { // Look back up to 50 lines
          const currentLine = lines[j].trim();
          
          // If we find a non-comment line after finding the end but before finding the start, then break
          if (foundJSDocEnd && !foundJSDocStart && 
              !currentLine.startsWith('//') && !currentLine.startsWith('/*') && 
              currentLine !== '' && !currentLine.includes('*/')) {
            break;
          }
          
          // Found the end marker of a JSDoc block
          if (currentLine.includes('*/')) {
            foundJSDocEnd = true;
            // Note the position where the JSDoc block ends
            const endPos = j;
            
            // Continue searching backward for the start marker
            while (j >= 0 && j >= i - 50) {
              if (lines[j].includes('/**')) {
                foundJSDocStart = true;
                
                // Extract the full JSDoc block from start to end
                jsdocContent = lines.slice(j, endPos + 1).join('\n');
                break;
              }
              j--;
            }
            
            // If we found the start, break the outer loop too
            if (foundJSDocStart) {
              break;
            }
          }
          
          j--;
        }
        
        // If no JSDoc was found, check for simple comments
        if (!jsdocContent) {
          // Look for a comment block starting with //
          j = i - 1;
          let commentBlock = [];
          
          while (j >= 0 && j >= i - 10) { // Look back up to 10 lines for // comments
            const currentLine = lines[j].trim();
            
            if (currentLine.startsWith('//')) {
              commentBlock.unshift(currentLine.substring(2).trim());
            } else if (currentLine !== '') {
              break; // Stop at non-comment, non-empty line
            }
            
            j--;
          }
          
          if (commentBlock.length > 0) {
            jsdocContent = '/**\n * ' + commentBlock.join('\n * ') + '\n */';
          }
        }
        
        const jsdoc = parseJSDoc(jsdocContent);
        
        // Look forward for createApiRoute configuration or parse function parameters for traditional endpoints
        let routeConfigLines = [];
        let k = i;
        let bracketCount = 0;
        let foundCreateApiRoute = false;
        let isTraditionalEndpoint = false;
        let handlerStartLine = -1;
        
        // First check if this is a traditional endpoint (non-createApiRoute)
        // Look ahead a bit to determine if this is a traditional endpoint or uses createApiRoute
        let lookAhead = i;
        let hasCreateApiRoute = false;
        const maxLookAhead = Math.min(i + 10, lines.length);
        
        while (lookAhead < maxLookAhead) {
          if (lines[lookAhead].includes('createApiRoute')) {
            hasCreateApiRoute = true;
            break;
          }
          lookAhead++;
        }
        
        // If no createApiRoute found, it's a traditional endpoint
        if (!hasCreateApiRoute) {
          isTraditionalEndpoint = true;
          handlerStartLine = i;
        }
        
        while (k < lines.length) {
          const currentLine = lines[k];
          
          if (currentLine.includes('createApiRoute')) {
            foundCreateApiRoute = true;
            isTraditionalEndpoint = false; // Override if we find createApiRoute
          }
          
          // For createApiRoute style
          if (foundCreateApiRoute) {
            routeConfigLines.push(currentLine);
            
            // Count brackets to find the end of the config
            for (const char of currentLine) {
              if (char === '{') bracketCount++;
              if (char === '}') bracketCount--;
            }
            
            if (bracketCount === 0 && currentLine.includes('})')) {
              break;
            }
          } 
          // For traditional endpoints, we need to ensure they have routes added
          else if (isTraditionalEndpoint) {
            // Just add a placeholder to ensure we have content for traditional endpoints
            routeConfigLines.push('// Traditional endpoint');
            // We need to break after adding one line to ensure we don't keep adding lines
            break;
          } if (isTraditionalEndpoint && k >= i + 5) {
            // Just break early, we're using JSDoc for params
            break;
          }
          
          k++;
        }
        
        const routeConfig = routeConfigLines.join('\n');
        
        // Extract required and optional parameters
        let requiredParams = [];
        let optionalParams = [];
        let inRequiredParams = false;
        let requiredBracketCount = 0;
        
        for (let m = 0; m < routeConfigLines.length; m++) {
          const line = routeConfigLines[m];
          
          if (line.includes('requiredParams:')) {
            inRequiredParams = true;
          }
          
          if (inRequiredParams) {
            // Count brackets
            for (const char of line) {
              if (char === '[') requiredBracketCount++;
              if (char === ']') requiredBracketCount--;
            }
            
            // Extract parameter objects and their comments
            if (line.includes('name:') && line.includes('from:') && line.includes('type:')) {
              // Extract parameter definition
              const nameMatch = line.match(/name:\s*['"]([^'"]+)['"]/);
              const typeMatch = line.match(/type:\s*['"]([^'"]+)['"]/);
              let fromMatch = line.match(/from:\s*\[([^\]]+)\]/);
              
              if (!fromMatch) {
                fromMatch = line.match(/from:\s*['"]([^'"]+)['"]/);
              }
              
              if (nameMatch && typeMatch) {
                const name = nameMatch[1];
                const type = typeMatch[1];
                const fromValue = fromMatch ? fromMatch[1] : '';
                const from = fromValue.includes(',') 
                  ? fromValue.split(',').map(s => s.trim().replace(/['"]/g, '')) 
                  : [fromValue.trim().replace(/['"]/g, '')];
                
                // Extract inline comment - use a more flexible pattern that works with the actual file format
                const commentIndex = line.indexOf('//');
                const description = commentIndex > -1 ? line.substring(commentIndex + 2).trim() : '';
                
                requiredParams.push({ name, type, from, description });
              }
            }
            
            if (requiredBracketCount === 0 && line.includes(']')) {
              inRequiredParams = false;
            }
          }
        }
        
        // Extract optional parameters from createApiRoute if not a traditional endpoint
        let inOptionalParams = false;
        let optionalBracketCount = 0;
        
        if (!isTraditionalEndpoint) {
          for (let m = 0; m < routeConfigLines.length; m++) {
          const line = routeConfigLines[m];
          
          if (line.includes('optionalParams:')) {
            inOptionalParams = true;
          }
          
          if (inOptionalParams) {
            // Count brackets
            for (const char of line) {
              if (char === '[') optionalBracketCount++;
              if (char === ']') optionalBracketCount--;
            }
            
            // Extract parameter objects and their comments
            if (line.includes('name:') && line.includes('from:') && line.includes('type:')) {
              // Extract parameter definition
              const nameMatch = line.match(/name:\s*['"]([^'"]+)['"]/);
              const typeMatch = line.match(/type:\s*['"]([^'"]+)['"]/);
              let fromMatch = line.match(/from:\s*\[([^\]]+)\]/);
              
              if (!fromMatch) {
                fromMatch = line.match(/from:\s*['"]([^'"]+)['"]/);
              }
              
              if (nameMatch && typeMatch) {
                const name = nameMatch[1];
                const type = typeMatch[1];
                const fromValue = fromMatch ? fromMatch[1] : '';
                const from = fromValue.includes(',') 
                  ? fromValue.split(',').map(s => s.trim().replace(/['"]/g, '')) 
                  : [fromValue.trim().replace(/['"]/g, '')];
                
                // Extract inline comment - use the same approach as with required params
                const commentIndex = line.indexOf('//');
                const description = commentIndex > -1 ? line.substring(commentIndex + 2).trim() : '';
                
                optionalParams.push({ name, type, from, description });
              }
            }
            
            if (optionalBracketCount === 0 && line.includes(']')) {
              inOptionalParams = false;
            }
          }
        } // End of if (!isTraditionalEndpoint) block
        
        // Extract type
        let type = '';
        for (const line of routeConfigLines) {
          const typeMatch = line.match(/type:\s*['"]([^'"]+)['"]/);
          if (typeMatch) {
            type = typeMatch[1];
            break;
          }
        }
        
        // Process routes differently based on whether they're traditional endpoints or createApiRoute style
        if (isTraditionalEndpoint) {
          // For traditional endpoints, transform JSDoc params into requiredParams and optionalParams format
          // This ensures consistent documentation format between both endpoint styles
          const processedRequiredParams = jsdoc.params
            .filter(param => !param.optional)
            .map(param => ({
              name: param.name,
              type: param.type || 'string',
              description: param.description || '',
              from: param.from || ['query'] // Default to query if not specified
            }));
          
          const processedOptionalParams = jsdoc.params
            .filter(param => param.optional)
            .map(param => ({
              name: param.name,
              type: param.type || 'string',
              description: param.description || '',
              from: param.from || ['query'] // Default to query if not specified
            }));
          
          // Always push the route for traditional endpoints if we have JSDoc content
          routes.push({
            method: method.toUpperCase(),
            path: path, // Use actual router path, not @route tag in case they don't match up
            route: jsdoc.route ? jsdoc.route.path : path, // Capture this for use in api-docs.json, so that path prefixes are captured.
            description: jsdoc.description,
            params: jsdoc.params,
            returns: jsdoc.returns,
            group: jsdoc.group,
            security: jsdoc.security,
            requiredParams: processedRequiredParams,
            optionalParams: processedOptionalParams,
            type: 'json' // Default type for traditional endpoints
          });
        } else {
          // createApiRoute style endpoints
          routes.push({
            method: method.toUpperCase(),
            path: path, // Use actual router path, not @route tag in case they don't match up
            route: jsdoc.route ? jsdoc.route.path : path, // Capture this for use in api-docs.json, so that path prefixes are captured.
            description: jsdoc.description,
            params: jsdoc.params,
            returns: jsdoc.returns,
            group: jsdoc.group,
            security: jsdoc.security,
            requiredParams,
            optionalParams,
            type
          });
        }
      } else if (isTraditionalEndpoint && jsdocContent) {
        // Process parameters from JSDoc for traditional endpoints
        const processedRequiredParams = jsdoc.params
          .filter(param => !param.optional)
          .map(param => ({
            name: param.name,
            type: param.type || 'string',
            description: param.description || '',
            from: param.from || ['query'] // Default to query if not specified
          }));
        
        const processedOptionalParams = jsdoc.params
          .filter(param => param.optional)
          .map(param => ({
            name: param.name,
            type: param.type || 'string',
            description: param.description || '',
            from: param.from || ['query'] // Default to query if not specified
          }));
        
        routes.push({
          method: method.toUpperCase(),
          path: path, // Use actual router path, not @route tag in case they don't match up
          route: jsdoc.route ? jsdoc.route.path : path, // Capture this for use in api-docs.json, so that path prefixes are captured.
          description: jsdoc.description,
          params: jsdoc.params,
          returns: jsdoc.returns,
          group: jsdoc.group,
          security: jsdoc.security,
          requiredParams: processedRequiredParams,
          optionalParams: processedOptionalParams,
          type: 'json' // Default type for traditional endpoints
        });
      }
    }
  }}
  
  return routes;
}

function extractParamsWithComments(paramsString) {
  const lines = paramsString.split('\n');
  const params = [];
  
  for (const line of lines) {
    // Match parameter definition
    const paramMatch = line.match(/{\s*name:\s*['"]([^'"]+)['"],\s*from:\s*(?:\[([^\]]+)\]|['"]([^'"]+)['"])\s*,\s*type:\s*['"]([^'"]+)['"]\s*}/);
    
    if (paramMatch) {
      const [, name, fromArray, fromString, type] = paramMatch;
      const from = fromArray ? fromArray.split(',').map(s => s.trim().replace(/['"]/g, '')) : [fromString];
      
      // Match inline comment if present
      const commentMatch = line.match(/\/\/\s*(.*?)$/);
      const description = commentMatch ? commentMatch[1].trim() : '';
      
      params.push({ name, from, type, description });
    }
  }
  
  return params;
}

function parseJSDoc(jsdocContent) {
  // Clean up the input content
  const lines = jsdocContent.split('\n')
    .map(line => line.replace(/^\s*\/\*\*\s*/, ''))  // Remove opening /**
    .map(line => line.replace(/^\s*\*\s?/, ''))      // Remove leading * 
    .map(line => line.replace(/\s*\*\/\s*$/, ''))    // Remove closing */
    .map(line => line.trim());
  
  const result = {
    description: '',
    params: [],
    returns: null,
    group: null,
    security: null,
    route: null
  };
  
  let currentSection = 'description';
  let descriptionLines = [];
  
  for (const line of lines) {
    if (line.startsWith('@param')) {
      currentSection = 'param';
      
      // Format: @param {type} name - [source,?] Description
      // This pattern specifically looks for the source in square brackets at the start of the description
      // Now also supports optional parameters marked with ? in the source part
      const paramMatch = line.match(/@param\s+(?:\{([^}]+)\})?\s*(\[?[\w\-\[\]\.]+\]?)(?:\s*-\s*)?\s*(?:\[([^\]]+)\])?\s*(.*)/);
      
      if (paramMatch) {
        const [, type, name, source, description] = paramMatch;
        
        // Check if parameter is optional based on name format [paramName]
        const isOptionalByName = name.startsWith('[') && name.endsWith(']');
        const cleanName = isOptionalByName ? name.slice(1, -1) : name;
        
        // Determine parameter source and check if marked as optional with ,?
        let from = ['query']; // Default source is query
        let isOptionalBySource = false;
        
        if (source) {
          // Check if any source contains ? to denote optional
          isOptionalBySource = source.includes('?');
          // Handle comma-separated sources like "query,body" or "query,?"
          from = source.split(',')
            .map(s => s.trim())
            .filter(s => s !== '?'); // Remove the ? marker from sources
        }
        
        // Parameter is optional if either notation is used
        const isOptional = isOptionalByName || isOptionalBySource;
        
        result.params.push({
          name: cleanName,
          type: type || 'string', // Default to string if type is not specified
          description: description.trim(),
          optional: isOptional,
          from: from
        });
      } else {
        // Try a simpler regex if the first one fails
        const simpleParamMatch = line.match(/@param\s+(?:\{([^}]+)\})?\s*(\[?[\w\-\[\]\.]+\]?)\s+(.*)/);
        if (simpleParamMatch) {
          const [, type, name, description] = simpleParamMatch;
          const isOptionalByName = name.startsWith('[') && name.endsWith(']');
          const cleanName = isOptionalByName ? name.slice(1, -1) : name;
          
          // Check for source in the description like "[query] Description", "[query,body] Description" or "[query,?] Description"
          let from = ['query']; // Default source
          let isOptionalBySource = false;
          const sourceMatch = description.match(/^\s*\[([^\]]+)\]\s*(.*)/);
          if (sourceMatch) {
            const [, source, restOfDescription] = sourceMatch;
            // Check if any source contains ? to denote optional
            isOptionalBySource = source.includes('?');
            from = source.split(',')
              .map(s => s.trim())
              .filter(s => s !== '?'); // Remove the ? marker from sources
            
            // Parameter is optional if either notation is used
            const isOptional = isOptionalByName || isOptionalBySource;
            
            result.params.push({
              name: cleanName,
              type: type || 'string', // Default to string
              description: restOfDescription.trim(),
              optional: isOptional,
              from: from
            });
          } else {
            // If no source is specified in the description, just use the name to determine if optional
            result.params.push({
              name: cleanName,
              type: type || 'string', // Default to string
              description: description.trim(),
              optional: isOptionalByName,
              from: from
            });
          }
        }
      }
    } else if (line.startsWith('@returns')) {
      const returnsMatch = line.match(/@returns\s+(?:\{([^}]+)\})?\s*(.*)/);
      if (returnsMatch) {
        const [, type, description] = returnsMatch;
        result.returns = { type: type || 'object', description: description.trim() };
      }
    } else if (line.startsWith('@route')) {
      const routeMatch = line.match(/@route\s+(GET|POST|PUT|DELETE)\s+([^\s]+)/i);
      if (routeMatch) {
        result.route = {
          method: routeMatch[1].toUpperCase(),
          path: routeMatch[2]
        };
      }
    } else if (line.startsWith('@group')) {
      result.group = line.replace('@group', '').trim();
    } else if (line.startsWith('@security')) {
      result.security = line.replace('@security', '').trim();
    } else if (currentSection === 'description' && line && !line.startsWith('@')) {
      descriptionLines.push(line);
    }
  }
  
  // Join the description lines and trim any whitespace
  let description = descriptionLines.join(' ').trim();
  
  // Remove trailing " /" which can appear from JSDoc parsing in createApiRoute
  if (description.endsWith(' /')) {
    description = description.slice(0, -2).trim();
  }
  
  result.description = description;
  return result;
}

function parseRouteConfig(configContent) {
  // Simple parsing of the route configuration
  const result = {
    type: null,
    requiredParams: [],
    optionalParams: []
  };
  
  // Extract type
  const typeMatch = configContent.match(/type:\s*['"]([^'"]+)['"]/);
  if (typeMatch) {
    result.type = typeMatch[1];
  }
  
  // Extract required parameters
  const requiredMatch = configContent.match(/requiredParams:\s*\[([\s\S]*?)\]/);
  if (requiredMatch) {
    result.requiredParams = parseParamArray(requiredMatch[1]);
  }
  
  // Extract optional parameters
  const optionalMatch = configContent.match(/optionalParams:\s*\[([\s\S]*?)\]/);
  if (optionalMatch) {
    result.optionalParams = parseParamArray(optionalMatch[1]);
  }
  
  return result;
}

function parseParamArray(paramArrayContent) {
  const params = [];
  // Match parameter objects and look for inline comments after them
  const paramLines = paramArrayContent.split('\n');
  
  for (let i = 0; i < paramLines.length; i++) {
    const line = paramLines[i].trim();
    const paramMatch = line.match(/\{\s*name:\s*['"]([^'"]+)['"],\s*from:\s*(?:\[([^\]]+)\]|['"]([^'"]+)['"])\s*,\s*type:\s*['"]([^'"]+)['"]\s*\}/);
    
    if (paramMatch) {
      const [, name, fromArray, fromString, type] = paramMatch;
      const from = fromArray ? fromArray.split(',').map(s => s.trim().replace(/['"]/g, '')) : [fromString];
      
      // Look for inline comment on the same line
      const commentMatch = line.match(/\/\/\s*(.*?)$/);
      const description = commentMatch ? commentMatch[1].trim() : '';
      
      params.push({ name, from, type, description });
    }
  }
  
  return params;
}

/**
 * Known parameterized route expansions.
 * When a route path contains a parameter (e.g., :documentType), and the parameter
 * has a known set of valid values, expand the single route into multiple doc headings.
 * This keeps the route code DRY while giving each concrete endpoint its own doc section.
 */
const ROUTE_PARAM_EXPANSIONS = {
  ':documentType': ['tokens', 'tiles', 'drawings', 'lights', 'sounds', 'notes', 'templates', 'walls'],
};

/**
 * Generate a single route section (heading, description, params, returns).
 * Used for both regular and expanded routes.
 */
function generateRouteSection(method, routePath, route) {
    let section = `## ${method} ${routePath}\n\n`;
    if (route.description) {
        section += `${route.description}\n\n`;
    }

    // Parameters table — omit the expanded param since it's now part of the path
    const expandedParamName = Object.keys(ROUTE_PARAM_EXPANSIONS)
        .find(p => route.route.includes(p))
        ?.replace(':', '');

    const hasRequiredParams = route.requiredParams && route.requiredParams.length > 0;
    const hasOptionalParams = route.optionalParams && route.optionalParams.length > 0;

    if (hasRequiredParams || hasOptionalParams) {
        const filteredRequired = (route.requiredParams || []).filter(p => p.name !== expandedParamName);
        const filteredOptional = (route.optionalParams || []).filter(p => p.name !== expandedParamName);

        if (filteredRequired.length > 0 || filteredOptional.length > 0) {
            section += `### Parameters\n\n`;
            section += `| Name | Type | Required | Source | Description |\n`;
            section += `|------|------|----------|--------|--------------|\n`;

            filteredRequired.forEach(param => {
                const source = param.from ? param.from.join(', ') : 'body, query';
                section += `| ${param.name} | ${param.type} | ✓ | ${source} | ${param.description || ''} |\n`;
            });
            filteredOptional.forEach(param => {
                const source = param.from ? param.from.join(', ') : 'body, query';
                section += `| ${param.name} | ${param.type} |  | ${source} | ${param.description || ''} |\n`;
            });

            section += '\n';
        }
    }

    if (route.returns) {
        section += `### Returns\n\n`;
        section += `**${route.returns.type}** - ${route.returns.description}\n\n`;
    }

    // ApiTester component
    section += `### Try It Out\n\n`;
    const allParameters = [
        ...(route.requiredParams || []).filter(p => p.name !== expandedParamName).map(p => ({
            name: p.name,
            type: p.type,
            required: true,
            source: Array.isArray(p.from) ? p.from[0] : (p.from || 'query')
        })),
        ...(route.optionalParams || []).filter(p => p.name !== expandedParamName).map(p => ({
            name: p.name,
            type: p.type,
            required: false,
            source: Array.isArray(p.from) ? p.from[0] : (p.from || 'query')
        }))
    ];
    const parametersJson = JSON.stringify(allParameters);
    section += `<ApiTester\n`;
    section += `  method="${method}"\n`;
    section += `  path="${routePath}"\n`;
    section += `  parameters={${parametersJson}}\n`;
    section += `/>\n\n`;

    return section;
}

function generateMarkdown(routes, moduleName) {
    let markdown = `---\n`;
    markdown += `tag: ${moduleName}\n`;
    markdown += `---\n\n`;
    markdown += `import ApiTester from '@site/src/components/ApiTester';\n\n`;
    markdown += `# ${moduleName}\n\n`;

    const sections = [];

    routes.forEach(route => {
        // Check if this route has an expandable parameter
        const paramKey = Object.keys(ROUTE_PARAM_EXPANSIONS).find(p => route.route.includes(p));

        if (paramKey) {
            // Expand into one section per concrete value
            const values = ROUTE_PARAM_EXPANSIONS[paramKey];
            values.forEach(value => {
                const concretePath = route.route.replace(paramKey, value);
                sections.push(generateRouteSection(route.method, concretePath, route));
            });
        } else {
            sections.push(generateRouteSection(route.method, route.route, route));
        }
    });

    markdown += sections.join('---\n\n');
    return markdown;
}

// Helper function to generate example values for different parameter types
function getExampleValue(type) {
  switch(type) {
    case 'string': return 'example-value';
    case 'number': return 123;
    case 'boolean': return true;
    case 'array': return ['example'];
    case 'object': return { key: 'value' };
    default: return 'value';
  }
}

/**
 * Generate an OpenAPI 3.0.3 spec from the extracted route data.
 */
function generateOpenApiSpec(allRoutes, packageJson) {
  const spec = {
    openapi: '3.0.3',
    info: {
      title: 'FoundryVTT REST API',
      description: 'REST API relay server for accessing Foundry VTT data remotely. Provides WebSocket connectivity and HTTP endpoints to interact with Foundry VTT worlds.',
      version: packageJson.version || '1.0.0',
      license: {
        name: 'MIT'
      }
    },
    servers: [
      {
        url: 'http://localhost:3010',
        description: 'Replaced dynamically at /openapi.json'
      }
    ],
    security: [{ apiKey: [] }],
    tags: [],
    paths: {},
    components: {
      securitySchemes: {
        apiKey: {
          type: 'apiKey',
          in: 'header',
          name: 'x-api-key',
          description: 'API key for authentication. Obtain one by registering at /auth/register.'
        }
      }
    }
  };

  // Derive tags from routes
  const tagSet = new Set();
  for (const route of allRoutes) {
    if (route.group) {
      tagSet.add(route.group);
    }
  }
  spec.tags = Array.from(tagSet).sort().map(name => ({ name }));

  // Type mapping helper
  function mapType(type) {
    switch (type) {
      case 'string': return { type: 'string' };
      case 'number': return { type: 'number' };
      case 'boolean': return { type: 'boolean' };
      case 'array': return { type: 'array', items: { type: 'string' } };
      case 'object': return { type: 'object' };
      default: return { type: 'string' };
    }
  }

  // Public endpoints that don't require auth
  const publicPaths = ['/api/status', '/api/docs', '/api/health'];

  for (const route of allRoutes) {
    // Use the full route path (from @route tag), converting Express :param to OpenAPI {param}
    let routePath = route.route || route.path;
    routePath = routePath.replace(/:(\w+)/g, '{$1}');

    // Ensure path starts with /
    if (!routePath.startsWith('/')) {
      routePath = '/' + routePath;
    }

    const method = route.method.toLowerCase();

    if (!spec.paths[routePath]) {
      spec.paths[routePath] = {};
    }

    const operation = {
      summary: route.description || '',
      operationId: `${method}_${routePath.replace(/[^a-zA-Z0-9]/g, '_').replace(/_+/g, '_')}`,
      responses: {
        '200': {
          description: route.returns ? route.returns.description : 'Successful response',
          content: {
            'application/json': {
              schema: { type: 'object' }
            }
          }
        },
        '400': {
          description: 'Bad request - missing or invalid parameters',
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: {
                  error: { type: 'string' }
                }
              }
            }
          }
        },
        '401': {
          description: 'Unauthorized - invalid or missing API key',
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: {
                  error: { type: 'string' }
                }
              }
            }
          }
        }
      }
    };

    // Add tag
    if (route.group) {
      operation.tags = [route.group];
    }

    // Handle security
    if (route.security === 'none' || publicPaths.includes(routePath)) {
      operation.security = [];
    }

    // Collect params by location
    const allParams = [
      ...(route.requiredParams || []).map(p => ({ ...p, required: true })),
      ...(route.optionalParams || []).map(p => ({ ...p, required: false }))
    ];

    const queryParams = [];
    const pathParams = [];
    const bodyProps = {};
    const requiredBodyProps = [];
    let hasBodyParams = false;

    for (const param of allParams) {
      const locations = Array.isArray(param.from) ? param.from : [param.from];

      // Check if this is a path parameter (appears in the route path as {param})
      if (routePath.includes(`{${param.name}}`)) {
        pathParams.push({
          name: param.name,
          in: 'path',
          required: true,
          description: param.description || '',
          schema: mapType(param.type)
        });
        continue;
      }

      if (locations.includes('params')) {
        pathParams.push({
          name: param.name,
          in: 'path',
          required: true,
          description: param.description || '',
          schema: mapType(param.type)
        });
      }

      if (locations.includes('query')) {
        queryParams.push({
          name: param.name,
          in: 'query',
          required: param.required && !locations.includes('body'),
          description: param.description || '',
          schema: mapType(param.type)
        });
      }

      if (locations.includes('body')) {
        hasBodyParams = true;
        bodyProps[param.name] = {
          ...mapType(param.type),
          description: param.description || ''
        };
        if (param.required && !locations.includes('query')) {
          requiredBodyProps.push(param.name);
        }
      }

      // If param has both body and query sources, add to both
      if (locations.includes('body') && locations.includes('query')) {
        // Already handled above - appears in both
      }

      // Default: if only 'query' or unspecified, already handled
    }

    // Add parameters
    const parameters = [...pathParams, ...queryParams];
    if (parameters.length > 0) {
      operation.parameters = parameters;
    }

    // Add request body for POST/PUT/DELETE with body params
    if (hasBodyParams && ['post', 'put', 'delete'].includes(method)) {
      // Special case: upload endpoint uses multipart
      if (routePath === '/upload') {
        operation.requestBody = {
          required: true,
          content: {
            'multipart/form-data': {
              schema: {
                type: 'object',
                properties: bodyProps,
                ...(requiredBodyProps.length > 0 ? { required: requiredBodyProps } : {})
              }
            },
            'application/json': {
              schema: {
                type: 'object',
                properties: bodyProps,
                ...(requiredBodyProps.length > 0 ? { required: requiredBodyProps } : {})
              }
            }
          }
        };
      } else {
        operation.requestBody = {
          required: requiredBodyProps.length > 0,
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: bodyProps,
                ...(requiredBodyProps.length > 0 ? { required: requiredBodyProps } : {})
              }
            }
          }
        };
      }
    }

    spec.paths[routePath][method] = operation;
  }

  return spec;
}

// Main execution
function main() {
  const apiDir = path.join(__dirname, '../src/routes/api');
  const markdownOutputDir = path.join(__dirname, '../docs/md/api');
  const jsonOutputDir = path.join(__dirname, '../public');
  const packageJsonPath = path.join(__dirname, '../package.json');

  // Ensure output directories exist
  if (!fs.existsSync(markdownOutputDir)) {
    fs.mkdirSync(markdownOutputDir, { recursive: true });
    console.log('Created markdown output directory:', markdownOutputDir);
  }
  if (!fs.existsSync(jsonOutputDir)) {
    fs.mkdirSync(jsonOutputDir, { recursive: true });
    console.log('Created JSON output directory:', jsonOutputDir);
  }

  const allRoutes = [];
  // Process each TypeScript file in the API directory
  const files = fs.readdirSync(apiDir).filter(file => file.endsWith('.ts'));
  
  // Debug - Process specific files first to troubleshoot
  const debugFiles = ['fileSystem.ts', 'session.ts', 'sheet.ts'];
  const otherFiles = files.filter(file => !debugFiles.includes(file));
  const sortedFiles = [...debugFiles, ...otherFiles].filter(file => files.includes(file));
  
  sortedFiles.forEach(file => {
    const filePath = path.join(apiDir, file);
    const moduleName = path.basename(file, '.ts');

    console.log(`Processing ${file}...`);

    try {
      const routes = extractApiInfo(filePath);

      if (routes.length > 0) {
        // Auto-assign group tag from module name if not set
        for (const route of routes) {
          if (!route.group) {
            route.group = moduleName.charAt(0).toUpperCase() + moduleName.slice(1);
          }
        }
        // Collect routes for the JSON file
        allRoutes.push(...routes);

        // Generate markdown for docusaurus
        let markdown = generateMarkdown(routes, moduleName);
        const outputPath = path.join(markdownOutputDir, `${moduleName}.md`);
        
        // Do NOT preserve code examples - they will be regenerated by updateDocsWithExamples.ts
        // This prevents duplication issues when running docs:generate multiple times
        
        fs.writeFileSync(outputPath, markdown);
        console.log(`Generated documentation for ${routes.length} routes in ${moduleName}.md`);
      } else {
        console.log(`No routes found in ${file}`);
      }
    } catch (error) {
      console.error(`Error processing ${file}:`, error.message);
    }
  });
  
  // Also parse auth routes
  const authFilePath = path.join(__dirname, '../src/routes/auth.ts');
  try {
    console.log('Processing auth.ts...');
    const authRoutes = extractApiInfo(authFilePath);
    if (authRoutes.length > 0) {
      allRoutes.push(...authRoutes);
      console.log(`Extracted ${authRoutes.length} auth routes`);
    } else {
      console.log('No routes found in auth.ts');
    }
  } catch (error) {
    console.error('Error processing auth.ts:', error.message);
  }

  console.log('Markdown documentation generation complete!');

  // Generate the single api-docs.json file
  try {
    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
    const apiDocs = {
      version: packageJson.version || "1.0.0",
      baseUrl: `https://foundryvtt-rest-api-relay.fly.dev`, // This will be replaced by the server dynamically
      authentication: {
        required: true,
        headerName: "x-api-key",
        description: "API key must be included in the x-api-key header for all endpoints except /api/status and /api/docs"
      },
      endpoints: allRoutes.flatMap(route => {
          const paramKey = Object.keys(ROUTE_PARAM_EXPANSIONS).find(p => route.route.includes(p));
          const expandedParamName = paramKey?.replace(':', '');
          const paths = paramKey
              ? ROUTE_PARAM_EXPANSIONS[paramKey].map(v => route.route.replace(paramKey, v))
              : [route.route];
          return paths.map(routePath => ({
              method: route.method.toUpperCase(),
              path: routePath,
              description: route.description,
              requiredParameters: route.requiredParams
                  .filter(p => p.name !== expandedParamName)
                  .map(p => ({ name: p.name, type: p.type, description: p.description, location: p.from.join(', ') })),
              optionalParameters: route.optionalParams
                  .filter(p => p.name !== expandedParamName)
                  .map(p => ({ name: p.name, type: p.type, description: p.description, location: p.from.join(', ') })),
          }));
      })
    };

    const outputPath = path.join(jsonOutputDir, `api-docs.json`);
    fs.writeFileSync(outputPath, JSON.stringify(apiDocs, null, 2));
    console.log(`Generated JSON documentation for ${allRoutes.length} routes in api-docs.json`);
  } catch (error) {
    console.error(`Error generating api-docs.json:`, error.message);
  }

  // Generate OpenAPI 3.0 spec
  try {
    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
    const openApiSpec = generateOpenApiSpec(allRoutes, packageJson);

    // Inject response examples from captured test data
    const examplesDir = path.join(__dirname, '../docs/examples');
    if (fs.existsSync(examplesDir)) {
      const { sanitizeExampleForDocs } = require('./shared/sanitize');
      const exampleFiles = fs.readdirSync(examplesDir).filter(f => f.endsWith('-examples.json'));
      let injectedCount = 0;

      for (const file of exampleFiles) {
        try {
          const examples = JSON.parse(fs.readFileSync(path.join(examplesDir, file), 'utf8'));

          for (const example of examples) {
            // Normalize endpoint path: convert Express :param to OpenAPI {param}
            let examplePath = (example.description || example.endpoint || '').replace(/:(\w+)/g, '{$1}');
            if (!examplePath.startsWith('/')) examplePath = '/' + examplePath;

            const method = (example.method || 'GET').toLowerCase();
            const pathEntry = openApiSpec.paths[examplePath];

            if (pathEntry && pathEntry[method]) {
              const operation = pathEntry[method];
              const statusCode = String(example.response?.status || 200);

              // Only inject if we don't already have an example for this status
              if (operation.responses[statusCode]?.content?.['application/json']) {
                const content = operation.responses[statusCode].content['application/json'];
                if (!content.example) {
                  const sanitized = sanitizeExampleForDocs(example);
                  content.example = sanitized.response.data;
                  injectedCount++;
                }
              }
            }
          }
        } catch (err) {
          console.warn(`  Warning: Could not process examples from ${file}: ${err.message}`);
        }
      }

      if (injectedCount > 0) {
        console.log(`Injected ${injectedCount} response examples into OpenAPI spec`);
      }
    } else {
      console.log('No examples directory found, skipping example injection');
    }

    const openApiOutputPath = path.join(jsonOutputDir, 'openapi.json');
    fs.writeFileSync(openApiOutputPath, JSON.stringify(openApiSpec, null, 2));
    console.log(`Generated OpenAPI 3.0 spec with ${Object.keys(openApiSpec.paths).length} paths in openapi.json`);

  } catch (error) {
    console.error('Error generating openapi.json:', error.message);
  }

  // Generate AsyncAPI spec for WebSocket API
  try {
    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
    const asyncApiSpec = generateAsyncApiSpec(allRoutes, packageJson);
    const asyncApiOutputPath = path.join(jsonOutputDir, 'asyncapi.json');
    fs.writeFileSync(asyncApiOutputPath, JSON.stringify(asyncApiSpec, null, 2));
    console.log(`Generated AsyncAPI spec with ${Object.keys(asyncApiSpec.channels).length} channels in asyncapi.json`);
  } catch (error) {
    console.error('Error generating asyncapi.json:', error.message);
  }

  // Generate WebSocket markdown docs
  try {
    const wsMarkdown = generateWebSocketMarkdown(allRoutes);
    const wsOutputPath = path.join(markdownOutputDir, 'websocket.md');
    fs.writeFileSync(wsOutputPath, wsMarkdown);
    console.log('Generated WebSocket documentation in websocket.md');
  } catch (error) {
    console.error('Error generating websocket.md:', error.message);
  }

  console.log('API documentation generation complete!');
}

/**
 * Generate an AsyncAPI 2.6.0 spec from the extracted route data.
 */
function generateAsyncApiSpec(allRoutes, packageJson) {
  // Type mapping helper (reused from OpenAPI)
  function mapType(type) {
    switch (type) {
      case 'string': return { type: 'string' };
      case 'number': return { type: 'number' };
      case 'boolean': return { type: 'boolean' };
      case 'array': return { type: 'array', items: { type: 'string' } };
      case 'object': return { type: 'object' };
      default: return { type: 'string' };
    }
  }

  const spec = {
    asyncapi: '2.6.0',
    info: {
      title: 'FoundryVTT WebSocket API',
      description: 'Client-facing WebSocket API for bidirectional communication with Foundry VTT through the relay server. Connect to /ws/api with a token and clientId to send requests and receive real-time events.',
      version: packageJson.version || '1.0.0',
      license: { name: 'MIT' }
    },
    servers: {
      production: {
        url: 'ws://localhost:3010/ws/api',
        protocol: 'ws',
        description: 'Replaced dynamically at /asyncapi.json',
        variables: {
          token: { description: 'API key for authentication' },
          clientId: { description: 'ID of the connected Foundry instance to target' }
        }
      }
    },
    channels: {},
    components: {
      schemas: {
        ErrorMessage: {
          type: 'object',
          properties: {
            type: { type: 'string', const: 'error' },
            error: { type: 'string' },
            requestId: { type: 'string' }
          },
          required: ['type', 'error']
        },
        ConnectedMessage: {
          type: 'object',
          properties: {
            type: { type: 'string', const: 'connected' },
            clientId: { type: 'string' },
            supportedTypes: { type: 'array', items: { type: 'string' } },
            eventChannels: { type: 'array', items: { type: 'string' } }
          }
        }
      }
    }
  };

  // Deduplicate routes by type (many REST endpoints share the same WS message type)
  const seenTypes = new Set();

  for (const route of allRoutes) {
    const msgType = route.type;
    if (!msgType || seenTypes.has(msgType)) continue;
    seenTypes.add(msgType);

    // Build request properties from params (exclude clientId — it's in the connection URL)
    const requestProps = {
      type: { type: 'string', const: msgType, description: 'Message type' },
      requestId: { type: 'string', description: 'Unique request ID for correlation' }
    };
    const requiredFields = ['type', 'requestId'];

    const allParams = [
      ...(route.requiredParams || []),
      ...(route.optionalParams || [])
    ].filter(p => p.name !== 'clientId');

    for (const param of allParams) {
      requestProps[param.name] = {
        ...mapType(param.type),
        description: param.description || ''
      };
    }

    // Add required params (minus clientId)
    for (const p of (route.requiredParams || [])) {
      if (p.name !== 'clientId') {
        requiredFields.push(p.name);
      }
    }

    // Response schema
    const responseProps = {
      type: { type: 'string', const: `${msgType}-result`, description: 'Response message type' },
      requestId: { type: 'string', description: 'Echoed request ID' },
      clientId: { type: 'string', description: 'Foundry client ID' },
      error: { type: 'string', description: 'Error message if request failed' }
    };

    const channelName = `request/${msgType}`;
    spec.channels[channelName] = {
      description: route.description || `${msgType} request/response`,
      publish: {
        operationId: `send_${msgType}`,
        summary: `Send ${msgType} request`,
        message: {
          name: msgType,
          summary: route.description || '',
          payload: {
            type: 'object',
            properties: requestProps,
            required: requiredFields
          }
        }
      },
      subscribe: {
        operationId: `receive_${msgType}_result`,
        summary: `Receive ${msgType} response`,
        message: {
          name: `${msgType}-result`,
          payload: {
            type: 'object',
            properties: responseProps,
            required: ['type', 'requestId']
          }
        }
      }
    };
  }

  // Event channels (subscribe-only)
  spec.channels['events/chat'] = {
    description: 'Real-time chat message events from Foundry. Subscribe with { type: "subscribe", channel: "chat-events" }.',
    subscribe: {
      operationId: 'receive_chat_event',
      summary: 'Receive chat events',
      message: {
        name: 'chat-event',
        payload: {
          type: 'object',
          properties: {
            type: { type: 'string', const: 'chat-event' },
            event: { type: 'string', enum: ['create', 'update', 'delete'], description: 'Event type' },
            data: { type: 'object', description: 'Chat message data' }
          },
          required: ['type', 'event', 'data']
        }
      }
    }
  };

  spec.channels['events/roll'] = {
    description: 'Real-time dice roll events from Foundry. Subscribe with { type: "subscribe", channel: "roll-events" }.',
    subscribe: {
      operationId: 'receive_roll_event',
      summary: 'Receive roll events',
      message: {
        name: 'roll-event',
        payload: {
          type: 'object',
          properties: {
            type: { type: 'string', const: 'roll-event' },
            data: { type: 'object', description: 'Roll data' }
          },
          required: ['type', 'data']
        }
      }
    }
  };

  // Subscribe/unsubscribe control channels
  spec.channels['control/subscribe'] = {
    description: 'Subscribe to event channels',
    publish: {
      operationId: 'subscribe',
      summary: 'Subscribe to an event channel',
      message: {
        name: 'subscribe',
        payload: {
          type: 'object',
          properties: {
            type: { type: 'string', const: 'subscribe' },
            channel: { type: 'string', enum: ['chat-events', 'roll-events'] },
            requestId: { type: 'string' },
            filters: {
              type: 'object',
              properties: {
                speaker: { type: 'string', description: 'Filter by speaker alias' },
                type: { type: 'number', description: 'Filter by chat message type' },
                whisperOnly: { type: 'boolean', description: 'Only receive whispered messages' },
                userId: { type: 'string', description: 'Filter by Foundry user ID or username' }
              }
            }
          },
          required: ['type', 'channel']
        }
      }
    }
  };

  spec.channels['control/unsubscribe'] = {
    description: 'Unsubscribe from event channels',
    publish: {
      operationId: 'unsubscribe',
      summary: 'Unsubscribe from an event channel',
      message: {
        name: 'unsubscribe',
        payload: {
          type: 'object',
          properties: {
            type: { type: 'string', const: 'unsubscribe' },
            channel: { type: 'string', enum: ['chat-events', 'roll-events'], description: 'Omit to unsubscribe from all channels' },
            requestId: { type: 'string' }
          },
          required: ['type']
        }
      }
    }
  };

  return spec;
}

/**
 * Generate WebSocket documentation markdown
 */
function generateWebSocketMarkdown(allRoutes) {
  let md = `---\ntag: WebSocket\n---\n\n`;
  md += `import Tabs from '@theme/Tabs';\nimport TabItem from '@theme/TabItem';\n\n`;
  md += `# WebSocket API\n\n`;
  md += `The WebSocket API provides bidirectional communication with Foundry VTT through the relay server. `;
  md += `It supports the same operations as the REST API, plus real-time event subscriptions.\n\n`;

  // Connection
  md += `## Connection\n\n`;
  md += `Connect to the WebSocket endpoint with your API key and target Foundry client ID:\n\n`;
  md += '```\n';
  md += `ws://<host>/ws/api?token=<apiKey>&clientId=<clientId>\n`;
  md += '```\n\n';
  md += `On successful connection, you will receive a \`connected\` message listing all supported message types and event channels.\n\n`;

  // Message format
  md += `## Message Format\n\n`;
  md += `All messages are JSON objects with a \`type\` field. Request messages must also include a \`requestId\` for correlation.\n\n`;
  md += `### Request\n\n`;
  md += '```json\n';
  md += `{\n  "type": "search",\n  "requestId": "my-unique-id",\n  "query": "dragon"\n}\n`;
  md += '```\n\n';
  md += `### Response\n\n`;
  md += '```json\n';
  md += `{\n  "type": "search-result",\n  "requestId": "my-unique-id",\n  "clientId": "abc123",\n  "results": [...]\n}\n`;
  md += '```\n\n';

  // Event subscriptions
  md += `## Event Subscriptions\n\n`;
  md += `Subscribe to real-time events from Foundry:\n\n`;
  md += '```json\n';
  md += `{\n  "type": "subscribe",\n  "channel": "chat-events",\n  "filters": { "speaker": "GM" }\n}\n`;
  md += '```\n\n';
  md += `Available channels: \`chat-events\`, \`roll-events\`\n\n`;
  md += `To unsubscribe:\n\n`;
  md += '```json\n';
  md += `{\n  "type": "unsubscribe",\n  "channel": "chat-events"\n}\n`;
  md += '```\n\n';

  // Supported message types
  md += `## Supported Message Types\n\n`;

  // Deduplicate by type
  const seenTypes = new Set();

  md += `| Type | Description | Required Params |\n`;
  md += `|------|-------------|-----------------|\n`;

  for (const route of allRoutes) {
    const msgType = route.type;
    if (!msgType || seenTypes.has(msgType)) continue;
    seenTypes.add(msgType);

    const requiredParams = (route.requiredParams || [])
      .filter(p => p.name !== 'clientId')
      .map(p => `\`${p.name}\``)
      .join(', ');

    md += `| \`${msgType}\` | ${(route.description || '').replace(/\|/g, '\\|').substring(0, 80)} | ${requiredParams || '—'} |\n`;
  }

  md += `\n`;

  // AsyncAPI spec link
  md += `## AsyncAPI Spec\n\n`;
  md += `The full AsyncAPI specification is available at [\`/asyncapi.json\`](/asyncapi.json).\n\n`;

  return md;
}

main();
