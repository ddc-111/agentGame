import { describe, it, expect, vi } from 'vitest';

vi.mock('@/stores/scene', () => ({
  useSceneStore: vi.fn()
}));

vi.mock('@/stores/npc', () => ({
  useNPCStore: vi.fn()
}));

vi.mock('@/stores/agent', () => ({
  useAgentStore: vi.fn()
}));

vi.mock('@/stores/llm', () => ({
  useLLMStore: vi.fn()
}));

vi.mock('@/stores/prompt', () => ({
  usePromptStore: vi.fn()
}));

vi.mock('@/stores/shop', () => ({
  useShopStore: vi.fn()
}));

vi.mock('@/stores/task', () => ({
  useTaskStore: vi.fn()
}));

vi.mock('@/stores/config', () => ({
  useConfigStore: vi.fn()
}));

vi.mock('@/stores/skill', () => ({
  useSkillStore: vi.fn()
}));

vi.mock('@/stores/achievement', () => ({
  useAchievementStore: vi.fn()
}));

vi.mock('@/stores/demo', () => ({
  useDemoStore: vi.fn()
}));

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
} from '@/stores/index';

describe('stores/index.js', () => {
  it('should export all stores', () => {
    expect(useSceneStore).toBeDefined();
    expect(useNPCStore).toBeDefined();
    expect(useAgentStore).toBeDefined();
    expect(useLLMStore).toBeDefined();
    expect(usePromptStore).toBeDefined();
    expect(useShopStore).toBeDefined();
    expect(useTaskStore).toBeDefined();
    expect(useConfigStore).toBeDefined();
    expect(useSkillStore).toBeDefined();
    expect(useAchievementStore).toBeDefined();
    expect(useDemoStore).toBeDefined();
  });

  it('should export mocked functions', () => {
    expect(typeof useSceneStore).toBe('function');
    expect(typeof useNPCStore).toBe('function');
    expect(typeof useAgentStore).toBe('function');
    expect(typeof useLLMStore).toBe('function');
    expect(typeof usePromptStore).toBe('function');
    expect(typeof useShopStore).toBe('function');
    expect(typeof useTaskStore).toBe('function');
    expect(typeof useConfigStore).toBe('function');
    expect(typeof useSkillStore).toBe('function');
    expect(typeof useAchievementStore).toBe('function');
    expect(typeof useDemoStore).toBe('function');
  });

  it('useSceneStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useSceneStore)).toBe(true);
  });

  it('useNPCStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useNPCStore)).toBe(true);
  });

  it('useAgentStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useAgentStore)).toBe(true);
  });

  it('useLLMStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useLLMStore)).toBe(true);
  });

  it('usePromptStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(usePromptStore)).toBe(true);
  });

  it('useShopStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useShopStore)).toBe(true);
  });

  it('useTaskStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useTaskStore)).toBe(true);
  });

  it('useConfigStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useConfigStore)).toBe(true);
  });

  it('useSkillStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useSkillStore)).toBe(true);
  });

  it('useAchievementStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useAchievementStore)).toBe(true);
  });

  it('useDemoStore should be a vi.fn mock', () => {
    expect(vi.isMockFunction(useDemoStore)).toBe(true);
  });
});
