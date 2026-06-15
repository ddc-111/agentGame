import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useLLMStore = defineStore('llm', () => {
  const providers = ref([
    {
      id: 'openai',
      name: 'OpenAI',
      baseUrl: 'https://api.openai.com/v1',
      apiKey: '',
      models: [
        { id: 'gpt-4', name: 'GPT-4', maxTokens: 8192 },
        { id: 'gpt-4-turbo', name: 'GPT-4 Turbo', maxTokens: 128000 },
        { id: 'gpt-3.5-turbo', name: 'GPT-3.5 Turbo', maxTokens: 4096 }
      ]
    },
    {
      id: 'anthropic',
      name: 'Anthropic',
      baseUrl: 'https://api.anthropic.com',
      apiKey: '',
      models: [
        { id: 'claude-3-opus', name: 'Claude 3 Opus', maxTokens: 200000 },
        { id: 'claude-3-sonnet', name: 'Claude 3 Sonnet', maxTokens: 200000 }
      ]
    },
    {
      id: 'local',
      name: '本地模型',
      baseUrl: 'http://localhost:11434',
      apiKey: '',
      models: [
        { id: 'qwen:7b', name: 'Qwen 7B', maxTokens: 4096 },
        { id: 'llama2:7b', name: 'Llama 2 7B', maxTokens: 4096 }
      ]
    }
  ]);

  const currentProvider = ref('openai');

  const addProvider = (provider) => {
    providers.value.push({
      id: `provider_${Date.now()}`,
      ...provider
    });
  };

  const updateProvider = (id, data) => {
    const index = providers.value.findIndex(p => p.id === id);
    if (index !== -1) {
      providers.value[index] = { ...providers.value[index], ...data };
    }
  };

  const deleteProvider = (id) => {
    providers.value = providers.value.filter(p => p.id !== id);
  };

  const getProviderById = (id) => {
    return providers.value.find(p => p.id === id);
  };

  return {
    providers,
    currentProvider,
    addProvider,
    updateProvider,
    deleteProvider,
    getProviderById
  };
});
