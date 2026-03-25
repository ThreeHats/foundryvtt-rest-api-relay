import { Request, Response, NextFunction } from 'express';
import { User } from '../models/user';
import { ApiKey } from '../models/apiKey';
import { ClientManager } from '../core/ClientManager';
import { sequelize } from '../sequelize';
import { log } from '../utils/logger';
import { apiKeyToSession } from '../routes/api';

// Helper function to update session activity timestamp
function updateSessionActivity(apiKey: string) {
  const session = apiKeyToSession.get(apiKey);
  if (session) {
    session.lastActivity = Date.now();
  }
}

// Flag to check if we're using memory store
const isMemoryStore = process.env.DB_TYPE === 'memory';

// Free tier request limit per month
const FREE_TIER_LIMIT = parseInt(process.env.FREE_API_REQUESTS_LIMIT || '100');

// Daily request limit for all users (configurable via environment variable)
const DAILY_REQUEST_LIMIT = parseInt(process.env.DAILY_REQUEST_LIMIT || '1000');

declare global {
  namespace Express {
    interface Request {
      user: any;
      subscriptionStatus?: string;
      masterApiKey?: string;
      scopedKey?: {
        id: number;
        scopedClientId?: string;
        scopedUserId?: string;
        dailyLimit?: number;
      };
    }
  }
}

export const authMiddleware = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  // If using memory store in local dev, bypass authentication
  if (isMemoryStore) {
    log.info('Using memory store - bypassing API key authentication');

    req.user = {
      id: 1,
      email: 'admin@example.com',
      apiKey: 'local-dev',
      requestsThisMonth: 0,
      subscriptionStatus: 'active'
    };
    next();
    return;
  }

  // Normal authentication flow for SQLite and PostgreSQL
  const apiKey = req.headers['x-api-key'] as string;
  const clientId = req.query.clientId as string;

  if (!apiKey) {
    res.status(401).json({ error: 'API key is required' });
    return;
  }

  try {
    // 1. Try master key lookup first (existing logic)
    const users = await User.findAll({ where: { apiKey } });

    if (users.length > 0) {
      // Master key auth — existing behavior
      const user = users[0];

      if (clientId) {
        const client = await ClientManager.getClient(clientId);
        if (!client) {
          res.status(404).json({ error: 'Invalid client ID' });
          return;
        }

        if (client.getApiKey() !== apiKey) {
          log.warn(`Client ID ${clientId} does not match API key ${apiKey.substring(0, 8)}...`);
          res.status(401).json({ error: 'Invalid API key for this client ID' });
          return;
        }
      }

      req.user = user;
      const subscriptionStatus = user.getDataValue ?
        user.getDataValue('subscriptionStatus') : user.subscriptionStatus;
      req.subscriptionStatus = subscriptionStatus || 'free';

      next();
      return;
    }

    // 2. Try scoped API key lookup
    const scopedKey = await ApiKey.findByKey(apiKey);

    if (!scopedKey) {
      res.status(401).json({ error: 'Invalid API key' });
      return;
    }

    // Validate scoped key is enabled
    const enabled = scopedKey.getDataValue ? scopedKey.getDataValue('enabled') : scopedKey.enabled;
    if (!enabled) {
      res.status(401).json({ error: 'API key is disabled' });
      return;
    }

    // Validate scoped key is not expired
    const expiresAt = scopedKey.getDataValue ? scopedKey.getDataValue('expiresAt') : scopedKey.expiresAt;
    if (expiresAt && new Date(expiresAt) < new Date()) {
      res.status(401).json({ error: 'API key has expired' });
      return;
    }

    // Look up parent user
    const userId = scopedKey.getDataValue ? scopedKey.getDataValue('userId') : scopedKey.userId;
    const parentUser = await User.findOne({ where: { id: userId } });
    if (!parentUser) {
      res.status(401).json({ error: 'Invalid API key' });
      return;
    }

    // Set request properties
    req.user = parentUser;
    const parentApiKey = parentUser.getDataValue ? parentUser.getDataValue('apiKey') : parentUser.apiKey;
    req.masterApiKey = parentApiKey;

    const scopedClientId = scopedKey.getDataValue ? scopedKey.getDataValue('scopedClientId') : scopedKey.scopedClientId;
    const scopedUserId = scopedKey.getDataValue ? scopedKey.getDataValue('scopedUserId') : scopedKey.scopedUserId;
    const dailyLimit = scopedKey.getDataValue ? scopedKey.getDataValue('dailyLimit') : scopedKey.dailyLimit;
    const scopedKeyId = scopedKey.getDataValue ? scopedKey.getDataValue('id') : scopedKey.id;

    req.scopedKey = {
      id: scopedKeyId,
      scopedClientId: scopedClientId || undefined,
      scopedUserId: scopedUserId || undefined,
      dailyLimit: dailyLimit || undefined
    };

    // Validate clientId if provided matches against parent user's master key
    if (clientId) {
      const client = await ClientManager.getClient(clientId);
      if (!client) {
        res.status(404).json({ error: 'Invalid client ID' });
        return;
      }

      if (client.getApiKey() !== parentApiKey) {
        log.warn(`Client ID ${clientId} does not match parent API key for scoped key`);
        res.status(401).json({ error: 'Invalid API key for this client ID' });
        return;
      }
    }

    const subscriptionStatus = parentUser.getDataValue ?
      parentUser.getDataValue('subscriptionStatus') : parentUser.subscriptionStatus;
    req.subscriptionStatus = subscriptionStatus || 'free';

    next();
  } catch (error) {
    log.error(`Auth error: ${error}`);
    res.status(500).json({ error: 'Authentication error' });
  }
};

