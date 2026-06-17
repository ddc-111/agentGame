import { describe, it, expect, beforeEach } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
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
  });

  describe('getConfig', () => {
    it('should return config section', () => {
      const gameSection = store.getConfig('game');
      expect(gameSection).toBeDefined();
    });

    it('should return undefined for non-existent section', () => {
      const result = store.getConfig('nonExistent');
      expect(result).toBeUndefined();
    });
  });

  describe('updateConfig', () => {
    it('should update config section', () => {
      const newData = { maxPlayers: 200 };
      store.updateConfig('game', newData);
      
      const game = store.getConfig('game');
      expect(game.maxPlayers).toBe(200);
    });
  });

  describe('exportConfig', () => {
    it('should export config as JSON string', () => {
      const exported = store.exportConfig();
      expect(typeof exported).toBe('string');
      
      const parsed = JSON.parse(exported);
      expect(parsed).toBeDefined();
    });
  });

  describe('importConfig', () => {
    it('should import valid JSON', () => {
      const importData = {
        game: { name: 'Test Game' }
      };
      const jsonString = JSON.stringify(importData);
      
      const result = store.importConfig(jsonString);
      expect(result).toBe(true);
    });

    it('should return false for invalid JSON', () => {
      const invalidJson = '{invalid json}';
      
      const result = store.importConfig(invalidJson);
      expect(result).toBe(false);
    });
  });
});
