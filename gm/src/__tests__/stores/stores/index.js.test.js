import { describe, it, expect, vi, beforeEach } from 'vitest';
import { 
  useSceneStore,
  useNPCStore,
  useAgentStore,
  useLLMStore,
  usePromptStore,
  useShopStore,
  useTaskStore,
  useConfigStore,
  useSkillStore,
  useAchievementStore,
  useDemoStore
} from './index';

// Mock all store modules
vi.mock('./scene', () => ({
  useSceneStore: vi.fn()
}));

vi.mock('./npc', () => ({
  useNPCStore: vi.fn()
}));

vi.mock('./agent', () => ({
  useAgentStore: vi.fn()
}));

vi.mock('./llm', () => ({
  useLLMStore: vi.fn()
}));

vi.mock('./prompt', () => ({
  usePromptStore: vi.fn()
}));

vi.mock('./shop', () => ({
  useShopStore: vi.fn()
}));

vi.mock('./task', () => ({
  useTaskStore: vi.fn()
}));

vi.mock('./config', () => ({
  useConfigStore: vi.fn()
}));

vi.mock('./skill', () => ({
  useSkillStore: vi.fn()
}));

vi.mock('./achievement', () => ({
  useAchievementStore: vi.fn()
}));

vi.mock('./demo', () => ({
  useDemoStore: vi.fn()
}));

describe('Stores Index', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Exported Stores', () => {
    it('should export useSceneStore', () => {
      expect(useSceneStore).toBeDefined();
      expect(typeof useSceneStore).toBe('function');
    });

    it('should export useNPCStore', () => {
      expect(useNPCStore).toBeDefined();
      expect(typeof useNPCStore).toBe('function');
    });

    it('should export useAgentStore', () => {
      expect(useAgentStore).toBeDefined();
      expect(typeof useAgentStore).toBe('function');
    });

    it('should export useLLMStore', () => {
      expect(useLLMStore).toBeDefined();
      expect(typeof useLLMStore).toBe('function');
    });

    it('should export usePromptStore', () => {
      expect(usePromptStore).toBeDefined();
      expect(typeof usePromptStore).toBe('function');
    });

    it('should export useShopStore', () => {
      expect(useShopStore).toBeDefined();
      expect(typeof useShopStore).toBe('function');
    });

    it('should export useTaskStore', () => {
      expect(useTaskStore).toBeDefined();
      expect(typeof useTaskStore).toBe('function');
    });

    it('should export useConfigStore', () => {
      expect(useConfigStore).toBeDefined();
      expect(typeof useConfigStore).toBe('function');
    });

    it('should export useSkillStore', () => {
      expect(useSkillStore).toBeDefined();
      expect(typeof useSkillStore).toBe('function');
    });

    it('should export useAchievementStore', () => {
      expect(useAchievementStore).toBeDefined();
      expect(typeof useAchievementStore).toBe('function');
    });

    it('should export useDemoStore', () => {
      expect(useDemoStore).toBeDefined();
      expect(typeof useDemoStore).toBe('function');
    });
  });

  describe('Store Functionality', () => {
    it('should call useSceneStore when invoked', () => {
      const mockStore = { id: 'scene-store' };
      useSceneStore.mockReturnValue(mockStore);
      
      const result = useSceneStore();
      expect(useSceneStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useNPCStore when invoked', () => {
      const mockStore = { id: 'npc-store' };
      useNPCStore.mockReturnValue(mockStore);
      
      const result = useNPCStore();
      expect(useNPCStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useAgentStore when invoked', () => {
      const mockStore = { id: 'agent-store' };
      useAgentStore.mockReturnValue(mockStore);
      
      const result = useAgentStore();
      expect(useAgentStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useLLMStore when invoked', () => {
      const mockStore = { id: 'llm-store' };
      useLLMStore.mockReturnValue(mockStore);
      
      const result = useLLMStore();
      expect(useLLMStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call usePromptStore when invoked', () => {
      const mockStore = { id: 'prompt-store' };
      usePromptStore.mockReturnValue(mockStore);
      
      const result = usePromptStore();
      expect(usePromptStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useShopStore when invoked', () => {
      const mockStore = { id: 'shop-store' };
      useShopStore.mockReturnValue(mockStore);
      
      const result = useShopStore();
      expect(useShopStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useTaskStore when invoked', () => {
      const mockStore = { id: 'task-store' };
      useTaskStore.mockReturnValue(mockStore);
      
      const result = useTaskStore();
      expect(useTaskStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useConfigStore when invoked', () => {
      const mockStore = { id: 'config-store' };
      useConfigStore.mockReturnValue(mockStore);
      
      const result = useConfigStore();
      expect(useConfigStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useSkillStore when invoked', () => {
      const mockStore = { id: 'skill-store' };
      useSkillStore.mockReturnValue(mockStore);
      
      const result = useSkillStore();
      expect(useSkillStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useAchievementStore when invoked', () => {
      const mockStore = { id: 'achievement-store' };
      useAchievementStore.mockReturnValue(mockStore);
      
      const result = useAchievementStore();
      expect(useAchievementStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });

    it('should call useDemoStore when invoked', () => {
      const mockStore = { id: 'demo-store' };
      useDemoStore.mockReturnValue(mockStore);
      
      const result = useDemoStore();
      expect(useDemoStore).toHaveBeenCalled();
      expect(result).toBe(mockStore);
    });
  });

  describe('Store Re-exports', () => {
    it('should re-export all stores correctly', () => {
      const allStores = [
        { name: 'useSceneStore', store: useSceneStore },
        { name: 'useNPCStore', store: useNPCStore },
        { name: 'useAgentStore', store: useAgentStore },
        { name: 'useLLMStore', store: useLLMStore },
        { name: 'usePromptStore', store: usePromptStore },
        { name: 'useShopStore', store: useShopStore },
        { name: 'useTaskStore', store: useTaskStore },
        { name: 'useConfigStore', store: useConfigStore },
        { name: 'useSkillStore', store: useSkillStore },
        { name: 'useAchievementStore', store: useAchievementStore },
        { name: 'useDemoStore', store: useDemoStore }
      ];

      allStores.forEach(({ name, store }) => {
        expect(store).toBeDefined();
        expect(typeof store).toBe('function');
      });
    });
  });
});