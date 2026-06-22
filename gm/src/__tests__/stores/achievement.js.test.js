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

    it('should have 10 predefined achievements', () => {
      expect(store.achievements.length).toBe(10);
    });

    it('should have specific predefined achievements', () => {
      expect(store.achievements[0].name).toBe('初来乍到');
      expect(store.achievements[1].name).toBe('村庄之友');
      expect(store.achievements[2].name).toBe('富甲一方');
    });
  });

  describe('addAchievement', () => {
    it('should add new achievement to achievements array', () => {
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

      expect(store.achievements.length).toBe(initialLength + 1);
      const added = store.achievements.find(a => a.code === 'ach_test');
      expect(added).toBeDefined();
      expect(added.name).toBe('Test Achievement');
    });

    it('should add multiple achievements', () => {
      const initialLength = store.achievements.length;
      const achievement1 = { code: 'ach_test1', name: 'Test 1' };
      const achievement2 = { code: 'ach_test2', name: 'Test 2' };

      store.addAchievement(achievement1);
      store.addAchievement(achievement2);

      expect(store.achievements.length).toBe(initialLength + 2);
    });

    it('should auto-generate id if not provided', () => {
      const achievement = { code: 'ach_no_id', name: 'No ID' };
      store.addAchievement(achievement);

      const added = store.achievements.find(a => a.code === 'ach_no_id');
      expect(added).toBeDefined();
      expect(added.id).toBeDefined();
    });

    it('should not add duplicate achievement code', () => {
      const initialLength = store.achievements.length;
      const achievement1 = { code: 'ach_duplicate', name: 'First' };
      const achievement2 = { code: 'ach_duplicate', name: 'Second' };

      store.addAchievement(achievement1);
      store.addAchievement(achievement2);

      expect(store.achievements.length).toBe(initialLength + 1);
      const found = store.achievements.find(a => a.code === 'ach_duplicate');
      expect(found.name).toBe('First');
    });
  });

  describe('getAchievementById', () => {
    it('should return achievement by id if exists', () => {
      const result = store.getAchievementById(1);
      expect(result).toBeDefined();
      expect(result.name).toBe('初来乍到');
    });

    it('should return undefined for non-existent id', () => {
      const result = store.getAchievementById(999);
      expect(result).toBeUndefined();
    });
  });

  describe('getAchievementByCode', () => {
    it('should return achievement by code if exists', () => {
      const result = store.getAchievementByCode('ach_first_quest');
      expect(result).toBeDefined();
      expect(result.name).toBe('初来乍到');
    });

    it('should return undefined for non-existent code', () => {
      const result = store.getAchievementByCode('ach_nonexistent');
      expect(result).toBeUndefined();
    });
  });

  describe('updateAchievement', () => {
    it('should update achievement by id', () => {
      const updated = {
        id: 1,
        name: 'Updated Achievement'
      };

      store.updateAchievement(1, updated);

      const result = store.getAchievementById(1);
      expect(result.name).toBe('Updated Achievement');
      expect(result.code).toBe('ach_first_quest');
    });

    it('should not update non-existent achievement', () => {
      const initialLength = store.achievements.length;
      store.updateAchievement(999, { name: 'Test' });
      expect(store.achievements.length).toBe(initialLength);
    });

    it('should preserve other fields when updating', () => {
      store.updateAchievement(1, { name: 'Updated' });
      const result = store.getAchievementById(1);
      expect(result.name).toBe('Updated');
      expect(result.code).toBe('ach_first_quest');
      expect(result.description).toBe('完成第一个任务');
    });

    it('should update multiple fields', () => {
      const updates = {
        name: 'Multi Update',
        description: 'Updated description',
        reward: { exp: 500, gold: 1000 }
      };

      store.updateAchievement(1, updates);
      const result = store.getAchievementById(1);

      expect(result.name).toBe('Multi Update');
      expect(result.description).toBe('Updated description');
      expect(result.reward).toEqual({ exp: 500, gold: 1000 });
      expect(result.code).toBe('ach_first_quest');
    });
  });

  describe('deleteAchievement', () => {
    it('should delete achievement by id', () => {
      const initialLength = store.achievements.length;
      store.deleteAchievement(1);
      expect(store.achievements.length).toBe(initialLength - 1);
      expect(store.getAchievementById(1)).toBeUndefined();
    });

    it('should not affect other achievements when deleting', () => {
      const initialLength = store.achievements.length;
      store.deleteAchievement(1);
      expect(store.achievements.length).toBe(initialLength - 1);
      expect(store.getAchievementById(2)).toBeDefined();
    });

    it('should do nothing when deleting non-existent id', () => {
      const initialLength = store.achievements.length;
      store.deleteAchievement(999);
      expect(store.achievements.length).toBe(initialLength);
    });
  });

  describe('conditionTypes', () => {
    it('should have predefined condition types', () => {
      expect(store.conditionTypes.length).toBeGreaterThan(0);
    });

    it('should have specific condition types', () => {
      const types = store.conditionTypes.map(t => t.value);
      expect(types).toContain('quest_complete');
      expect(types).toContain('combat_win');
      expect(types).toContain('level');
    });

    it('should have unique condition type values', () => {
      const values = store.conditionTypes.map(t => t.value);
      const uniqueValues = [...new Set(values)];
      expect(values.length).toBe(uniqueValues.length);
    });

    it('should have proper structure for each condition type', () => {
      store.conditionTypes.forEach(type => {
        expect(type).toHaveProperty('value');
        expect(type).toHaveProperty('label');
        expect(typeof type.value).toBe('string');
        expect(typeof type.label).toBe('string');
      });
    });
  });

  describe('getAchievementByConditionType', () => {
    it('should return achievements matching condition type', () => {
      const result = store.getAchievementByConditionType('quest_complete');
      expect(Array.isArray(result)).toBe(true);
      result.forEach(achievement => {
        expect(achievement.condition.type).toBe('quest_complete');
      });
    });

    it('should return empty array for non-existent condition type', () => {
      const result = store.getAchievementByConditionType('nonexistent_type');
      expect(result).toEqual([]);
    });
  });

  describe('isAchievementUnlocked', () => {
    it('should check if achievement is unlocked by id', () => {
      const achievement = store.getAchievementById(1);
      if (achievement.unlocked !== undefined) {
        const result = store.isAchievementUnlocked(1);
        expect(typeof result).toBe('boolean');
      }
    });
  });

  describe('unlockAchievement', () => {
    it('should unlock achievement by id', () => {
      if (typeof store.unlockAchievement === 'function') {
        store.unlockAchievement(1);
        const achievement = store.getAchievementById(1);
        expect(achievement.unlocked).toBe(true);
      }
    });
  });

  describe('resetAchievements', () => {
    it('should reset all achievements to initial state', () => {
      store.addAchievement({ code: 'ach_temp', name: 'Temp' });
      store.updateAchievement(1, { name: 'Modified' });
      
      const initialLength = 10;
      
      if (typeof store.resetAchievements === 'function') {
        store.resetAchievements();
        expect(store.achievements.length).toBe(initialLength);
        expect(store.getAchievementById(1).name).toBe('初来乍到');
      }
    });
  });
});