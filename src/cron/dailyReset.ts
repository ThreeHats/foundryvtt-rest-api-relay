import { User } from '../models/user';
import { ApiKey } from '../models/apiKey';
import { log } from '../utils/logger';
import { getRedisClient } from '../config/redis';

/**
 * Reset the requestsToday counter for all users
 * This is run daily at midnight via a cron job
 */
export async function resetDailyRequests(): Promise<void> {
  const redis = getRedisClient();
  const lockKey = 'daily_reset_lock';
  const lockValue = `${process.env.FLY_ALLOC_ID || 'local'}_${Date.now()}`;
  const lockTTL = 300; // 5 minutes lock timeout
  
  try {
    // Try to acquire distributed lock
    if (redis) {
      // Use Redis SET with NX (not exists) and PX (expire in milliseconds) options
      const lockAcquired = await redis.set(lockKey, lockValue, { 
        NX: true, 
        PX: lockTTL * 1000 
      });
      
      if (!lockAcquired) {
        log.info('Daily reset already running on another instance - skipping');
        return;
      }
      
      log.info(`Acquired daily reset lock on instance ${process.env.FLY_ALLOC_ID || 'local'}`);
    } else {
      log.warn('Redis not available - proceeding with reset (single instance mode)');
    }
    
    const startTime = Date.now();
    log.info(`Starting daily API request count reset for all users at ${new Date().toISOString()}`);
    
    // For SQL databases, we can do a bulk update
    const [updatedCount] = await User.update(
      { 
        requestsToday: 0,
        lastRequestDate: new Date()
      },
      { where: {} }  // Empty where clause updates all records
    );
    
    // Get total count of users for verification
    const totalUsers = await User.count();
    
    // Calculate execution time
    const executionTime = Date.now() - startTime;
    
    // Also reset ApiKey daily counters
    try {
      await ApiKey.resetDailyCounters();
      log.info('Reset daily counters for all scoped API keys');
    } catch (apiKeyError) {
      log.error(`Error resetting ApiKey daily counters: ${apiKeyError}`);
    }

    log.info(`Successfully reset daily request count for ${updatedCount} of ${totalUsers} users (took ${executionTime}ms)`);
    
    // Verify the reset worked by checking a sample
    const sampleUser = await User.findOne({ where: {} });
    if (sampleUser) {
      const requestCount = sampleUser.getDataValue ? 
        sampleUser.getDataValue('requestsToday') : sampleUser.requestsToday;
      
      log.info(`Verification - Sample user daily request count: ${requestCount}`);
      
      if (requestCount !== 0) {
        log.warn('Daily reset verification failed - some users may still have non-zero daily request counts');
      }
    }
    
    // Store completion timestamp in Redis for monitoring
    if (redis) {
      await redis.set('last_daily_reset', new Date().toISOString());
      await redis.expire('last_daily_reset', 86400 * 2); // Keep for 2 days
    }
    
  } catch (error) {
    log.error(`Error resetting daily request counts: ${error}`);
    
    // Try again with a different approach if the first method fails
    try {
      log.info('Attempting alternate daily reset method using findAll + individual updates');
      
      // Get all users and update them individually
      const users = await User.findAll({});
      let successCount = 0;
      
      for (const user of users) {
        if ('setDataValue' in user && typeof user.setDataValue === 'function') {
          user.setDataValue('requestsToday', 0);
          user.setDataValue('lastRequestDate', new Date());
          await user.save();
          successCount++;
        } else if ('requestsToday' in user) {
          user.requestsToday = 0;
          user.lastRequestDate = new Date();
          if ('save' in user && typeof user.save === 'function') {
            await user.save();
            successCount++;
          }
        }
      }
      
      log.info(`Daily recovery method successful - reset ${successCount} users individually`);
    } catch (recoveryError) {
      log.error(`Daily recovery attempt also failed: ${recoveryError}`);
    }
    
  } finally {
    // Release the lock
    if (redis) {
      try {
        // Only release if we still own the lock
        const script = `
          if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
          else
            return 0
          end
        `;
        await redis.eval(script, { keys: [lockKey], arguments: [lockValue] });
        log.info('Released daily reset lock');
      } catch (lockError) {
        log.warn(`Error releasing daily lock: ${lockError}`);
      }
    }
  }
}
