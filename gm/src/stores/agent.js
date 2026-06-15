import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAgentStore = defineStore('agent', () => {
  const agents = ref([
    {
      id: 'agent_001',
      name: '李掌柜智能体',
      description: '杂货铺老板的AI智能体，负责与玩家进行智能对话',
      llmConfig: {
        provider: 'openai',
        model: 'gpt-4',
        temperature: 0.7,
        maxTokens: 500
      },
      systemPrompt: `你是李掌柜，一位古风小镇杂货铺的老板。
你性格精明但诚实，对顾客热情周到。
你经营各种日常用品、药材、食材等。
当顾客询问商品时，你会详细介绍商品的用途和价格。
你偶尔会讲一些镇上的趣事。

当前店铺状态：
- 营业中
- 库存充足

与顾客对话时，保持古风语气，使用"客官"、"小店"等称呼。`,
      memoryConfig: {
        type: 'sliding_window',
        maxMessages: 20,
        summaryEnabled: true,
        summaryThreshold: 50
      },
      knowledgeBase: [
        { id: 'kb_001', name: '商品目录', type: 'text' },
        { id: 'kb_002', name: '价格表', type: 'table' }
      ],
      tools: [
        { name: 'query_inventory', description: '查询库存' },
        { name: 'make_deal', description: '进行交易' }
      ]
    },
    {
      id: 'agent_002',
      name: '王大娘智能体',
      description: '茶摊老板娘的AI智能体，喜欢聊天',
      llmConfig: {
        provider: 'openai',
        model: 'gpt-4',
        temperature: 0.8,
        maxTokens: 300
      },
      systemPrompt: `你是王大娘，一位古风小镇茶摊的老板娘。
你性格热情开朗，喜欢与人聊天。
你卖各种茶水，也提供一些小点心。
你知道镇上很多八卦和故事。
经常主动与路人搭话，分享镇上的新鲜事。

与客人对话时，使用亲切的语气，如"哎呀"、"客官"等。`,
      memoryConfig: {
        type: 'sliding_window',
        maxMessages: 15,
        summaryEnabled: false
      },
      knowledgeBase: [
        { id: 'kb_003', name: '茶品目录', type: 'text' },
        { id: 'kb_004', name: '镇上故事', type: 'text' }
      ],
      tools: []
    }
  ]);

  const currentAgent = ref(null);

  const addAgent = (agent) => {
    agents.value.push({
      id: `agent_${Date.now()}`,
      ...agent
    });
  };

  const updateAgent = (id, data) => {
    const index = agents.value.findIndex(a => a.id === id);
    if (index !== -1) {
      agents.value[index] = { ...agents.value[index], ...data };
    }
  };

  const deleteAgent = (id) => {
    agents.value = agents.value.filter(a => a.id !== id);
  };

  const getAgentById = (id) => {
    return agents.value.find(a => a.id === id);
  };

  return {
    agents,
    currentAgent,
    addAgent,
    updateAgent,
    deleteAgent,
    getAgentById
  };
});
