import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useDemoStore = defineStore('demo', () => {
  const demos = ref([
    {
      id: 'beginner_village',
      name: '新手村完整流程',
      description: '展示从角色创建到完成第一个任务的完整游戏流程',
      category: 'gameplay',
      duration: '5分钟',
      steps: 14,
      tags: ['新手', '主线', '全流程']
    },
    {
      id: 'combat_demo',
      name: '战斗系统演示',
      description: '展示回合制战斗的完整机制',
      category: 'combat',
      duration: '3分钟',
      steps: 12,
      tags: ['战斗', '技能', '道具']
    },
    {
      id: 'npc_ai_demo',
      name: 'AI对话演示',
      description: '展示AI驱动的NPC对话系统',
      category: 'ai',
      duration: '4分钟',
      steps: 13,
      tags: ['AI', '对话', 'NPC']
    }
  ]);

  const activeDemo = ref(null);
  const currentStep = ref(0);
  const isPlaying = ref(false);
  const playbackSpeed = ref(1);

  const demoResults = ref([]);

  const selectDemo = (demoId) => {
    activeDemo.value = demos.value.find(d => d.id === demoId) || null;
    currentStep.value = 0;
    demoResults.value = [];
    isPlaying.value = false;
  };

  const nextStep = () => {
    if (activeDemo.value && currentStep.value < activeDemo.value.steps - 1) {
      currentStep.value++;
      return true;
    }
    return false;
  };

  const prevStep = () => {
    if (currentStep.value > 0) {
      currentStep.value--;
      return true;
    }
    return false;
  };

  const goToStep = (step) => {
    if (activeDemo.value && step >= 0 && step < activeDemo.value.steps) {
      currentStep.value = step;
    }
  };

  const togglePlay = () => {
    isPlaying.value = !isPlaying.value;
  };

  const stop = () => {
    isPlaying.value = false;
    currentStep.value = 0;
  };

  const setSpeed = (speed) => {
    playbackSpeed.value = speed;
  };

  const addResult = (result) => {
    demoResults.value.push({
      step: currentStep.value,
      timestamp: Date.now(),
      ...result
    });
  };

  const resetResults = () => {
    demoResults.value = [];
  };

  return {
    demos,
    activeDemo,
    currentStep,
    isPlaying,
    playbackSpeed,
    demoResults,
    selectDemo,
    nextStep,
    prevStep,
    goToStep,
    togglePlay,
    stop,
    setSpeed,
    addResult,
    resetResults
  };
});