export const trackApiUsage = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  // Skip usage tracking in memory store mode
  if (isMemoryStore) {
    return next();
  }

  // Normal API usage tracking
  try {
    // For scoped keys, use the parent user's master key for user-level tracking
    const rawApiKey = req.headers['x-api-key'] as string;
    const masterKey = req.masterApiKey || rawApiKey;

    if (masterKey) {
      // Use the User.findOne method that works with both sequelize and memory store
      const user = await User.findOne({ where: { apiKey: masterKey } });

      if (user) {
        // Always track api usage regardless of subscription status
        if ('getDataValue' in user && typeof user.getDataValue === 'function') {
          // Check if it's a new day - reset daily counter if needed
          const today = new Date().toISOString().split('T')[0]; // YYYY-MM-DD format
          const lastRequestDate = user.getDataValue('lastRequestDate');

          // Safely handle lastRequestDate - it might be a string, Date, or null
          let lastRequestDateStr = null;
          if (lastRequestDate) {
            if (lastRequestDate instanceof Date) {
              lastRequestDateStr = lastRequestDate.toISOString().split('T')[0];
            } else if (typeof lastRequestDate === 'string') {
              lastRequestDateStr = new Date(lastRequestDate).toISOString().split('T')[0];
            }
          }

          if (lastRequestDateStr !== today) {
            // New day - reset daily counter
            user.setDataValue('requestsToday', 0);
            user.setDataValue('lastRequestDate', new Date());
          }

          // Get current request counts
          const currentMonthlyRequests = user.getDataValue('requestsThisMonth') || 0;
          const currentDailyRequests = user.getDataValue('requestsToday') || 0;

          // Check daily rate limit
          if (currentDailyRequests >= DAILY_REQUEST_LIMIT) {
            // Calculate midnight of the next day for reset time
            const tomorrow = new Date();
            tomorrow.setDate(tomorrow.getDate() + 1);
            tomorrow.setHours(0, 0, 0, 0); // Set to midnight

            log.warn(`Daily rate limit hit for user ID ${user.getDataValue('id')} - ${currentDailyRequests}/${DAILY_REQUEST_LIMIT} requests`);

            res.status(429).json({
              error: 'Daily API request limit reached',
              dailyLimit: DAILY_REQUEST_LIMIT,
              currentRequests: currentDailyRequests,
              message: `You have reached the daily limit of ${DAILY_REQUEST_LIMIT} requests. Please try again tomorrow.`,
              resetsAt: tomorrow.toISOString()
            });
            return;
          }

          // Increment both counters
          user.setDataValue('requestsThisMonth', currentMonthlyRequests + 1);
          user.setDataValue('requestsToday', currentDailyRequests + 1);
          user.setDataValue('lastRequestDate', new Date());

          // Log with proper data access
          log.info(`Incrementing requests for user ID ${user.getDataValue('id')} - Monthly: ${user.getDataValue('requestsThisMonth')}, Daily: ${user.getDataValue('requestsToday')}`);

          // Save the updated user
          if ('save' in user && typeof user.save === 'function') {
            await user.save();
          }

          updateSessionActivity(masterKey);
        } else if ('requestsThisMonth' in user) {
          // Fallback for memory store
          const today = new Date().toISOString().split('T')[0];

          // Safely handle lastRequestDate for memory store too
          let lastRequestDateStr = null;
          if (user.lastRequestDate) {
            if (user.lastRequestDate instanceof Date) {
              lastRequestDateStr = user.lastRequestDate.toISOString().split('T')[0];
            } else if (typeof user.lastRequestDate === 'string') {
              lastRequestDateStr = new Date(user.lastRequestDate).toISOString().split('T')[0];
            }
          }

          if (lastRequestDateStr !== today) {
            user.requestsToday = 0;
            user.lastRequestDate = new Date();
          }

          // Check daily rate limit
          if (user.requestsToday >= DAILY_REQUEST_LIMIT) {
            // Calculate midnight of the next day for reset time
            const tomorrow = new Date();
            tomorrow.setDate(tomorrow.getDate() + 1);
            tomorrow.setHours(0, 0, 0, 0); // Set to midnight

            res.status(429).json({
              error: 'Daily API request limit reached',
              dailyLimit: DAILY_REQUEST_LIMIT,
              currentRequests: user.requestsToday,
              message: `You have reached the daily limit of ${DAILY_REQUEST_LIMIT} requests. Please try again tomorrow.`,
              resetsAt: tomorrow.toISOString()
            });
            return;
          }

          user.requestsThisMonth += 1;
          user.requestsToday += 1;
          user.lastRequestDate = new Date();
          updateSessionActivity(masterKey);
        }

        // Enforce monthly limits only for free tier users
        const subscriptionStatus = user.getDataValue ?
          user.getDataValue('subscriptionStatus') : user.subscriptionStatus;

        if (subscriptionStatus !== 'active') {
          const requestCount = user.getDataValue ?
            user.getDataValue('requestsThisMonth') : user.requestsThisMonth;

          if (requestCount >= FREE_TIER_LIMIT) {
            res.status(429).json({
              error: 'Monthly API request limit reached',
              limit: FREE_TIER_LIMIT,
              message: 'Please upgrade to a paid subscription for unlimited monthly API access',
              upgradeUrl: '/api/subscriptions/create-checkout-session'
            });
            return;
          }
        }

        // Per-key daily rate limit check for scoped keys
        if (req.scopedKey && req.scopedKey.dailyLimit) {
          const scopedKeyRecord = await ApiKey.findOne({ where: { id: req.scopedKey.id } });
          if (scopedKeyRecord) {
            const today = new Date().toISOString().split('T')[0];
            const keyLastDate = scopedKeyRecord.getDataValue ? scopedKeyRecord.getDataValue('lastRequestDate') : scopedKeyRecord.lastRequestDate;

            let keyLastDateStr: string | null = null;
            if (keyLastDate) {
              if (keyLastDate instanceof Date) {
                keyLastDateStr = keyLastDate.toISOString().split('T')[0];
              } else if (typeof keyLastDate === 'string') {
                keyLastDateStr = new Date(keyLastDate).toISOString().split('T')[0];
              }
            }

            if (keyLastDateStr !== today) {
              if (scopedKeyRecord.setDataValue) {
                scopedKeyRecord.setDataValue('requestsToday', 0);
                scopedKeyRecord.setDataValue('lastRequestDate', new Date());
              } else {
                scopedKeyRecord.requestsToday = 0;
                scopedKeyRecord.lastRequestDate = new Date();
              }
            }

            const keyDailyRequests = scopedKeyRecord.getDataValue ?
              scopedKeyRecord.getDataValue('requestsToday') : scopedKeyRecord.requestsToday;

            if (keyDailyRequests >= req.scopedKey.dailyLimit) {
              const tomorrow = new Date();
              tomorrow.setDate(tomorrow.getDate() + 1);
              tomorrow.setHours(0, 0, 0, 0);

              res.status(429).json({
                error: 'API key daily limit reached',
                dailyLimit: req.scopedKey.dailyLimit,
                currentRequests: keyDailyRequests,
                message: `This API key has reached its daily limit of ${req.scopedKey.dailyLimit} requests.`,
                resetsAt: tomorrow.toISOString()
              });
              return;
            }

            // Increment per-key counter
            if (scopedKeyRecord.setDataValue) {
              scopedKeyRecord.setDataValue('requestsToday', (keyDailyRequests || 0) + 1);
              scopedKeyRecord.setDataValue('lastRequestDate', new Date());
            } else {
              scopedKeyRecord.requestsToday = (keyDailyRequests || 0) + 1;
              scopedKeyRecord.lastRequestDate = new Date();
            }

            if (scopedKeyRecord.save && typeof scopedKeyRecord.save === 'function') {
              await scopedKeyRecord.save();
            }
          }
        }

        next();
      } else {
        // If using a scoped key, the master key lookup already passed in authMiddleware
        // so this shouldn't happen, but handle gracefully
        log.warn(`API key not found for usage tracking: ${masterKey.substring(0, 8)}...`);
        res.status(401).json({ error: 'Invalid API key' });
        return;
      }
    } else {
      log.warn('API key is required for usage tracking');
      res.status(401).json({ error: 'API key is required' });
      return;
    }
  } catch (error) {
    log.error(`Error tracking API usage: ${error}`);
    res.status(500).json({ error: 'Internal server error' });
    return;
  }
};
