import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useConfigStore = defineStore('config', () => {
  const gameConfig = ref({
    game: {
      name: '古风RPG',
      version: '1.0.0',
      description: '一个AI驱动的古风RPG游戏',
      maxPlayers: 100,
      tickRate: 20
    },
    player: {
      startScene: 'scene_001',
      startPosition: { x: 400, y: 300 },
      startGold: 1000,
      startLevel: 1,
      maxLevel: 100,
      baseHP: 100,
      baseMP: 50,
      baseAttack: 10,
      baseDefense: 5
    },
    world: {
      dayNightCycle: true,
      dayDuration: 1200,
      weatherEnabled: true,
      weatherTypes: ['sunny', 'cloudy', 'rainy', 'snowy']
    },
    combat: {
      enabled: true,
      turnBased: true,
      maxTurns: 20,
      criticalRate: 0.1,
      criticalMultiplier: 1.5
    },
    economy: {
      inflationRate: 0.01,
      taxRate: 0.05,
      maxGold: 999999
    }
  });

  const updateConfig = (section, data) => {
    if (gameConfig.value[section]) {
      gameConfig.value[section] = { ...gameConfig.value[section], ...data };
    }
  };

  const getConfig = (section) => {
    return gameConfig.value[section];
  };

  const exportConfig = () => {
    return JSON.stringify(gameConfig.value, null, 2);
  };

  const importConfig = (json) => {
    try {
      const data = JSON.parse(json);
      gameConfig.value = { ...gameConfig.value, ...data };
      return true;
    } catch (e) {
      return false;
    }
  };

  return {
    gameConfig,
    updateConfig,
    getConfig,
    exportConfig,
    importConfig
  };
});
