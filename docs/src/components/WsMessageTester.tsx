/**
 * WsMessageTester — per-type try-it-out panel for the WebSocket API reference page.
 *
 * Uses a module-level singleton WebSocket so all WsMessageTester instances on the page
 * share one connection rather than opening 90+ individual connections.
 *
 * WsConnectionBar renders the connect/disconnect controls; it must appear once per page.
 */
import React, { useState, useCallback, useEffect, useRef } from 'react';
import InteractiveJson from './InteractiveJson';

// ---------------------------------------------------------------------------
// Module-level shared connection singleton
// ---------------------------------------------------------------------------

type ConnectionListener = (connected: boolean) => void;

let sharedWs: WebSocket | null = null;
let sharedConnected = false;
const connectionListeners = new Set<ConnectionListener>();
const responseWaiters = new Map<string, (data: any) => void>();

function notifyConnectionListeners(connected: boolean) {
  sharedConnected = connected;
  for (const cb of connectionListeners) cb(connected);
}

export function connectShared(serverUrl: string, apiKey: string, clientId: string) {
  if (sharedWs) {
    sharedWs.close();
    sharedWs = null;
  }
  // Reject any pending waiters
  for (const resolve of responseWaiters.values()) {
    resolve({ error: 'Connection reset' });
  }
  responseWaiters.clear();

  // Token must NOT appear in the URL (server logs, proxies). Use auth-via-first-message.
  const base = serverUrl.replace(/\/$/, '').replace(/^http/, 'ws');
  const url = `${base}/ws/api`;

  try {
    const ws = new WebSocket(url);
    sharedWs = ws;

    ws.onopen = () => {
      // Send credentials as the first WebSocket message — never in the URL.
      ws.send(JSON.stringify({ type: 'auth', token: apiKey, clientId }));
    };

    ws.onmessage = (evt) => {
      let data: any = evt.data;
      try { data = JSON.parse(evt.data); } catch {}
      // Server sends {"type":"connected"} after successful auth.
      if (data && data.type === 'connected') {
        notifyConnectionListeners(true);
      }
      if (data && typeof data === 'object' && data.requestId) {
        const waiter = responseWaiters.get(data.requestId);
        if (waiter) {
          responseWaiters.delete(data.requestId);
          waiter(data);
        }
      }
    };

    ws.onerror = () => notifyConnectionListeners(false);

    ws.onclose = () => {
      for (const resolve of responseWaiters.values()) {
        resolve({ error: 'Connection closed' });
      }
      responseWaiters.clear();
      notifyConnectionListeners(false);
    };
  } catch {
    notifyConnectionListeners(false);
  }
}

export function disconnectShared() {
  if (sharedWs) {
    sharedWs.close();
    sharedWs = null;
  }
  notifyConnectionListeners(false);
}

export function sendShared(msg: any): boolean {
  if (!sharedWs || sharedWs.readyState !== WebSocket.OPEN) return false;
  sharedWs.send(JSON.stringify(msg));
  return true;
}

export function waitForResponse(requestId: string, timeoutMs: number): Promise<any> {
  return new Promise((resolve) => {
    const timer = setTimeout(() => {
      responseWaiters.delete(requestId);
      resolve({ error: `Timed out after ${timeoutMs / 1000}s waiting for ${requestId}` });
    }, timeoutMs);

    responseWaiters.set(requestId, (data) => {
      clearTimeout(timer);
      resolve(data);
    });
  });
}

function useSharedConnection(): boolean {
  const [connected, setConnected] = useState(sharedConnected);

  useEffect(() => {
    const cb: ConnectionListener = (c) => setConnected(c);
    connectionListeners.add(cb);
    return () => { connectionListeners.delete(cb); };
  }, []);

  return connected;
}

// ---------------------------------------------------------------------------
// WsConnectionBar — renders once per page
// ---------------------------------------------------------------------------

