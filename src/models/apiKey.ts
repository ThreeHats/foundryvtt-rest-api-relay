import { Model, DataTypes, Sequelize, Op } from 'sequelize';
import { sequelize } from '../sequelize';
import { log } from '../utils/logger';
import crypto from 'crypto';

const isMemoryStore = process.env.DB_TYPE === 'memory';

// In-memory stores for API keys
const memoryApiKeys: Map<string, any> = new Map(); // key -> record
const memoryApiKeysByUser: Map<number, any[]> = new Map(); // userId -> records[]
let memoryIdCounter = 0;

function createMemoryRecord(data: any): any {
  const record: any = {
    id: ++memoryIdCounter,
    userId: data.userId,
    key: data.key || crypto.randomBytes(32).toString('hex'),
    name: data.name,
    scopedClientId: data.scopedClientId || null,
    scopedUserId: data.scopedUserId || null,
    dailyLimit: data.dailyLimit ?? null,
    requestsToday: data.requestsToday || 0,
    lastRequestDate: data.lastRequestDate || null,
    foundryUrl: data.foundryUrl || null,
    foundryUsername: data.foundryUsername || null,
    encryptedFoundryPassword: data.encryptedFoundryPassword || null,
    passwordIv: data.passwordIv || null,
    passwordAuthTag: data.passwordAuthTag || null,
    expiresAt: data.expiresAt || null,
    enabled: data.enabled !== undefined ? data.enabled : true,
    createdAt: new Date(),
    updatedAt: new Date(),
    getDataValue(key: string): any { return (this as any)[key]; },
    setDataValue(key: string, value: any): void { (this as any)[key] = value; },
    save: async function () {
      this.updatedAt = new Date();
      memoryApiKeys.set(this.key, this);
      return this;
    },
    update: async function (values: any) {
      Object.assign(this, values, { updatedAt: new Date() });
      memoryApiKeys.set(this.key, this);
      return this;
    },
    destroy: async function () {
      memoryApiKeys.delete(this.key);
      const userKeys = memoryApiKeysByUser.get(this.userId);
      if (userKeys) {
        const idx = userKeys.findIndex((k: any) => k.id === this.id);
        if (idx !== -1) userKeys.splice(idx, 1);
        if (userKeys.length === 0) memoryApiKeysByUser.delete(this.userId);
      }
    }
  };
  return record;
}

export class ApiKey extends Model {
  declare id: number;
  declare userId: number;
  declare key: string;
  declare name: string;
  declare scopedClientId: string | null;
  declare scopedUserId: string | null;
  declare dailyLimit: number | null;
  declare requestsToday: number;
  declare lastRequestDate: Date | null;
  declare foundryUrl: string | null;
  declare foundryUsername: string | null;
  declare encryptedFoundryPassword: string | null;
  declare passwordIv: string | null;
  declare passwordAuthTag: string | null;
  declare expiresAt: Date | null;
  declare enabled: boolean;
  declare createdAt: Date;
  declare updatedAt: Date;

  static async findOne(options: any): Promise<any> {
    if (isMemoryStore) {
      const tokens = Array.from(memoryApiKeys.values());
      if (options.where) {
        return tokens.find(t => {
          for (const k of Object.keys(options.where)) {
            if (t[k] !== options.where[k]) return false;
          }
          return true;
        }) || null;
      }
      return null;
    }
    return super.findOne(options);
  }

  static async create(data: any): Promise<any> {
    if (isMemoryStore) {
      const record = createMemoryRecord(data);
      memoryApiKeys.set(record.key, record);
      if (!memoryApiKeysByUser.has(record.userId)) {
        memoryApiKeysByUser.set(record.userId, []);
      }
      memoryApiKeysByUser.get(record.userId)!.push(record);
      return record;
    }
    return super.create(data);
  }

  static async findAll(options: any): Promise<any[]> {
    if (isMemoryStore) {
      if (options?.where?.userId !== undefined) {
        return memoryApiKeysByUser.get(options.where.userId) || [];
      }
      return Array.from(memoryApiKeys.values());
    }
    return super.findAll(options);
  }

  /**
   * Look up a scoped API key by its token value
   */
  static async findByKey(key: string): Promise<any> {
    if (isMemoryStore) {
      return memoryApiKeys.get(key) || null;
    }
    return ApiKey.findOne({ where: { key } });
  }

  /**
   * Find all API keys owned by a user
   */
  static async findAllByUser(userId: number): Promise<any[]> {
    if (isMemoryStore) {
      return memoryApiKeysByUser.get(userId) || [];
    }
    return ApiKey.findAll({ where: { userId } });
  }

  /**
   * Delete all API keys owned by a user (for master key regen cascade)
   */
  static async deleteAllByUser(userId: number): Promise<number> {
    if (isMemoryStore) {
      const userKeys = memoryApiKeysByUser.get(userId) || [];
      const count = userKeys.length;
      for (const key of userKeys) {
        memoryApiKeys.delete(key.key);
      }
      memoryApiKeysByUser.delete(userId);
      return count;
    }
    return ApiKey.destroy({ where: { userId } });
  }

  /**
   * Reset requestsToday for all API keys (cron job)
   */
  static async resetDailyCounters(): Promise<void> {
    if (isMemoryStore) {
      for (const record of memoryApiKeys.values()) {
        record.requestsToday = 0;
        record.lastRequestDate = null;
      }
      return;
    }
    await ApiKey.update(
      { requestsToday: 0, lastRequestDate: null },
      { where: {} }
    );
  }
}

if (!isMemoryStore) {
  ApiKey.init({
    id: {
      type: DataTypes.INTEGER,
      autoIncrement: true,
      primaryKey: true
    },
    userId: {
      type: DataTypes.INTEGER,
      allowNull: false
    },
    key: {
      type: DataTypes.STRING(64),
      allowNull: false,
      unique: true,
      defaultValue: () => crypto.randomBytes(32).toString('hex')
    },
    name: {
      type: DataTypes.STRING,
      allowNull: false
    },
    scopedClientId: {
      type: DataTypes.STRING,
      allowNull: true
    },
    scopedUserId: {
      type: DataTypes.STRING,
      allowNull: true
    },
    dailyLimit: {
      type: DataTypes.INTEGER,
      allowNull: true
    },
    requestsToday: {
      type: DataTypes.INTEGER,
      defaultValue: 0
    },
    lastRequestDate: {
      type: DataTypes.DATEONLY,
      allowNull: true
    },
    foundryUrl: {
      type: DataTypes.STRING,
      allowNull: true
    },
    foundryUsername: {
      type: DataTypes.STRING,
      allowNull: true
    },
    encryptedFoundryPassword: {
      type: DataTypes.TEXT,
      allowNull: true
    },
    passwordIv: {
      type: DataTypes.STRING,
      allowNull: true
    },
    passwordAuthTag: {
      type: DataTypes.STRING,
      allowNull: true
    },
    expiresAt: {
      type: DataTypes.DATE,
      allowNull: true
    },
    enabled: {
      type: DataTypes.BOOLEAN,
      defaultValue: true
    }
  }, {
    sequelize: sequelize as Sequelize,
    modelName: 'ApiKey',
    tableName: 'ApiKeys'
  });
}

export default ApiKey;
