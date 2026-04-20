import React, { useState, useCallback, useRef } from 'react';
import InteractiveJson from './InteractiveJson';

interface WsTesterProps {
  channels?: string[];
}

interface WsMessage {
  id: number;
  direction: 'sent' | 'received';
  data: any;
  time: string;
}

const CHANNELS = ['chat-events', 'roll-events', 'hooks', 'combat-events', 'actor-events', 'scene-events'];

export default function WsTester({ channels = CHANNELS }: WsTesterProps) {
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
  const [clientId, setClientId] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('apiTester_clientId') || '';
    }
    return '';
  });
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState<WsMessage[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [expandedMsg, setExpandedMsg] = useState<number | null>(null);
  const [subscriptions, setSubscriptions] = useState<Set<string>>(new Set());
  const [customMsg, setCustomMsg] = useState('{\n  "type": "",\n  "requestId": "test-1"\n}');
  const wsRef = useRef<WebSocket | null>(null);
  const msgIdRef = useRef(0);

  const handleServerUrlChange = useCallback((value: string) => {
    setServerUrl(value);
    if (typeof window !== 'undefined') localStorage.setItem('apiTester_serverUrl', value);
  }, []);

  const handleApiKeyChange = useCallback((value: string) => {
    setApiKey(value);
    if (typeof window !== 'undefined') localStorage.setItem('apiTester_apiKey', value);
  }, []);

  const handleClientIdChange = useCallback((value: string) => {
    setClientId(value);
    if (typeof window !== 'undefined') localStorage.setItem('apiTester_clientId', value);
  }, []);

  const addMessage = useCallback((direction: 'sent' | 'received', data: any) => {
    const id = ++msgIdRef.current;
    setMessages(prev => {
      const next = [...prev, { id, direction, data, time: new Date().toLocaleTimeString() }];
      return next.length > 200 ? next.slice(-200) : next;
    });
  }, []);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setConnected(false);
    setSubscriptions(new Set());
  }, []);

  const connect = useCallback(() => {
    disconnect();
    setError(null);
    setMessages([]);
    msgIdRef.current = 0;

    // Token must NOT appear in the URL — it would be logged by every proxy and
    // server in the path. Connect without credentials, then send them as the
    // first WebSocket message (auth-via-first-message protocol).
    const base = serverUrl.replace(/\/$/, '').replace(/^http/, 'ws');
    const url = `${base}/ws/api`;

    try {
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        // Send credentials as the first message — never in the URL.
        ws.send(JSON.stringify({ type: 'auth', token: apiKey, clientId }));
      };

      ws.onmessage = (evt) => {
        let parsed: any = evt.data;
        try { parsed = JSON.parse(evt.data); } catch {}

        // Server sends {"type":"connected"} after successful auth.
        if (parsed?.type === 'connected') {
          setConnected(true);
        }

        addMessage('received', parsed);
      };

      ws.onerror = () => setError('WebSocket error');

      ws.onclose = (evt) => {
        setConnected(false);
        if (evt.code !== 1000) {
          setError(`Closed: ${evt.reason || `code ${evt.code}`}`);
        }
      };
    } catch (err: any) {
      setError(err.message || 'Connection failed');
    }
  }, [serverUrl, apiKey, clientId, disconnect, addMessage]);

  const sendMessage = useCallback((msg: any) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return;
    const str = typeof msg === 'string' ? msg : JSON.stringify(msg);
    wsRef.current.send(str);
    let parsed = msg;
    if (typeof msg === 'string') { try { parsed = JSON.parse(msg); } catch {} }
    addMessage('sent', parsed);
  }, [addMessage]);

  const toggleSubscription = useCallback((channel: string) => {
    const isSubscribed = subscriptions.has(channel);
    const msg = {
      type: isSubscribed ? 'unsubscribe' : 'subscribe',
      channel,
      requestId: `${isSubscribed ? 'unsub' : 'sub'}-${channel}-${Date.now()}`,
    };
    sendMessage(msg);
    setSubscriptions(prev => {
      const next = new Set(prev);
      if (isSubscribed) next.delete(channel);
      else next.add(channel);
      return next;
    });
  }, [subscriptions, sendMessage]);

  const sendCustomMessage = useCallback(() => {
    try {
      const parsed = JSON.parse(customMsg);
      sendMessage(parsed);
    } catch {
      setError('Invalid JSON');
    }
  }, [customMsg, sendMessage]);

  return (
    <div className="api-tester">
      <button
        className="api-tester__toggle"
        onClick={() => setIsOpen(!isOpen)}
        type="button"
      >
        <span className="api-tester__toggle-icon">{isOpen ? '\u25BC' : '\u25B6'}</span>
        <span className="api-tester__toggle-method" style={{ backgroundColor: '#f59e0b' }}>WS</span>
        <span className="api-tester__toggle-path">/ws/api</span>
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
          <div className="api-tester__field">
            <label className="api-tester__label">Client ID</label>
            <input className="api-tester__input" type="text" value={clientId}
              onChange={(e) => handleClientIdChange(e.target.value)} placeholder="your-client-id" />
          </div>

          <div className="stream-tester__controls">
            {!connected ? (
              <button className="api-tester__send" onClick={connect} type="button">Connect</button>
            ) : (
              <button className="stream-tester__disconnect" onClick={disconnect} type="button">Disconnect</button>
            )}
            <span className={`stream-tester__status ${connected ? 'stream-tester__status--connected' : ''}`}>
              {connected ? 'Connected' : 'Disconnected'}
            </span>
            {messages.length > 0 && (
              <button className="stream-tester__clear" onClick={() => setMessages([])} type="button">Clear</button>
            )}
          </div>

          {error && <div className="api-tester__error">{error}</div>}

          {connected && (
            <>
              <div className="ws-tester__channels">
                <h4 className="api-tester__params-title">Event Subscriptions</h4>
                <div className="ws-tester__channel-grid">
                  {channels.map((ch) => (
                    <button
                      key={ch}
                      className={`ws-tester__channel-btn ${subscriptions.has(ch) ? 'ws-tester__channel-btn--active' : ''}`}
                      onClick={() => toggleSubscription(ch)}
                      type="button"
                    >
                      {subscriptions.has(ch) ? '\u2713 ' : ''}{ch}
                    </button>
                  ))}
                </div>
              </div>

              <div className="ws-tester__custom">
                <h4 className="api-tester__params-title">Send Message</h4>
                <textarea
                  className="ws-tester__textarea"
                  value={customMsg}
                  onChange={(e) => setCustomMsg(e.target.value)}
                  rows={4}
                  spellCheck={false}
                />
                <button className="api-tester__send" onClick={sendCustomMessage} type="button">Send</button>
              </div>
            </>
          )}

          {messages.length > 0 && (
            <div className="stream-tester__log">
              <div className="stream-tester__log-header">
                Messages ({messages.length})
              </div>
              <div className="stream-tester__log-body">
                {messages.map((msg) => (
                  <div key={msg.id} className={`stream-tester__event stream-tester__event--${msg.direction}`}>
                    <button
                      className="stream-tester__event-header"
                      onClick={() => setExpandedMsg(expandedMsg === msg.id ? null : msg.id)}
                      type="button"
                    >
                      <span className={`stream-tester__event-badge stream-tester__event-badge--${msg.direction}`}>
                        {msg.direction === 'sent' ? '\u2191' : '\u2193'}{' '}
                        {typeof msg.data === 'object' ? msg.data.type || msg.direction : msg.direction}
                      </span>
                      <span className="stream-tester__event-time">{msg.time}</span>
                      <span className="stream-tester__event-toggle">
                        {expandedMsg === msg.id ? '\u25BC' : '\u25B6'}
                      </span>
                    </button>
                    {expandedMsg === msg.id && (
                      <InteractiveJson data={msg.data} className="stream-tester__event-data" />
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
