<template>
  <div class="memory-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>记忆配置 - {{ agent?.name || '未选择智能体' }}</span>
          <el-button type="primary" @click="handleSave">保存</el-button>
        </div>
      </template>

      <el-form :model="form" label-width="150px">
        <el-divider content-position="left">短期记忆</el-divider>

        <el-form-item label="记忆类型">
          <el-radio-group v-model="form.type">
            <el-radio label="sliding_window">滑动窗口</el-radio>
            <el-radio label="token_limit">Token限制</el-radio>
            <el-radio label="conversation_length">对话长度</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.type === 'sliding_window'" label="最大消息数">
          <el-input-number v-model="form.maxMessages" :min="5" :max="100" :step="5" />
          <span class="form-tip">保留最近N条消息作为上下文</span>
        </el-form-item>

        <el-form-item v-if="form.type === 'token_limit'" label="最大Token数">
          <el-input-number v-model="form.maxTokens" :min="1000" :max="100000" :step="1000" />
          <span class="form-tip">当上下文超过Token限制时，自动压缩早期对话</span>
        </el-form-item>

        <el-divider content-position="left">长期记忆（摘要）</el-divider>

        <el-form-item label="启用摘要">
          <el-switch v-model="form.summaryEnabled" />
          <span class="form-tip">启用后，系统会自动总结对话内容形成长期记忆</span>
        </el-form-item>

        <template v-if="form.summaryEnabled">
          <el-form-item label="摘要触发阈值">
            <el-input-number v-model="form.summaryThreshold" :min="10" :max="200" :step="10" />
            <span class="form-tip">当对话消息数超过此值时触发摘要生成</span>
          </el-form-item>

          <el-form-item label="摘要模型">
            <el-select v-model="form.summaryModel" placeholder="选择摘要模型">
              <el-option label="使用相同模型" value="same" />
              <el-option label="GPT-3.5 Turbo" value="gpt-3.5-turbo" />
              <el-option label="本地模型" value="local" />
            </el-select>
            <span class="form-tip">建议使用较便宜的模型生成摘要以节省成本</span>
          </el-form-item>

          <el-form-item label="摘要提示词">
            <el-input
              v-model="form.summaryPrompt"
              type="textarea"
              :rows="6"
              placeholder="输入用于生成摘要的提示词"
            />
          </el-form-item>
        </template>

        <el-divider content-position="left">记忆检索</el-divider>

        <el-form-item label="启用向量检索">
          <el-switch v-model="form.vectorSearchEnabled" />
          <span class="form-tip">启用后，可以从历史对话中检索相关信息</span>
        </el-form-item>

        <template v-if="form.vectorSearchEnabled">
          <el-form-item label="向量维度">
            <el-input-number v-model="form.vectorDimension" :min="128" :max="4096" :step="128" />
          </el-form-item>

          <el-form-item label="检索数量">
            <el-input-number v-model="form.topK" :min="1" :max="20" :step="1" />
            <span class="form-tip">每次检索返回的相关记忆条数</span>
          </el-form-item>

          <el-form-item label="相似度阈值">
            <el-slider v-model="form.similarityThreshold" :min="0" :max="1" :step="0.05" show-input />
            <span class="form-tip">只有相似度高于此值的记忆才会被返回</span>
          </el-form-item>
        </template>

        <el-divider content-position="left">记忆预览</el-divider>

        <el-form-item>
          <el-button @click="handleTestMemory">测试记忆</el-button>
          <el-button @click="handleClearMemory">清空记忆</el-button>
        </el-form-item>

        <el-card v-if="memoryPreview.length > 0" class="memory-preview-card">
          <div v-for="(memory, index) in memoryPreview" :key="index" class="memory-item">
            <div class="memory-header">
              <el-tag :type="memory.type === 'summary' ? 'warning' : 'info'">
                {{ memory.type === 'summary' ? '摘要' : '消息' }}
              </el-tag>
              <span class="memory-time">{{ memory.time }}</span>
            </div>
            <div class="memory-content">{{ memory.content }}</div>
          </div>
        </el-card>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAgentStore } from '@/stores';
import { ElMessage } from 'element-plus';

const route = useRoute();
const agentStore = useAgentStore();

const agent = computed(() => agentStore.getAgentById(route.params.id));

const form = ref({
  type: 'sliding_window',
  maxMessages: 20,
  maxTokens: 4000,
  summaryEnabled: true,
  summaryThreshold: 50,
  summaryModel: 'same',
  summaryPrompt: `请将以下对话总结为简短的摘要，保留关键信息：
1. 对话主题
2. 达成的共识或交易
3. 需要记住的重要信息

对话内容：
{{conversation}}`,
  vectorSearchEnabled: false,
  vectorDimension: 1536,
  topK: 5,
  similarityThreshold: 0.7
});

const memoryPreview = ref([]);

onMounted(() => {
  if (agent.value?.memoryConfig) {
    form.value = { ...form.value, ...agent.value.memoryConfig };
  }
});

const handleSave = () => {
  if (agent.value) {
    agentStore.updateAgent(agent.value.id, {
      memoryConfig: { ...form.value }
    });
    ElMessage.success('保存成功');
  }
};

const handleTestMemory = () => {
  memoryPreview.value = [
    {
      type: 'summary',
      time: '2024-01-15 10:30',
      content: '玩家询问了草药的价格，购买了3份草药，花费300文。'
    },
    {
      type: 'message',
      time: '2024-01-15 10:28',
      content: '客官，草药100文一份，要几份？'
    },
    {
      type: 'message',
      time: '2024-01-15 10:27',
      content: '我想买些草药。'
    }
  ];
};

const handleClearMemory = () => {
  memoryPreview.value = [];
  ElMessage.success('记忆已清空');
};
</script>

<style scoped>
.memory-config {
  max-width: 1200px;
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

.memory-preview-card {
  margin-top: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.memory-item {
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
}

.memory-item:last-child {
  border-bottom: none;
}

.memory-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.memory-time {
  font-size: 12px;
  color: #909399;
}

.memory-content {
  font-size: 14px;
  color: #303133;
  line-height: 1.6;
}
</style>
