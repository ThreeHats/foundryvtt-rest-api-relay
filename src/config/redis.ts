import { createClient, RedisClientType } from 'redis';
import { log } from '../utils/logger';

// Global Redis client
let redisClient: RedisClientType | null = null;
let redisConnected = false;
let redisReconnecting = false;
let lastRedisError: Error | null = null;

// Redis connection details from environment
const REDIS_URL = process.env.REDIS_URL || '';
const REDIS_ENABLED = process.env.REDIS_ENABLED !== 'false';

// Local mode detection (memory or sqlite DB types don't use Redis)
export const isLocalMode = process.env.DB_TYPE === 'memory' || process.env.DB_TYPE === 'sqlite';

// Initialize Redis client
export async function initRedis(): Promise<boolean> {
  if (!REDIS_ENABLED || !REDIS_URL) {
    log.info('Redis is disabled by configuration');
    return false;
  }

  try {
    // Hide credentials when logging
    const sanitizedUrl = REDIS_URL.replace(/redis:\/\/.*@/, 'redis://***@');
    log.info(`Initializing Redis client at ${sanitizedUrl}`);
    
    // Upstash requires TLS to be enabled
    const url = new URL(REDIS_URL);
    const isUpstash = url.hostname.includes('upstash.io');
    
    redisClient = createClient({
      url: REDIS_URL,
      socket: {
        tls: isUpstash, // Enable TLS for Upstash connections
        rejectUnauthorized: false, // For self-signed certificates
        connectTimeout: 15000, // Increase connect timeout to 15s
        reconnectStrategy: (retries) => {
          redisReconnecting = true;
          const maxRetries = 20; // Increased max retries
          
          if (retries >= maxRetries) {
            log.error(`Maximum Redis reconnection attempts (${maxRetries}) reached`);
            return false; // Stop retrying
          }
          
          const delay = Math.min(retries * 500, 10000);
          log.warn(`Redis reconnect attempt ${retries}, next attempt in ${delay}ms`);
          return delay;
        }
      },
      readonly: false,
      legacyMode: false,
      disableOfflineQueue: false
      // Connection retry behavior is handled by socket.reconnectStrategy above
    });

    // Set up event listeners for logging and state management
    redisClient.on('connect', () => {
      log.info('Redis client connected');
    });

    redisClient.on('ready', () => {
      redisConnected = true;
      redisReconnecting = false;
      lastRedisError = null;
      log.info('Redis client ready');
    });

    redisClient.on('error', (err) => {
      redisConnected = false;
      lastRedisError = err;
      log.error(`Redis client error: ${err.message}`);
    });

    redisClient.on('end', () => {
      redisConnected = false;
      log.warn('Redis client connection closed');
    });

    redisClient.on('reconnecting', () => {
      redisReconnecting = true;
      log.warn('Redis client reconnecting...');
    });

    // Connect to Redis
    await redisClient.connect();
    log.info('Redis client connected successfully');
    
    // Verify connection with a simple ping
    const pingResult = await redisClient.ping();
    log.info(`Redis ping result: ${pingResult}`);
    
    redisConnected = true;
    return true;
  } catch (error) {
    redisConnected = false;
    lastRedisError = error as Error;
    log.error(`Failed to initialize Redis: ${error}`);
    return false;
  }
}

// Get Redis client with safety checks
export function getRedisClient(): RedisClientType | null {
  if (!REDIS_ENABLED) {
    return null;
  }
  
  if (!redisConnected && !redisReconnecting && redisClient) {
    // Try to reconnect if client exists but is disconnected
    log.warn('Redis disconnected, attempting to reconnect...');
    redisClient.connect().catch(err => {
      log.error(`Failed to reconnect to Redis: ${err}`);
    });
  }
  
  return redisClient;
}

// Check Redis health
export function checkRedisHealth(): { healthy: boolean, error?: string, skipped?: boolean } {
  if (isLocalMode || !REDIS_ENABLED || !REDIS_URL) {
    return { healthy: true, skipped: true };
  }
  
  if (!redisClient) {
    return { healthy: false, error: 'Redis client not initialized' };
  }
  
  if (!redisConnected) {
    return { 
      healthy: false, 
      error: lastRedisError ? `Redis disconnected: ${lastRedisError.message}` : 'Redis disconnected' 
    };
  }
  
  return { healthy: true };
}

// Safely execute Redis operations with error handling
export async function safeRedisOperation<T>(operation: () => Promise<T>, fallback: T): Promise<T> {
  if (!REDIS_ENABLED || !redisClient || !redisConnected) {
    return fallback;
  }
  
  try {
    return await operation();
  } catch (error) {
    log.error(`Redis operation failed: ${error}`);
    return fallback;
  }
}

// Close Redis connection
export async function closeRedis(): Promise<void> {
  if (redisClient) {
    log.info('Closing Redis connection');
    await redisClient.quit().catch(err => {
      log.error(`Error closing Redis connection: ${err}`);
    });
    redisConnected = false;
    redisClient = null;
  }
}