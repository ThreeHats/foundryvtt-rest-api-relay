import { log } from '../utils/logger';
import { getRedisClient } from '../config/redis';

// Store active browser sessions locally
export const browserSessions = new Map<string, any>(); // Using 'any' for puppeteer.Browser

// Check if a client ID is from a headless session
export function isHeadlessClient(clientId: string): boolean {
  return clientId.startsWith('foundry-');
}

// Generate a placeholder client ID for headless sessions (actual ID is determined by the module)
export function getHeadlessClientId(userId: string): string {
  return `foundry-pending-${userId}`;
}

// Track pending headless sessions
export async function registerHeadlessSession(sessionId: string, userId: string, apiKey: string): Promise<void> {
  const redis = getRedisClient();
  if (!redis) return;

  const clientId = getHeadlessClientId(userId);
  const instanceId = process.env.FLY_ALLOC_ID || 'local';
  
  try {
    // Store the session mapping
    await redis.hSet(`headless_session:${sessionId}`, {
      clientId,
      apiKey,
      instanceId,
      created: Date.now()
    });
    
    // Set expiration (3 hours)
    await redis.expire(`headless_session:${sessionId}`, 10800);
    
    // Store reverse lookup from clientId to sessionId
    await redis.set(`headless_client:${clientId}`, sessionId);
    await redis.expire(`headless_client:${clientId}`, 10800);
    
    // Store apiKey to instanceId mapping - CRITICAL FOR REQUEST FORWARDING
    await redis.set(`apikey:${apiKey}:instance`, instanceId);
    await redis.expire(`apikey:${apiKey}:instance`, 10800);
    
    // Also store client to instance mapping for socket lookups
    await redis.set(`client:${clientId}:instance`, instanceId);  
    await redis.expire(`client:${clientId}:instance`, 10800);
    
    log.info(`Registered headless session ${sessionId} for client ${clientId} on instance ${instanceId}`);
  } catch (error) {
    log.error(`Failed to register headless session: ${error}`);
  }
}

// Validate client connections - MODIFIED TO HANDLE INSTANCE MIGRATIONS
export async function validateHeadlessSession(clientId: string, token: string): Promise<boolean> {
  // Skip non-headless clients
  if (!isHeadlessClient(clientId)) {
    return true;
  }
  
  try {
    const redis = getRedisClient();
    if (!redis) return true; // Allow if Redis is not available
    
    // Get the session ID for this client
    const sessionId = await redis.get(`headless_client:${clientId}`);
    if (!sessionId) {
      log.warn(`No session found for headless client ${clientId}`);
      return true; // Allow connection but log warning
    }
    
    // Get session data
    const sessionData = await redis.hGetAll(`headless_session:${sessionId}`);
    if (!sessionData || !sessionData.apiKey) {
      log.warn(`Session data not found for session ${sessionId}`);
      return true;
    }
    
    // Check if API key matches
    if (sessionData.apiKey !== token) {
      log.warn(`API key mismatch for headless client ${clientId}`);
      return false; // Reject the connection
    }
    
    // Update the instance ID in case client connects to a different instance
    const currentInstanceId = process.env.FLY_ALLOC_ID || 'local';
    
    // Always update the instance location when validating a session
    await redis.hSet(`headless_session:${sessionId}`, "instanceId", currentInstanceId);
    await redis.set(`client:${clientId}:instance`, currentInstanceId);
    await redis.set(`apikey:${sessionData.apiKey}:instance`, currentInstanceId);
    
    // Touch all keys to refresh TTL
    await redis.expire(`headless_session:${sessionId}`, 10800);
    await redis.expire(`headless_client:${clientId}`, 10800);
    await redis.expire(`apikey:${sessionData.apiKey}:instance`, 10800);
    await redis.expire(`client:${clientId}:instance`, 10800);
    
    log.info(`Headless client ${clientId} validated successfully`);
    return true;
  } catch (error) {
    log.error(`Error validating headless session: ${error}`);
    return true; // Allow connection on error to avoid blocking legitimate connections
  }
}

// Update Redis with the actual clientId once the module connects
export async function updateHeadlessSessionClientId(sessionId: string, actualClientId: string, apiKey: string): Promise<void> {
  const redis = getRedisClient();
  if (!redis) return;

  const instanceId = process.env.FLY_ALLOC_ID || 'local';

  try {
    // Update the session record
    await redis.hSet(`headless_session:${sessionId}`, 'clientId', actualClientId);

    // Set the reverse lookup for the actual clientId
    await redis.set(`headless_client:${actualClientId}`, sessionId);
    await redis.expire(`headless_client:${actualClientId}`, 10800);

    // Update client-to-instance mapping
    await redis.set(`client:${actualClientId}:instance`, instanceId);
    await redis.expire(`client:${actualClientId}:instance`, 10800);

    log.info(`Updated headless session ${sessionId} with actual clientId: ${actualClientId}`);
  } catch (error) {
    log.error(`Failed to update headless session clientId: ${error}`);
  }
}

// Get session ID from client ID
export async function getSessionForClient(clientId: string): Promise<string | null> {
  try {
    const redis = getRedisClient();
    if (!redis) return null;
    
    return await redis.get(`headless_client:${clientId}`);
  } catch (error) {
    log.error(`Error getting session for client: ${error}`);
    return null;
  }
}

// Process pending session requests that were forwarded from other instances
export async function checkPendingHeadlessSessions() {
  try {
    const redis = getRedisClient();
    if (!redis) return;
    
    const instanceId = process.env.FLY_ALLOC_ID || 'local';
    
    // Find all handshakes that belong to this instance
    const handshakeKeys = await redis.keys('handshake:*');
    
    for (const key of handshakeKeys) {
      const handshakeData = await redis.hGetAll(key);
      
      if (handshakeData.instanceId === instanceId) {
        const handshakeToken = key.split(':')[1];
        
        // Check if there's a pending session request for this handshake
        const pendingSessionKey = `pending_session:${handshakeToken}`;
        const pendingSessionExists = await redis.exists(pendingSessionKey);
        
        if (pendingSessionExists) {
          const sessionData = await redis.hGetAll(pendingSessionKey);
          
          log.info(`Processing pending session request for handshake ${handshakeToken.substring(0, 8)}...`);
          
          // Process the session request here
          // This should trigger your session creation logic without going through the API endpoint
          
          // Clean up the pending request
          await redis.del(pendingSessionKey);
        }
      }
    }
  } catch (error) {
    log.error(`Error checking pending headless sessions: ${error}`);
  }
}

// Schedule this to run regularly
export function scheduleHeadlessSessionsCheck() {
  setInterval(checkPendingHeadlessSessions, 5000); // Check every 5 seconds
}