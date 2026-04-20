import React, { useState, useCallback, useRef } from 'react';
import InteractiveJson from './InteractiveJson';

interface Parameter {
  name: string;
  type: string;
  required: boolean;
  source: string;
}

interface SseTesterProps {
  path: string;
  parameters?: Parameter[];
}

interface SseEvent {
  id: number;
  event: string;
  data: any;
  time: string;
}

export default function SseTester({ path, parameters = [] }: SseTesterProps) {
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
  const [connected, setConnected] = useState(false);
  const [events, setEvents] = useState<SseEvent[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [expandedEvent, setExpandedEvent] = useState<number | null>(null);
  const abortRef = useRef<AbortController | null>(null);
  const eventIdRef = useRef(0);

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

  const disconnect = useCallback(() => {
    if (abortRef.current) {
      abortRef.current.abort();
      abortRef.current = null;
    }
    setConnected(false);
  }, []);

  const connect = useCallback(async () => {
    disconnect();
    setError(null);
    setEvents([]);
    eventIdRef.current = 0;

    let url = `${serverUrl.replace(/\/$/, '')}${path}`;
    const queryParams = parameters
      .filter(p => p.source === 'query' && paramValues[p.name])
      .map(p => `${encodeURIComponent(p.name)}=${encodeURIComponent(paramValues[p.name])}`)
      .join('&');
    if (queryParams) url += `?${queryParams}`;

    const abort = new AbortController();
    abortRef.current = abort;

    try {
      const headers: Record<string, string> = { 'Accept': 'text/event-stream' };
      if (apiKey) headers['x-api-key'] = apiKey;

      const res = await fetch(url, { headers, signal: abort.signal });

      if (!res.ok) {
        let msg = `HTTP ${res.status}`;
        try { const body = await res.json(); msg = body.error || msg; } catch {}
        setError(msg);
        return;
      }

      setConnected(true);
      const reader = res.body?.getReader();
      if (!reader) { setError('No response body'); return; }

      const decoder = new TextDecoder();
      let buffer = '';
      let currentEvent = '';
      let currentData = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        buffer = lines.pop() || '';

        for (const line of lines) {
          if (line.startsWith('event: ')) {
            currentEvent = line.slice(7).trim();
          } else if (line.startsWith('data: ')) {
            currentData += (currentData ? '\n' : '') + line.slice(6);
          } else if (line === '' && (currentEvent || currentData)) {
            let parsed: any = currentData;
            try { parsed = JSON.parse(currentData); } catch {}
            const id = ++eventIdRef.current;
            setEvents(prev => {
              const next = [...prev, {
                id,
                event: currentEvent || 'message',
                data: parsed,
                time: new Date().toLocaleTimeString(),
              }];
              return next.length > 100 ? next.slice(-100) : next;
            });
            currentEvent = '';
            currentData = '';
          }
        }
      }
    } catch (err: any) {
      if (err.name !== 'AbortError') {
        setError(err.message || 'Connection failed');
      }
    } finally {
      setConnected(false);
    }
  }, [serverUrl, apiKey, path, parameters, paramValues, disconnect]);

  return (
    <div className="api-tester">
      <button
        className="api-tester__toggle"
        onClick={() => setIsOpen(!isOpen)}
        type="button"
      >
        <span className="api-tester__toggle-icon">{isOpen ? '\u25BC' : '\u25B6'}</span>
        <span className="api-tester__toggle-method" style={{ backgroundColor: '#8b5cf6' }}>SSE</span>
        <span className="api-tester__toggle-path">{path}</span>
      </button>

      {isOpen && (
        <div className="api-tester__panel">
          <div className="api-tester__field">
            <label className="api-tester__label">Server URL</label>
            <input className="api-tester__input" type="text" value={serverUrl}
              onChange={(e) => handleServerUrlChange(e.target.value)} placeholder="http://localhost:3010" />
          </div>
          <div className="api-tester__field">
            <label className="api-tester__label">API Key</label>
            <input className="api-tester__input" type="password" value={apiKey}
              onChange={(e) => handleApiKeyChange(e.target.value)} placeholder="your-api-key" />
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
                  <input className="api-tester__input" type="text"
                    value={paramValues[param.name] || ''}
                    onChange={(e) => handleParamChange(param.name, e.target.value)}
                    placeholder={`Enter ${param.name}`} />
                </div>
              ))}
            </div>
          )}

          <div className="stream-tester__controls">
            {!connected ? (
              <button className="api-tester__send" onClick={connect} type="button">Connect</button>
            ) : (
              <button className="stream-tester__disconnect" onClick={disconnect} type="button">Disconnect</button>
            )}
            <span className={`stream-tester__status ${connected ? 'stream-tester__status--connected' : ''}`}>
              {connected ? 'Connected' : 'Disconnected'}
            </span>
            {events.length > 0 && (
              <button className="stream-tester__clear" onClick={() => setEvents([])} type="button">Clear</button>
            )}
          </div>

          {error && <div className="api-tester__error">{error}</div>}

          {events.length > 0 && (
            <div className="stream-tester__log">
              <div className="stream-tester__log-header">
                Events ({events.length})
              </div>
              <div className="stream-tester__log-body">
                {events.map((evt) => (
                  <div key={evt.id} className="stream-tester__event">
                    <button
                      className="stream-tester__event-header"
                      onClick={() => setExpandedEvent(expandedEvent === evt.id ? null : evt.id)}
                      type="button"
                    >
                      <span className="stream-tester__event-badge">{evt.event}</span>
                      <span className="stream-tester__event-time">{evt.time}</span>
                      <span className="stream-tester__event-toggle">
                        {expandedEvent === evt.id ? '\u25BC' : '\u25B6'}
                      </span>
                    </button>
                    {expandedEvent === evt.id && (
                      <InteractiveJson data={evt.data} className="stream-tester__event-data" />
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
