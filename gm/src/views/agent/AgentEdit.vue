<template>
  <div class="agent-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑智能体' : '新建智能体' }}</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="智能体ID">
              <el-input v-model="form.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="名称" required>
              <el-input v-model="form.name" placeholder="请输入智能体名称" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="描述">
              <el-input v-model="form.description" placeholder="请输入描述" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">大模型配置</el-divider>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="模型提供商">
              <el-select v-model="form.llmConfig.provider" placeholder="选择提供商">
                <el-option
                  v-for="provider in llmStore.providers"
                  :key="provider.id"
                  :label="provider.name"
                  :value="provider.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="模型">
              <el-select v-model="form.llmConfig.model" placeholder="选择模型">
                <el-option
                  v-for="model in availableModels"
                  :key="model.id"
                  :label="model.name"
                  :value="model.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Temperature">
              <el-slider v-model="form.llmConfig.temperature" :min="0" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="最大Token">
              <el-input-number v-model="form.llmConfig.maxTokens" :min="100" :max="8000" :step="100" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">系统提示词</el-divider>

        <el-form-item>
          <el-input
            v-model="form.systemPrompt"
            type="textarea"
            :rows="10"
            placeholder="输入系统提示词，定义NPC的人设、行为规则等"
          />
          <div class="prompt-tools">
            <el-button size="small" @click="insertVariable('npc_name')">插入NPC名称</el-button>
            <el-button size="small" @click="insertVariable('current_scene')">插入当前场景</el-button>
            <el-button size="small" @click="insertVariable('player_name')">插入玩家名称</el-button>
            <el-button size="small" @click="insertVariable('current_time')">插入当前时间</el-button>
          </div>
        </el-form-item>

        <el-divider content-position="left">知识库</el-divider>

        <el-form-item>
          <el-button @click="handleAddKnowledge">添加知识库</el-button>
          <el-table :data="form.knowledgeBase" style="margin-top: 10px;">
            <el-table-column prop="id" label="ID" width="120" />
            <el-table-column prop="name" label="名称" width="200" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag>{{ row.type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemoveKnowledge($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>

        <el-divider content-position="left">工具</el-divider>

        <el-form-item>
          <el-button @click="handleAddTool">添加工具</el-button>
          <el-table :data="form.tools" style="margin-top: 10px;">
            <el-table-column prop="name" label="工具名称" width="200" />
            <el-table-column prop="description" label="描述" />
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemoveTool($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAgentStore, useLLMStore } from '@/stores';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();
const agentStore = useAgentStore();
const llmStore = useLLMStore();

const isEdit = computed(() => !!route.params.id);

const form = ref({
  id: '',
  name: '',
  description: '',
  llmConfig: {
    provider: 'openai',
    model: 'gpt-4',
    temperature: 0.7,
    maxTokens: 500
  },
  systemPrompt: '',
  memoryConfig: {
    type: 'sliding_window',
    maxMessages: 20,
    summaryEnabled: true,
    summaryThreshold: 50
  },
  knowledgeBase: [],
  tools: []
});

const availableModels = computed(() => {
  const provider = llmStore.providers.find(p => p.id === form.value.llmConfig.provider);
  return provider?.models || [];
});

onMounted(() => {
  if (route.params.id) {
    const agent = agentStore.getAgentById(route.params.id);
    if (agent) {
      form.value = { ...agent };
    }
  } else {
    form.value.id = `agent_${Date.now()}`;
  }
});

const insertVariable = (varName) => {
  const textarea = document.querySelector('.agent-edit textarea');
  if (textarea) {
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const text = form.value.systemPrompt;
    form.value.systemPrompt = text.substring(0, start) + `{{${varName}}}` + text.substring(end);
  }
};

const handleAddKnowledge = () => {
  form.value.knowledgeBase.push({
    id: `kb_${Date.now()}`,
    name: '新知识库',
    type: 'text'
  });
};

const handleRemoveKnowledge = (index) => {
  form.value.knowledgeBase.splice(index, 1);
};

const handleAddTool = () => {
  form.value.tools.push({
    name: '',
    description: ''
  });
};

const handleRemoveTool = (index) => {
  form.value.tools.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.name) {
    ElMessage.warning('请输入智能体名称');
    return;
  }

  if (isEdit.value) {
    agentStore.updateAgent(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    agentStore.addAgent(form.value);
    ElMessage.success('创建成功');
  }
  router.push('/agent/list');
};

const handleCancel = () => {
  router.push('/agent/list');
};
</script>

<style scoped>
.agent-edit {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.prompt-tools {
  margin-top: 10px;
  display: flex;
  gap: 8px;
}
</style>
