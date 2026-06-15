<template>
  <div class="scene-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>场景列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建场景
          </el-button>
        </div>
      </template>

      <el-table :data="sceneStore.scenes" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column prop="name" label="场景名称" width="150" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="background" label="背景" width="150" />
        <el-table-column label="尺寸" width="120">
          <template #default="{ row }">
            {{ row.width }} x {{ row.height }}
          </template>
        </el-table-column>
        <el-table-column label="NPC数量" width="100">
          <template #default="{ row }">
            {{ row.npcs?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" text @click="handlePreview(row)">预览</el-button>
            <el-button type="danger" text @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useSceneStore } from '@/stores';
import { ElMessageBox, ElMessage } from 'element-plus';

const router = useRouter();
const sceneStore = useSceneStore();

const handleAdd = () => {
  router.push('/scene/edit');
};

const handleEdit = (scene) => {
  router.push(`/scene/edit/${scene.id}`);
};

const handlePreview = (scene) => {
  ElMessage.info(`预览场景: ${scene.name}`);
};

const handleDelete = async (scene) => {
  try {
    await ElMessageBox.confirm('确定要删除该场景吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    sceneStore.deleteScene(scene.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};
</script>

<style scoped>
.scene-list {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
