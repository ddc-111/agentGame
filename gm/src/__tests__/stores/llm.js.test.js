import { describe, it, expect, beforeEach, vi } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import { useLLMStore } from '../../stores/llm';

describe('useLLMStore', () => {
  let store;
  let pinia;

  beforeEach(() => {
    pinia = createPinia();
    setActivePinia(pinia);
    store = useLLMStore();
  });

  describe('初始状态', () => {
    it('应该有正确的providers初始数据', () => {
      expect(store.providers).toHaveLength(3);
      expect(store.providers[0].id).toBe('openai');
      expect(store.providers[1].id).toBe('anthropic');
      expect(store.providers[2].id).toBe('local');
    });

    it('应该有正确的currentProvider初始值', () => {
      expect(store.currentProvider).toBe('openai');
    });

    it('每个provider应该包含必需的字段', () => {
      store.providers.forEach(provider => {
        expect(provider).toHaveProperty('id');
        expect(provider).toHaveProperty('name');
        expect(provider).toHaveProperty('baseUrl');
        expect(provider).toHaveProperty('apiKey');
        expect(provider).toHaveProperty('models');
        expect(Array.isArray(provider.models)).toBe(true);
      });
    });

    it('每个模型应该包含必需的字段', () => {
      store.providers.forEach(provider => {
        provider.models.forEach(model => {
          expect(model).toHaveProperty('id');
          expect(model).toHaveProperty('name');
          expect(model).toHaveProperty('maxTokens');
          expect(typeof model.maxTokens).toBe('number');
        });
      });
    });
  });

  describe('addProvider', () => {
    it('应该正确添加新的provider', () => {
      const newProvider = {
        name: 'Test Provider',
        baseUrl: 'https://test.com',
        apiKey: 'test-key',
        models: [
          { id: 'test-model-1', name: 'Test Model 1', maxTokens: 1000 }
        ]
      };

      const initialLength = store.providers.length;
      store.addProvider(newProvider);

      expect(store.providers).toHaveLength(initialLength + 1);

      const addedProvider = store.providers[store.providers.length - 1];
      expect(addedProvider.name).toBe(newProvider.name);
      expect(addedProvider.baseUrl).toBe(newProvider.baseUrl);
      expect(addedProvider.apiKey).toBe(newProvider.apiKey);
      expect(addedProvider.models).toEqual(newProvider.models);
      expect(addedProvider.id).toContain('provider_');
    });
  });

  describe('updateProvider', () => {
    it('应该正确更新现有provider', () => {
      const updateData = {
        name: 'Updated OpenAI',
        apiKey: 'new-api-key',
        baseUrl: 'https://new-openai.com/v2'
      };

      store.updateProvider('openai', updateData);

      const updatedProvider = store.getProviderById('openai');
      expect(updatedProvider.name).toBe('Updated OpenAI');
      expect(updatedProvider.apiKey).toBe('new-api-key');
      expect(updatedProvider.baseUrl).toBe('https://new-openai.com/v2');
      expect(updatedProvider.id).toBe('openai');
    });

    it('当id不存在时不应该修改任何provider', () => {
      const initialProviders = JSON.parse(JSON.stringify(store.providers));
      store.updateProvider('non-existent-id', { name: 'Test' });
      expect(store.providers).toEqual(initialProviders);
    });
  });

  describe('deleteProvider', () => {
    it('应该正确删除指定的provider', () => {
      const initialLength = store.providers.length;
      store.deleteProvider('openai');
      expect(store.providers).toHaveLength(initialLength - 1);
      expect(store.getProviderById('openai')).toBeUndefined();
    });
  });
});
