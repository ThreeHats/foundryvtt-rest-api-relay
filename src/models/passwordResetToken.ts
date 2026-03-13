import { Model, DataTypes, Sequelize, Op } from 'sequelize';
import { sequelize } from '../sequelize';
import { log } from '../utils/logger';

const isMemoryStore = process.env.DB_TYPE === 'memory';

// In-memory store for password reset tokens
const memoryTokens: Map<string, any> = new Map();

export class PasswordResetToken extends Model {
  declare id: number;
  declare userId: number;
  declare tokenHash: string;
  declare expiresAt: Date;
  declare used: boolean;
  declare createdAt: Date;
  declare updatedAt: Date;

  static async findOne(options: any): Promise<any> {
    if (isMemoryStore) {
      const tokens = Array.from(memoryTokens.values());
      if (options.where) {
        return tokens.find(t => {
          for (const key of Object.keys(options.where)) {
            if (key === 'expiresAt' && options.where[key]?.[Op.gt]) {
              if (!(new Date(t.expiresAt) > options.where[key][Op.gt])) return false;
            } else if (t[key] !== options.where[key]) {
              return false;
            }
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
      const token: any = {
        id: memoryTokens.size + 1,
        userId: data.userId,
        tokenHash: data.tokenHash,
        expiresAt: data.expiresAt,
        used: data.used || false,
        createdAt: new Date(),
        updatedAt: new Date(),
        getDataValue: function (key: string): any { return (this as any)[key]; },
        setDataValue: function (key: string, value: any): void { (this as any)[key] = value; },
        save: async function () {
          this.updatedAt = new Date();
          memoryTokens.set(String(this.id), this);
          return this;
        },
        update: async function (values: any) {
          Object.assign(this, values, { updatedAt: new Date() });
          memoryTokens.set(String(this.id), this);
          return this;
        }
      };
      memoryTokens.set(String(token.id), token);
      return token;
    }
    return super.create(data);
  }

  /**
   * Invalidate all unused tokens for a given user
   */
  static async invalidateForUser(userId: number): Promise<void> {
    if (isMemoryStore) {
      for (const token of memoryTokens.values()) {
        if (token.userId === userId && !token.used) {
          token.used = true;
          token.updatedAt = new Date();
        }
      }
      return;
    }
    await PasswordResetToken.update(
      { used: true },
      { where: { userId, used: false } }
    );
  }

  /**
   * Clean up expired/used tokens older than 24 hours
   */
  static async cleanupExpired(): Promise<number> {
    const cutoff = new Date(Date.now() - 24 * 60 * 60 * 1000);
    if (isMemoryStore) {
      let count = 0;
      for (const [key, token] of memoryTokens.entries()) {
        if ((token.used || new Date(token.expiresAt) < new Date()) && new Date(token.createdAt) < cutoff) {
          memoryTokens.delete(key);
          count++;
        }
      }
      return count;
    }
    const deleted = await PasswordResetToken.destroy({
      where: {
        [Op.or]: [
          { used: true },
          { expiresAt: { [Op.lt]: new Date() } }
        ],
        createdAt: { [Op.lt]: cutoff }
      }
    });
    return deleted;
  }
}

if (!isMemoryStore) {
  PasswordResetToken.init({
    id: {
      type: DataTypes.INTEGER,
      autoIncrement: true,
      primaryKey: true
    },
    userId: {
      type: DataTypes.INTEGER,
      allowNull: false
    },
    tokenHash: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true
    },
    expiresAt: {
      type: DataTypes.DATE,
      allowNull: false
    },
    used: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: false
    }
  }, {
    sequelize: sequelize as Sequelize,
    modelName: 'PasswordResetToken',
    tableName: 'PasswordResetTokens'
  });
}

export default PasswordResetToken;
