/**
 * Test environment variables
 * Centralized place for all test configuration
 */

export const testVariables = {
  baseUrl: process.env.TEST_BASE_URL || 'http://localhost:3010',
  apiKey: process.env.TEST_API_KEY || 'test-api-key',
  
  // These will be set during test execution
  clientId: '', // Set from global session
  foundryUrl: 'http://localhost:30013',
  worldName: 'testing',
  username: 'Gamemaster',
  
  // Test data UUIDs (will be populated as entities are created)
  testActorUuid: '',
  testItemUuid: '',
  testSceneUuid: '',
  testMacroUuid: ''
};

/**
 * Update a variable value
 */
export function setVariable(key: keyof typeof testVariables, value: string): void {
  testVariables[key] = value;
}

/**
 * Get all variables as a record
 */
export function getVariables(): Record<string, string> {
  return { ...testVariables };
}
