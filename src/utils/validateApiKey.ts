import { User } from '../models/user';
import { ApiKey } from '../models/apiKey';
import { log } from './logger';

const isMemoryStore = process.env.DB_TYPE === 'memory';

export interface ApiKeyValidationResult {
  valid: boolean;
  masterApiKey?: string;
  scopedClientId?: string | null;
  scopedUserId?: string | null;
}

/**
 * Validate API key against database.
 * Shared between the Foundry module WebSocket (/relay) and the client-facing WebSocket (/ws/api).
 * Supports both master keys and scoped API keys.
 * @param apiKey The API key to validate
 * @returns validation result with master key info if scoped
 */
export async function validateApiKey(apiKey: string): Promise<boolean> {
  const result = await validateApiKeyDetailed(apiKey);
  return result.valid;
}

/**
 * Detailed API key validation that returns scoped key info for WebSocket auth.
 */
export async function validateApiKeyDetailed(apiKey: string): Promise<ApiKeyValidationResult> {
  if (isMemoryStore) {
    return { valid: true };
  }

  try {
    // 1. Try master key lookup
    const user = await User.findOne({ where: { apiKey } });
    if (user) {
      return { valid: true, masterApiKey: apiKey };
    }

    // 2. Try scoped API key lookup
    const scopedKey = await ApiKey.findByKey(apiKey);
    if (!scopedKey) {
      return { valid: false };
    }

    // Check enabled
    const enabled = scopedKey.getDataValue ? scopedKey.getDataValue('enabled') : scopedKey.enabled;
    if (!enabled) {
      return { valid: false };
    }

    // Check expiry
    const expiresAt = scopedKey.getDataValue ? scopedKey.getDataValue('expiresAt') : scopedKey.expiresAt;
    if (expiresAt && new Date(expiresAt) < new Date()) {
      return { valid: false };
    }

    // Look up parent user to get master key
    const userId = scopedKey.getDataValue ? scopedKey.getDataValue('userId') : scopedKey.userId;
    const parentUser = await User.findOne({ where: { id: userId } });
    if (!parentUser) {
      return { valid: false };
    }

    const masterApiKey = parentUser.getDataValue ? parentUser.getDataValue('apiKey') : parentUser.apiKey;
    const scopedClientId = scopedKey.getDataValue ? scopedKey.getDataValue('scopedClientId') : scopedKey.scopedClientId;
    const scopedUserId = scopedKey.getDataValue ? scopedKey.getDataValue('scopedUserId') : scopedKey.scopedUserId;

    return { valid: true, masterApiKey, scopedClientId, scopedUserId };
  } catch (error) {
    log.error(`Error validating API key: ${error instanceof Error ? error.message : 'Unknown error'}`);
    return { valid: false };
  }
}

/**
 * Track per-message API usage for a WebSocket connection.
 * Lightweight version of trackApiUsage middleware for WS messages.
 */
export async function trackWsApiUsage(apiKey: string): Promise<{ allowed: boolean; error?: string }> {
  if (isMemoryStore) {
    return { allowed: true };
  }

  try {
    const user = await User.findOne({ where: { apiKey } });
    if (!user) {
      return { allowed: false, error: 'Invalid API key' };
    }

    // Daily rate limit check
    const DAILY_REQUEST_LIMIT = parseInt(process.env.DAILY_REQUEST_LIMIT || '1000');
    const today = new Date().toISOString().split('T')[0];

    if ('getDataValue' in user && typeof user.getDataValue === 'function') {
      const lastRequestDate = user.getDataValue('lastRequestDate');
      let lastRequestDateStr: string | null = null;
      if (lastRequestDate) {
        if (lastRequestDate instanceof Date) {
          lastRequestDateStr = lastRequestDate.toISOString().split('T')[0];
        } else if (typeof lastRequestDate === 'string') {
          lastRequestDateStr = new Date(lastRequestDate).toISOString().split('T')[0];
        }
      }

      if (lastRequestDateStr !== today) {
        user.setDataValue('requestsToday', 0);
        user.setDataValue('lastRequestDate', new Date());
      }

      const currentDailyRequests = user.getDataValue('requestsToday') || 0;
      if (currentDailyRequests >= DAILY_REQUEST_LIMIT) {
        return { allowed: false, error: 'Daily API request limit reached' };
      }

      const currentMonthlyRequests = user.getDataValue('requestsThisMonth') || 0;
      user.setDataValue('requestsThisMonth', currentMonthlyRequests + 1);
      user.setDataValue('requestsToday', currentDailyRequests + 1);
      user.setDataValue('lastRequestDate', new Date());

      if ('save' in user && typeof user.save === 'function') {
        await user.save();
      }
    }

    return { allowed: true };
  } catch (error) {
    log.error(`Error tracking WS API usage: ${error}`);
    return { allowed: true }; // Don't block on tracking errors
  }
}
