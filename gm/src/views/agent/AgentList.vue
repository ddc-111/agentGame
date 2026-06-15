<template>
  <div class="agent-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>智能体列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建智能体
          </el-button>
        </div>
      </template>

      <el-table :data="agentStore.agents" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="模型配置" width="200">
          <template #default="{ row }">
            <el-tag>{{ row.llmConfig?.provider }}</el-tag>
            <el-tag type="info" style="margin-left: 4px">{{ row.llmConfig?.model }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="记忆配置" width="120">
          <template #default="{ row }">
            <el-tag :type="row.memoryConfig?.summaryEnabled ? 'success' : 'info'">
              {{ row.memoryConfig?.summaryEnabled ? '已启用摘要' : '基础模式' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="知识库" width="100">
          <template #default="{ row }">
            {{ row.knowledgeBase?.length || 0 }} 个
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" text @click="handleMemory(row)">记忆</el-button>
            <el-button type="primary" text @click="handleTest(row)">测试</el-button>
            <el-button type="danger" text @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useAgentStore } from '@/stores';
import { ElMessageBox, ElMessage } from 'element-plus';

const router = useRouter();
const agentStore = useAgentStore();

const handleAdd = () => {
  router.push('/agent/edit');
};

const handleEdit = (agent) => {
  router.push(`/agent/edit/${agent.id}`);
};

const handleMemory = (agent) => {
  router.push(`/agent/memory/${agent.id}`);
};

const handleTest = (agent) => {
  ElMessage.info(`测试智能体: ${agent.name}`);
};

const handleDelete = async (agent) => {
  try {
    await ElMessageBox.confirm('确定要删除该智能体吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    agentStore.deleteAgent(agent.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};
</script>

<style scoped>
.agent-list {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
