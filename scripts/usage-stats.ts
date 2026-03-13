// Script to check current API usage statistics
// This can be used for monitoring and debugging request tracking
import { User } from '../src/models/user';
import { log } from '../src/utils/logger';
import { sequelize } from '../src/sequelize';
import { config } from 'dotenv';

// Load environment variables
config();

async function main() {
  log.info('Checking API usage statistics...');
  
  try {
    // Initialize database connection
    await sequelize.sync();
    log.info('Database connection established');
    
    // Get total user count
    const totalUsers = await User.count();
    
    // Get users with request counts
    const usersWithRequests = await User.findAll({
      attributes: ['id', 'email', 'requestsThisMonth', 'subscriptionStatus'],
      order: [['requestsThisMonth', 'DESC']]
    });
    
    // Calculate statistics
    const totalRequests = usersWithRequests.reduce((sum, user) => {
      const requests = user.getDataValue ? 
        user.getDataValue('requestsThisMonth') : user.requestsThisMonth;
      return sum + (requests || 0);
    }, 0);
    
    const usersWithActivity = usersWithRequests.filter(user => {
      const requests = user.getDataValue ? 
        user.getDataValue('requestsThisMonth') : user.requestsThisMonth;
      return requests > 0;
    });
    
    const avgRequestsPerActiveUser = usersWithActivity.length > 0 ? 
      totalRequests / usersWithActivity.length : 0;
    
    // Display statistics
    log.info('=== API Usage Statistics ===');
    log.info(`Total users: ${totalUsers}`);
    log.info(`Users with activity this month: ${usersWithActivity.length}`);
    log.info(`Total requests this month: ${totalRequests}`);
    log.info(`Average requests per active user: ${avgRequestsPerActiveUser.toFixed(2)}`);
    
    // Show top 10 users by request count
    const top10 = usersWithRequests.slice(0, 10);
    if (top10.length > 0) {
      log.info('\n=== Top 10 Users by Request Count ===');
      top10.forEach((user, index) => {
        const id = user.getDataValue ? user.getDataValue('id') : user.id;
        const requests = user.getDataValue ? user.getDataValue('requestsThisMonth') : user.requestsThisMonth;
        const status = user.getDataValue ? user.getDataValue('subscriptionStatus') : user.subscriptionStatus;
        log.info(`${index + 1}. User ${id}: ${requests} requests (${status || 'free'})`);
      });
    }
    
    // Check if reset is needed (if it's the 1st of the month and there are non-zero counts)
    const now = new Date();
    if (now.getDate() === 1 && totalRequests > 0) {
      log.info('\n⚠️  It\'s the 1st of the month but there are still non-zero request counts.');
      log.info('You may want to run: npm run reset-requests');
    }
    
  } catch (error) {
    log.error(`Error checking API usage: ${error}`);
    process.exit(1);  } finally {
    // Close database connection if it's a real Sequelize instance
    if ('close' in sequelize && typeof sequelize.close === 'function') {
      await sequelize.close();
      log.info('Database connection closed');
    }
  }
  
  process.exit(0);
}

main();
