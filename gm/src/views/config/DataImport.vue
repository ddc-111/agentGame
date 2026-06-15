<template>
  <div class="data-import">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>数据导入</span>
        </div>
      </template>

      <el-alert
        title="上传JSON格式的游戏数据文件进行导入"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />

      <el-form label-width="120px">
        <el-form-item label="导入模式">
          <el-radio-group v-model="importMode">
            <el-radio label="merge">合并（保留现有数据）</el-radio>
            <el-radio label="overwrite">覆盖（替换现有数据）</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="选择文件">
          <el-upload
            class="upload-area"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            accept=".json"
          >
            <el-icon class="el-icon--upload"><Upload /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">只能上传 JSON 文件</div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleImport" :disabled="!importData">
            <el-icon><Upload /></el-icon>
            导入数据
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider content-position="left">数据预览</el-divider>

      <el-input
        v-model="previewContent"
        type="textarea"
        :rows="15"
        readonly
        placeholder="上传文件后查看预览"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useSceneStore, useNPCStore, useAgentStore, usePromptStore, useShopStore, useTaskStore, useConfigStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const sceneStore = useSceneStore();
const npcStore = useNPCStore();
const agentStore = useAgentStore();
const promptStore = usePromptStore();
const shopStore = useShopStore();
const taskStore = useTaskStore();
const configStore = useConfigStore();

const importMode = ref('merge');
const importData = ref(null);
const previewContent = ref('');

const handleFileChange = (file) => {
  const reader = new FileReader();
  reader.onload = (e) => {
    try {
      const data = JSON.parse(e.target.result);
      importData.value = data;
      previewContent.value = JSON.stringify(data, null, 2);
      ElMessage.success('文件解析成功');
    } catch (error) {
      ElMessage.error('文件格式错误，请上传有效的JSON文件');
    }
  };
  reader.readAsText(file.raw);
};

const handleImport = async () => {
  if (!importData.value) {
    ElMessage.warning('请先上传文件');
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确定要${importMode.value === 'merge' ? '合并' : '覆盖'}导入数据吗？`,
      '确认导入',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );

    const data = importData.value;

    if (data.scenes) {
      if (importMode.value === 'overwrite') {
        sceneStore.scenes = data.scenes;
      } else {
        data.scenes.forEach(scene => {
          if (!sceneStore.getSceneById(scene.id)) {
            sceneStore.addScene(scene);
          }
        });
      }
    }

    if (data.npcs) {
      if (importMode.value === 'overwrite') {
        npcStore.npcs = data.npcs;
      } else {
        data.npcs.forEach(npc => {
          if (!npcStore.getNPCById(npc.id)) {
            npcStore.addNPC(npc);
          }
        });
      }
    }

    if (data.agents) {
      if (importMode.value === 'overwrite') {
        agentStore.agents = data.agents;
      } else {
        data.agents.forEach(agent => {
          if (!agentStore.getAgentById(agent.id)) {
            agentStore.addAgent(agent);
          }
        });
      }
    }

    ElMessage.success('导入成功');
  } catch {
    // 取消操作
  }
};
</script>

<style scoped>
.data-import {
  max-width: 1000px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload-area {
  width: 100%;
}
</style>
