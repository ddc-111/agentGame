```javascript
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useLLMStore } from './stores/llm.js';

// Mock Date.now for predictable provider IDs
vi.spyOn(Date, 'now').mockReturnValue(1234567890);

describe('useLLMStore', () => {
  let store;

  beforeEach(() => {
    // 创建新的Pinia实例并激活
    setActivePinia(createPinia());
    // 创建store实例
    store = useLLMStore();
  });

  describe('initial state', () => {
    it('should have 3 default providers', () => {
      expect(store.providers).toHaveLength(3);
      expect(store.providers.map(p => p.id)).toEqual(['openai', 'anthropic', 'local']);
    });

    it('should have currentProvider set to openai', () => {
      expect(store.currentProvider).toBe('openai');
    });

    it('should have correct openai provider structure', () => {
      const openai = store.providers.find(p => p.id === 'openai');
      expect(openai).toEqual({
        id: 'openai',
        name: 'OpenAI',
        baseUrl: 'https://api.openai.com/v1',
        apiKey: '',
        models: [
          { id: 'gpt-4', name: 'GPT-4', maxTokens: 8192 },
          { id: 'gpt-4-turbo', name: 'GPT-4 Turbo', maxTokens: 128000 },
          { id: 'gpt-3.5-turbo', name: 'GPT-3.5 Turbo', maxTokens: 4096 }
        ]
      });
    });
  });

  describe('addProvider', () => {
    it('should add a new provider with generated id', () => {
      const initialLength = store.providers.length;
      const newProvider = {
        name: 'Test Provider',
        baseUrl: 'https://test.com',
        apiKey: 'test-key',
        models: [{ id: 'test-model', name: 'Test Model', maxTokens: 1000 }]
      };

      store.addProvider(newProvider);

      expect(store.providers).toHaveLength(initialLength + 1);
      const addedProvider = store.providers[store.providers.length - 1];
      expect(addedProvider.id).toBe('provider_1234567890');
      expect(addedProvider.name).toBe('Test Provider');
      expect(addedProvider.baseUrl).toBe('https://test.com');
      expect(addedProvider.apiKey).toBe('test-key');
      expect(addedProvider.models).toEqual([
        { id: 'test-model', name: 'Test Model', maxTokens: 1000 }
      ]);
    });

    it('should preserve existing providers after adding', () => {
      const initialProviders = [...store.providers];
      store.addProvider({ name: 'New Provider' });

      expect(store.providers.slice(0, 3)).toEqual(initialProviders);
    });
  });

  describe('updateProvider', () => {
    it('should update provider fields by id', () => {
      store.updateProvider('openai', {
        apiKey: 'new-api-key',
        name: 'Updated OpenAI'
      });

      const openai = store.providers.find(p => p.id === 'openai');
      expect(openai.apiKey).toBe('new-api-key');
      expect(openai.name).toBe('Updated OpenAI');
    });

    it('should not modify other providers when updating', () => {
      const anthropicBefore = { ...store.providers.find(p => p.id === 'anthropic') };
      
      store.updateProvider('openai', { apiKey: 'test-key' });
      
      const anthropicAfter = store.providers.find(p => p.id === 'anthropic');
      expect(anthropicAfter).toEqual(anthropicBefore);
    });

    it('should not throw when updating non-existent provider', () => {
      const initialProviders = [...store.providers];
      
      store.updateProvider('non-existent', { apiKey: 'test-key' });
      
      expect(store.providers).toEqual(initialProviders);
    });

    it('should merge data without removing existing fields', () => {
      store.updateProvider('openai', { apiKey: 'test-key' });
      
      const openai = store.providers.find(p => p.id === 'openai');
      expect(openai.id).toBe('openai');
      expect(openai.name).toBe('OpenAI');
      expect(openai.baseUrl).toBe('https://api.openai.com/v1');
      expect(openai.apiKey).toBe('test-key');
      expect(openai.models).toEqual([
        { id: 'gpt-4', name: 'GPT-4', maxTokens: 8192 },
        { id: 'gpt-4-turbo', name: 'GPT-4 Turbo', maxTokens: 128000 },
        { id: 'gpt-3.5-turbo', name: 'GPT-3.5 Turbo', maxTokens: 4096 }
      ]);
    });
  });

  describe('deleteProvider', () => {
    it('should remove provider by id', () => {
      const initialLength = store.providers.length;
      
      store.deleteProvider('openai');
      
      expect(store.providers).toHaveLength(initialLength - 1);
      expect(store.providers.find(p => p.id === 'openai')).toBeUndefined();
    });

    it('should not affect other providers when deleting', () => {
      const anthropicBefore = { ...store.providers.find(p => p.id === 'anthropic') };
      const localBefore = { ...store.providers.find(p => p.id === 'local') };
      
      store.deleteProvider('openai');
      
      const anthropicAfter = store.providers.find(p => p.id === 'anthropic');
      const localAfter = store.providers.find(p => p.id === 'local');
      expect(anthropicAfter).toEqual(anthropicBefore);
      expect(localAfter).toEqual(localBefore);
    });

    it('should not throw when deleting non-existent provider', () => {
      const initialProviders = [...store.providers];
      
      store.deleteProvider('non-existent');
      
      expect(store.providers).toEqual(initialProviders);
    });
  });

  describe('getProviderById', () => {
    it('should return provider when found', () => {
      const openai = store.getProviderById('openai');
      expect(openai).toEqual({
        id: 'openai',
        name: 'OpenAI',
        baseUrl: 'https://api.openai.com/v1',
        apiKey: '',
        models: [
          { id: 'gpt-4', name: 'GPT-4', maxTokens: 8192 },
          { id: 'gpt-4-turbo', name: 'GPT-4 Turbo', maxTokens: 128000 },
          { id: 'gpt-3.5-turbo', name: 'GPT-3.5 Turbo', maxTokens: 4096 }
        ]
      });
    });

    it('should return undefined when not found', () => {
      const result = store.getProviderById('non-existent');
      expect(result).toBeUndefined();
    });
  });
});
```