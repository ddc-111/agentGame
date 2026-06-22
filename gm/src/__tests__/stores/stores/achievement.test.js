import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAchievementStore } from './achievement';

// Mock Date.now for consistent testing
const mockDateNow = 1700000000000;
vi.spyOn(Date, 'now').mockImplementation(() => mockDateNow);

describe('useAchievementStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useAchievementStore();
  });

  describe('initial state', () => {
    it('should have initial achievements array with 10 items', () => {
      expect(store.achievements).toHaveLength(10);
    });

    it('should have correct structure for first achievement', () => {
      const first = store.achievements[0];
      expect(first).toHaveProperty('id', 1);
      expect(first).toHaveProperty('name', '初来乍到');
      expect(first).toHaveProperty('code', 'ach_first_quest');
      expect(first).toHaveProperty('description', '完成第一个任务');
      expect(first).toHaveProperty('condition', { type: 'quest_complete', value: 1 });
      expect(first).toHaveProperty('reward', { exp: 50, gold: 100 });
      expect(first).toHaveProperty('icon', '⭐');
    });

    it('should have all condition types', () => {
      expect(store.conditionTypes).toHaveLength(10);
      const conditionValues = store.conditionTypes.map(t => t.value);
      expect(conditionValues).toContain('quest_complete');
      expect(conditionValues).toContain('talk_all_npcs');
      expect(conditionValues).toContain('gold');
      expect(conditionValues).toContain('combat_win');
      expect(conditionValues).toContain('explore');
      expect(conditionValues).toContain('collect');
      expect(conditionValues).toContain('skill_use');
      expect(conditionValues).toContain('level');
      expect(conditionValues).toContain('item_use');
      expect(conditionValues).toContain('death');
    });
  });

  describe('addAchievement', () => {
    it('should add a new achievement', () => {
      const initialLength = store.achievements.length;
      const newAchievement = {
        name: '测试成就',
        code: 'ach_test',
        description: '这是一个测试成就',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '🧪'
      };

      store.addAchievement(newAchievement);

      expect(store.achievements).toHaveLength(initialLength + 1);
      
      const added = store.achievements.find(a => a.name === '测试成就');
      expect(added).toBeDefined();
      expect(added.id).toBe(mockDateNow);
      expect(added.code).toBe('ach_test');
      expect(added.description).toBe('这是一个测试成就');
      expect(added.condition).toEqual({ type: 'test', value: 1 });
      expect(added.reward).toEqual({ exp: 100, gold: 200 });
      expect(added.icon).toBe('🧪');
    });

    it('should generate unique id using Date.now', () => {
      const newAchievement = {
        name: '唯一ID测试',
        code: 'ach_unique',
        description: '测试ID唯一性',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '🔢'
      };

      store.addAchievement(newAchievement);
      
      const added = store.achievements.find(a => a.name === '唯一ID测试');
      expect(added.id).toBe(mockDateNow);
    });
  });

  describe('updateAchievement', () => {
    it('should update an existing achievement', () => {
      const achievementId = 1;
      const updates = {
        name: '更新后的名称',
        description: '更新后的描述'
      };

      store.updateAchievement(achievementId, updates);

      const updated = store.achievements.find(a => a.id === achievementId);
      expect(updated.name).toBe('更新后的名称');
      expect(updated.description).toBe('更新后的描述');
      // Original fields should remain
      expect(updated.code).toBe('ach_first_quest');
      expect(updated.condition).toEqual({ type: 'quest_complete', value: 1 });
    });

    it('should do nothing if achievement id does not exist', () => {
      const initialLength = store.achievements.length;
      const nonExistentId = 999;

      store.updateAchievement(nonExistentId, { name: '不存在' });

      expect(store.achievements).toHaveLength(initialLength);
      // No achievement should have the name '不存在'
      expect(store.achievements.find(a => a.name === '不存在')).toBeUndefined();
    });

    it('should partially update achievement', () => {
      const achievementId = 2;
      const updates = { icon: '👤' };

      store.updateAchievement(achievementId, updates);

      const updated = store.achievements.find(a => a.id === achievementId);
      expect(updated.icon).toBe('👤');
      // Other fields should remain unchanged
      expect(updated.name).toBe('村庄之友');
      expect(updated.code).toBe('ach_talk_all_npcs');
    });
  });

  describe('deleteAchievement', () => {
    it('should delete an existing achievement', () => {
      const initialLength = store.achievements.length;
      const achievementIdToDelete = 3;

      store.deleteAchievement(achievementIdToDelete);

      expect(store.achievements).toHaveLength(initialLength - 1);
      expect(store.achievements.find(a => a.id === achievementIdToDelete)).toBeUndefined();
    });

    it('should do nothing if achievement id does not exist', () => {
      const initialLength = store.achievements.length;
      const nonExistentId = 999;

      store.deleteAchievement(nonExistentId);

      expect(store.achievements).toHaveLength(initialLength);
    });

    it('should delete multiple achievements one by one', () => {
      const idsToDelete = [4, 5];
      
      idsToDelete.forEach(id => {
        store.deleteAchievement(id);
      });

      expect(store.achievements.find(a => a.id === 4)).toBeUndefined();
      expect(store.achievements.find(a => a.id === 5)).toBeUndefined();
      // Other achievements should still exist
      expect(store.achievements.find(a => a.id === 1)).toBeDefined();
      expect(store.achievements.find(a => a.id === 6)).toBeDefined();
    });
  });

  describe('getAchievementById', () => {
    it('should return an achievement by id', () => {
      const achievementId = 1;
      const achievement = store.getAchievementById(achievementId);
      
      expect(achievement).toBeDefined();
      expect(achievement.id).toBe(achievementId);
      expect(achievement.name).toBe('初来乍到');
    });

    it('should return undefined for non-existent id', () => {
      const nonExistentId = 999;
      const achievement = store.getAchievementById(nonExistentId);
      
      expect(achievement).toBeUndefined();
    });

    it('should return the correct achievement among multiple', () => {
      const achievementId = 5;
      const achievement = store.getAchievementById(achievementId);
      
      expect(achievement).toBeDefined();
      expect(achievement.id).toBe(achievementId);
      expect(achievement.name).toBe('探索者');
      expect(achievement.code).toBe('ach_explorer');
    });
  });

  describe('reactivity', () => {
    it('should be reactive when adding achievements', () => {
      const newAchievement = {
        name: '响应式测试',
        code: 'ach_reactive',
        description: '测试响应式',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '⚡'
      };

      store.addAchievement(newAchievement);
      
      expect(store.achievements.value).toBeDefined();
      expect(store.achievements.value.find(a => a.name === '响应式测试')).toBeDefined();
    });

    it('should be reactive when deleting achievements', () => {
      const initialLength = store.achievements.value.length;
      
      store.deleteAchievement(1);
      
      expect(store.achievements.value.length).toBe(initialLength - 1);
      expect(store.achievements.value.find(a => a.id === 1)).toBeUndefined();
    });

    it('should be reactive when updating achievements', () => {
      store.updateAchievement(1, { name: '响应式更新' });
      
      const updated = store.achievements.value.find(a => a.id === 1);
      expect(updated.name).toBe('响应式更新');
    });
  });
});