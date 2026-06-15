<template>
  <div class="task-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>任务列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建任务
          </el-button>
        </div>
      </template>

      <el-table :data="taskStore.tasks" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'main' ? 'primary' : 'info'">
              {{ row.type === 'main' ? '主线' : '支线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="目标" width="100">
          <template #default="{ row }">
            {{ row.objectives?.length || 0 }} 个
          </template>
        </el-table-column>
        <el-table-column label="奖励" width="200">
          <template #default="{ row }">
            <span v-if="row.rewards">
              经验: {{ row.rewards.exp }}, 金币: {{ row.rewards.gold }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" text @click="handleFlow(row)">流程</el-button>
            <el-button type="danger" text @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useTaskStore } from '@/stores';
import { ElMessageBox, ElMessage } from 'element-plus';

const router = useRouter();
const taskStore = useTaskStore();

const getStatusTag = (status) => {
  const map = {
    active: 'success',
    locked: 'info',
    completed: 'warning',
    failed: 'danger'
  };
  return map[status] || '';
};

const handleAdd = () => {
  router.push('/task/edit');
};

const handleEdit = (task) => {
  router.push(`/task/edit/${task.id}`);
};

const handleFlow = (task) => {
  router.push(`/task/flow/${task.id}`);
};

const handleDelete = async (task) => {
  try {
    await ElMessageBox.confirm('确定要删除该任务吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    taskStore.deleteTask(task.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};
</script>

<style scoped>
.task-list {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
