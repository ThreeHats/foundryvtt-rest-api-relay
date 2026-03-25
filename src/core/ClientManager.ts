// src/core/ClientManager.ts
import { WebSocket } from "ws";
import { log } from "../utils/logger";
import { Client } from "./Client";
import { WSCloseCodes } from "../lib/constants";
import { getRedisClient } from "../config/redis";
import { SheetSessionManager } from "./SheetSessionManager";

type MessageHandler = (client: Client, message: any) => void;

const INSTANCE_ID = process.env.FLY_ALLOC_ID || 'local';
const CLIENT_EXPIRY = 60 * 60 * 2; // 2 hours expiry for Redis keys
const isMemoryStore = process.env.DB_TYPE === 'memory';

export class ClientManager {
  private static clients = new Map<string, Client>();
  private static tokenGroups = new Map<string, Set<string>>(); 
  private static messageHandlers = new Map<string, MessageHandler[]>();

  /**
   * Add a new client to the manager
   */
  static async addClient(
    ws: WebSocket, 
    id: string, 
    token: string, 
    worldId: string | null, 
    worldTitle: string | null,
    foundryVersion: string | null = null,
    systemId: string | null = null,
    systemTitle: string | null = null,
    systemVersion: string | null = null,
    customName: string | null = null
  ): Promise<Client | null> {
    // If a client with this ID is already connected, replace it.
    // This happens when a headless session takes over from a browser tab,
    // or when a module reconnects after a brief disconnect.
    if (this.clients.has(id)) {
      log.info(`Client ${id} already exists, replacing old connection`);
      const oldClient = this.clients.get(id)!;
      const oldToken = oldClient.getApiKey();
      // Remove from old token group
      this.tokenGroups.get(oldToken)?.delete(id);
      if (this.tokenGroups.get(oldToken)?.size === 0) {
        this.tokenGroups.delete(oldToken);
      }
      oldClient.disconnect();
      this.clients.delete(id);
    }

    // Create new client
    const client = new Client(ws, id, token, worldId, worldTitle, foundryVersion, systemId, systemTitle, systemVersion, customName);
    this.clients.set(id, client);

    // Add client to token group
    if (!this.tokenGroups.has(token)) {
      this.tokenGroups.set(token, new Set());
    }
    this.tokenGroups.get(token)?.add(id);
    
    // Store clientId->instanceId and apiKey->clientId mappings in Redis
    try {
      const redis = getRedisClient();
      
      if (redis) {
        // Store the mapping between API key (token) and this instance
        await redis.set(`apikey:${token}:instance`, INSTANCE_ID, { EX: CLIENT_EXPIRY });
        
        // Store the client ID -> instance ID mapping
        await redis.set(`client:${id}:instance`, INSTANCE_ID, { EX: CLIENT_EXPIRY });
        
        // Store client ID -> API key mapping for lookup
        await redis.set(`client:${id}:apikey`, token, { EX: CLIENT_EXPIRY });

        if (worldId) {
            await redis.set(`client:${id}:worldId`, worldId, { EX: CLIENT_EXPIRY });
        }
        if (worldTitle) {
            await redis.set(`client:${id}:worldTitle`, worldTitle, { EX: CLIENT_EXPIRY });
        }
        if (foundryVersion) {
            await redis.set(`client:${id}:foundryVersion`, foundryVersion, { EX: CLIENT_EXPIRY });
        }
        if (systemId) {
            await redis.set(`client:${id}:systemId`, systemId, { EX: CLIENT_EXPIRY });
        }
        if (systemTitle) {
            await redis.set(`client:${id}:systemTitle`, systemTitle, { EX: CLIENT_EXPIRY });
        }
        if (systemVersion) {
            await redis.set(`client:${id}:systemVersion`, systemVersion, { EX: CLIENT_EXPIRY });
        }
        if (customName) {
            await redis.set(`client:${id}:customName`, customName, { EX: CLIENT_EXPIRY });
        }
        
        // Add client ID to the list of clients for this token/API key
        await redis.sAdd(`apikey:${token}:clients`, id);
        await redis.expire(`apikey:${token}:clients`, CLIENT_EXPIRY);
        
        const tokenTrunicated = `${token.substring(0, 8)}...`;
        log.info(`Client ${id} registered in Redis with token ${tokenTrunicated}`);
      }
    } catch (error) {
      log.error(`Failed to register client in Redis: ${error}`);
      // Continue even if Redis fails - local operation will still work
    }

    const tokenTrunicated = `${token.substring(0, 8)}...`;
    log.info(`Client ${id} connected with token ${tokenTrunicated}`);
    return client;
  }

