/**
 * Global teardown - runs ONCE after all test files
 */
import fs from 'fs/promises';
import path from 'path';

const GLOBAL_VARS_FILE = path.join(__dirname, '.global-vars.json');

export default async function globalTeardown() {
  console.log('\n🧹 Global Teardown: Cleaning up...\n');
  
  // Clean up global variables file
  try {
    await fs.unlink(GLOBAL_VARS_FILE);
    console.log('   Removed global variables file');
  } catch (error) {
    // Ignore if file doesn't exist
  }
  
  console.log('\n✅ Global Teardown Complete\n');
}
