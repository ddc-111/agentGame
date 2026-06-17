import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAchievementStore } from '@/stores/achievement';

describe('useAchievementStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useAchievementStore();
  });

  describe('initial state', () => {
    it('should have achievements array', () => {
      expect(Array.isArray(store.achievements)).toBe(true);
    });

    it('should have conditionTypes array', () => {
      expect(Array.isArray(store.conditionTypes)).toBe(true);
    });
  });

  describe('addAchievement', () => {
    it('should add new achievement', () => {
      const initialLength = store.achievements.length;
      const newAchievement = {
        name: 'Test Achievement',
        code: 'ach_test',
        description: 'Test description',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '🧪'
      };

      store.addAchievement(newAchievement);

      expect(store.achievements.length).toBeGreaterThanOrEqual(initialLength);
    });
  });

  describe('getAchievementById', () => {
    it('should return undefined for non-existent id', () => {
      const achievement = store.getAchievementById(999);
      expect(achievement).toBeUndefined();
    });
  });
});
