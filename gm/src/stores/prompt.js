import { defineStore } from 'pinia';
import { ref } from 'vue';

export const usePromptStore = defineStore('prompt', () => {
  const templates = ref([
    {
      id: 'template_001',
      name: 'NPC基础人设',
      category: 'system',
      content: `你是{{npc_name}}，{{npc_title}}。
{{npc_description}}

当前状态：
- 位置：{{current_scene}}
- 时间：{{current_time}}
- 心情：{{mood}}

{{#if has_shop}}
你经营着一家{{shop_type}}店。
{{/if}}

{{#if has_task}}
你有任务要交给冒险者：{{task_description}}
{{/if}}

请保持角色设定，用古风语气与玩家对话。`,
      variables: [
        { name: 'npc_name', type: 'string', description: 'NPC名称' },
        { name: 'npc_title', type: 'string', description: 'NPC称号' },
        { name: 'npc_description', type: 'string', description: 'NPC描述' },
        { name: 'current_scene', type: 'string', description: '当前场景' },
        { name: 'current_time', type: 'string', description: '当前时间' },
        { name: 'mood', type: 'string', description: 'NPC心情' },
        { name: 'has_shop', type: 'boolean', description: '是否有商店' },
        { name: 'shop_type', type: 'string', description: '商店类型' },
        { name: 'has_task', type: 'boolean', description: '是否有任务' },
        { name: 'task_description', type: 'string', description: '任务描述' }
      ]
    },
    {
      id: 'template_002',
      name: '商人对话模板',
      category: 'system',
      content: `你是{{npc_name}}，一位{{personality}}的商人。

你售卖以下商品：
{{#each items}}
- {{this.name}}: {{this.description}} ({{this.price}}文)
{{/each}}

交易规则：
1. 可以适当还价，但不能低于成本价
2. 购买超过3件可打9折
3. 老顾客可以赊账

当前库存状态：
{{#each items}}
- {{this.name}}: {{this.stock}}件
{{/each}}`,
      variables: [
        { name: 'npc_name', type: 'string', description: 'NPC名称' },
        { name: 'personality', type: 'string', description: '性格特点' },
        { name: 'items', type: 'array', description: '商品列表' }
      ]
    },
    {
      id: 'template_003',
      name: '对话摘要',
      category: 'summary',
      content: `请将以下对话总结为简短的摘要，保留关键信息：

对话内容：
{{conversation}}

请总结：
1. 对话主题
2. 达成的共识或交易
3. 需要记住的重要信息`,
      variables: [
        { name: 'conversation', type: 'string', description: '对话内容' }
      ]
    }
  ]);

  const variables = ref([
    { name: 'npc_name', type: 'string', description: 'NPC名称', source: 'npc' },
    { name: 'npc_title', type: 'string', description: 'NPC称号', source: 'npc' },
    { name: 'npc_description', type: 'string', description: 'NPC描述', source: 'npc' },
    { name: 'current_scene', type: 'string', description: '当前场景', source: 'game' },
    { name: 'current_time', type: 'string', description: '当前时间', source: 'game' },
    { name: 'player_name', type: 'string', description: '玩家名称', source: 'player' },
    { name: 'player_level', type: 'number', description: '玩家等级', source: 'player' },
    { name: 'mood', type: 'string', description: 'NPC心情', source: 'agent' },
    { name: 'has_shop', type: 'boolean', description: '是否有商店', source: 'npc' },
    { name: 'shop_type', type: 'string', description: '商店类型', source: 'shop' }
  ]);

  const addTemplate = (template) => {
    templates.value.push({
      id: `template_${Date.now()}`,
      ...template
    });
  };

  const updateTemplate = (id, data) => {
    const index = templates.value.findIndex(t => t.id === id);
    if (index !== -1) {
      templates.value[index] = { ...templates.value[index], ...data };
    }
  };

  const deleteTemplate = (id) => {
    templates.value = templates.value.filter(t => t.id !== id);
  };

  const addVariable = (variable) => {
    variables.value.push(variable);
  };

  const updateVariable = (name, data) => {
    const index = variables.value.findIndex(v => v.name === name);
    if (index !== -1) {
      variables.value[index] = { ...variables.value[index], ...data };
    }
  };

  const deleteVariable = (name) => {
    variables.value = variables.value.filter(v => v.name !== name);
  };

  return {
    templates,
    variables,
    addTemplate,
    updateTemplate,
    deleteTemplate,
    addVariable,
    updateVariable,
    deleteVariable
  };
});
