<template>
  <div class="npc-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>NPC列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建NPC
          </el-button>
        </div>
      </template>

      <el-table :data="npcStore.npcs" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :size="40" :src="row.avatar">
              {{ row.name?.charAt(0) }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="120" />
        <el-table-column prop="title" label="称号" width="120" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="智能体" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.agentId" type="success">已配置</el-tag>
            <el-tag v-else type="info">未配置</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="商店" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.shopId" type="warning">有商店</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" text @click="handleDialogue(row)">对话</el-button>
            <el-button type="danger" text @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useNPCStore } from '@/stores';
import { ElMessageBox, ElMessage } from 'element-plus';

const router = useRouter();
const npcStore = useNPCStore();

const handleAdd = () => {
  router.push('/npc/edit');
};

const handleEdit = (npc) => {
  router.push(`/npc/edit/${npc.id}`);
};

const handleDialogue = (npc) => {
  router.push(`/npc/dialogue/${npc.id}`);
};

const handleDelete = async (npc) => {
  try {
    await ElMessageBox.confirm('确定要删除该NPC吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    npcStore.deleteNPC(npc.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};
</script>

<style scoped>
.npc-list {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
