// Jest globals and types
import { jest } from '@jest/globals';
import dotenv from 'dotenv';
import path from 'path';

// Load environment variables for testing
dotenv.config({ path: path.join(__dirname, '..', '.env.test') });

// Increase timeout for all tests
jest.setTimeout(120000);

// Verify the relay server is reachable before running any tests
const baseUrl = process.env.TEST_BASE_URL || 'http://localhost:3010';

beforeAll(async () => {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 5000);
    await fetch(`${baseUrl}/health`, { signal: controller.signal });
    clearTimeout(timeout);
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw new Error(
      `\n\nCould not connect to the relay server at ${baseUrl}\n` +
      `   Error: ${message}\n\n` +
      `   Make sure the relay server is running and TEST_BASE_URL is correct in .env.test\n`
    );
  }
}, 10000);

