<template>
  <div class="demo-player">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>演示场景列表</span>
          </template>
          <div class="demo-list">
            <div
              v-for="demo in demoStore.demos"
              :key="demo.id"
              class="demo-item"
              :class="{ active: demoStore.activeDemo?.id === demo.id }"
              @click="selectDemo(demo.id)"
            >
              <div class="demo-icon">{{ getDemoIcon(demo.category) }}</div>
              <div class="demo-info">
                <div class="demo-name">{{ demo.name }}</div>
                <div class="demo-desc">{{ demo.description }}</div>
                <div class="demo-meta">
                  <el-tag size="small">{{ demo.duration }}</el-tag>
                  <el-tag size="small" type="info">{{ demo.steps }}步</el-tag>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card v-if="demoStore.activeDemo">
          <template #header>
            <div class="card-header">
              <span>{{ demoStore.activeDemo.name }}</span>
              <div class="playback-controls">
                <el-button-group>
                  <el-button @click="demoStore.prevStep()" :disabled="demoStore.currentStep === 0">
                    <el-icon><ArrowLeft /></el-icon>
                  </el-button>
                  <el-button @click="demoStore.togglePlay()" :type="demoStore.isPlaying ? 'warning' : 'primary'">
                    <el-icon>
                      <component :is="demoStore.isPlaying ? 'VideoPause' : 'VideoPlay'" />
                    </el-icon>
                  </el-button>
                  <el-button @click="demoStore.nextStep()" :disabled="demoStore.currentStep >= demoStore.activeDemo.steps - 1">
                    <el-icon><ArrowRight /></el-icon>
                  </el-button>
                </el-button-group>
                <el-select v-model="playbackSpeed" style="width: 100px; margin-left: 8px">
                  <el-option label="0.5x" :value="0.5" />
                  <el-option label="1x" :value="1" />
                  <el-option label="2x" :value="2" />
                </el-select>
              </div>
            </div>
          </template>

          <div class="step-progress">
            <el-progress
              :percentage="(demoStore.currentStep / (demoStore.activeDemo.steps - 1)) * 100"
              :format="() => `${demoStore.currentStep + 1}/${demoStore.activeDemo.steps}`"
            />
          </div>

          <div class="step-content" v-if="currentStepData">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="步骤名称" :span="2">
                <span class="step-name">{{ currentStepData.name }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="描述" :span="2">{{ currentStepData.description }}</el-descriptions-item>
              <el-descriptions-item label="操作" :span="2">
                <div class="action-box">{{ currentStepData.action }}</div>
              </el-descriptions-item>
              <el-descriptions-item label="预期结果" :span="2">
                <div class="expected-box">{{ currentStepData.expected }}</div>
              </el-descriptions-item>
              <el-descriptions-item label="场景" v-if="currentStepData.scene">{{ currentStepData.scene }}</el-descriptions-item>
              <el-descriptions-item label="NPC" v-if="currentStepData.npc">{{ currentStepData.npc }}</el-descriptions-item>
              <el-descriptions-item label="预计时长">{{ currentStepData.duration }}秒</el-descriptions-item>
            </el-descriptions>

            <div class="step-actions">
              <el-button type="success" @click="markResult('pass')">
                <el-icon><CircleCheck /></el-icon> 通过
              </el-button>
              <el-button type="danger" @click="markResult('fail')">
                <el-icon><CircleClose /></el-icon> 失败
              </el-button>
              <el-button type="warning" @click="markResult('skip')">
                <el-icon><DArrowRight /></el-icon> 跳过
              </el-button>
              <el-input
                v-model="stepNote"
                placeholder="备注..."
                style="width: 200px; margin-left: 16px"
              />
            </div>
          </div>

          <el-divider />

          <div class="results-section">
            <h4>测试结果</h4>
            <div class="results-grid">
              <div
                v-for="(result, index) in demoStore.demoResults"
                :key="index"
                class="result-item"
                :class="result.status"
              >
                <span class="result-step">步骤{{ result.step + 1 }}</span>
                <el-tag :type="result.status === 'pass' ? 'success' : result.status === 'fail' ? 'danger' : 'warning'" size="small">
                  {{ result.status === 'pass' ? '通过' : result.status === 'fail' ? '失败' : '跳过' }}
                </el-tag>
                <span class="result-note" v-if="result.note">{{ result.note }}</span>
              </div>
            </div>
            <div class="results-summary" v-if="demoStore.demoResults.length > 0">
              <el-tag type="success">通过: {{ passCount }}</el-tag>
              <el-tag type="danger">失败: {{ failCount }}</el-tag>
              <el-tag type="warning">跳过: {{ skipCount }}</el-tag>
            </div>
          </div>
        </el-card>

        <el-card v-else>
          <el-empty description="选择一个演示场景开始" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useDemoStore } from '@/stores';

