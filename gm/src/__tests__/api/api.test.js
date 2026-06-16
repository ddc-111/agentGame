import { describe, it, expect, beforeEach, vi } from 'vitest';
import api, { sceneApi, npcApi, agentApi, llmApi, configApi, generatorApi } from '@/api/index';

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should create axios instance with correct config', () => {
    expect(api.defaults.baseURL).toBe('/api');
    expect(api.defaults.timeout).toBe(10000);
    expect(api.defaults.headers['Content-Type']).toBe('application/json');
  });

  it('should add auth token to request if exists', async () => {
    localStorage.getItem.mockReturnValue('test-token');
    const config = { headers: {} };
    const interceptor = api.interceptors.request.handlers[0].fulfilled;
    const result = interceptor(config);
    expect(result.headers.Authorization).toBe('Bearer test-token');
  });

  it('should not add auth token if not exists', async () => {
    localStorage.getItem.mockReturnValue(null);
    const config = { headers: {} };
    const interceptor = api.interceptors.request.handlers[0].fulfilled;
    const result = interceptor(config);
    expect(result.headers.Authorization).toBeUndefined();
  });

  it('should handle request error', async () => {
    const error = new Error('Request failed');
    const interceptor = api.interceptors.request.handlers[0].rejected;
    await expect(interceptor(error)).rejects.toThrow('Request failed');
  });

  it('should return response data on success', () => {
    const response = { data: { id: 1, name: 'test' } };
    const interceptor = api.interceptors.response.handlers[0].fulfilled;
    const result = interceptor(response);
    expect(result).toEqual({ id: 1, name: 'test' });
  });

  it('should handle response error', async () => {
    const error = new Error('Response failed');
    const interceptor = api.interceptors.response.handlers[0].rejected;
    await expect(interceptor(error)).rejects.toThrow('Response failed');
  });

  describe('Scene API', () => {
    it('should have all methods', () => {
      expect(sceneApi.getAll).toBeDefined();
      expect(sceneApi.getById).toBeDefined();
      expect(sceneApi.create).toBeDefined();
      expect(sceneApi.update).toBeDefined();
      expect(sceneApi.delete).toBeDefined();
    });
  });

  describe('NPC API', () => {
    it('should have all methods', () => {
      expect(npcApi.getAll).toBeDefined();
      expect(npcApi.getById).toBeDefined();
      expect(npcApi.create).toBeDefined();
      expect(npcApi.update).toBeDefined();
      expect(npcApi.delete).toBeDefined();
    });
  });

  describe('Agent API', () => {
    it('should have all methods', () => {
      expect(agentApi.getAll).toBeDefined();
      expect(agentApi.getById).toBeDefined();
      expect(agentApi.create).toBeDefined();
      expect(agentApi.update).toBeDefined();
      expect(agentApi.delete).toBeDefined();
      expect(agentApi.chat).toBeDefined();
    });
  });

  describe('LLM API', () => {
    it('should have all methods', () => {
      expect(llmApi.getProviders).toBeDefined();
      expect(llmApi.testConnection).toBeDefined();
    });
  });

  describe('Config API', () => {
    it('should have all methods', () => {
      expect(configApi.get).toBeDefined();
      expect(configApi.update).toBeDefined();
      expect(configApi.export).toBeDefined();
      expect(configApi.import).toBeDefined();
    });
  });

  describe('Generator API', () => {
    it('should have all methods', () => {
      expect(generatorApi.generate).toBeDefined();
      expect(generatorApi.getStatus).toBeDefined();
      expect(generatorApi.test).toBeDefined();
    });
  });
});
