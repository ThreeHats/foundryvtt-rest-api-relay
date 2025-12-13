/**
 * Convert captured test examples to Docusaurus markdown
 * This script reads captured examples from Jest tests and merges them with existing API docs
 */

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

/**
 * Sanitize sensitive values in examples for documentation
 * Replaces real API keys and client IDs with placeholders
 * This is the single point of sanitization - captureExample stores raw values
 */
function sanitizeExampleForDocs(example: CapturedExample): CapturedExample {
  // Deep clone to avoid mutating original
  const sanitized = JSON.parse(JSON.stringify(example));
  
  // Helper to sanitize a string
  const sanitizeString = (str: string): string => {
    return str
      // Replace 32-char hex API keys
      .replace(/[a-f0-9]{32}/gi, 'your-api-key-here')
      // Replace client IDs
      .replace(/foundry-[a-zA-Z0-9]{16}/g, 'your-client-id');
  };
  
  // Sanitize request URL
  if (sanitized.request?.url) {
    sanitized.request.url = sanitizeString(sanitized.request.url);
  }
  
  // Sanitize request headers (especially x-api-key)
  if (sanitized.request?.headers) {
    for (const key of Object.keys(sanitized.request.headers)) {
      if (key.toLowerCase() === 'x-api-key') {
        sanitized.request.headers[key] = 'your-api-key-here';
      } else {
        sanitized.request.headers[key] = sanitizeString(sanitized.request.headers[key]);
      }
    }
  }
  
  // Sanitize response data - replace client IDs and API keys
  if (sanitized.response?.data) {
    const dataStr = JSON.stringify(sanitized.response.data);
    sanitized.response.data = JSON.parse(sanitizeString(dataStr));
  }
  
  // Sanitize code examples - replace API keys and client IDs
  if (sanitized.codeExamples) {
    for (const lang of Object.keys(sanitized.codeExamples) as Array<keyof typeof sanitized.codeExamples>) {
      if (sanitized.codeExamples[lang]) {
        sanitized.codeExamples[lang] = sanitizeString(sanitized.codeExamples[lang]);
      }
    }
  }
  
  return sanitized;
}

/**
 * Generate markdown for a single code example
 */
function generateSingleExample(example: CapturedExample): string {
  let md = '\n### Code Examples\n\n';

  // Add tabs for different languages
  md += '<Tabs groupId="programming-language">\n';
  
  // JavaScript tab
  md += '<TabItem value="javascript" label="JavaScript">\n\n';
  md += '```javascript\n';
  md += example.codeExamples.javascript;
  md += '\n```\n\n';
  md += '</TabItem>\n';

  // cURL tab
  md += '<TabItem value="curl" label="cURL">\n\n';
  md += '```bash\n';
  md += example.codeExamples.curl;
  md += '\n```\n\n';
  md += '</TabItem>\n';

  // Python tab
  md += '<TabItem value="python" label="Python">\n\n';
  md += '```python\n';
  md += example.codeExamples.python;
  md += '\n```\n\n';
  md += '</TabItem>\n';

  // TypeScript tab
  md += '<TabItem value="typescript" label="TypeScript">\n\n';
  md += '```typescript\n';
  md += example.codeExamples.typescript;
  md += '\n```\n\n';
  md += '</TabItem>\n';

  // Emojicode tab
  md += '<TabItem value="emojicode" label="Emojicode">\n\n';
  md += '```emojicode\n';
  md += example.codeExamples.emojicode;
  md += '\n```\n\n';
  md += '</TabItem>\n';
  md += '</Tabs>\n\n';

  // Add response section
  md += '#### Response\n\n';
  md += `**Status:** ${example.response.status}\n\n`;
  md += '```json\n';
  md += JSON.stringify(example.response.data, null, 2);
  md += '\n```\n\n';

  return md;
}

/**
 * Update markdown file with code examples
 */
