<template>
  <div class="llm-test">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>大模型连接测试</span>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>测试配置</template>
            <el-form :model="testConfig" label-width="100px">
              <el-form-item label="提供商">
                <el-select v-model="testConfig.provider" placeholder="选择提供商">
                  <el-option
                    v-for="provider in llmStore.providers"
                    :key="provider.id"
                    :label="provider.name"
                    :value="provider.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="模型">
                <el-select v-model="testConfig.model" placeholder="选择模型">
                  <el-option
                    v-for="model in availableModels"
                    :key="model.id"
                    :label="model.name"
                    :value="model.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="测试消息">
                <el-input
                  v-model="testConfig.message"
                  type="textarea"
                  :rows="4"
                  placeholder="输入测试消息"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleTest" :loading="loading">
                  发送测试
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>

        <el-col :span="12">
          <el-card>
            <template #header>测试结果</template>
            <div v-if="testResult" class="test-result">
              <div class="result-header">
                <el-tag :type="testResult.success ? 'success' : 'danger'">
                  {{ testResult.success ? '成功' : '失败' }}
                </el-tag>
                <span class="result-time">耗时: {{ testResult.time }}ms</span>
              </div>
              <div class="result-content">
                <pre>{{ testResult.content }}</pre>
              </div>
              <div v-if="testResult.usage" class="result-usage">
                <span>Token使用: {{ testResult.usage.total }}</span>
                <span>提示: {{ testResult.usage.prompt }}</span>
                <span>补全: {{ testResult.usage.completion }}</span>
              </div>
            </div>
            <el-empty v-else description="点击发送测试查看结果" />
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useLLMStore } from '@/stores';
import { ElMessage } from 'element-plus';

const llmStore = useLLMStore();

const testConfig = ref({
  provider: 'openai',
  model: 'gpt-4',
  message: '你好，请介绍一下你自己。'
});

const loading = ref(false);
const testResult = ref(null);

const availableModels = computed(() => {
  const provider = llmStore.providers.find(p => p.id === testConfig.value.provider);
  return provider?.models || [];
});

const handleTest = async () => {
  loading.value = true;
  testResult.value = null;

  // 模拟测试
  setTimeout(() => {
    testResult.value = {
      success: true,
      time: 1234,
      content: '你好！我是一个AI助手，很高兴为你服务。我可以帮助你解答问题、提供信息和进行对话。有什么我可以帮助你的吗？',
      usage: {
        total: 156,
        prompt: 45,
        completion: 111
      }
    };
    loading.value = false;
    ElMessage.success('测试完成');
  }, 1500);
};
</script>

<style scoped>
.llm-test {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.test-result {
  padding: 16px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.result-time {
  color: #909399;
  font-size: 14px;
}

.result-content {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 16px;
}

.result-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.result-usage {
  display: flex;
  gap: 20px;
  color: #606266;
  font-size: 14px;
}
</style>
