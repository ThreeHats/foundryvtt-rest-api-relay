import { Request, Response, NextFunction } from 'express';
import { getRedisClient, safeRedisOperation } from '../config/redis';
import { log } from '../utils/logger';

// Store headless session data in Redis
export async function storeHeadlessSession(
  sessionId: string, 
  userId: string,
  apiKey: string, 
  instanceId: string
): Promise<boolean> {
  const redis = getRedisClient();
  if (!redis) return false;
  
  const clientId = `foundry-${userId}`;
  
  try {
    // Store session data with 3 hour expiry
    await redis.hSet(`headless_session:${sessionId}`, {
      clientId,
      apiKey,
      instanceId,
      created: Date.now().toString()
    });
    
    await redis.expire(`headless_session:${sessionId}`, 10800);
    
    // Store clientId to sessionId mapping
    await redis.set(`headless_client:${clientId}`, sessionId);
    await redis.expire(`headless_client:${clientId}`, 10800);
    
    // Store apiKey to sessionId mapping for quick lookups
    await redis.set(`headless_apikey:${apiKey}:session`, sessionId);
    await redis.expire(`headless_apikey:${apiKey}:session`, 10800);
    
    log.info(`Stored headless session ${sessionId} for client ${clientId} in Redis`);
    return true;
  } catch (error) {
    log.error(`Failed to store headless session in Redis: ${error}`);
    return false;
  }
}

// Get session ID for a client
export async function getSessionForClient(clientId: string): Promise<string | null> {
  return await safeRedisOperation(
    async () => {
      const redis = getRedisClient();
      if (!redis) return null;
      return await redis.get(`headless_client:${clientId}`);
    },
    null
  );
}

// Get session data by session ID
export async function getSessionData(sessionId: string): Promise<any | null> {
  return await safeRedisOperation(
    async () => {
      const redis = getRedisClient();
      if (!redis) return null;
      
      const exists = await redis.exists(`headless_session:${sessionId}`);
      if (!exists) return null;
      
      return await redis.hGetAll(`headless_session:${sessionId}`);
    },
    null
  );
}

// Get active sessions for API key
export async function getSessionForApiKey(apiKey: string): Promise<string | null> {
  return await safeRedisOperation(
    async () => {
      const redis = getRedisClient();
      if (!redis) return null;
      return await redis.get(`headless_apikey:${apiKey}:session`);
    },
    null
  );
}

// Clean up session data
export async function cleanupHeadlessSession(sessionId: string): Promise<boolean> {
  return await safeRedisOperation(
    async () => {
      const redis = getRedisClient();
      if (!redis) return false;
      
      // Get session data first to find clientId and apiKey
      const sessionData = await redis.hGetAll(`headless_session:${sessionId}`);
      if (!sessionData) return false;
      
      const { clientId, apiKey } = sessionData;
      
      // Clean up all related Redis keys
      await redis.del(`headless_session:${sessionId}`);
      if (clientId) await redis.del(`headless_client:${clientId}`);
      if (apiKey) await redis.del(`headless_apikey:${apiKey}:session`);
      
      log.info(`Cleaned up Redis data for session ${sessionId}`);
      return true;
    },
    false
  );
}

// Redis middleware to ensure cross-instance awareness
export function redisSessionMiddleware(req: Request, res: Response, next: NextFunction): void {
  const clientId = req.query.clientId as string;
  
  // If no client ID or not a valid string, continue with the request
  if (!clientId || typeof clientId !== 'string' || !clientId.startsWith('foundry-')) {
    return next();
  }
  
  // Check if this client has a session in Redis
  getSessionForClient(clientId)
    .then(sessionId => {
      if (sessionId) {
        // Add session ID to request for other middleware/handlers
        (req as any).headlessSessionId = sessionId;
        
        // Get full session data
        return getSessionData(sessionId);
      }
      return null;
    })
    .then(sessionData => {
      if (sessionData) {
        // Add session data to request
        (req as any).headlessSessionData = sessionData;
      }
      next();
    })
    .catch(error => {
      log.error(`Error in Redis session middleware: ${error}`);
      next();
    });
}