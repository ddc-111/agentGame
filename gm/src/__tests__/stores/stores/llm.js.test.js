import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useLLMStore } from './llm.js';

// Mock Date.now for consistent testing
const mockDateNow = 1234567890000;
vi.spyOn(Date, 'now').mockReturnValue(mockDateNow);

describe('LLM Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useLLMStore();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('Initial State', () => {
    it('should have initial providers', () => {
      expect(store.providers).toHaveLength(3);
      expect(store.providers[0].id).toBe('openai');
      expect(store.providers[1].id).toBe('anthropic');
      expect(store.providers[2].id).toBe('local');
    });

    it('should have correct models for each provider', () => {
      const openaiProvider = store.providers.find(p => p.id === 'openai');
      expect(openaiProvider.models).toHaveLength(3);
      expect(openaiProvider.models[0].id).toBe('gpt-4');

      const anthropicProvider = store.providers.find(p => p.id === 'anthropic');
      expect(anthropicProvider.models).toHaveLength(2);
      expect(anthropicProvider.models[0].id).toBe('claude-3-opus');

      const localProvider = store.providers.find(p => p.id === 'local');
      expect(localProvider.models).toHaveLength(2);
      expect(localProvider.models[0].id).toBe('qwen:7b');
    });

    it('should have currentProvider set to openai by default', () => {
      expect(store.currentProvider).toBe('openai');
    });
  });

  describe('addProvider', () => {
    it('should add a new provider with generated id', () => {
      const newProvider = {
        name: 'Test Provider',
        baseUrl: 'https://test.com',
        apiKey: 'test-key',
        models: [
          { id: 'test-model', name: 'Test Model', maxTokens: 1000 }
        ]
      };

      const initialLength = store.providers.length;
      store.addProvider(newProvider);

      expect(store.providers).toHaveLength(initialLength + 1);

      const addedProvider = store.providers[store.providers.length - 1];
      expect(addedProvider.id).toBe(`provider_${mockDateNow}`);
      expect(addedProvider.name).toBe('Test Provider');
      expect(addedProvider.baseUrl).toBe('https://test.com');
      expect(addedProvider.apiKey).toBe('test-key');
      expect(addedProvider.models).toHaveLength(1);
    });

    it('should handle provider with missing optional fields', () => {
      const minimalProvider = {
        name: 'Minimal Provider'
      };

      store.addProvider(minimalProvider);

      const addedProvider = store.providers[store.providers.length - 1];
      expect(addedProvider.id).toBe(`provider_${mockDateNow}`);
      expect(addedProvider.name).toBe('Minimal Provider');
      expect(addedProvider.baseUrl).toBeUndefined();
      expect(addedProvider.apiKey).toBeUndefined();
      expect(addedProvider.models).toBeUndefined();
    });
  });

  describe('updateProvider', () => {
    it('should update an existing provider', () => {
      const updateData = {
        name: 'Updated OpenAI',
        apiKey: 'new-api-key',
        baseUrl: 'https://new.openai.com/v1'
      };

      store.updateProvider('openai', updateData);

      const updatedProvider = store.providers.find(p => p.id === 'openai');
      expect(updatedProvider.name).toBe('Updated OpenAI');
      expect(updatedProvider.apiKey).toBe('new-api-key');
      expect(updatedProvider.baseUrl).toBe('https://new.openai.com/v1');
    });

    it('should only update specified fields', () => {
      const originalProvider = store.providers.find(p => p.id === 'openai');
      const originalModels = originalProvider.models;

      store.updateProvider('openai', { name: 'Updated Name' });

      const updatedProvider = store.providers.find(p => p.id === 'openai');
      expect(updatedProvider.name).toBe('Updated Name');
      expect(updatedProvider.baseUrl).toBe('https://api.openai.com/v1');
      expect(updatedProvider.apiKey).toBe('');
      expect(updatedProvider.models).toBe(originalModels);
    });

    it('should do nothing for non-existent provider', () => {
      const originalProviders = [...store.providers];
      
      store.updateProvider('non-existent-id', { name: 'Test' });
      
      expect(store.providers).toEqual(originalProviders);
    });
  });

  describe('deleteProvider', () => {
    it('should delete a provider by id', () => {
      const initialLength = store.providers.length;
      
      store.deleteProvider('anthropic');
      
      expect(store.providers).toHaveLength(initialLength - 1);
      expect(store.providers.find(p => p.id === 'anthropic')).toBeUndefined();
    });

    it('should not modify providers for non-existent id', () => {
      const initialProviders = [...store.providers];
      
      store.deleteProvider('non-existent-id');
      
      expect(store.providers).toEqual(initialProviders);
    });

    it('should allow deleting all providers', () => {
      store.deleteProvider('openai');
      store.deleteProvider('anthropic');
      store.deleteProvider('local');
      
      expect(store.providers).toHaveLength(0);
    });
  });

  describe('getProviderById', () => {
    it('should return the correct provider when it exists', () => {
      const provider = store.getProviderById('openai');
      
      expect(provider).toBeDefined();
      expect(provider.id).toBe('openai');
      expect(provider.name).toBe('OpenAI');
    });

    it('should return undefined for non-existent provider', () => {
      const provider = store.getProviderById('non-existent-id');
      
      expect(provider).toBeUndefined();
    });

    it('should return providers with correct structure', () => {
      const provider = store.getProviderById('anthropic');
      
      expect(provider).toHaveProperty('id');
      expect(provider).toHaveProperty('name');
      expect(provider).toHaveProperty('baseUrl');
      expect(provider).toHaveProperty('apiKey');
      expect(provider).toHaveProperty('models');
    });
  });

  describe('Integration tests', () => {
    it('should maintain provider list after multiple operations', () => {
      // Add a provider
      store.addProvider({
        name: 'New Provider',
        baseUrl: 'https://new.com',
        apiKey: 'key123',
        models: [{ id: 'model1', name: 'Model 1', maxTokens: 5000 }]
      });

      // Update a provider
      store.updateProvider('openai', { apiKey: 'updated-key' });

      // Delete a provider
      store.deleteProvider('local');

      // Verify the final state
      expect(store.providers).toHaveLength(3); // Added one, deleted one
      expect(store.providers.find(p => p.id === 'openai').apiKey).toBe('updated-key');
      expect(store.providers.find(p => p.id === 'local')).toBeUndefined();
      expect(store.providers.find(p => p.name === 'New Provider')).toBeDefined();
    });

    it('should handle concurrent operations correctly', () => {
      // Test multiple operations in sequence
      store.updateProvider('openai', { apiKey: 'first-update' });
      store.updateProvider('openai', { apiKey: 'second-update' });
      
      expect(store.getProviderById('openai').apiKey).toBe('second-update');
    });
  });
});