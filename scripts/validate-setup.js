#!/usr/bin/env node

/**
 * @file validate-setup.js
 * @description Validates your test environment setup before running tests
 */

const fs = require('fs');
const path = require('path');

async function validateSetup() {
  console.log('Validating Foundry REST API Test Setup...\n');

  // Check for .env.test file
  const envPath = path.join(process.cwd(), '.env.test');
  if (!fs.existsSync(envPath)) {
    console.error('Missing .env.test file');
    console.log('   → Copy .env.test.example to .env.test and configure your settings');
    process.exit(1);
  }
  console.log('Found .env.test file');

  // Load environment variables
  require('dotenv').config({ path: envPath });

  // Validate required environment variables
  const required = ['TEST_BASE_URL', 'TEST_API_KEY'];
  for (const env of required) {
    if (!process.env[env]) {
      console.error(`Missing required environment variable: ${env}`);
      process.exit(1);
    }
  }
  console.log('Required environment variables present');

  // Check authentication credentials
  const username = process.env.FOUNDRY_USERNAME;
  const password = process.env.FOUNDRY_PASSWORD;
  
  if (!username) {
    console.error('Missing Foundry username: set FOUNDRY_USERNAME');
    process.exit(1);
  }
  if (!password) {
    console.error('Missing Foundry password: set FOUNDRY_PASSWORD');
    process.exit(1);
  }
  console.log(`Foundry credentials configured for user: ${username}`);

  // Check for Foundry instance URLs
  const foundryUrls = [];
  for (const [key, value] of Object.entries(process.env)) {
    if (key.startsWith('FOUNDRY_V') && key.endsWith('_URL')) {
      foundryUrls.push({ version: key.replace('FOUNDRY_V', '').replace('_URL', ''), url: value });
    }
  }

  if (foundryUrls.length === 0) {
    console.error('No Foundry instance URLs found');
    console.log('   → Set FOUNDRY_V12_URL, FOUNDRY_V13_URL, etc. in your .env.test');
    process.exit(1);
  }

  console.log(`Found ${foundryUrls.length} Foundry instance(s):`);
  foundryUrls.forEach(({ version, url }) => {
    const worldKey = `FOUNDRY_V${version}_WORLD`;
    const worldName = process.env[worldKey] || process.env.TEST_DEFAULT_WORLD || 'test-world';
    console.log(`   → v${version}: ${url}`);
    console.log(`      World: ${worldName}${!process.env[worldKey] ? ' (using default)' : ''}`);
  });

  // Test relay server connection
  console.log('\nTesting relay server connection...');
  try {
    const response = await fetch(`${process.env.TEST_BASE_URL}/api/status`);
    if (response.ok) {
      const data = await response.json();
      console.log(`Relay server responding (version: ${data.version})`);
    } else {
      console.error(`Relay server returned status ${response.status}`);
      console.log('   → Make sure your relay server is running: pnpm run dev');
      process.exit(1);
    }
  } catch (error) {
    console.error('Cannot connect to relay server:', error.message);
    console.log('   → Make sure your relay server is running: pnpm run dev');
    process.exit(1);
  }

  // Test Foundry instances
  console.log('\nTesting Foundry instances...');
  for (const { version, url } of foundryUrls) {
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 5000);
      const response = await fetch(url, { signal: controller.signal });
      clearTimeout(timeoutId);
      if (response.ok) {
        console.log(`Foundry v${version} responding at ${url}`);
      } else {
        console.log(`Foundry v${version} returned status ${response.status} at ${url}`);
      }
    } catch (error) {
      console.log(`Cannot connect to Foundry v${version} at ${url}: ${error.message}`);
      console.log(`   → Make sure Foundry v${version} is running`);
    }
  }

  // Test authentication to relay
  console.log('\nTesting API authentication...');
  try {
    const response = await fetch(`${process.env.TEST_BASE_URL}/clients`, {
      headers: {
        'x-api-key': process.env.TEST_API_KEY
      }
    });

    if (response.ok) {
      console.log('API authentication successful');
    } else if (response.status === 401) {
      console.error('API authentication failed - check your TEST_API_KEY');
      process.exit(1);
    } else {
      console.log(`API returned status ${response.status}`);
    }
  } catch (error) {
    console.error('Error testing API authentication:', error.message);
  }

  console.log('\nSetup validation complete!');
  console.log('\nNext steps for testing:');
  console.log('1. Make sure all Foundry instances are running');
  console.log('2. Ensure you have worlds created for each version (or use defaults)');
  console.log('3. Verify credentials can log into Foundry instances');
  console.log('4. Run tests: pnpm test');
  console.log('\nOptional: Set version-specific worlds in .env.test:');
  foundryUrls.forEach(({ version }) => {
    if (!process.env[`FOUNDRY_V${version}_WORLD`]) {
      console.log(`   FOUNDRY_V${version}_WORLD=my-world-name`);
    }
  });
}

// Add fetch polyfill for older Node.js versions
if (!global.fetch) {
  global.fetch = require('node-fetch');
}

validateSetup().catch(console.error);