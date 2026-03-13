import { sequelize } from '../sequelize';
import { User } from './user';
import { PasswordResetToken } from './passwordResetToken';
import crypto from 'crypto';
import { log } from '../utils/logger';

async function initializeDatabase() {
  try {
    log.info('Starting database initialization...');
    log.info('Using database', { databaseUrl: process.env.DATABASE_URL });
    
    // Test the connection first
    await sequelize.authenticate();
    log.info('Database connection has been established successfully.');
    
    // Sync all models - this creates or alters tables to match models
    log.info('Syncing database models...');
    await sequelize.sync({ alter: true });
    log.info('Database models synchronized.');
    
    // Only create admin user if ADMIN_EMAIL and ADMIN_PASSWORD are provided
    const adminEmail = process.env.ADMIN_EMAIL;
    const adminPassword = process.env.ADMIN_PASSWORD;
    
    if (adminEmail && adminPassword) {
      log.info('Creating admin user from environment variables...');
      
      const user = await User.create({
        email: adminEmail,
        password: adminPassword, // Will be hashed by the beforeCreate hook
        apiKey: crypto.randomBytes(32).toString('hex'),
        requestsThisMonth: 0
      });
      
      log.info('Admin user created successfully');
    } else {
      log.info('No ADMIN_EMAIL/ADMIN_PASSWORD provided - skipping admin user creation');
      log.info('Users can register via /auth/register endpoint');
    }
    
    log.info('Database initialization complete!');
    return true;
  } catch (error) {
    log.error('Database initialization failed', { error });
    return false;
  }
}

// Run the function if this script is executed directly
if (require.main === module) {
  initializeDatabase()
    .then((result) => {
      log.info(`Initialization ${result ? 'succeeded' : 'failed'}`);
      process.exit(result ? 0 : 1);
    })
    .catch(error => {
      log.error('Failed to initialize database', { error });
      process.exit(1);
    });
}

export { initializeDatabase };