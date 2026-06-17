```javascript
import { describe, it, expect, beforeEach } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import { useSkillStore } from './skill.js';

describe('useSkillStore', () => {
  let store;

  beforeEach(() => {
    const pinia = createPinia();
    setActivePinia(pinia);
    store = useSkillStore();
  });

  it('should have initial skills array', () => {
    expect(store.skills.value).toHaveLength(7);
    expect(store.skills.value[0]).toEqual({
      id: 1,
      name: '基础斩击',
      code: 'skill_basic_slash',
      description: '用武器进行强力斩击，造成1.5倍武器伤害',
      type: 'attack',
      mpCost: 5,
      damage: 15,
      heal: 0,
      cooldown: 0,
      level: 1,
      effect: {}
    });
  });

  it('should have initial skillTree object', () => {
    expect(store.skillTree.value.nodes).toHaveLength(8);
    expect(store.skillTree.value.edges).toHaveLength(7);
    expect(store.skillTree.value.nodes[0]).toEqual({
      id: 'root',
      x: 400,
      y: 50,
      label: '技能树',
      type: 'root'
    });
  });

  it('addSkill should add a new skill', () => {
    const initialLength = store.skills.value.length;
    const newSkill = {
      name: '新技能',
      code: 'new_skill',
      description: '一个新技能',
      type: 'attack',
      mpCost: 10,
      damage: 20,
      heal: 0,
      cooldown: 1,
      level: 1,
      effect: {}
    };

    store.addSkill(newSkill);

    expect(store.skills.value).toHaveLength(initialLength + 1);
    const addedSkill = store.skills.value[store.skills.value.length - 1];
    expect(addedSkill.name).toBe('新技能');
    expect(addedSkill.code).toBe('new_skill');
    expect(addedSkill.id).toBeDefined();
  });

  it('updateSkill should update an existing skill', () => {
    const id = 1;
    const newData = { name: '更新后的斩击', description: '描述已更新' };

    store.updateSkill(id, newData);

    const updatedSkill = store.getSkillById(id);
    expect(updatedSkill.name).toBe('更新后的斩击');
    expect(updatedSkill.description).toBe('描述已更新');
    expect(updatedSkill.id).toBe(id);
  });

  it('updateSkill should do nothing if id does not exist', () => {
    const nonExistentId = 999;
    const initialSkills = [...store.skills.value];

    store.updateSkill(nonExistentId, { name: '不存在' });

    expect(store.skills.value).toEqual(initialSkills);
  });

  it('deleteSkill should remove a skill by id', () => {
    const id = 2;
    const initialLength = store.skills.value.length;

    store.deleteSkill(id);

    expect(store.skills.value).toHaveLength(initialLength - 1);
    expect(store.getSkillById(id)).toBeUndefined();
  });

  it('deleteSkill should do nothing if id does not exist', () => {
    const nonExistentId = 999;
    const initialSkills = [...store.skills.value];

    store.deleteSkill(nonExistentId);

    expect(store.skills.value).toEqual(initialSkills);
  });

  it('getSkillById should return the skill with the given id', () => {
    const id = 3;
    const skill = store.getSkillById(id);

    expect(skill).toBeDefined();
    expect(skill.id).toBe(id);
    expect(skill.name).toBe('火球术');
  });

  it('getSkillById should return undefined for non-existent id', () => {
    const nonExistentId = 999;
    const skill = store.getSkillById(nonExistentId);

    expect(skill).toBeUndefined();
  });

  it('getSkillByCode should return the skill with the given code', () => {
    const code = 'skill_heal';
    const skill = store.getSkillByCode(code);

    expect(skill).toBeDefined();
    expect(skill.code).toBe(code);
    expect(skill.name).toBe('治疗术');
  });

  it('getSkillByCode should return undefined for non-existent code', () => {
    const nonExistentCode = 'non_existent_code';
    const skill = store.getSkillByCode(nonExistentCode);

    expect(skill).toBeUndefined();
  });

  it('updateSkillTree should update the skillTree', () => {
    const newTree = {
      nodes: [{ id: 'new_node', x: 100, y: 100, label: '新节点', type: 'skill' }],
      edges: []
    };

    store.updateSkillTree(newTree);

    expect(store.skillTree.value).toEqual(newTree);
  });
});
```