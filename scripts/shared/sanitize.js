/**
 * Shared sanitization logic for API documentation examples.
 * Replaces real API keys, client IDs, and other sensitive values with placeholders.
 *
 * Used by both generateApiDocs.js (OpenAPI spec examples) and updateDocsWithExamples.ts (markdown docs).
 */

/**
 * Sanitize a string by replacing sensitive values with placeholders
 * @param {string} str - The string to sanitize
 * @returns {string} The sanitized string
 */
function sanitizeString(str) {
  return str
    // Replace 32+ char hex API keys
    .replace(/[a-f0-9]{32,}/gi, 'your-api-key-here')
    // Replace client IDs
    .replace(/foundry-[a-zA-Z0-9]{16}/g, 'your-client-id');
}

/**
 * Deep sanitize an object by replacing sensitive string values
 * @param {any} obj - The object to sanitize
 * @returns {any} A new sanitized object
 */
function sanitizeObject(obj) {
  if (obj === null || obj === undefined) return obj;
  if (typeof obj === 'string') return sanitizeString(obj);
  if (typeof obj !== 'object') return obj;

  const str = JSON.stringify(obj);
  return JSON.parse(sanitizeString(str));
}

/**
 * Sanitize a full captured example for documentation
 * @param {object} example - A captured example object
 * @returns {object} A sanitized copy of the example
 */
function sanitizeExampleForDocs(example) {
  // Deep clone to avoid mutating original
  const sanitized = JSON.parse(JSON.stringify(example));

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

  // Sanitize response data
  if (sanitized.response?.data) {
    sanitized.response.data = sanitizeObject(sanitized.response.data);
  }

  // Sanitize code examples
  if (sanitized.codeExamples) {
    for (const lang of Object.keys(sanitized.codeExamples)) {
      if (sanitized.codeExamples[lang]) {
        sanitized.codeExamples[lang] = sanitizeString(sanitized.codeExamples[lang]);
      }
    }
  }

  return sanitized;
}

module.exports = { sanitizeString, sanitizeObject, sanitizeExampleForDocs };