  /**
   * Remove a client from the manager
   */
  static async removeClient(id: string): Promise<void> {
    const client = this.clients.get(id);
    if (client) {
      const token = client.getApiKey();

      // Terminate any active sheet sessions for this client
      SheetSessionManager.terminateSessionsForClient(id);

      // Clean up local state
      this.clients.delete(id);
      this.tokenGroups.get(token)?.delete(id);
      if (this.tokenGroups.get(token)?.size === 0) {
        this.tokenGroups.delete(token);
      }
      
      // Clean up Redis state
      try {
        const redis = getRedisClient();
        if (redis) {
          // FIXED: Consistent key naming
          await redis.del(`client:${id}:instance`);
          await redis.del(`client:${id}:apikey`);
          await redis.sRem(`apikey:${token}:clients`, id);
          
          // If this was the last client for this token, remove instance mapping
          const remainingClients = await redis.sCard(`apikey:${token}:clients`);
          if (remainingClients === 0) {
            await redis.del(`apikey:${token}:instance`);
            await redis.del(`apikey:${token}:clients`);
          }
        }
      } catch (error) {
        log.error(`Failed to clean up client in Redis: ${error}`);
        // Continue even if Redis fails
      }
      
      log.info(`Client ${id} disconnected`);
    }
  }

  /**
   * Get a client by ID
   */
  static async getClient(id: string): Promise<Client | null> {
    // First check local clients
    const client = this.clients.get(id);
    if (client) {
      return client;
    }
    
    // If not found locally, check if this client should be on a different instance
    try {
      const redis = getRedisClient();
      if (redis) {
        // FIXED: Consistent key naming
        const instanceId = await redis.get(`client:${id}:instance`);
        
        if (instanceId && instanceId !== INSTANCE_ID) {
          // This client exists but is connected to a different instance
          log.info(`Client ${id} is connected to instance ${instanceId}, not to this instance ${INSTANCE_ID}`);
          return null;
        }
      }
    } catch (error) {
      log.error(`Error checking Redis for client: ${error}`);
    }
    
    return null;
  }

  /**
   * Get the instance ID for a client by ID
   */
  static async getClientInstance(id: string): Promise<string | null> {
    try {
      const redis = getRedisClient();
      if (redis) {
        return await redis.get(`client:id:${id}:instance`);
      } else {
        return null;
      }
    } catch (error) {
      log.error(`Error getting client instance from Redis: ${error}`);
      return null;
    }
  }

  /**
   * Get the instance ID for an API token
   */
  static async getInstanceForToken(token: string): Promise<string | null> {
    try {
      const redis = getRedisClient();
      if (redis) {
        return await redis.get(`client:${token}:instance`);
      } else {
        return null;
      }
    } catch (error) {
      log.error(`Error getting instance for token from Redis: ${error}`);
      return null;
    }
  }

  /**
   * Get the instance ID for an API key
   */
  static async getInstanceForApiKey(apiKey: string): Promise<string | null> {
    try {
      const redis = getRedisClient();
      if (redis) {
        // Directly look up the instance for this API key
        return await redis.get(`apikey:${apiKey}:instance`);
      }
      return null;
    } catch (error) {
      log.error(`Error getting instance for API key from Redis: ${error}`);
      return null;
    }
  }

  /**
   * Get all connected clients for an API key
   */
  static async getConnectedClients(apiKey: string): Promise<string[]> {
    const localClients = Array.from(this.tokenGroups.get(apiKey) || [])
      .filter(id => this.clients.has(id) && this.clients.get(id)!.isAlive())
      .map(id => id); // Just return the client IDs
    
    // In a distributed setup, we'd need to query other instances
    // For now, just return local clients
    return localClients;
  }

