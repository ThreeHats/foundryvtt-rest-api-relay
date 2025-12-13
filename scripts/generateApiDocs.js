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
      const methodMatch = line.match(/(\w+[Rr]outer)\.(get|post|put|delete)\(\s*["']([^"']+)["']/);
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

function generateMarkdown(routes, moduleName) {
    let markdown = `---\n`;
    markdown += `tag: ${moduleName}\n`;
    markdown += `---\n\n`;
    markdown += `# ${moduleName}\n\n`;

    routes.forEach((route, index) => {
        markdown += `## ${route.method} ${route.path}\n\n`;
        if (route.description) {
            markdown += `${route.description}\n\n`;
        }
        
        // Parameters table
        const hasRequiredParams = route.requiredParams && route.requiredParams.length > 0;
        const hasOptionalParams = route.optionalParams && route.optionalParams.length > 0;
        
        if (hasRequiredParams || hasOptionalParams) {
            markdown += `### Parameters\n\n`;
            markdown += `| Name | Type | Required | Source | Description |\n`;
            markdown += `|------|------|----------|--------|--------------|\n`;
            
            // Required parameters first
            if (hasRequiredParams) {
                route.requiredParams.forEach(param => {
                    const source = param.from ? param.from.join(', ') : 'body, query';
                    markdown += `| ${param.name} | ${param.type} | ✓ | ${source} | ${param.description || ''} |\n`;
                });
            }
            
            // Then optional parameters
            if (hasOptionalParams) {
                route.optionalParams.forEach(param => {
                    const source = param.from ? param.from.join(', ') : 'body, query';
                    markdown += `| ${param.name} | ${param.type} |  | ${source} | ${param.description || ''} |\n`;
                });
            }
            
            markdown += '\n';
        }
        
        // Returns
        if (route.returns) {
            markdown += `### Returns\n\n`;
            markdown += `**${route.returns.type}** - ${route.returns.description}\n\n`;
        }
        
        // Add separator between routes, but not after the last one
        if (index < routes.length - 1) {
            markdown += '---\n\n';
        }
    });

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
  
  console.log('Markdown documentation generation complete!');

  // Generate the single api-docs.json file
  try {
    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
    const apiDocs = {
      version: packageJson.version || "1.0.0",
      baseUrl: `https://your-relay-server.com`, // This will be replaced by the server dynamically
      authentication: {
        required: true,
        headerName: "x-api-key",
        description: "API key must be included in the x-api-key header for all endpoints except /api/status and /api/docs"
      },
      endpoints: allRoutes.map(route => ({
          method: route.method.toUpperCase(),
          path: route.route, // Use the full path from @route JSDoc tag which includes any path prefixes (e.g., /dnd5e)
          description: route.description,
          requiredParameters: route.requiredParams.map(p => ({ name: p.name, type: p.type, description: p.description, location: p.from.join(', ') })),
          optionalParameters: route.optionalParams.map(p => ({ name: p.name, type: p.type, description: p.description, location: p.from.join(', ') })),
      }))
    };

    const outputPath = path.join(jsonOutputDir, `api-docs.json`);
    fs.writeFileSync(outputPath, JSON.stringify(apiDocs, null, 2));
    console.log(`Generated JSON documentation for ${allRoutes.length} routes in api-docs.json`);
  } catch (error) {
    console.error(`Error generating api-docs.json:`, error.message);
  }

  console.log('API documentation generation complete!');
}

main();
