#!/usr/bin/env tsx
/**
 * Run Jest tests in a specific order
 * Runs all tests in a SINGLE Jest process to preserve state between test files
 * Uses a custom test sequencer to enforce the order
 */

import { spawn } from 'child_process';
import path from 'path';
import fs from 'fs';
import { TEST_ORDER } from '../tests/helpers/testSequencer';

/** Recursively find all test files under a directory, indexed by basename */
function indexTestFiles(dir: string): Map<string, string> {
  const index = new Map<string, string>();
  function walk(d: string) {
    for (const entry of fs.readdirSync(d, { withFileTypes: true })) {
      const full = path.join(d, entry.name);
      if (entry.isDirectory()) walk(full);
      else if (entry.name.endsWith('.test.ts')) index.set(entry.name, full);
    }
  }
  walk(dir);
  return index;
}

async function runAllTestsInOrder(): Promise<void> {
  console.log('🚀 Starting ordered test execution (single Jest process)...\n');

  // Index all test files by filename (supports subdirectories)
  const fileIndex = indexTestFiles(path.join('tests', 'integration'));
  const testFiles = TEST_ORDER.map(name => {
    const fullPath = fileIndex.get(name);
    if (!fullPath) {
      console.warn(` Test file not found: ${name}`);
      return null;
    }
    return fullPath;
  }).filter((f): f is string => f !== null);
  
  console.log(`Running ${testFiles.length} test files in order:\n`);
  testFiles.forEach((file, i) => {
    console.log(`  ${i + 1}. ${file.split('/').pop()}`);
  });
  console.log('');
  
  // Path to custom sequencer (use TS file - Jest will handle via ts-jest)
  const sequencerPath = path.resolve(__dirname, '../tests/helpers/testSequencer.ts');
  
  return new Promise((resolve, reject) => {
    // Run Jest with custom sequencer to enforce test order
    const jestProcess = spawn('jest', [
      '--runInBand',                          // Run tests sequentially in one process
      '--no-cache',
      `--testSequencer=${sequencerPath}`,     // Use our custom sequencer
      ...testFiles                            // Pass all test files
    ], {
      stdio: 'inherit',
      shell: true
    });

    jestProcess.on('close', (code) => {
      console.log('\n' + '─'.repeat(60));
      if (code !== 0) {
        console.log(`\n❌ Tests failed with exit code ${code}\n`);
        reject(new Error(`Tests failed with exit code ${code}`));
      } else {
        console.log('\n✨ All tests completed successfully!\n');
        resolve();
      }
    });

    jestProcess.on('error', (error) => {
      console.error('Error running Jest:', error);
      reject(error);
    });
  });
}

runAllTestsInOrder().catch(error => {
  console.error('Error running tests:', error);
  process.exit(1);
});
