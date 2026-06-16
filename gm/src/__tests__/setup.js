import { vi } from 'vitest';

globalThis.fetch = vi.fn();

globalThis.localStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
};

globalThis.navigator = {
  clipboard: {
    writeText: vi.fn()
  }
};