const demoStore = useDemoStore();
const playbackSpeed = ref(1);
const stepNote = ref('');

const currentStepData = computed(() => {
  if (!demoStore.activeDemo) return null;
  const steps = getDemoSteps(demoStore.activeDemo.id);
  return steps[demoStore.currentStep] || null;
});

const passCount = computed(() => demoStore.demoResults.filter(r => r.status === 'pass').length);
const failCount = computed(() => demoStore.demoResults.filter(r => r.status === 'fail').length);
const skipCount = computed(() => demoStore.demoResults.filter(r => r.status === 'skip').length);

const getDemoIcon = (category) => {
  const map = { gameplay: '🎮', combat: '⚔️', ai: '🤖' };
  return map[category] || '📋';
};

const getDemoSteps = (demoId) => {
  const stepsMap = {
    beginner_village: [
      { name: '进入游戏', description: '玩家创建角色，出现在青石村村口', action: '创建新角色，选择默认外观，点击开始游戏', expected: '角色出现在村口场景，显示欢迎提示，触发「初来乍到」任务', scene: 'scene_village_entrance', duration: 5 },
      { name: '认识老村长', description: '与村口的老村长对话，了解村庄情况', action: '靠近老村长NPC，点击对话按钮', expected: '老村长用慈祥的语气回应，介绍青石村的情况', scene: 'scene_village_entrance', npc: 'npc_chief_chen', duration: 15 },
      { name: '前往村中心', description: '通过传送点前往村中心广场', action: '走到村口右侧的传送点', expected: '场景切换到村中心广场', scene: 'scene_village_center', duration: 5 },
      { name: '拜访王大娘茶摊', description: '在茶摊与王大娘聊天', action: '走到茶摊区域，与王大娘对话', expected: '王大娘热情招呼，分享村里的故事', scene: 'scene_tea_stand', npc: 'npc_tea_wang', duration: 15 },
      { name: '前往杂货铺', description: '通过传送点进入杂货铺', action: '从村中心走到杂货铺入口传送点', expected: '场景切换到杂货铺', scene: 'scene_general_store', duration: 5 },
      { name: '购买补给品', description: '与李掌柜对话，购买草药和馒头', action: '与李掌柜对话，选择购买3份草药和5个馒头', expected: '完成交易，金币减少，背包增加道具', scene: 'scene_general_store', npc: 'npc_merchant_li', duration: 20 },
      { name: '前往铁匠铺', description: '进入铁匠铺购买武器', action: '从杂货铺返回村中心，走到铁匠铺入口', expected: '场景切换到铁匠铺', scene: 'scene_blacksmith', duration: 5 },
      { name: '购买铁剑', description: '与张铁匠对话，购买一把铁剑', action: '与张铁匠对话，选择购买铁剑', expected: '完成交易，装备铁剑后攻击力提升', scene: 'scene_blacksmith', npc: 'npc_blacksmith_zhang', duration: 20 },
      { name: '前往村外小路', description: '前往村外寻找猎户老周', action: '从村中心右侧传送点前往村外小路', expected: '场景切换到村外小路', scene: 'scene_village_path', duration: 5 },
      { name: '与猎户老周对话', description: '了解狼群的情况，接受委托', action: '与猎户老周对话', expected: '老周说明情况，玩家接受任务', scene: 'scene_village_path', npc: 'npc_hunter_zhou', duration: 15 },
      { name: '遭遇野狼', description: '在村外小路遭遇野狼，进入战斗', action: '在小路上探索，触发随机遭遇', expected: '进入回合制战斗界面', scene: 'scene_village_path', duration: 30 },
      { name: '使用道具恢复', description: '战斗后使用草药恢复生命值', action: '打开背包，使用草药', expected: '生命值恢复20点', duration: 5 },
      { name: '返回村庄交任务', description: '回到村中心向村长汇报', action: '使用回城符或走回村庄', expected: '完成任务链，获得奖励', scene: 'scene_village_entrance', npc: 'npc_chief_chen', duration: 15 },
      { name: '成就展示', description: '查看获得的成就', action: '打开成就面板', expected: '显示已获得的成就', duration: 10 }
    ],
    combat_demo: [
      { name: '准备阶段', description: '确保角色已装备武器，携带草药', action: '打开背包确认装备和道具', expected: '角色装备铁剑，背包有草药', duration: 5 },
      { name: '进入战斗', description: '在村外小路遭遇野狼', action: '触发随机战斗', expected: '进入战斗界面，显示敌方信息', scene: 'scene_village_path', duration: 5 },
      { name: '普通攻击', description: '使用普通攻击造成基础伤害', action: '点击攻击按钮', expected: '造成17点伤害', duration: 5 },
      { name: '使用技能', description: '使用基础斩击技能', action: '选择基础斩击', expected: '造成25点伤害，消耗5MP', duration: 8 },
      { name: '使用增益技能', description: '使用战吼提升攻击力', action: '选择战吼', expected: '攻击力提升30%', duration: 8 },
      { name: '使用道具', description: '生命值较低时使用草药', action: '选择草药使用', expected: 'HP恢复20点', duration: 5 },
      { name: '暴击演示', description: '触发暴击造成高额伤害', action: '继续攻击', expected: '触发暴击，伤害翻倍', duration: 5 },
      { name: '击败敌人', description: '将野狼HP降为0', action: '继续攻击', expected: '胜利，获得奖励', duration: 10 },
      { name: '逃跑演示', description: '展示逃跑机制', action: '选择逃跑', expected: '有概率逃跑成功', duration: 10 },
      { name: '战斗失败', description: '展示HP归零的处理', action: '让HP降为0', expected: '自动回城恢复', duration: 10 }
    ],
    npc_ai_demo: [
      { name: '初次对话', description: '与老村长进行初次对话', action: '说「你好，我是新来的」', expected: '老村长用慈祥语气回应', npc: 'npc_chief_chen', duration: 15 },
      { name: '追问村庄历史', description: '询问村庄的历史背景', action: '问「这个村子有什么故事？」', expected: '讲述青石村的由来', npc: 'npc_chief_chen', duration: 15 },
      { name: '请求任务指引', description: '询问有什么可以帮忙的', action: '说「有什么我能帮忙的？」', expected: '引导去找其他NPC', npc: 'npc_chief_chen', duration: 15 },
      { name: '智能商品推荐', description: '与商人对话获取推荐', action: '说「有什么推荐的？」', expected: '根据等级推荐商品', npc: 'npc_merchant_li', duration: 15 },
      { name: '讨价还价', description: '尝试与商人讨价还价', action: '说「能便宜点吗？」', expected: '可能给予折扣', npc: 'npc_merchant_li', duration: 15 },
      { name: '时间感知', description: '根据时间触发不同对话', action: '与王大娘对话', expected: '根据时间讲故事', npc: 'npc_tea_wang', duration: 15 },
      { name: '隐藏信息获取', description: '通过关键词获取隐藏信息', action: '问「哪里有宝藏？」', expected: '透露山洞秘密', npc: 'npc_tea_wang', duration: 15 },
      { name: '知识问答', description: '回答铁匠的问题', action: '回答兵器知识问题', expected: '答对获得折扣', npc: 'npc_blacksmith_zhang', duration: 20 },
      { name: '动态任务生成', description: '根据等级获取不同任务', action: '与猎户对话', expected: '给出适合等级的任务', npc: 'npc_hunter_zhou', duration: 15 }
    ]
  };
  return stepsMap[demoId] || [];
};

