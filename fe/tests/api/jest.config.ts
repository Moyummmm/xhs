/** @type {import('ts-jest').JestConfigWithTsJest} */
export default {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>'],
  testMatch: ['**/tests/api/tests/**/*.test.ts'],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  collectCoverageFrom: [
    'src/api/**/*.ts',
    '!src/**/*.d.ts',
  ],
  coverageDirectory: 'coverage/api',
  setupFilesAfterEnv: ['<rootDir>/tests/api/setup.ts'],
  testTimeout: 30000,
};
