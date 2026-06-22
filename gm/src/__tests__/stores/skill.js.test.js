import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useSkillStore } from './skill.js';

// Mock Date.now for consistent testing
vi.spyOn(Date, 'now').mockReturnValue(1234567890);

describe('Skill Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useSkillStore();
  });

  describe('Initial State', () => {
    it('should initialize with correct skills array', () => {
      expect(store.skills).toHaveLength(7);
      expect(store.skills[0].code).toBe('skill_basic_slash');
      expect(store.skills[2].code).toBe('skill_fireball');
    });

    it('should initialize with correct skill tree', () => {
      expect(store.skillTree.nodes).toHaveLength(8);
      expect(store.skillTree.edges).toHaveLength(7);
      expect(store.skillTree.nodes[0].id).toBe('root');
    });
  });

  describe('addSkill', () => {
    it('should add a new skill to the store', () => {
      const newSkill = {
        name: '新技能',
        code: 'skill_new',
        description: '测试技能',
        type: 'attack',
        mpCost: 10,
        damage: 20,
        heal: 0,
        cooldown: 1,
        level: 1,
        effect: {}
      };

      store.addSkill(newSkill);

      expect(store.skills).toHaveLength(8);
      expect(store.skills[7].id).toBe(1234567890);
      expect(store.skills[7].name).toBe('新技能');
      expect(store.skills[7].code).toBe('skill_new');
    });

    it('should preserve existing skills when adding new one', () => {
      const existingSkillsCount = store.skills.length;
      store.addSkill({ name: '测试', code: 'test' });

      expect(store.skills).toHaveLength(existingSkillsCount + 1);
      expect(store.skills[0].code).toBe('skill_basic_slash');
    });
  });

  describe('updateSkill', () => {
    it('should update an existing skill', () => {
      const updateData = { name: '升级斩击', damage: 25 };
      store.updateSkill(1, updateData);

      const updatedSkill = store.skills.find(s => s.id === 1);
      expect(updatedSkill.name).toBe('升级斩击');
      expect(updatedSkill.damage).toBe(25);
      expect(updatedSkill.code).toBe('skill_basic_slash'); // should preserve other properties
    });

    it('should do nothing if skill id does not exist', () => {
      const originalLength = store.skills.length;
      store.updateSkill(999, { name: '不存在' });

      expect(store.skills).toHaveLength(originalLength);
    });
  });

  describe('deleteSkill', () => {
    it('should remove skill by id', () => {
      store.deleteSkill(1);

      expect(store.skills).toHaveLength(6);
      expect(store.skills.find(s => s.id === 1)).toBeUndefined();
    });

    it('should preserve other skills when deleting', () => {
      const secondSkillId = store.skills[1].id;
      store.deleteSkill(1);

      expect(store.skills.find(s => s.id === secondSkillId)).toBeDefined();
      expect(store.skills[0].id).toBe(secondSkillId);
    });
  });

  describe('getSkillById', () => {
    it('should return skill by id', () => {
      const skill = store.getSkillById(1);
      expect(skill.name).toBe('基础斩击');
      expect(skill.code).toBe('skill_basic_slash');
    });

    it('should return undefined if skill id does not exist', () => {
      const skill = store.getSkillById(999);
      expect(skill).toBeUndefined();
    });
  });

  describe('getSkillByCode', () => {
    it('should return skill by code', () => {
      const skill = store.getSkillByCode('skill_fireball');
      expect(skill.name).toBe('火球术');
      expect(skill.id).toBe(3);
    });

    it('should return undefined if skill code does not exist', () => {
      const skill = store.getSkillByCode('skill_nonexistent');
      expect(skill).toBeUndefined();
    });
  });

  describe('updateSkillTree', () => {
    it('should update the entire skill tree', () => {
      const newTree = {
        nodes: [{ id: 'new_root', x: 100, y: 100, label: '新树', type: 'root' }],
        edges: []
      };

      store.updateSkillTree(newTree);

      expect(store.skillTree).toEqual(newTree);
      expect(store.skillTree.nodes).toHaveLength(1);
      expect(store.skillTree.nodes[0].id).toBe('new_root');
    });

    it('should completely replace the previous tree', () => {
      const originalNodeCount = store.skillTree.nodes.length;
      store.updateSkillTree({ nodes: [], edges: [] });

      expect(store.skillTree.nodes).toHaveLength(0);
      expect(store.skillTree.edges).toHaveLength(0);
    });
  });
});