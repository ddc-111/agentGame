<template>
  <div class="home">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon" style="background-color: #409eff">
            <el-icon size="24"><Picture /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ sceneStore.scenes.length }}</div>
            <div class="stat-label">场景数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon" style="background-color: #67c23a">
            <el-icon size="24"><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ npcStore.npcs.length }}</div>
            <div class="stat-label">NPC数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon" style="background-color: #e6a23c">
            <el-icon size="24"><ChatDotRound /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ agentStore.agents.length }}</div>
            <div class="stat-label">智能体数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon" style="background-color: #f56c6c">
            <el-icon size="24"><List /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ taskStore.tasks.length }}</div>
            <div class="stat-label">任务数量</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快速操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/scene/edit')">
              <el-icon><Plus /></el-icon>
              新建场景
            </el-button>
            <el-button type="success" @click="$router.push('/npc/edit')">
              <el-icon><Plus /></el-icon>
              新建NPC
            </el-button>
            <el-button type="warning" @click="$router.push('/agent/edit')">
              <el-icon><Plus /></el-icon>
              新建智能体
            </el-button>
            <el-button type="danger" @click="$router.push('/task/edit')">
              <el-icon><Plus /></el-icon>
              新建任务
            </el-button>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统状态</span>
            </div>
          </template>
          <div class="system-status">
            <div class="status-item">
              <span class="status-label">服务端状态</span>
              <el-tag type="success">运行中</el-tag>
            </div>
            <div class="status-item">
              <span class="status-label">数据库连接</span>
              <el-tag type="success">正常</el-tag>
            </div>
            <div class="status-item">
              <span class="status-label">AI服务</span>
              <el-tag type="warning">未配置</el-tag>
            </div>
            <div class="status-item">
              <span class="status-label">最后保存</span>
              <span>{{ lastSaveTime }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>NPC购物流程预览</span>
              <el-button type="primary" text @click="$router.push('/task/flow/flow_001')">编辑流程</el-button>
            </div>
          </template>
          <div class="flow-preview">
            <div class="flow-steps">
              <div class="flow-step" v-for="(step, index) in shoppingFlow" :key="index">
                <div class="step-icon" :class="step.type">
                  <el-icon size="20">
                    <component :is="step.icon" />
                  </el-icon>
                </div>
                <div class="step-content">
                  <div class="step-title">{{ step.title }}</div>
                  <div class="step-desc">{{ step.description }}</div>
                </div>
                <el-icon v-if="index < shoppingFlow.length - 1" class="step-arrow"><ArrowRight /></el-icon>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useSceneStore, useNPCStore, useAgentStore, useTaskStore } from '@/stores';

const sceneStore = useSceneStore();
const npcStore = useNPCStore();
const agentStore = useAgentStore();
const taskStore = useTaskStore();

const lastSaveTime = ref(new Date().toLocaleString());

const shoppingFlow = [
  {
    type: 'start',
    icon: 'HomeFilled',
    title: 'NPC在家',
    description: '李掌柜在家中准备出门'
  },
  {
    type: 'action',
    icon: 'Walk',
    title: '前往商店',
    description: '离开家，前往杂货铺'
  },
  {
    type: 'condition',
    icon: 'QuestionFilled',
    title: '检查商店',
    description: '确认商店是否开门'
  },
  {
    type: 'action',
    icon: 'ShoppingCart',
    title: '购买物品',
    description: '购买草药、馒头等物资'
  },
  {
    type: 'action',
    icon: 'House',
    title: '返回家中',
    description: '带着购买的物品回家'
  },
  {
    type: 'end',
    icon: 'CircleCheck',
    title: '流程结束',
    description: 'NPC完成购物，更新库存'
  }
];
</script>

<style scoped>
.home {
  max-width: 1400px;
  margin: 0 auto;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
}

.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  width: 100%;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 16px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.content-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quick-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.system-status {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-label {
  color: #606266;
}

.flow-preview {
  overflow-x: auto;
  padding: 20px 0;
}

.flow-steps {
  display: flex;
  align-items: center;
  gap: 16px;
  min-width: max-content;
}

.flow-step {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.step-icon.start {
  background-color: #909399;
}

.step-icon.action {
  background-color: #409eff;
}

.step-icon.condition {
  background-color: #e6a23c;
}

.step-icon.end {
  background-color: #67c23a;
}

.step-content {
  min-width: 120px;
}

.step-title {
  font-weight: bold;
  color: #303133;
  font-size: 14px;
}

.step-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.step-arrow {
  color: #c0c4cc;
  font-size: 20px;
}
</style>
