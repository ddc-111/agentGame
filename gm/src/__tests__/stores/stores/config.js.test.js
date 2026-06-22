import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('vue', () => ({
  ref: vi.fn(initialValue => ({ value: initialValue }))
}));

vi.mock('pinia', () => ({
  defineStore: vi.fn((id, setup) => setup)
}));

import { useConfigStore } from './config.js';

describe('useConfigStore', () => {
  let store;

  beforeEach(() => {
    store = useConfigStore();
  });

  describe('initialization', () => {
    it('should have default gameConfig', () => {
      expect(store.gameConfig.value).toBeDefined();
      expect(store.gameConfig.value.game.name).toBe('古风RPG');
      expect(store.gameConfig.value.game.version).toBe('1.0.0');
      expect(store.gameConfig.value.game.description).toBe('一个AI驱动的古风RPG游戏');
      expect(store.gameConfig.value.game.maxPlayers).toBe(100);
      expect(store.gameConfig.value.game.tickRate).toBe(20);
      expect(store.gameConfig.value.player.startScene).toBe('scene_001');
      expect(store.gameConfig.value.player.startPosition).toEqual({ x: 400, y: 300 });
      expect(store.gameConfig.value.player.startGold).toBe(1000);
      expect(store.gameConfig.value.player.startLevel).toBe(1);
      expect(store.gameConfig.value.player.maxLevel).toBe(100);
      expect(store.gameConfig.value.player.baseHP).toBe(100);
      expect(store.gameConfig.value.player.baseMP).toBe(50);
      expect(store.gameConfig.value.player.baseAttack).toBe(10);
      expect(store.gameConfig.value.player.baseDefense).toBe(5);
      expect(store.gameConfig.value.world.dayNightCycle).toBe(true);
      expect(store.gameConfig.value.world.dayDuration).toBe(1200);
      expect(store.gameConfig.value.world.weatherEnabled).toBe(true);
      expect(store.gameConfig.value.world.weatherTypes).toEqual(['sunny', 'cloudy', 'rainy', 'snowy']);
      expect(store.gameConfig.value.combat.enabled).toBe(true);
      expect(store.gameConfig.value.combat.turnBased).toBe(true);
      expect(store.gameConfig.value.combat.maxTurns).toBe(20);
      expect(store.gameConfig.value.combat.criticalRate).toBe(0.1);
      expect(store.gameConfig.value.combat.criticalMultiplier).toBe(1.5);
      expect(store.gameConfig.value.economy.inflationRate).toBe(0.01);
      expect(store.gameConfig.value.economy.taxRate).toBe(0.05);
      expect(store.gameConfig.value.economy.maxGold).toBe(999999);
    });
  });

  describe('updateConfig', () => {
    it('should update existing section', () => {
      const newData = { name: '新游戏名', maxPlayers: 200 };
      store.updateConfig('game', newData);
      expect(store.gameConfig.value.game.name).toBe('新游戏名');
      expect(store.gameConfig.value.game.maxPlayers).toBe(200);
      expect(store.gameConfig.value.game.version).toBe('1.0.0');
    });

    it('should add new properties to section', () => {
      const newData = { newProperty: 'newValue' };
      store.updateConfig('player', newData);
      expect(store.gameConfig.value.player.newProperty).toBe('newValue');
      expect(store.gameConfig.value.player.startScene).toBe('scene_001');
    });

    it('should not update non-existing section', () => {
      const newData = { name: 'test' };
      store.updateConfig('nonExisting', newData);
      expect(store.gameConfig.value.nonExisting).toBeUndefined();
    });

    it('should update nested properties', () => {
      const newData = { startPosition: { x: 500, y: 400 } };
      store.updateConfig('player', newData);
      expect(store.gameConfig.value.player.startPosition).toEqual({ x: 500, y: 400 });
    });
  });

  describe('getConfig', () => {
    it('should return existing section', () => {
      const gameSection = store.getConfig('game');
      expect(gameSection).toBeDefined();
      expect(gameSection.name).toBe('古风RPG');
      expect(gameSection.version).toBe('1.0.0');
    });

    it('should return player section', () => {
      const playerSection = store.getConfig('player');
      expect(playerSection).toBeDefined();
      expect(playerSection.startScene).toBe('scene_001');
      expect(playerSection.startGold).toBe(1000);
    });

    it('should return undefined for non-existing section', () => {
      const result = store.getConfig('nonExisting');
      expect(result).toBeUndefined();
    });
  });

  describe('exportConfig', () => {
    it('should return JSON string', () => {
      const json = store.exportConfig();
      expect(typeof json).toBe('string');
      const parsed = JSON.parse(json);
      expect(parsed).toEqual(store.gameConfig.value);
    });

    it('should format JSON with indentation', () => {
      const json = store.exportConfig();
      expect(json).toContain('\n');
      expect(json).toContain('  ');
    });

    it('should export current config after updates', () => {
      store.updateConfig('game', { name: '更新游戏