function updateMarkdownWithExamples(mdPath: string, examples: CapturedExample[]): void {
  if (!fs.existsSync(mdPath)) {
    console.warn(`Markdown file not found: ${mdPath}`);
    return;
  }

  let content = fs.readFileSync(mdPath, 'utf8');

  // Check if we need to add imports for tabs at the top
  if (!content.includes('import Tabs') && examples.length > 0) {
    // Add after the frontmatter - look for the closing --- of frontmatter
    const lines = content.split('\n');
    let frontmatterEndLine = -1;
    
    if (lines[0] === '---') {
      // Find the closing ---
      for (let i = 1; i < lines.length; i++) {
        if (lines[i] === '---') {
          frontmatterEndLine = i;
          break;
        }
      }
    }
    
    if (frontmatterEndLine !== -1) {
      const imports = [
        "import Tabs from '@theme/Tabs';",
        "import TabItem from '@theme/TabItem';",
        ''
      ];
      lines.splice(frontmatterEndLine + 1, 0, ...imports);
      content = lines.join('\n');
    }
  }

  // Remove ALL existing code examples sections
  content = content.replace(/\n### Code Examples\n[\s\S]*?(?=\n## |$)/g, '');

  // For each example, find the corresponding endpoint and insert the example
  for (const example of examples) {
    // Use description for matching if it contains parameter placeholders (e.g., ":uuid"),
    // otherwise fall back to converting the endpoint to a pattern
    let matchEndpoint = example.description;
    
    // If description doesn't look like a path pattern, convert the endpoint
    if (!matchEndpoint.startsWith('/')) {
      matchEndpoint = example.endpoint
        // Convert UUIDs like "Macro.3MHryjHf5ahMAGVU" to ":uuid"
        .replace(/\/[A-Z][a-z]+\.[a-zA-Z0-9]{16}/g, '/:uuid')
        // Convert other ID-like segments (16+ char alphanumeric) to ":id"
        .replace(/\/[a-zA-Z0-9]{16,}/g, '/:id');
    }
    
    const escapedEndpoint = matchEndpoint.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    
    // Match the entire endpoint section - from ## METHOD /endpoint to just before the next --- or ## or EOF
    const endpointSectionRegex = new RegExp(
      `(## ${example.method} ${escapedEndpoint}\\n[\\s\\S]*?)(?=\\n---|\\n## |$)`,
      'i'
    );

    const match = content.match(endpointSectionRegex);
    if (match) {
      const endpointSection = match[1];
      
      // Check if this section already has code examples
      if (endpointSection.includes('### Code Examples')) {
        // Skip - already has examples (shouldn't happen since we removed them all)
        continue;
      }
      
      // Find the best insertion point:
      // 1. After ### Returns section if it exists
      // 2. After ### Parameters table if no Returns
      let insertionPoint = endpointSection.length;
      
      // Try to find ### Returns section
      const returnsMatch = endpointSection.match(/### Returns\n\n\*\*[^\n]*?\n/);
      if (returnsMatch) {
        insertionPoint = returnsMatch.index! + returnsMatch[0].length;
      } else {
        // No Returns section, try to find end of Parameters table
        const paramsMatch = endpointSection.match(/### Parameters[\s\S]*?\n\n/);
        if (paramsMatch) {
          insertionPoint = paramsMatch.index! + paramsMatch[0].length;
        }
      }
      
      // Insert the example at the insertion point
      const exampleSection = generateSingleExample(example);
      const before = endpointSection.substring(0, insertionPoint);
      const after = endpointSection.substring(insertionPoint);
      const newEndpointSection = before + exampleSection + after;
      
      content = content.replace(endpointSectionRegex, newEndpointSection);
    } else {
      console.warn(`  Warning: Could not find endpoint ${example.method} ${matchEndpoint} in markdown`);
    }
  }

  fs.writeFileSync(mdPath, content);
  console.log(`Updated ${mdPath} with ${examples.length} examples`);
}

/**
 * Main execution
 */
async function main() {
  const examplesDir = path.join(__dirname, '../docs/examples');
  const docsDir = path.join(__dirname, '../docs/md/api');

  if (!fs.existsSync(examplesDir)) {
    console.error('Examples directory not found. Run tests first to capture examples.');
    process.exit(1);
  }

  // Read all captured example files
  const exampleFiles = fs.readdirSync(examplesDir)
    .filter(file => file.endsWith('-examples.json'));

  console.log(`Found ${exampleFiles.length} example files`);

  for (const file of exampleFiles) {
    const examplesPath = path.join(examplesDir, file);
    const examples: CapturedExample[] = JSON.parse(fs.readFileSync(examplesPath, 'utf8'));

    // Determine the corresponding markdown file
    // e.g., roll-endpoints-examples.json -> roll.md
    // e.g., roll-examples.json -> roll.md
    const baseName = file.replace('-endpoints-examples.json', '').replace('-examples.json', '');
    const mdPath = path.join(docsDir, `${baseName}.md`);

    console.log(`Processing ${file} (${examples.length} examples) -> ${baseName}.md`);
    
    if (!fs.existsSync(mdPath)) {
      console.warn(`  Warning: ${baseName}.md not found, skipping`);
      continue;
    }
    
    // Sanitize examples before updating docs
    const sanitizedExamples = examples.map(sanitizeExampleForDocs);
    updateMarkdownWithExamples(mdPath, sanitizedExamples);
  }

  console.log('\nDocumentation update complete!');
  console.log('\nNext steps:');
  console.log('1. Review the updated markdown files in docs/md/api/');
  console.log('2. Run: cd docs && pnpm run build');
  console.log('3. Preview: cd docs && pnpm run start');
}

main().catch(console.error);
