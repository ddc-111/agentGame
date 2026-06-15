import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useNPCStore = defineStore('npc', () => {
  const npcs = ref([
    {
      id: 'npc_001',
      name: '李掌柜',
      title: '杂货铺老板',
      description: '一位精明的中年商人，经营着镇上最大的杂货铺',
      avatar: 'merchant_li.png',
      sprite: 'merchant_li_sprite.png',
      position: { x: 400, y: 300, scene: 'scene_001' },
      agentId: 'agent_001',
      dialogues: ['dialogue_001'],
      shopId: 'shop_001',
      behaviors: ['idle', 'greet', 'sell'],
      schedule: [
        { time: '06:00', action: 'open_shop', position: { x: 400, y: 300, scene: 'scene_002' } },
        { time: '22:00', action: 'close_shop', position: { x: 200, y: 200, scene: 'scene_001' } }
      ]
    },
    {
      id: 'npc_002',
      name: '王大娘',
      title: '茶摊老板娘',
      description: '热情好客的中年妇人，卖的茶远近闻名',
      avatar: 'tea_wang.png',
      sprite: 'tea_wang_sprite.png',
      position: { x: 600, y: 400, scene: 'scene_001' },
      agentId: 'agent_002',
      dialogues: ['dialogue_002'],
      shopId: null,
      behaviors: ['idle', 'chat', 'serve_tea'],
      schedule: []
    },
    {
      id: 'npc_003',
      name: '张铁匠',
      title: '铁匠铺师傅',
      description: '技艺精湛的老铁匠，打造的兵器远近闻名',
      avatar: 'blacksmith_zhang.png',
      sprite: 'blacksmith_zhang_sprite.png',
      position: { x: 300, y: 350, scene: 'scene_002' },
      agentId: 'agent_003',
      dialogues: ['dialogue_003'],
      shopId: 'shop_002',
      behaviors: ['idle', 'forge', 'sell'],
      schedule: []
    }
  ]);

  const currentNPC = ref(null);

  const addNPC = (npc) => {
    npcs.value.push({
      id: `npc_${Date.now()}`,
      ...npc
    });
  };

  const updateNPC = (id, data) => {
    const index = npcs.value.findIndex(n => n.id === id);
    if (index !== -1) {
      npcs.value[index] = { ...npcs.value[index], ...data };
    }
  };

  const deleteNPC = (id) => {
    npcs.value = npcs.value.filter(n => n.id !== id);
  };

  const getNPCById = (id) => {
    return npcs.value.find(n => n.id === id);
  };

  return {
    npcs,
    currentNPC,
    addNPC,
    updateNPC,
    deleteNPC,
    getNPCById
  };
});
