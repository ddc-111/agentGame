import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAgentStore } from '@/stores/agent';

describe('Agent Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useAgentStore();
  });

  it('should have initial agents', () => {
    expect(store.agents).toHaveLength(2);
    expect(store.agents[0].id).toBe('agent_001');
    expect(store.agents[1].id).toBe('agent_002');
  });

  it('should add a new agent', () => {
    const initialLength = store.agents.length;
    store.addAgent({
      name: '新智能体',
      description: '测试智能体',
      llmConfig: {
        provider: 'openai',
        model: 'gpt-4',
        temperature: 0.7,
        maxTokens: 500
      },
      systemPrompt: '测试提示词',
      memoryConfig: {
        type: 'sliding_window',
        maxMessages: 20,
        summaryEnabled: true,
        summaryThreshold: 50
      },
      knowledgeBase: [],
      tools: []
    });
    expect(store.agents).toHaveLength(initialLength + 1);
    expect(store.agents[store.agents.length - 1].name).toBe('新智能体');
    expect(store.agents[store.agents.length - 1].id).toMatch(/^agent_\d+$/);
  });

  it('should update an existing agent', () => {
    store.updateAgent('agent_001', { name: '更新后的智能体' });
    const agent = store.getAgentById('agent_001');
    expect(agent.name).toBe('更新后的智能体');
  });

  it('should not update non-existent agent', () => {
    const initialLength = store.agents.length;
    store.updateAgent('non_existent', { name: 'test' });
    expect(store.agents).toHaveLength(initialLength);
  });

  it('should delete an agent', () => {
    const initialLength = store.agents.length;
    store.deleteAgent('agent_001');
    expect(store.agents).toHaveLength(initialLength - 1);
    expect(store.getAgentById('agent_001')).toBeUndefined();
  });

  it('should get agent by id', () => {
    const agent = store.getAgentById('agent_002');
    expect(agent).toBeDefined();
    expect(agent.name).toBe('王大娘智能体');
  });

  it('should return undefined for non-existent agent', () => {
    const agent = store.getAgentById('non_existent');
    expect(agent).toBeUndefined();
  });

  it('should have currentAgent initially null', () => {
    expect(store.currentAgent).toBeNull();
  });
});
