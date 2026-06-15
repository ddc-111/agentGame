<template>
  <div class="data-export">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>数据导出</span>
        </div>
      </template>

      <el-alert
        title="选择要导出的数据类型，然后点击导出按钮"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />

      <el-form label-width="120px">
        <el-form-item label="导出内容">
          <el-checkbox-group v-model="exportTypes">
            <el-checkbox label="scenes">场景数据</el-checkbox>
            <el-checkbox label="npcs">NPC数据</el-checkbox>
            <el-checkbox label="agents">智能体数据</el-checkbox>
            <el-checkbox label="prompts">提示词模板</el-checkbox>
            <el-checkbox label="shops">商店数据</el-checkbox>
            <el-checkbox label="items">道具数据</el-checkbox>
            <el-checkbox label="tasks">任务数据</el-checkbox>
            <el-checkbox label="config">游戏配置</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="导出格式">
          <el-radio-group v-model="exportFormat">
            <el-radio label="json">JSON</el-radio>
            <el-radio label="yaml">YAML</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出数据
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider content-position="left">导出预览</el-divider>

      <el-input
        v-model="exportPreview"
        type="textarea"
        :rows="15"
        readonly
        placeholder="点击导出按钮查看预览"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useSceneStore, useNPCStore, useAgentStore, usePromptStore, useShopStore, useTaskStore, useConfigStore } from '@/stores';
import { ElMessage } from 'element-plus';

const sceneStore = useSceneStore();
const npcStore = useNPCStore();
const agentStore = useAgentStore();
const promptStore = usePromptStore();
const shopStore = useShopStore();
const taskStore = useTaskStore();
const configStore = useConfigStore();

const exportTypes = ref(['scenes', 'npcs', 'agents']);
const exportFormat = ref('json');
const exportPreview = ref('');

const handleExport = () => {
  const data = {};

  if (exportTypes.value.includes('scenes')) {
    data.scenes = sceneStore.scenes;
  }
  if (exportTypes.value.includes('npcs')) {
    data.npcs = npcStore.npcs;
  }
  if (exportTypes.value.includes('agents')) {
    data.agents = agentStore.agents;
  }
  if (exportTypes.value.includes('prompts')) {
    data.prompts = promptStore.templates;
  }
  if (exportTypes.value.includes('shops')) {
    data.shops = shopStore.shops;
  }
  if (exportTypes.value.includes('items')) {
    data.items = shopStore.items;
  }
  if (exportTypes.value.includes('tasks')) {
    data.tasks = taskStore.tasks;
  }
  if (exportTypes.value.includes('config')) {
    data.config = configStore.gameConfig;
  }

  const jsonStr = JSON.stringify(data, null, 2);
  exportPreview.value = jsonStr;

  // 下载文件
  const blob = new Blob([jsonStr], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `game_data_${new Date().toISOString().slice(0, 10)}.json`;
  a.click();
  URL.revokeObjectURL(url);

  ElMessage.success('导出成功');
};
</script>

<style scoped>
.data-export {
  max-width: 1000px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
