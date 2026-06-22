import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useConfigStore } from '@/stores/config';

describe('useConfigStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useConfigStore();
  });

  describe('initial state', () => {
    it('should have gameConfig', () => {
      expect(store.gameConfig).toBeDefined();
    });

    it('should have default gameConfig values', () => {
      expect(store.gameConfig.game).toBeDefined();
      expect(store.gameConfig.game.maxPlayers).toBe(100);
      expect(store.gameConfig.game.name).toBe('古风RPG');
    });
  });

  describe('getConfig', () => {
    it('should return config section when exists', () => {
      const gameSection = store.getConfig('game');
      expect(gameSection).toBeDefined();
      expect(gameSection.name).toBe('古风RPG');
    });

    it('should return undefined for non-existent section', () => {
      const result = store.getConfig('nonExistent');
      expect(result).toBeUndefined();
    });

    it('should return different sections correctly', () => {
      const gameSection = store.getConfig('game');
      const playerSection = store.getConfig('player');
      expect(gameSection).toBeDefined();
      expect(playerSection).toBeDefined();
      expect(gameSection).not.toEqual(playerSection);
    });
  });

  describe('updateConfig', () => {
    it('should update config section with new data', () => {
      store.updateConfig('game', { name: '新游戏名' });
      expect(store.gameConfig.game.name).toBe('新游戏名');
    });

    it('should merge data with existing config', () => {
      store.updateConfig('game', { name: '新游戏名' });
      expect(store.gameConfig.game.maxPlayers).toBe(100);
    });

    it('should not affect other sections when updating', () => {
      const originalPlayer = { ...store.gameConfig.player };
      store.updateConfig('game', { name: '新游戏名' });
      expect(store.gameConfig.player).toEqual(originalPlayer);
    });
  });

  describe('exportConfig', () => {
    it('should export config as valid JSON string', () => {
      const exported = store.exportConfig();
      expect(typeof exported).toBe('string');
      expect(() => JSON.parse(exported)).not.toThrow();
    });

    it('should export current state of config', () => {
      const exported = store.exportConfig();
      const parsed = JSON.parse(exported);
      expect(parsed.game.name).toBe('古风RPG');
    });
  });

  describe('importConfig', () => {
    it('should import valid JSON and return true', () => {
      const result = store.importConfig('{"game": {"name": "导入的游戏"}}');
      expect(result).toBe(true);
    });

    it('should update store after successful import', () => {
      store.importConfig('{"game": {"name": "导入的游戏"}}');
      expect(store.gameConfig.game.name).toBe('导入的游戏');
    });

    it('should return false for invalid JSON', () => {
      const result = store.importConfig('invalid json');
      expect(result).toBe(false);
    });

    it('should return false for empty string', () => {
      const result = store.importConfig('');
      expect(result).toBe(false);
    });
  });
});