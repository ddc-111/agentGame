import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { createPinia, setActivePinia, defineStore } from 'pinia';

// Mock localStorage using vi.stubGlobal for better compatibility
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};

// 在外部定义store，避免在测试内部重复定义
const useTestStore = defineStore('test', {
  state: () => ({
    count: 0,
  }),
  actions: {
    increment() {
      this.count++;
    },
  },
});

describe('Store Tests', () => {
  let pinia;

  beforeEach(() => {
    // Create a new Pinia instance and set it as active before each test
    pinia = createPinia();
    setActivePinia(pinia);

    // Mock localStorage globally
    vi.stubGlobal('localStorage', localStorageMock);
    // Clear all mocks before each test
    vi.clearAllMocks();
  });

  afterEach(() => {
    // Unstub localStorage to clean up
    vi.unstubAllGlobals();
    // Reset the active Pinia instance to avoid contamination
    setActivePinia(null);
    // Reset Pinia (optional, but good practice)
    pinia = null;
  });

  it('should increment count in test store', () => {
    // 直接使用在外部定义的store
    const store = useTestStore();
    expect(store.count).toBe(0);
    store.increment();
    expect(store.count).toBe(1);
  });

  it('should handle localStorage mock', () => {
    // Example test that uses localStorage
    localStorage.getItem('testKey');
    expect(localStorage.getItem).toHaveBeenCalledWith('testKey');
  });
});
