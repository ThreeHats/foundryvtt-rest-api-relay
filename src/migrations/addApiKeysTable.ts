import { sequelize } from '../sequelize';
import { log } from '../utils/logger';

/**
 * Migration to ensure the ApiKeys table exists.
 *
 * The table is normally created by sequelize.sync() from the model definition.
 * This migration handles the case where the model is added to an existing database
 * that was synced before the ApiKey model existed.
 *
 * For SQLite: sequelize.sync({ alter: true }) handles table creation from the model.
 * For PostgreSQL: we use raw SQL as a safety net if sync didn't create it.
 */
export async function migrateApiKeysTable(): Promise<void> {
  try {
    const isMemoryStore = process.env.DB_TYPE === 'memory';
    if (isMemoryStore) {
      log.info('Using memory store - skipping ApiKeys table migration');
      return;
    }

    if (!('query' in sequelize)) {
      log.warn('Database does not support migrations - skipping');
      return;
    }

    log.info('Starting migration to create ApiKeys table');

    const dialect = (sequelize as any).getDialect?.() || process.env.DB_TYPE || 'postgres';

    if (dialect === 'sqlite') {
      // SQLite: sequelize.sync() already created the table from the model.
      // Just verify it exists.
      try {
        await (sequelize as any).query(`SELECT 1 FROM "ApiKeys" LIMIT 0;`);
        log.info('ApiKeys table already exists (SQLite)');
      } catch {
        // Table doesn't exist yet — sync should have created it.
        // Force a model sync for just this table.
        const { ApiKey } = await import('../models/apiKey');
        await (ApiKey as any).sync();
        log.info('ApiKeys table created via model sync (SQLite)');
      }
    } else {
      // PostgreSQL: use dialect-specific SQL
      try {
        await (sequelize as any).query(`
          CREATE TABLE IF NOT EXISTS "ApiKeys" (
            "id" SERIAL PRIMARY KEY,
            "userId" INTEGER NOT NULL,
            "key" VARCHAR(64) NOT NULL UNIQUE,
            "name" VARCHAR(255) NOT NULL,
            "scopedClientId" VARCHAR(255),
            "scopedUserId" VARCHAR(255),
            "dailyLimit" INTEGER,
            "requestsToday" INTEGER DEFAULT 0,
            "lastRequestDate" DATE,
            "foundryUrl" VARCHAR(255),
            "foundryUsername" VARCHAR(255),
            "encryptedFoundryPassword" TEXT,
            "passwordIv" VARCHAR(255),
            "passwordAuthTag" VARCHAR(255),
            "expiresAt" TIMESTAMP WITH TIME ZONE,
            "enabled" BOOLEAN DEFAULT true,
            "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
          );
        `);
        log.info('ApiKeys table created (or already exists)');
      } catch (error: any) {
        if (error.message.includes('already exists')) {
          log.info('ApiKeys table already exists - skipping creation');
        } else {
          throw error;
        }
      }
    }

    log.info('ApiKeys table migration completed successfully');
  } catch (error) {
    log.error('ApiKeys table migration failed', { error });
    throw error;
  }
}