export function WsConnectionBar() {
  const connected = useSharedConnection();

  const [serverUrl, setServerUrl] = useState(() =>
    typeof window !== 'undefined'
      ? localStorage.getItem('apiTester_serverUrl') || 'http://localhost:3010'
      : 'http://localhost:3010'
  );
  const [apiKey, setApiKey] = useState(() =>
    typeof window !== 'undefined' ? localStorage.getItem('apiTester_apiKey') || '' : ''
  );
  const [clientId, setClientId] = useState(() =>
    typeof window !== 'undefined' ? localStorage.getItem('apiTester_clientId') || '' : ''
  );

  const handleServerUrlChange = (v: string) => {
    setServerUrl(v);
    if (typeof window !== 'undefined') localStorage.setItem('apiTester_serverUrl', v);
  };
  const handleApiKeyChange = (v: string) => {
    setApiKey(v);
    if (typeof window !== 'undefined') localStorage.setItem('apiTester_apiKey', v);
  };
  const handleClientIdChange = (v: string) => {
    setClientId(v);
    if (typeof window !== 'undefined') localStorage.setItem('apiTester_clientId', v);
  };

  const connect = () => connectShared(serverUrl, apiKey, clientId);
  const disconnect = disconnectShared;

  return (
    <div className="api-tester ws-connection-bar">
      <div className="api-tester__panel" style={{ marginBottom: '1rem' }}>
        <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', alignItems: 'flex-end' }}>
          <div className="api-tester__field" style={{ flex: '2', minWidth: '160px' }}>
            <label className="api-tester__label">Server URL</label>
            <input
              className="api-tester__input"
              type="text"
              value={serverUrl}
              onChange={(e) => handleServerUrlChange(e.target.value)}
              placeholder="http://localhost:3010"
            />
          </div>
          <div className="api-tester__field" style={{ flex: '2', minWidth: '140px' }}>
            <label className="api-tester__label">API Key</label>
            <input
              className="api-tester__input"
              type="password"
              value={apiKey}
              onChange={(e) => handleApiKeyChange(e.target.value)}
              placeholder="your-api-key"
            />
          </div>
          <div className="api-tester__field" style={{ flex: '1', minWidth: '120px' }}>
            <label className="api-tester__label">Client ID</label>
            <input
              className="api-tester__input"
              type="text"
              value={clientId}
              onChange={(e) => handleClientIdChange(e.target.value)}
              placeholder="your-client-id"
            />
          </div>
          <div style={{ paddingBottom: '0.25rem' }}>
            {!connected ? (
              <button className="api-tester__send" onClick={connect} type="button">
                Connect
              </button>
            ) : (
              <button className="stream-tester__disconnect" onClick={disconnect} type="button">
                Disconnect
              </button>
            )}
            <span
              className={`stream-tester__status ${connected ? 'stream-tester__status--connected' : ''}`}
              style={{ marginLeft: '0.5rem' }}
            >
              {connected ? 'Connected' : 'Disconnected'}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// WsMessageTester — per-type try-it-out panel
// ---------------------------------------------------------------------------

export interface WsParameter {
  name: string;
  type: string;
  required: boolean;
  description?: string;
}

interface WsMessageTesterProps {
  messageType: string;
  resultType?: string;
  parameters?: WsParameter[];
}

function genRequestId(messageType: string): string {
  return `${messageType}_${Date.now()}_${Math.random().toString(36).slice(2, 7)}`;
}

export default function WsMessageTester({
  messageType,
  resultType,
  parameters = [],
}: WsMessageTesterProps) {
  const connected = useSharedConnection();
  const [isOpen, setIsOpen] = useState(false);
  const [paramValues, setParamValues] = useState<Record<string, string>>({});
  const [response, setResponse] = useState<any | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleParamChange = useCallback((name: string, value: string) => {
    setParamValues((prev) => ({ ...prev, [name]: value }));
  }, []);

  const sendMessage = useCallback(async () => {
    if (!connected) {
      setError('Not connected — use the connection bar above to connect first.');
      return;
    }
    setLoading(true);
    setError(null);
    setResponse(null);

    const requestId = genRequestId(messageType);
    const msg: Record<string, any> = { type: messageType, requestId };

    for (const param of parameters) {
      const raw = paramValues[param.name];
      if (raw === undefined || raw === '') continue;
      if (param.type === 'number') {
        msg[param.name] = Number(raw);
      } else if (param.type === 'boolean') {
        msg[param.name] = raw === 'true' || raw === '1';
      } else if (param.type === 'object' || param.type === 'array') {
        try {
          msg[param.name] = JSON.parse(raw);
        } catch {
          msg[param.name] = raw;
        }
      } else {
        msg[param.name] = raw;
      }
    }

    const sent = sendShared(msg);
    if (!sent) {
      setLoading(false);
      setError('Failed to send — connection may have dropped.');
      return;
    }

    const data = await waitForResponse(requestId, 60000);
    setResponse(data);
    setLoading(false);
  }, [connected, messageType, parameters, paramValues]);

  const effectiveResultType = resultType || `${messageType}-result`;

  return (
    <div className="api-tester">
      <button
        className="api-tester__toggle"
        onClick={() => setIsOpen(!isOpen)}
        type="button"
      >
        <span className="api-tester__toggle-icon">{isOpen ? '\u25BC' : '\u25B6'}</span>
        <span className="api-tester__toggle-method" style={{ backgroundColor: '#f59e0b' }}>WS</span>
        <span className="api-tester__toggle-path">
          {messageType} → {effectiveResultType}
        </span>
      </button>

      {isOpen && (
        <div className="api-tester__panel">
          {!connected && (
            <div className="api-tester__error" style={{ marginBottom: '0.75rem' }}>
              Connect using the bar above to try this message type.
            </div>
          )}

          {parameters.length > 0 && (
            <div className="api-tester__params">
              <h4 className="api-tester__params-title">Parameters</h4>
              {parameters.map((param) => (
                <div key={param.name} className="api-tester__field">
                  <label className="api-tester__label">
                    {param.name}
                    <span className="api-tester__param-meta">
                      {' '}({param.type}
                      {param.required && <span className="api-tester__required"> *</span>})
                    </span>
                    {param.description && (
                      <span className="api-tester__param-meta"> — {param.description}</span>
                    )}
                  </label>
                  {param.type === 'object' || param.type === 'array' ? (
                    <textarea
                      className="ws-tester__textarea"
                      value={paramValues[param.name] || ''}
                      onChange={(e) => handleParamChange(param.name, e.target.value)}
                      rows={3}
                      placeholder={param.type === 'object' ? '{ }' : '[ ]'}
                      spellCheck={false}
                    />
                  ) : param.type === 'boolean' ? (
                    <select
                      className="api-tester__input"
                      value={paramValues[param.name] || ''}
                      onChange={(e) => handleParamChange(param.name, e.target.value)}
                    >
                      <option value="">(not set)</option>
                      <option value="true">true</option>
                      <option value="false">false</option>
                    </select>
                  ) : (
                    <input
                      className="api-tester__input"
                      type={param.type === 'number' ? 'number' : 'text'}
                      value={paramValues[param.name] || ''}
                      onChange={(e) => handleParamChange(param.name, e.target.value)}
                      placeholder={`Enter ${param.name}`}
                    />
                  )}
                </div>
              ))}
            </div>
          )}

          <button
            className="api-tester__send"
            onClick={sendMessage}
            disabled={loading || !connected}
            type="button"
          >
            {loading ? 'Waiting for response...' : `Send ${messageType}`}
          </button>

          {error && <div className="api-tester__error">{error}</div>}

          {response && (
            <div className="api-tester__response">
              <div className="api-tester__response-header">
                <span
                  className={`api-tester__status api-tester__status--${
                    response.success === false ? '4xx' : '2xx'
                  }`}
                >
                  {response.success === false ? 'error' : 'ok'}
                </span>
                <span className="api-tester__response-label">Response</span>
              </div>
              <InteractiveJson data={response} className="api-tester__response-body" />
            </div>
          )}
        </div>
      )}
    </div>
  );
}
