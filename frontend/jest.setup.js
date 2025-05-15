// Import jest-dom to add custom jest matchers for asserting on DOM nodes
import '@testing-library/jest-dom';

// Mock the next/router
jest.mock('next/router', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    prefetch: jest.fn(),
    back: jest.fn(),
    pathname: '/',
    query: {},
  }),
}));

// Mock for styleMock
global.__mocks__ = {
  styleMock: {},
  fileMock: 'test-file-stub',
};

// Setup mock directories
if (typeof window === 'undefined') {
  const { mkdir, writeFile } = require('fs').promises;
  const path = require('path');
  
  const createDir = async (dirPath) => {
    try {
      await mkdir(dirPath, { recursive: true });
    } catch (error) {
      // Directory might already exist
    }
  };
  
  // Create mock directories and files
  Promise.all([
    createDir(path.join(process.cwd(), '__mocks__')),
    writeFile(path.join(process.cwd(), '__mocks__/styleMock.js'), 'module.exports = {};'),
    writeFile(path.join(process.cwd(), '__mocks__/fileMock.js'), 'module.exports = "test-file-stub";'),
  ]).catch(console.error);
}