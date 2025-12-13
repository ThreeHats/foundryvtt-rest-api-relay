// Jest globals and types
import { jest } from '@jest/globals';
import dotenv from 'dotenv';
import path from 'path';

// Load environment variables for testing
dotenv.config({ path: path.join(__dirname, '..', '.env.test') });

// Increase timeout for all tests
jest.setTimeout(120000);

