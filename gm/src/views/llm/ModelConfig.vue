<template>
  <div class="model-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>模型配置</span>
        </div>
      </template>

      <el-alert
        title="在此页面可以为不同场景配置使用不同的模型"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />

      <el-form :model="config" label-width="150px">
        <el-divider content-position="left">默认配置</el-divider>

        <el-form-item label="默认提供商">
          <el-select v-model="config.defaultProvider" placeholder="选择默认提供商">
            <el-option
              v-for="provider in llmStore.providers"
              :key="provider.id"
              :label="provider.name"
              :value="provider.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="默认模型">
          <el-select v-model="config.defaultModel" placeholder="选择默认模型">
            <el-option
              v-for="model in defaultModels"
              :key="model.id"
              :label="model.name"
              :value="model.id"
            />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">场景配置</el-divider>

        <el-form-item label="NPC对话">
          <el-select v-model="config.npcChatModel" placeholder="选择模型">
            <el-option label="使用默认" value="default" />
            <el-option
              v-for="model in allModels"
              :key="model.id"
              :label="model.name"
              :value="model.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="摘要生成">
          <el-select v-model="config.summaryModel" placeholder="选择模型">
            <el-option label="使用默认" value="default" />
            <el-option
              v-for="model in allModels"
              :key="model.id"
              :label="model.name"
              :value="model.id"
            />
          </el-select>
          <span class="form-tip">建议使用较便宜的模型</span>
        </el-form-item>

        <el-form-item label="任务生成">
          <el-select v-model="config.taskModel" placeholder="选择模型">
            <el-option label="使用默认" value="default" />
            <el-option
              v-for="model in allModels"
              :key="model.id"
              :label="model.name"
              :value="model.id"
            />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">高级配置</el-divider>

        <el-form-item label="请求超时(秒)">
          <el-input-number v-model="config.timeout" :min="5" :max="120" :step="5" />
        </el-form-item>

        <el-form-item label="最大重试次数">
          <el-input-number v-model="config.maxRetries" :min="0" :max="10" :step="1" />
        </el-form-item>

        <el-form-item label="流式输出">
          <el-switch v-model="config.streamEnabled" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave">保存配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useLLMStore } from '@/stores';
import { ElMessage } from 'element-plus';

const llmStore = useLLMStore();

const config = ref({
  defaultProvider: 'openai',
  defaultModel: 'gpt-4',
  npcChatModel: 'default',
  summaryModel: 'gpt-3.5-turbo',
  taskModel: 'default',
  timeout: 30,
  maxRetries: 3,
  streamEnabled: true
});

const defaultModels = computed(() => {
  const provider = llmStore.providers.find(p => p.id === config.value.defaultProvider);
  return provider?.models || [];
});

const allModels = computed(() => {
  return llmStore.providers.flatMap(p => p.models);
});

const handleSave = () => {
  ElMessage.success('配置已保存');
};
</script>

<style scoped>
.model-config {
  max-width: 1000px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
</style>
