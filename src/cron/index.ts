import * as cron from 'node-cron';
import { resetMonthlyRequests } from './monthlyReset';
import { resetDailyRequests } from './dailyReset';
import { PasswordResetToken } from '../models/passwordResetToken';
import { log } from '../utils/logger';

// Track scheduled jobs - use the correct type
let monthlyResetJob: cron.ScheduledTask | null = null;
let dailyResetJob: cron.ScheduledTask | null = null;
let tokenCleanupJob: cron.ScheduledTask | null = null;

/**
 * Set up all cron jobs for the application
 */
export function setupCronJobs(): void {
  if (!monthlyResetJob) {
    // Reset request counts at midnight on the first day of each month
    // Cron format: minute hour day month day-of-week
    monthlyResetJob = cron.schedule('0 0 1 * *', async () => {
      log.info('Running scheduled monthly request count reset');
      try {
        await resetMonthlyRequests();
        log.info('Monthly request count reset completed successfully via cron job');
      } catch (error) {
        log.error(`Error in monthly reset cron job: ${error}`);
        // Add retry logic
        setTimeout(async () => {
          log.info('Retrying monthly request count reset after failure');
          try {
            await resetMonthlyRequests();
            log.info('Monthly request count reset retry successful');
          } catch (retryError) {
            log.error(`Monthly reset retry also failed: ${retryError}`);
          }
        }, 5 * 60 * 1000); // Retry after 5 minutes
      }
    }, {
      timezone: 'UTC'
    });
    
    monthlyResetJob.start();
    log.info('Monthly reset cron job scheduled');
    
    // Also run immediately when starting the server if it's the 1st day of month
    const now = new Date();
    if (now.getDate() === 1) {
      log.info('Today is the 1st day of the month - running request reset immediately');
      resetMonthlyRequests().catch(error => {
        log.error(`Error running immediate monthly reset: ${error}`);
      });
    }
  } else {
    log.info('Monthly reset cron job already scheduled');
  }
  
  if (!dailyResetJob) {
    // Reset daily request counts at midnight every day
    // Cron format: minute hour day month day-of-week
    dailyResetJob = cron.schedule('0 0 * * *', async () => {
      log.info('Running scheduled daily request count reset');
      try {
        await resetDailyRequests();
        log.info('Daily request count reset completed successfully via cron job');
      } catch (error) {
        log.error(`Error in daily reset cron job: ${error}`);
        // Add retry logic
        setTimeout(async () => {
          log.info('Retrying daily request count reset after failure');
          try {
            await resetDailyRequests();
            log.info('Daily request count reset retry successful');
          } catch (retryError) {
            log.error(`Daily reset retry also failed: ${retryError}`);
          }
        }, 5 * 60 * 1000); // Retry after 5 minutes
      }
    }, {
      timezone: 'UTC'
    });
    
    dailyResetJob.start();
    log.info('Daily reset cron job scheduled');
  } else {
    log.info('Daily reset cron job already scheduled');
  }

  if (!tokenCleanupJob) {
    // Clean up expired/used password reset tokens every 6 hours
    tokenCleanupJob = cron.schedule('0 */6 * * *', async () => {
      log.info('Running scheduled password reset token cleanup');
      try {
        const deleted = await PasswordResetToken.cleanupExpired();
        log.info(`Password reset token cleanup completed, removed ${deleted} tokens`);
      } catch (error) {
        log.error(`Error in token cleanup cron job: ${error}`);
      }
    }, {
      timezone: 'UTC'
    });

    tokenCleanupJob.start();
    log.info('Token cleanup cron job scheduled');
  } else {
    log.info('Token cleanup cron job already scheduled');
  }

  log.info('Cron jobs setup completed');
}

/**
 * Stop all cron jobs (useful for graceful shutdown)
 */
export function stopCronJobs(): void {
  if (monthlyResetJob) {
    monthlyResetJob.stop();
    monthlyResetJob = null;
    log.info('Monthly reset cron job stopped');
  }
  
  if (dailyResetJob) {
    dailyResetJob.stop();
    dailyResetJob = null;
    log.info('Daily reset cron job stopped');
  }

  if (tokenCleanupJob) {
    tokenCleanupJob.stop();
    tokenCleanupJob = null;
    log.info('Token cleanup cron job stopped');
  }
}

/**
 * Get status of cron jobs
 */
export function getCronJobStatus(): {
  monthlyReset: { scheduled: boolean; active: boolean };
  dailyReset: { scheduled: boolean; active: boolean };
  tokenCleanup: { scheduled: boolean; active: boolean };
} {
  return {
    monthlyReset: {
      scheduled: monthlyResetJob !== null,
      active: monthlyResetJob !== null
    },
    dailyReset: {
      scheduled: dailyResetJob !== null,
      active: dailyResetJob !== null
    },
    tokenCleanup: {
      scheduled: tokenCleanupJob !== null,
      active: tokenCleanupJob !== null
    }
  };
}

/**
 * Manually trigger the monthly reset (for testing or emergency use)
 */
export async function triggerMonthlyReset(): Promise<void> {
  log.info('Manually triggering monthly request count reset');
  await resetMonthlyRequests();
}

/**
 * Manually trigger the daily reset (for testing or emergency use)
 */
export async function triggerDailyReset(): Promise<void> {
  log.info('Manually triggering daily request count reset');
  await resetDailyRequests();
}
