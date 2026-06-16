import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('gm_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export const sceneApi = {
  getAll: () => api.get('/scenes'),
  getById: (id) => api.get(`/scenes/${id}`),
  create: (data) => api.post('/scenes', data),
  update: (id, data) => api.put(`/scenes/${id}`, data),
  delete: (id) => api.delete(`/scenes/${id}`)
};

export const npcApi = {
  getAll: () => api.get('/npcs'),
  getById: (id) => api.get(`/npcs/${id}`),
  create: (data) => api.post('/npcs', data),
  update: (id, data) => api.put(`/npcs/${id}`, data),
  delete: (id) => api.delete(`/npcs/${id}`)
};

export const agentApi = {
  getAll: () => api.get('/agents'),
  getById: (id) => api.get(`/agents/${id}`),
  create: (data) => api.post('/agents', data),
  update: (id, data) => api.put(`/agents/${id}`, data),
  delete: (id) => api.delete(`/agents/${id}`),
  chat: (id, message) => api.post(`/agents/${id}/chat`, { message })
};

export const llmApi = {
  getProviders: () => api.get('/llm/providers'),
  testConnection: (providerId) => api.post(`/llm/test`, { providerId })
};

export const configApi = {
  get: () => api.get('/config'),
  update: (data) => api.put('/config', data),
  export: () => api.get('/config/export'),
  import: (data) => api.post('/config/import', data)
};

export const generatorApi = {
  generate: (data) => api.post('/generator/generate', data),
  getStatus: () => api.get('/generator/status'),
  test: () => api.post('/generator/test')
};

export default api;
