import React, { useState, useCallback } from 'react';

interface Parameter {
  name: string;
  type: string;
  required: boolean;
  source: string;
}

interface ApiTesterProps {
  method: string;
  path: string;
  parameters?: Parameter[];
}

const METHOD_COLORS: Record<string, string> = {
  GET: '#10b981',
  POST: '#3b82f6',
  PUT: '#f59e0b',
  DELETE: '#ef4444',
  PATCH: '#8b5cf6',
};

export default function ApiTester({ method, path, parameters = [] }: ApiTesterProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [serverUrl, setServerUrl] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('apiTester_serverUrl') || 'http://localhost:3010';
    }
    return 'http://localhost:3010';
  });
  const [apiKey, setApiKey] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('apiTester_apiKey') || '';
    }
    return '';
  });
  const [paramValues, setParamValues] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {};
    if (typeof window !== 'undefined') {
      const savedClientId = localStorage.getItem('apiTester_clientId');
      if (savedClientId && parameters.some(p => p.name === 'clientId')) {
        initial.clientId = savedClientId;
      }
    }
    return initial;
  });
  const [response, setResponse] = useState<{ status: number; data: any } | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleServerUrlChange = useCallback((value: string) => {
    setServerUrl(value);
    if (typeof window !== 'undefined') {
      localStorage.setItem('apiTester_serverUrl', value);
    }
  }, []);

  const handleApiKeyChange = useCallback((value: string) => {
    setApiKey(value);
    if (typeof window !== 'undefined') {
      localStorage.setItem('apiTester_apiKey', value);
    }
  }, []);

  const handleParamChange = useCallback((name: string, value: string) => {
    setParamValues(prev => ({ ...prev, [name]: value }));
    if (name === 'clientId' && typeof window !== 'undefined') {
      localStorage.setItem('apiTester_clientId', value);
    }
  }, []);

  const sendRequest = useCallback(async () => {
    setLoading(true);
    setError(null);
    setResponse(null);

    try {
      // Build URL with path params replaced
      let url = `${serverUrl.replace(/\/$/, '')}${path}`;
      for (const param of parameters) {
        if (param.source === 'params' || param.source === 'path') {
          url = url.replace(`{${param.name}}`, encodeURIComponent(paramValues[param.name] || ''));
          url = url.replace(`:${param.name}`, encodeURIComponent(paramValues[param.name] || ''));
        }
      }

      // Build query params
      const queryParams = parameters
        .filter(p => p.source === 'query' && paramValues[p.name])
        .map(p => `${encodeURIComponent(p.name)}=${encodeURIComponent(paramValues[p.name])}`)
        .join('&');

      if (queryParams) {
        url += `?${queryParams}`;
      }

      // Build body for POST/PUT/DELETE
      const bodyParams = parameters.filter(p => p.source === 'body' && paramValues[p.name]);
      let body: string | undefined;
      if (['POST', 'PUT', 'DELETE', 'PATCH'].includes(method) && bodyParams.length > 0) {
        const bodyObj: Record<string, any> = {};
        for (const p of bodyParams) {
          const val = paramValues[p.name];
          if (p.type === 'number') {
            bodyObj[p.name] = Number(val);
          } else if (p.type === 'boolean') {
            bodyObj[p.name] = val === 'true';
          } else if (p.type === 'object' || p.type === 'array') {
            try { bodyObj[p.name] = JSON.parse(val); } catch { bodyObj[p.name] = val; }
          } else {
            bodyObj[p.name] = val;
          }
        }
        body = JSON.stringify(bodyObj);
      }

      const headers: Record<string, string> = {};
      if (apiKey) {
        headers['x-api-key'] = apiKey;
      }
      if (body) {
        headers['Content-Type'] = 'application/json';
      }

      // Add header params
      for (const param of parameters) {
        if (param.source === 'header' && paramValues[param.name]) {
          headers[param.name] = paramValues[param.name];
        }
      }

      const fetchOptions: RequestInit = {
        method,
        headers,
      };
      if (body) {
        fetchOptions.body = body;
      }

      const res = await fetch(url, fetchOptions);
      let data: any;
      const contentType = res.headers.get('content-type') || '';
      if (contentType.includes('application/json')) {
        data = await res.json();
      } else {
        data = await res.text();
      }

      setResponse({ status: res.status, data });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  }, [serverUrl, apiKey, method, path, parameters, paramValues]);

  const methodColor = METHOD_COLORS[method] || '#6b7280';

  return (
    <div className="api-tester">
      <button
        className="api-tester__toggle"
        onClick={() => setIsOpen(!isOpen)}
        type="button"
      >
        <span className="api-tester__toggle-icon">{isOpen ? '\u25BC' : '\u25B6'}</span>
        <span className="api-tester__toggle-method" style={{ backgroundColor: methodColor }}>
          {method}
        </span>
        <span className="api-tester__toggle-path">{path}</span>
      </button>

      {isOpen && (
        <div className="api-tester__panel">
          <div className="api-tester__field">
            <label className="api-tester__label">Server URL</label>
            <input
              className="api-tester__input"
              type="text"
              value={serverUrl}
              onChange={(e) => handleServerUrlChange(e.target.value)}
              placeholder="http://localhost:3010"
            />
          </div>

          <div className="api-tester__field">
            <label className="api-tester__label">API Key</label>
            <input
              className="api-tester__input"
              type="password"
              value={apiKey}
              onChange={(e) => handleApiKeyChange(e.target.value)}
              placeholder="your-api-key"
            />
          </div>

          {parameters.length > 0 && (
            <div className="api-tester__params">
              <h4 className="api-tester__params-title">Parameters</h4>
              {parameters.map((param) => (
                <div key={param.name} className="api-tester__field">
                  <label className="api-tester__label">
                    {param.name}
                    <span className="api-tester__param-meta">
                      {' '}({param.type}, {param.source})
                      {param.required && <span className="api-tester__required">*</span>}
                    </span>
                  </label>
                  <input
                    className="api-tester__input"
                    type="text"
                    value={paramValues[param.name] || ''}
                    onChange={(e) => handleParamChange(param.name, e.target.value)}
                    placeholder={`Enter ${param.name}`}
                  />
                </div>
              ))}
            </div>
          )}

          <button
            className="api-tester__send"
            onClick={sendRequest}
            disabled={loading}
            type="button"
          >
            {loading ? 'Sending...' : 'Send Request'}
          </button>

          {error && (
            <div className="api-tester__error">
              {error}
            </div>
          )}

          {response && (
            <div className="api-tester__response">
              <div className="api-tester__response-header">
                <span className={`api-tester__status api-tester__status--${Math.floor(response.status / 100)}xx`}>
                  {response.status}
                </span>
                <span className="api-tester__response-label">Response</span>
              </div>
              <pre className="api-tester__response-body">
                {typeof response.data === 'string' ? response.data : JSON.stringify(response.data, null, 2)}
              </pre>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
