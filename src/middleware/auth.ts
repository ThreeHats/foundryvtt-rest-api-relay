import { Request, Response, NextFunction } from 'express';
import { User } from '../models/user';
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
    // Find all users with the matching API key
    const users = await User.findAll({ where: { apiKey } });
    
    if (users.length === 0) {
      res.status(401).json({ error: 'Invalid API key' });
      return;
    }
    
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
    
    const user = users[0];
    req.user = user;
    
    const subscriptionStatus = user.getDataValue ? 
      user.getDataValue('subscriptionStatus') : user.subscriptionStatus;
    
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
    const apiKey = req.headers['x-api-key'] as string;
    
    if (apiKey) {
      // Use the User.findOne method that works with both sequelize and memory store
      const user = await User.findOne({ where: { apiKey } });
      
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

          updateSessionActivity(apiKey);
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
          updateSessionActivity(apiKey);
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
        
        next();
      } else {
        log.warn(`API key not found: ${apiKey.substring(0, 8)}...`);
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