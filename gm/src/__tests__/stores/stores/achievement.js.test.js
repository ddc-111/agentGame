import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAchievementStore } from './stores/achievement.js';

// Mock Date.now to make tests deterministic
const mockDateNow = 1234567890;
vi.spyOn(Date, 'now').mockImplementation(() => mockDateNow);

describe('Achievement Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useAchievementStore();
    Date.now.mockClear();
  });

  describe('Initial State', () => {
    it('should have 10 achievements', () => {
      expect(store.achievements).toHaveLength(10);
    });

    it('should have correct conditionTypes', () => {
      expect(store.conditionTypes).toHaveLength(10);
      expect(store.conditionTypes).toEqual(
        expect.arrayContaining([
          expect.objectContaining({ value: 'quest_complete' }),
          expect.objectContaining({ value: 'talk_all_npcs' }),
          expect.objectContaining({ value: 'gold' }),
          expect.objectContaining({ value: 'combat_win' }),
          expect.objectContaining({ value: 'explore' }),
          expect.objectContaining({ value: 'collect' }),
          expect.objectContaining({ value: 'skill_use' }),
          expect.objectContaining({ value: 'level' }),
          expect.objectContaining({ value: 'item_use' }),
          expect.objectContaining({ value: 'death' })
        ])
      );
    });

    it('should have correct achievement structure', () => {
      const firstAchievement = store.achievements[0];
      expect(firstAchievement).toEqual({
        id: 1,
        name: '初来乍到',
        code: 'ach_first_quest',
        description: '完成第一个任务',
        condition: { type: 'quest_complete', value: 1 },
        reward: { exp: 50, gold: 100 },
        icon: '⭐'
      });
    });
  });

  describe('addAchievement', () => {
    it('should add a new achievement', () => {
      const newAchievement = {
        name: '测试成就',
        code: 'ach_test',
        description: '这是一个测试成就',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '🎮'
      };

      const initialLength = store.achievements.length;
      store.addAchievement(newAchievement);

      expect(store.achievements).toHaveLength(initialLength + 1);
      expect(Date.now).toHaveBeenCalled();
      
      const addedAchievement = store.achievements[store.achievements.length - 1];
      expect(addedAchievement).toEqual({
        id: mockDateNow,
        ...newAchievement
      });
    });

    it('should add multiple achievements', () => {
      const achievement1 = { name: '成就1', code: 'ach_1' };
      const achievement2 = { name: '成就2', code: 'ach_2' };

      store.addAchievement(achievement1);
      store.addAchievement(achievement2);

      expect(store.achievements).toHaveLength(12);
      expect(store.achievements[10].name).toBe('成就1');
      expect(store.achievements[11].name).toBe('成就2');
    });
  });

  describe('updateAchievement', () => {
    it('should update an existing achievement', () => {
      const updatedData = {
        name: '更新后的成就',
        description: '这是一个更新后的描述'
      };

      store.updateAchievement(1, updatedData);

      const achievement = store.achievements.find(a => a.id === 1);
      expect(achievement.name).toBe('更新后的成就');
      expect(achievement.description).toBe('这是一个更新后的描述');
      expect(achievement.code).toBe('ach_first_quest'); // 其他字段应该保持不变
    });

    it('should not update if id does not exist', () => {
      const initialAchievements = [...store.achievements];
      store.updateAchievement(999, { name: '不存在的成就' });

      expect(store.achievements).toEqual(initialAchievements);
    });

    it('should handle partial updates', () => {
      store.updateAchievement(1, { icon: '⭐⭐' });

      const achievement = store.achievements.find(a => a.id === 1);
      expect(achievement.icon).toBe('⭐⭐');
      expect(achievement.name).toBe('初来乍到'); // name should remain unchanged
    });
  });

  describe('deleteAchievement', () => {
    it('should delete an achievement', () => {
      const initialLength = store.achievements.length;
      store.deleteAchievement(1);

      expect(store.achievements).toHaveLength(initialLength - 1);
      expect(store.achievements.find(a => a.id === 1)).toBeUndefined();
    });

    it('should not delete if id does not exist', () => {
      const initialLength = store.achievements.length;
      store.deleteAchievement(999);

      expect(store.achievements).toHaveLength(initialLength);
    });

    it('should delete correct achievement', () => {
      store.deleteAchievement(3);
      const achievement = store.achievements.find(a => a.code === 'ach_rich');
      expect(achievement).toBeUndefined();
    });
  });

  describe('getAchievementById', () => {
    it('should return achievement by id', () => {
      const achievement = store.getAchievementById(1);
      expect(achievement).toEqual({
        id: 1,
        name: '初来乍到',
        code: 'ach_first_quest',
        description: '完成第一个任务',
        condition: { type: 'quest_complete', value: 1 },
        reward: { exp: 50, gold: 100 },
        icon: '⭐'
      });
    });

    it('should return undefined for non-existent id', () => {
      const achievement = store.getAchievementById(999);
      expect(achievement).toBeUndefined();
    });

    it('should return correct achievement when multiple exist', () => {
      const achievement = store.getAchievementById(5);
      expect(achievement.name).toBe('探索者');
      expect(achievement.code).toBe('ach_explorer');
    });
  });

  describe('Store Methods Integration', () => {
    it('should handle full CRUD operations', () => {
      // Add
      const newAchievement = {
        name: '集成测试成就',
        code: 'ach_integration',
        description: '测试集成',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '🧪'
      };
      store.addAchievement(newAchievement);
      const addedAchievement = store.achievements[store.achievements.length - 1];

      // Read
      const fetchedAchievement = store.getAchievementById(addedAchievement.id);
      expect(fetchedAchievement).toEqual(addedAchievement);

      // Update
      store.updateAchievement(addedAchievement.id, { name: '更新后的集成测试' });
      const updatedAchievement = store.getAchievementById(addedAchievement.id);
      expect(updatedAchievement.name).toBe('更新后的集成测试');

      // Delete
      store.deleteAchievement(addedAchievement.id);
      expect(store.getAchievementById(addedAchievement.id)).toBeUndefined();
    });
  });
});