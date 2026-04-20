import WebSocket from 'ws';
import { EventEmitter } from 'events';

interface WsClientOptions {
  timeout?: number;
}

export class WsTestClient extends EventEmitter {
  private ws: WebSocket | null = null;
  private responseHandlers = new Map<string, { resolve: (data: any) => void; reject: (err: Error) => void }>();
  private eventHandlers = new Map<string, ((data: any) => void)[]>();

  /**
   * Connect to the client-facing WebSocket endpoint
   */
  async connect(url: string, token: string, clientId: string): Promise<any> {
    // clientId is non-secret and may be passed in the URL; token must NOT be in the URL
    const fullUrl = clientId ? `${url}?clientId=${encodeURIComponent(clientId)}` : url;

    return new Promise<any>((resolve, reject) => {
      let isConnected = false;

      const timeout = setTimeout(() => {
        reject(new Error('Connection timeout'));
      }, 10000);

      this.ws = new WebSocket(fullUrl);

      this.ws.on('open', () => {
        // Send auth message — token must not appear in URL params
        this.ws!.send(JSON.stringify({ type: 'auth', token }));
      });

      this.ws.on('message', (raw: Buffer | string) => {
        try {
          const data = JSON.parse(typeof raw === 'string' ? raw : raw.toString('utf8'));

          // Handle initial connected message
          if (data.type === 'connected') {
            isConnected = true;
            clearTimeout(timeout);
            resolve(data);
          }

          // Route response to pending request
          if (data.requestId && this.responseHandlers.has(data.requestId)) {
            const handler = this.responseHandlers.get(data.requestId)!;
            this.responseHandlers.delete(data.requestId);
            handler.resolve(data);
          }

          // Route events to registered handlers
          const eventTypes = ['chat-event', 'roll-event', 'hook-event', 'combat-event', 'actor-event', 'scene-event'];
          if (eventTypes.includes(data.type)) {
            const handlers = this.eventHandlers.get(data.type) || [];
            for (const handler of handlers) {
              handler(data);
            }
          }

          // Also emit for subscribed/unsubscribed confirmations
          if (data.type === 'subscribed' || data.type === 'unsubscribed') {
            if (data.requestId && this.responseHandlers.has(data.requestId)) {
              const handler = this.responseHandlers.get(data.requestId)!;
              this.responseHandlers.delete(data.requestId);
              handler.resolve(data);
            }
          }

          // Error messages
          if (data.type === 'error') {
            if (data.requestId && this.responseHandlers.has(data.requestId)) {
              const handler = this.responseHandlers.get(data.requestId)!;
              this.responseHandlers.delete(data.requestId);
              handler.resolve(data); // Resolve (not reject) so tests can inspect the error
            }
          }

          // Emit raw message for advanced usage
          this.emit('message', data);
        } catch {
          // Ignore parse errors
        }
      });

      this.ws.on('error', (err) => {
        clearTimeout(timeout);
        reject(err);
      });

      this.ws.on('close', (code, reason) => {
        clearTimeout(timeout);
        // Reject connect() promise if auth hasn't completed yet
        if (!isConnected) {
          reject(new Error(`WebSocket closed before auth completed (code: ${code}, reason: ${reason})`));
        }
        // Reject all pending requests
        for (const [, handler] of this.responseHandlers) {
          handler.reject(new Error(`WebSocket closed (code: ${code}, reason: ${reason})`));
        }
        this.responseHandlers.clear();
      });
    });
  }

  /**
   * Send a message and wait for the matching response
   */
  async sendAndWait(message: Record<string, any>, timeout = 15000): Promise<any> {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket is not connected');
    }

    const requestId = message.requestId || `test_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
    const msg = { ...message, requestId };

    return new Promise<any>((resolve, reject) => {
      const timer = setTimeout(() => {
        this.responseHandlers.delete(requestId);
        reject(new Error(`Request ${requestId} timed out after ${timeout}ms`));
      }, timeout);

      this.responseHandlers.set(requestId, {
        resolve: (data) => {
          clearTimeout(timer);
          resolve(data);
        },
        reject: (err) => {
          clearTimeout(timer);
          reject(err);
        },
      });

      this.ws!.send(JSON.stringify(msg));
    });
  }

  /**
   * Subscribe to an event channel
   */
  async subscribe(
    channel: 'chat-events' | 'roll-events' | 'hooks' | 'combat-events' | 'actor-events' | 'scene-events',
    filters: Record<string, any> = {}
  ): Promise<any> {
    return this.sendAndWait({ type: 'subscribe', channel, filters });
  }

  /**
   * Unsubscribe from an event channel
   */
  async unsubscribe(channel?: string): Promise<any> {
    return this.sendAndWait({ type: 'unsubscribe', channel });
  }

  /**
   * Register a handler for event messages
   */
  onEvent(
    eventType: 'chat-event' | 'roll-event' | 'hook-event' | 'combat-event' | 'actor-event' | 'scene-event',
    handler: (data: any) => void
  ): void {
    if (!this.eventHandlers.has(eventType)) {
      this.eventHandlers.set(eventType, []);
    }
    this.eventHandlers.get(eventType)!.push(handler);
  }

  /**
   * Send a raw message without waiting for response
   */
  send(message: Record<string, any>): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket is not connected');
    }
    this.ws.send(JSON.stringify(message));
  }

  /**
   * Close the connection
   */
  close(): void {
    if (this.ws) {
      this.ws.close(1000, 'test complete');
      this.ws = null;
    }
    this.responseHandlers.clear();
    this.eventHandlers.clear();
  }

  get readyState(): number | undefined {
    return this.ws?.readyState;
  }
}
