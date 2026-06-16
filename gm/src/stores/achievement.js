import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAchievementStore = defineStore('achievement', () => {
  const achievements = ref([
    {
      id: 1,
      name: '初来乍到',
      code: 'ach_first_quest',
      description: '完成第一个任务',
      condition: { type: 'quest_complete', value: 1 },
      reward: { exp: 50, gold: 100 },
      icon: '⭐'
    },
    {
      id: 2,
      name: '村庄之友',
      code: 'ach_talk_all_npcs',
      description: '与所有NPC对话过',
      condition: { type: 'talk_all_npcs', value: 1 },
      reward: { exp: 200, gold: 300 },
      icon: '👥'
    },
    {
      id: 3,
      name: '富甲一方',
      code: 'ach_rich',
      description: '累计获得10000金币',
      condition: { type: 'gold', value: 10000 },
      reward: { exp: 500, gold: 1000 },
      icon: '💰'
    },
    {
      id: 4,
      name: '百战百胜',
      code: 'ach_combat_100',
      description: '赢得100场战斗',
      condition: { type: 'combat_win', value: 100 },
      reward: { exp: 1000, gold: 2000 },
      icon: '⚔️'
    },
    {
      id: 5,
      name: '探索者',
      code: 'ach_explorer',
      description: '探索所有场景',
      condition: { type: 'explore', value: 6 },
      reward: { exp: 300, gold: 500 },
      icon: '🗺️'
    },
    {
      id: 6,
      name: '收藏家',
      code: 'ach_collector',
      description: '拥有50种不同的道具',
      condition: { type: 'collect', value: 50 },
      reward: { exp: 800, gold: 1500 },
      icon: '🎒'
    },
    {
      id: 7,
      name: '初试牛刀',
      code: 'ach_first_combat',
      description: '赢得第一场战斗',
      condition: { type: 'combat_win', value: 1 },
      reward: { exp: 30, gold: 50 },
      icon: '🗡️'
    },
    {
      id: 8,
      name: '技能大师',
      code: 'ach_skill_master',
      description: '使用100次技能',
      condition: { type: 'skill_use', value: 100 },
      reward: { exp: 600, gold: 800 },
      icon: '✨'
    },
    {
      id: 9,
      name: '等级10',
      code: 'ach_level_10',
      description: '达到10级',
      condition: { type: 'level', value: 10 },
      reward: { exp: 200, gold: 500 },
      icon: '🏆'
    },
    {
      id: 10,
      name: '等级20',
      code: 'ach_level_20',
      description: '达到20级',
      condition: { type: 'level', value: 20 },
      reward: { exp: 500, gold: 1000 },
      icon: '👑'
    }
  ]);

  const conditionTypes = [
    { value: 'quest_complete', label: '完成任务数', description: '完成指定数量的任务' },
    { value: 'talk_all_npcs', label: '对话所有NPC', description: '与所有NPC至少对话一次' },
    { value: 'gold', label: '累计金币', description: '累计获得指定数量的金币' },
    { value: 'combat_win', label: '战斗胜利', description: '赢得指定场数的战斗' },
    { value: 'explore', label: '探索场景', description: '探索指定数量的场景' },
    { value: 'collect', label: '收集道具', description: '拥有指定种类的道具' },
    { value: 'skill_use', label: '使用技能', description: '使用指定次数的技能' },
    { value: 'level', label: '达到等级', description: '达到指定等级' },
    { value: 'item_use', label: '使用道具', description: '使用指定次数的道具' },
    { value: 'death', label: '死亡次数', description: '累计死亡指定次数' }
  ];

  const addAchievement = (achievement) => {
    achievements.value.push({
      id: Date.now(),
      ...achievement
    });
  };

  const updateAchievement = (id, data) => {
    const index = achievements.value.findIndex(a => a.id === id);
    if (index !== -1) {
      achievements.value[index] = { ...achievements.value[index], ...data };
    }
  };

  const deleteAchievement = (id) => {
    achievements.value = achievements.value.filter(a => a.id !== id);
  };

  const getAchievementById = (id) => {
    return achievements.value.find(a => a.id === id);
  };

  return {
    achievements,
    conditionTypes,
    addAchievement,
    updateAchievement,
    deleteAchievement,
    getAchievementById
  };
});
