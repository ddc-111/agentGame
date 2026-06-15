import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSceneStore = defineStore('scene', () => {
  const scenes = ref([
    {
      id: 'scene_001',
      name: '小镇中心',
      description: '古风小镇的中心广场，人来人往',
      background: 'town_center.png',
      width: 1920,
      height: 1080,
      npcs: ['npc_001', 'npc_002'],
      portals: [
        { x: 100, y: 500, targetScene: 'scene_002', targetX: 800, targetY: 400 }
      ],
      tiles: []
    },
    {
      id: 'scene_002',
      name: '杂货铺',
      description: '售卖各种日常用品的店铺',
      background: 'shop_general.png',
      width: 1280,
      height: 720,
      npcs: ['npc_003'],
      portals: [
        { x: 800, y: 600, targetScene: 'scene_001', targetX: 150, targetY: 500 }
      ],
      tiles: []
    }
  ]);

  const currentScene = ref(null);

  const addScene = (scene) => {
    scenes.value.push({
      id: `scene_${Date.now()}`,
      ...scene
    });
  };

  const updateScene = (id, data) => {
    const index = scenes.value.findIndex(s => s.id === id);
    if (index !== -1) {
      scenes.value[index] = { ...scenes.value[index], ...data };
    }
  };

  const deleteScene = (id) => {
    scenes.value = scenes.value.filter(s => s.id !== id);
  };

  const getSceneById = (id) => {
    return scenes.value.find(s => s.id === id);
  };

  return {
    scenes,
    currentScene,
    addScene,
    updateScene,
    deleteScene,
    getSceneById
  };
});