  /**
   * Update client's last seen timestamp
   */
  static updateClientLastSeen(id: string): void {
    const client = this.clients.get(id);
    if (client) {
      client.updateLastSeen();
      
      // Also update expiry in Redis
      try {
        const redis = getRedisClient();
        if (redis) {
          const token = client.getApiKey();
          
          // FIXED: Refresh expiry for all Redis keys with consistent naming
          redis.expire(`client:${id}:instance`, CLIENT_EXPIRY);
          redis.expire(`client:${id}:apikey`, CLIENT_EXPIRY);
          redis.expire(`apikey:${token}:clients`, CLIENT_EXPIRY);
          redis.expire(`apikey:${token}:instance`, CLIENT_EXPIRY);
        }
      } catch (error) {
        // Ignore Redis errors for heartbeat updates
      }
    }
  }

  /**
   * Broadcast a message to all clients in the same token group
   */
  static async broadcastToGroup(senderId: string, message: any): Promise<void> {
    const sender = this.clients.get(senderId);
    if (!sender) return;

    const token = sender.getApiKey();
    
    // Broadcast to local clients
    const groupClients = this.tokenGroups.get(token);
    if (groupClients) {
      for (const clientId of groupClients) {
        if (clientId !== senderId) {
          const client = this.clients.get(clientId);
          if (client && client.isAlive()) {
            client.send(message);
          }
        }
      }
    }
    
    // In a distributed setup with pub/sub, we'd publish to a Redis channel here
    // For Fly.io private network, we'd need to implement a pub/sub system
    // This is beyond the scope of this implementation
  }

  /**
   * Register a handler for a specific message type
   */
  static onMessageType(type: string, handler: MessageHandler): void {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, []);
    }
    this.messageHandlers.get(type)!.push(handler);
  }

  /**
   * Process an incoming message
   */
  static handleIncomingMessage(clientId: string, message: any): void {
    try {
      const client = this.clients.get(clientId);
      if (!client) return;

      // Update last seen timestamp
      client.updateLastSeen();

      // Handle ping messages specially
      if (message.type === "ping") {
        client.send({ type: "pong" });
        return;
      }
      
      // Handle other message types with registered handlers
      if (message.type && this.messageHandlers.has(message.type)) {
        for (const handler of this.messageHandlers.get(message.type)!) {
          try {
            handler(client, message);
          } catch (handlerError) {
            log.error('Error in message handler', {
              clientId,
              messageType: message.type,
              requestId: message.requestId,
              error: handlerError instanceof Error ? {
                name: handlerError.name,
                message: handlerError.message,
                stack: handlerError.stack
              } : String(handlerError)
            });
          }
        }
        return;
      }

      // Broadcast other messages
      this.broadcastToGroup(clientId, message);
    } catch (error) {
      log.error('Error handling message', {
        clientId,
        messageType: message?.type,
        requestId: message?.requestId,
        error: error instanceof Error ? {
          name: error.name,
          message: error.message,
          stack: error.stack
        } : String(error)
      });
    }
  }

  /**
   * Clean up inactive clients
   */
  static async cleanupInactiveClients(): Promise<void> {
    const toRemove: string[] = [];
    
    // Check all clients - only use isAlive() which now incorporates the grace period
    for (const [id, client] of this.clients.entries()) {
      if (!client.isAlive()) {
        toRemove.push(id);
      }
    }
    
    // Remove inactive clients
    for (const id of toRemove) {
      log.info(`Removing inactive client ${id}`);
      await this.removeClient(id);
    }
  }
}

export class WebSocketManager {
  // other properties...
  private messageHandlers: Map<string, (data: any) => void> = new Map();

  // other methods...
  
  onMessageType(type: string, handler: (data: any) => void): void {
    this.messageHandlers.set(type, handler);
  }

  private onMessage(event: MessageEvent): void {
    try {
      const data = JSON.parse(event.data);
      
      if (data.type && this.messageHandlers.has(data.type)) {
        this.messageHandlers.get(data.type)!(data);
      }
    } catch (error) {
      log.error('Error processing message', { error });
    }
  }
}