const selectDemo = (id) => {
  demoStore.selectDemo(id);
  stepNote.value = '';
};

const markResult = (status) => {
  demoStore.addResult({ status, note: stepNote.value });
  stepNote.value = '';
  if (status !== 'skip' && demoStore.currentStep < demoStore.activeDemo.steps - 1) {
    demoStore.nextStep();
  }
};

watch(playbackSpeed, (val) => {
  demoStore.setSpeed(val);
});
</script>

<style scoped>
.demo-player {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.playback-controls {
  display: flex;
  align-items: center;
}

.demo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.demo-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.demo-item:hover {
  border-color: #409eff;
  background: #f0f9ff;
}

.demo-item.active {
  border-color: #409eff;
  background: #ecf5ff;
}

.demo-icon {
  font-size: 32px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 8px;
}

.demo-info {
  flex: 1;
}

.demo-name {
  font-weight: 500;
  color: #303133;
}

.demo-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.demo-meta {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.step-progress {
  margin-bottom: 20px;
}

.step-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.action-box {
  padding: 8px 12px;
  background: #f0f9ff;
  border-left: 3px solid #409eff;
  border-radius: 0 4px 4px 0;
}

.expected-box {
  padding: 8px 12px;
  background: #f0f9eb;
  border-left: 3px solid #67c23a;
  border-radius: 0 4px 4px 0;
}

.step-actions {
  display: flex;
  align-items: center;
  margin-top: 16px;
}

.results-section h4 {
  margin: 0 0 12px;
}

.results-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.result-step {
  font-size: 12px;
  color: #909399;
}

.result-note {
  font-size: 12px;
  color: #606266;
}

.results-summary {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}
</style>
