import { sequelize } from '../sequelize';
import { log } from '../utils/logger';

/**
 * Migration to add maxHeadlessSessions column to Users table.
 * sequelize.sync({ alter: true }) should handle this automatically,
 * but this migration is a safety net for production databases.
 */
export async function migrateMaxHeadlessSessions(): Promise<void> {
  try {
    const isMemoryStore = process.env.DB_TYPE === 'memory';
    if (isMemoryStore) {
      log.info('Using memory store - skipping maxHeadlessSessions migration');
      return;
    }

    if (!('query' in sequelize)) {
      log.warn('Database does not support migrations - skipping');
      return;
    }

    log.info('Starting migration to add maxHeadlessSessions column');

    try {
      await (sequelize as any).query(`
        ALTER TABLE "Users"
        ADD COLUMN "maxHeadlessSessions" INTEGER;
      `);
      log.info('Added maxHeadlessSessions column');
    } catch (error: any) {
      if (error.message.includes('already exists') ||
          error.message.includes('duplicate column name') ||
          error.message.includes('duplicate column')) {
        log.info('maxHeadlessSessions column already exists - skipping');
      } else {
        throw error;
      }
    }

    log.info('maxHeadlessSessions migration completed successfully');
  } catch (error) {
    log.error('maxHeadlessSessions migration failed', { error });
    throw error;
  }
}
