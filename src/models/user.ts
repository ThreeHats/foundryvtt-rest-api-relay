import { Model, DataTypes, Sequelize } from 'sequelize';
import { sequelize } from '../sequelize';
import bcrypt from 'bcryptjs';
import crypto from 'crypto';
import { log } from '../utils/logger';

// Check if we're using the memory store
const isMemoryStore = process.env.DB_TYPE === 'memory';

export class User extends Model {
  // Declare types for TypeScript without public fields to avoid Sequelize conflicts
  declare id: number;
  declare email: string;
  declare password: string;
  declare apiKey: string;
  declare requestsThisMonth: number;
  declare requestsToday: number;
  declare lastRequestDate: Date;
  declare createdAt: Date;
  declare updatedAt: Date;
  declare stripeCustomerId?: string;
  declare subscriptionStatus?: string;
  declare subscriptionId?: string;
  declare subscriptionEndsAt?: Date;
  declare maxHeadlessSessions?: number;

  // Memory store methods
  static async findOne(options: any): Promise<any> {
    if (isMemoryStore) {
      if (options.where && options.where.apiKey) {
        return (sequelize as any).getUser(options.where.apiKey);
      }
      if (options.where && options.where.email) {
        const users = Array.from((sequelize as any).users.values());
        return users.find(u => (u as User).email === options.where.email) || null;
      }
      if (options.where && options.where.id) {
        const users = Array.from((sequelize as any).users.values());
        return users.find(u => (u as any).id === options.where.id) || null;
      }
      return null;
    }
    return super.findOne(options);
  }

  static async create(data: any): Promise<any> {
    if (isMemoryStore) {
      const memoryStore = sequelize as any;
      
      if (memoryStore.users.has(data.email)) {
        throw new Error('User already exists');
      }
      
      const user: any = {
        id: memoryStore.users.size + 1,
        email: data.email,
        password: data.password,
        apiKey: data.apiKey || crypto.randomBytes(32).toString('hex'),
        requestsThisMonth: data.requestsThisMonth || 0,
        requestsToday: data.requestsToday || 0,
        lastRequestDate: data.lastRequestDate || null,
        subscriptionStatus: data.subscriptionStatus || 'free',
        maxHeadlessSessions: data.maxHeadlessSessions || null,
        createdAt: new Date(),
        updatedAt: new Date(),
        getDataValue: function(key: string): any { 
          return (this as any)[key]; 
        },
        setDataValue: function(key: string, value: any): void { 
          (this as any)[key] = value; 
        }
      };
      
      memoryStore.users.set(data.email, user);
      memoryStore.apiKeys.set(user.apiKey, data.email);
      
      return user;
    }
    
    return super.create(data);
  }

  static async findAll(options: any): Promise<any[]> {
    if (isMemoryStore) {
      const memoryStore = sequelize as any;
      if (options.where && options.where.apiKey) {
        const user = memoryStore.getUser(options.where.apiKey);
        return user ? [user] : [];
      }
      return Array.from(memoryStore.users.values());
    }
    return super.findAll(options);
  }
}

// Initialize with Sequelize if not using memory store
if (!isMemoryStore) {
  User.init({
    id: {
      type: DataTypes.INTEGER,
      autoIncrement: true,
      primaryKey: true
    },
    email: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true,
      validate: {
        isEmail: true
      }
    },
    password: {
      type: DataTypes.STRING,
      allowNull: false
    },
    apiKey: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true,
      defaultValue: () => crypto.randomBytes(32).toString('hex')
    },
    requestsThisMonth: {
      type: DataTypes.INTEGER,
      defaultValue: 0
    },
    requestsToday: {
      type: DataTypes.INTEGER,
      defaultValue: 0
    },
    lastRequestDate: {
      type: DataTypes.DATEONLY,
      allowNull: true
    },
    stripeCustomerId: {
      type: DataTypes.STRING,
      allowNull: true
    },
    subscriptionStatus: {
      type: DataTypes.STRING,
      allowNull: true,
      defaultValue: 'free'
    },
    subscriptionId: {
      type: DataTypes.STRING,
      allowNull: true
    },
    subscriptionEndsAt: {
      type: DataTypes.DATE,
      allowNull: true
    },
    maxHeadlessSessions: {
      type: DataTypes.INTEGER,
      allowNull: true
    }
  }, {
    sequelize: sequelize as Sequelize,
    modelName: 'User',
    tableName: 'Users',
    hooks: {
      beforeCreate: async (user) => {
        if (user.getDataValue('password')) {
          const salt = await bcrypt.genSalt(10);
          const hashedPassword = await bcrypt.hash(user.getDataValue('password'), salt);
          user.setDataValue('password', hashedPassword);
        }
      },
      beforeUpdate: async (user) => {
        if (user.changed('password')) {
          log.debug('Updating password for user', { id: user.getDataValue('id') });
          const salt = await bcrypt.genSalt(10);
          const hashedPassword = await bcrypt.hash(user.getDataValue('password'), salt);
          user.setDataValue('password', hashedPassword);
        }
      }
    }
  });
}

export default User;
