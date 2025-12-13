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
  // If raw URL has query params, use it as-is; otherwise build from parts
  let fullUrl = config.url.raw;
  
  // If raw URL doesn't have query params but config does, append them
  if (!fullUrl.includes('?') && config.url.query && config.url.query.length > 0) {
    const queryParams = config.url.query
      ?.filter(q => !q.disabled)
      .map(q => `${q.key}=${encodeURIComponent(q.value)}`)
      .join('&');
    fullUrl = `${fullUrl}?${queryParams}`;
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
 * Replace variables in a config object
 */
export function replaceVariables(
  config: ApiRequestConfig, 
  variables: Record<string, string>
): ApiRequestConfig {
  const configStr = JSON.stringify(config);
  let replaced = configStr;
  
  for (const [key, value] of Object.entries(variables)) {
    replaced = replaced.replace(new RegExp(`{{${key}}}`, 'g'), value);
  }
  
  return JSON.parse(replaced);
}
