<template>
  <div class="shop-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>商店列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建商店
          </el-button>
        </div>
      </template>

      <el-table :data="shopStore.shops" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column label="店主" width="120">
          <template #default="{ row }">
            {{ getNPCName(row.owner) }}
          </template>
        </el-table-column>
        <el-table-column label="商品数量" width="100">
          <template #default="{ row }">
            {{ row.items?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="营业时间" width="150">
          <template #default="{ row }">
            {{ row.openTime }} - {{ row.closeTime }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" text @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useShopStore, useNPCStore } from '@/stores';
import { ElMessageBox, ElMessage } from 'element-plus';

const router = useRouter();
const shopStore = useShopStore();
const npcStore = useNPCStore();

const getNPCName = (npcId) => {
  const npc = npcStore.getNPCById(npcId);
  return npc?.name || '-';
};

const handleAdd = () => {
  router.push('/shop/edit');
};

const handleEdit = (shop) => {
  router.push(`/shop/edit/${shop.id}`);
};

const handleDelete = async (shop) => {
  try {
    await ElMessageBox.confirm('确定要删除该商店吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    shopStore.deleteShop(shop.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};
</script>

<style scoped>
.shop-list {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
