module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/tests'],
  testMatch: [
    '**/__tests__/**/*.ts',
    '**/?(*.)+(spec|test).ts'
  ],
  transform: {
    '^.+\\.ts$': ['ts-jest', {
      tsconfig: 'tsconfig.json'
    }],
  },
  collectCoverageFrom: [
    'src/**/*.ts',
    '!src/**/*.d.ts',
    '!src/**/*.test.ts',
    '!src/**/index.ts'
  ],
  coverageDirectory: 'coverage',
  coverageReporters: ['text', 'lcov', 'html', 'json'],
  globalTeardown: '<rootDir>/tests/globalTeardown.ts',
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
  testTimeout: 120000, // 2 minutes timeout for integration tests
  maxWorkers: 1, // Run test FILES sequentially
  maxConcurrency: 1, // Run individual TESTS sequentially within each file
  verbose: true,
  forceExit: true, // Required for integration tests with Puppeteer/WebSocket connections
  detectOpenHandles: false, // Axios keep-alive warnings are expected, not leaks
  reporters: [
    'default',
    ['jest-junit', {
      outputDirectory: './test-results',
      outputName: 'junit.xml',
      classNameTemplate: '{classname}',
      titleTemplate: '{title}',
      ancestorSeparator: ' › ',
      usePathForSuiteName: true
    }],
    ['jest-html-reporter', {
      pageTitle: 'Foundry REST API Test Report',
      outputPath: './test-results/test-report.html',
      includeFailureMsg: true,
      includeConsoleLog: true,
      theme: 'darkTheme',
      logo: '',
      dateFormat: 'yyyy-mm-dd HH:MM:ss',
      sort: 'status',
      executionTimeWarningThreshold: 5
    }]
  ],
  // Add types for Jest globals
  moduleFileExtensions: ['ts', 'tsx', 'js']
};