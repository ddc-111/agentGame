import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSkillStore = defineStore('skill', () => {
  const skills = ref([
    {
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
    },
    {
      id: 2,
      name: '治疗术',
      code: 'skill_heal',
      description: '使用法力恢复生命值',
      type: 'heal',
      mpCost: 8,
      damage: 0,
      heal: 30,
      cooldown: 1,
      level: 3,
      effect: {}
    },
    {
      id: 3,
      name: '火球术',
      code: 'skill_fireball',
      description: '发射火球造成魔法伤害',
      type: 'attack',
      mpCost: 15,
      damage: 25,
      heal: 0,
      cooldown: 1,
      level: 5,
      effect: { type: 'burn', duration: 2 }
    },
    {
      id: 4,
      name: '铁壁',
      code: 'skill_iron_wall',
      description: '提升防御力3回合',
      type: 'buff',
      mpCost: 10,
      damage: 0,
      heal: 0,
      cooldown: 2,
      level: 4,
      effect: { type: 'defense_up', duration: 3, value: 50 }
    },
    {
      id: 5,
      name: '疾风步',
      code: 'skill_wind_step',
      description: '提升速度3回合',
      type: 'buff',
      mpCost: 12,
      damage: 0,
      heal: 0,
      cooldown: 2,
      level: 6,
      effect: { type: 'speed_up', duration: 3, value: 50 }
    },
    {
      id: 6,
      name: '毒击',
      code: 'skill_poison_strike',
      description: '附带毒素的攻击，使敌人持续掉血',
      type: 'debuff',
      mpCost: 10,
      damage: 10,
      heal: 0,
      cooldown: 2,
      level: 3,
      effect: { type: 'poison', duration: 3, value: 5 }
    },
    {
      id: 7,
      name: '战吼',
      code: 'skill_warcry',
      description: '发出战吼，提升攻击力2回合',
      type: 'buff',
      mpCost: 8,
      damage: 0,
      heal: 0,
      cooldown: 2,
      level: 2,
      effect: { type: 'attack_up', duration: 2, value: 30 }
    }
  ]);

  const skillTree = ref({
    nodes: [
      { id: 'root', x: 400, y: 50, label: '技能树', type: 'root' },
      { id: 'atk_1', x: 200, y: 150, label: '基础斩击', type: 'skill', skillCode: 'skill_basic_slash' },
      { id: 'atk_2', x: 100, y: 280, label: '火球术', type: 'skill', skillCode: 'skill_fireball', requires: ['atk_1'] },
      { id: 'atk_3', x: 300, y: 280, label: '毒击', type: 'skill', skillCode: 'skill_poison_strike', requires: ['atk_1'] },
      { id: 'buf_1', x: 600, y: 150, label: '战吼', type: 'skill', skillCode: 'skill_warcry' },
      { id: 'buf_2', x: 500, y: 280, label: '铁壁', type: 'skill', skillCode: 'skill_iron_wall', requires: ['buf_1'] },
      { id: 'buf_3', x: 700, y: 280, label: '疾风步', type: 'skill', skillCode: 'skill_wind_step', requires: ['buf_1'] },
      { id: 'heal_1', x: 400, y: 380, label: '治疗术', type: 'skill', skillCode: 'skill_heal' }
    ],
    edges: [
      { source: 'root', target: 'atk_1' },
      { source: 'root', target: 'buf_1' },
      { source: 'root', target: 'heal_1' },
      { source: 'atk_1', target: 'atk_2' },
      { source: 'atk_1', target: 'atk_3' },
      { source: 'buf_1', target: 'buf_2' },
      { source: 'buf_1', target: 'buf_3' }
    ]
  });

  const addSkill = (skill) => {
    skills.value.push({
      id: Date.now(),
      ...skill
    });
  };

  const updateSkill = (id, data) => {
    const index = skills.value.findIndex(s => s.id === id);
    if (index !== -1) {
      skills.value[index] = { ...skills.value[index], ...data };
    }
  };

  const deleteSkill = (id) => {
    skills.value = skills.value.filter(s => s.id !== id);
  };

  const getSkillById = (id) => {
    return skills.value.find(s => s.id === id);
  };

  const getSkillByCode = (code) => {
    return skills.value.find(s => s.code === code);
  };

  const updateSkillTree = (tree) => {
    skillTree.value = tree;
  };

  return {
    skills,
    skillTree,
    addSkill,
    updateSkill,
    deleteSkill,
    getSkillById,
    getSkillByCode,
    updateSkillTree
  };
});
