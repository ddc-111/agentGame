import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useTaskStore = defineStore('task', () => {
  const tasks = ref([
    {
      id: 'task_001',
      name: '初来乍到',
      type: 'main',
      description: '新来的冒险者，先去杂货铺买些必需品吧',
      status: 'active',
      trigger: {
        type: 'auto',
        conditions: [
          { type: 'player_level', operator: '==', value: 1 }
        ]
      },
      objectives: [
        {
          id: 'obj_001',
          type: 'dialogue',
          target: 'npc_001',
          description: '与李掌柜对话',
          completed: false
        },
        {
          id: 'obj_002',
          type: 'collect',
          target: 'item_001',
          count: 3,
          description: '购买3份草药',
          completed: false
        },
        {
          id: 'obj_003',
          type: 'collect',
          target: 'item_003',
          count: 5,
          description: '购买5个馒头',
          completed: false
        }
      ],
      rewards: {
        exp: 100,
        gold: 500,
        items: []
      },
      nextTask: 'task_002',
      dialogue: 'dialogue_task_001'
    },
    {
      id: 'task_002',
      name: '装备自己',
      type: 'main',
      description: '有了补给，接下来去铁匠铺买把武器吧',
      status: 'locked',
      trigger: {
        type: 'task_complete',
        conditions: [
          { type: 'task_id', value: 'task_001' }
        ]
      },
      objectives: [
        {
          id: 'obj_004',
          type: 'dialogue',
          target: 'npc_003',
          description: '与张铁匠对话',
          completed: false
        },
        {
          id: 'obj_005',
          type: 'collect',
          target: 'item_010',
          count: 1,
          description: '购买一把铁剑',
          completed: false
        }
      ],
      rewards: {
        exp: 200,
        gold: 300,
        items: []
      },
      nextTask: 'task_003',
      dialogue: 'dialogue_task_002'
    }
  ]);

  const flows = ref([
    {
      id: 'flow_001',
      name: 'NPC出门购物流程',
      description: 'NPC从家里出发，前往商店购买物品的完整流程',
      nodes: [
        {
          id: 'node_1',
          type: 'start',
          position: { x: 100, y: 200 },
          data: { label: '开始' }
        },
        {
          id: 'node_2',
          type: 'action',
          position: { x: 250, y: 200 },
          data: {
            label: 'NPC离开家',
            action: 'move',
            params: {
              npcId: 'npc_001',
              from: 'home',
              to: 'scene_001',
              animation: 'walk'
            }
          }
        },
        {
          id: 'node_3',
          type: 'action',
          position: { x: 400, y: 200 },
          data: {
            label: '前往商店',
            action: 'move',
            params: {
              npcId: 'npc_001',
              from: 'scene_001',
              to: 'scene_002',
              path: [
                { x: 400, y: 300 },
                { x: 500, y: 350 },
                { x: 600, y: 400 }
              ]
            }
          }
        },
        {
          id: 'node_4',
          type: 'condition',
          position: { x: 550, y: 200 },
          data: {
            label: '商店开门?',
            condition: 'shop.isOpen(npc.shopId)'
          }
        },
        {
          id: 'node_5',
          type: 'action',
          position: { x: 700, y: 100 },
          data: {
            label: '等待商店开门',
            action: 'wait',
            params: {
              duration: 60,
              until: 'shop.open'
            }
          }
        },
        {
          id: 'node_6',
          type: 'action',
          position: { x: 700, y: 300 },
          data: {
            label: '进入商店',
            action: 'enter',
            params: {
              npcId: 'npc_001',
              shopId: 'shop_001'
            }
          }
        },
        {
          id: 'node_7',
          type: 'action',
          position: { x: 850, y: 200 },
          data: {
            label: '查看商品',
            action: 'browse',
            params: {
              npcId: 'npc_001',
              shopId: 'shop_001',
              duration: 30
            }
          }
        },
        {
          id: 'node_8',
          type: 'action',
          position: { x: 1000, y: 200 },
          data: {
            label: '购买物品',
            action: 'purchase',
            params: {
              npcId: 'npc_001',
              shopId: 'shop_001',
              items: [
                { itemId: 'item_001', count: 5 },
                { itemId: 'item_003', count: 10 }
              ]
            }
          }
        },
        {
          id: 'node_9',
          type: 'action',
          position: { x: 1150, y: 200 },
          data: {
            label: '离开商店',
            action: 'leave',
            params: {
              npcId: 'npc_001',
              shopId: 'shop_001'
            }
          }
        },
        {
          id: 'node_10',
          type: 'action',
          position: { x: 1300, y: 200 },
          data: {
            label: '返回家中',
            action: 'move',
            params: {
              npcId: 'npc_001',
              from: 'scene_002',
              to: 'home'
            }
          }
        },
        {
          id: 'node_11',
          type: 'end',
          position: { x: 1450, y: 200 },
          data: { label: '结束' }
        }
      ],
      edges: [
        { id: 'e1-2', source: 'node_1', target: 'node_2' },
        { id: 'e2-3', source: 'node_2', target: 'node_3' },
        { id: 'e3-4', source: 'node_3', target: 'node_4' },
        { id: 'e4-5', source: 'node_4', target: 'node_5', label: '否', sourceHandle: 'no' },
        { id: 'e4-6', source: 'node_4', target: 'node_6', label: '是', sourceHandle: 'yes' },
        { id: 'e5-6', source: 'node_5', target: 'node_6' },
        { id: 'e6-7', source: 'node_6', target: 'node_7' },
        { id: 'e7-8', source: 'node_7', target: 'node_8' },
        { id: 'e8-9', source: 'node_8', target: 'node_9' },
        { id: 'e9-10', source: 'node_9', target: 'node_10' },
        { id: 'e10-11', source: 'node_10', target: 'node_11' }
      ]
    }
  ]);

  const addTask = (task) => {
    tasks.value.push({
      id: `task_${Date.now()}`,
      ...task
    });
  };

  const updateTask = (id, data) => {
    const index = tasks.value.findIndex(t => t.id === id);
    if (index !== -1) {
      tasks.value[index] = { ...tasks.value[index], ...data };
    }
  };

  const deleteTask = (id) => {
    tasks.value = tasks.value.filter(t => t.id !== id);
  };

  const getTaskById = (id) => {
    return tasks.value.find(t => t.id === id);
  };

  const addFlow = (flow) => {
    flows.value.push({
      id: `flow_${Date.now()}`,
      ...flow
    });
  };

  const updateFlow = (id, data) => {
    const index = flows.value.findIndex(f => f.id === id);
    if (index !== -1) {
      flows.value[index] = { ...flows.value[index], ...data };
    }
  };

  const deleteFlow = (id) => {
    flows.value = flows.value.filter(f => f.id !== id);
  };

  return {
    tasks,
    flows,
    addTask,
    updateTask,
    deleteTask,
    getTaskById,
    addFlow,
    updateFlow,
    deleteFlow
  };
});
