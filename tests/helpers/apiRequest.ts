/**
 * API Request Helper for Tests
 * Provides a clean interface for making API requests with full visibility of request structure
 * Easy to convert to documentation examples in different languages
 */

import axios, { AxiosRequestConfig, Method } from 'axios';

export interface ApiRequestConfig {
  url: {
    raw: string;
    host: string[];
    path: string[];
    query?: Array<{
      key: string;
      value: string;
      description?: string;
      disabled?: boolean;
    }>;
  };
  method: Method;
  header?: Array<{
    key: string;
    value: string;
    type?: string;
  }>;
  body?: {
    mode: 'raw';
    raw?: string;
  };
  data?: any; // Alternative to body for direct data
}

export interface ApiResponse<T = any> {
  data: T;
  status: number;
  statusText: string;
  headers: any;
  config: AxiosRequestConfig;
}

/**
 * Make an API request using the a config structure
 */
export async function makeRequest(config: ApiRequestConfig): Promise<ApiResponse> {
  // Use the raw URL which should already have variables replaced
  let fullUrl = config.url.raw;
  
  // If config has query params, append them properly
  if (config.url.query && config.url.query.length > 0) {
    try {
      const url = new URL(fullUrl);
      const activeParams = config.url.query.filter(q => !q.disabled);
      
      // Only modify URL if there are active params
      if (activeParams.length > 0) {
        activeParams.forEach(q => {
          url.searchParams.append(q.key, q.value);
        });
        fullUrl = url.toString();
      }
    } catch (error) {
      // If URL parsing fails (e.g., relative URL), fall back to manual construction
      const activeParams = config.url.query.filter(q => !q.disabled);
      if (activeParams.length > 0) {
        const separator = fullUrl.includes('?') ? '&' : '?';
        const queryParams = activeParams
          .map(q => `${encodeURIComponent(q.key)}=${encodeURIComponent(q.value)}`)
          .join('&');
        fullUrl = `${fullUrl}${separator}${queryParams}`;
      }
    }
  }

  // Build headers
  const headers: Record<string, string> = {};
  if (config.header) {
    config.header.forEach(h => {
      headers[h.key] = h.value;
    });
  }

  // Build request body
  let data: any = undefined;
  if (config.body?.mode === 'raw' && config.body.raw) {
    try {
      data = JSON.parse(config.body.raw);
    } catch {
      data = config.body.raw;
    }
  } else if (config.data) {
    data = config.data;
  }

  // Make request
  const axiosConfig: AxiosRequestConfig = {
    url: fullUrl,
    method: config.method,
    headers,
    data,
    timeout: 120000,
    validateStatus: () => true // Don't throw on any status code
  };

  const response = await axios.request(axiosConfig);

  return {
    data: response.data,
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
    config: response.config
  };
}

/**
 * Escape regex metacharacters in a string
 */
function escapeRegex(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Replace variables in a config object
 */
export function replaceVariables(
  config: ApiRequestConfig, 
  variables: Record<string, string>
): ApiRequestConfig {
  const configStr = JSON.stringify(config);
  let replaced = configStr;
  
  for (const [key, value] of Object.entries(variables)) {
    const escapedKey = escapeRegex(key);
    replaced = replaced.replace(new RegExp(`{{${escapedKey}}}`, 'g'), value);
  }
  
  return JSON.parse(replaced);
}
