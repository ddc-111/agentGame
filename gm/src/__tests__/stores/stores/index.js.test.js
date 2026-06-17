```javascript
import { describe, it, expect, vi } from 'vitest';
import { useSceneStore } from './scene';
import { useNPCStore } from './npc';
import { useAgentStore } from './agent';
import { useLLMStore } from './llm';
import { usePromptStore } from './prompt';
import { useShopStore } from './shop';
import { useTaskStore } from './task';
import { useConfigStore } from './config';
import { useSkillStore } from './skill';
import { useAchievementStore } from './achievement';
import { useDemoStore } from './demo';

vi.mock('./scene', () => ({
  useSceneStore: vi.fn(),
}));
vi.mock('./npc', () => ({
  useNPCStore: vi.fn(),
}));
vi.mock('./agent', () => ({
  useAgentStore: vi.fn(),
}));
vi.mock('./llm', () => ({
  useLLMStore: vi.fn(),
}));
vi.mock('./prompt', () => ({
  usePromptStore: vi.fn(),
}));
vi.mock('./shop', () => ({
  useShopStore: vi.fn(),
}));
vi.mock('./task', () => ({
  useTaskStore: vi.fn(),
}));
vi.mock('./config', () => ({
  useConfigStore: vi.fn(),
}));
vi.mock('./skill', () => ({
  useSkillStore: vi.fn(),
}));
vi.mock('./achievement', () => ({
  useAchievementStore: vi.fn(),
}));
vi.mock('./demo', () => ({
  useDemoStore: vi.fn(),
}));

describe('stores/index', () => {
  it('should export useSceneStore from scene module', () => {
    expect(useSceneStore).toBeDefined();
    expect(useSceneStore).toBe(vi.mocked(useSceneStore));
  });

  it('should export useNPCStore from npc module', () => {
    expect(useNPCStore).toBeDefined();
    expect(useNPCStore).toBe(vi.mocked(useNPCStore));
  });

  it('should export useAgentStore from agent module', () => {
    expect(useAgentStore).toBeDefined();
    expect(useAgentStore).toBe(vi.mocked(useAgentStore));
  });

  it('should export useLLMStore from llm module', () => {
    expect(useLLMStore).toBeDefined();
    expect(useLLMStore).toBe(vi.mocked(useLLMStore));
  });

  it('should export usePromptStore from prompt module', () => {
    expect(usePromptStore).toBeDefined();
    expect(usePromptStore).toBe(vi.mocked(usePromptStore));
  });

  it('should export useShopStore from shop module', () => {
    expect(useShopStore).toBeDefined();
    expect(useShopStore).toBe(vi.mocked(useShopStore));
  });

  it('should export useTaskStore from task module', () => {
    expect(useTaskStore).toBeDefined();
    expect(useTaskStore).toBe(vi.mocked(useTaskStore));
  });

  it('should export useConfigStore from config module', () => {
    expect(useConfigStore).toBeDefined();
    expect(useConfigStore).toBe(vi.mocked(useConfigStore));
  });

  it('should export useSkillStore from skill module', () => {
    expect(useSkillStore).toBeDefined();
    expect(useSkillStore).toBe(vi.mocked(useSkillStore));
  });

  it('should export useAchievementStore from achievement module', () => {
    expect(useAchievementStore).toBeDefined();
    expect(useAchievementStore).toBe(vi.mocked(useAchievementStore));
  });

  it('should export useDemoStore from demo module', () => {
    expect(useDemoStore).toBeDefined();
    expect(useDemoStore).toBe(vi.mocked(useDemoStore));
  });
});
```