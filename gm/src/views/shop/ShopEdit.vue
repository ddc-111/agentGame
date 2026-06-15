<template>
  <div class="shop-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑商店' : '新建商店' }}</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="商店ID">
              <el-input v-model="form.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="商店名称" required>
              <el-input v-model="form.name" placeholder="请输入商店名称" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="商店类型">
              <el-select v-model="form.type">
                <el-option label="杂货铺" value="general" />
                <el-option label="铁匠铺" value="blacksmith" />
                <el-option label="药店" value="pharmacy" />
                <el-option label="酒楼" value="restaurant" />
                <el-option label="布庄" value="cloth" />
                <el-option label="书店" value="bookstore" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入商店描述" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="店主NPC">
              <el-select v-model="form.owner" placeholder="选择NPC">
                <el-option
                  v-for="npc in npcStore.npcs"
                  :key="npc.id"
                  :label="npc.name"
                  :value="npc.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="所在场景">
              <el-select v-model="form.scene" placeholder="选择场景">
                <el-option
                  v-for="scene in sceneStore.scenes"
                  :key="scene.id"
                  :label="scene.name"
                  :value="scene.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="开门时间">
              <el-time-picker v-model="form.openTime" format="HH:mm" value-format="HH:mm" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="关门时间">
              <el-time-picker v-model="form.closeTime" format="HH:mm" value-format="HH:mm" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">折扣配置</el-divider>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="折扣门槛">
              <el-input-number v-model="form.discount.threshold" :min="1" :max="100" />
              <span class="form-tip">件</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="折扣率">
              <el-slider v-model="form.discount.rate" :min="0.1" :max="1" :step="0.05" show-input />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">商品配置</el-divider>

        <el-form-item>
          <el-button @click="handleAddItem">添加商品</el-button>
          <el-table :data="form.items" style="margin-top: 10px;">
            <el-table-column label="商品" width="250">
              <template #default="{ row }">
                <el-select v-model="row.itemId" placeholder="选择商品">
                  <el-option
                    v-for="item in shopStore.items"
                    :key="item.id"
                    :label="item.name"
                    :value="item.id"
                  />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="价格" width="150">
              <template #default="{ row }">
                <el-input-number v-model="row.price" :min="0" :step="10" />
              </template>
            </el-table-column>
            <el-table-column label="库存" width="150">
              <template #default="{ row }">
                <el-input-number v-model="row.stock" :min="0" :step="10" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemoveItem($index)">删除</el-button>
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
import { useShopStore, useNPCStore, useSceneStore } from '@/stores';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();
const shopStore = useShopStore();
const npcStore = useNPCStore();
const sceneStore = useSceneStore();

const isEdit = computed(() => !!route.params.id);

const form = ref({
  id: '',
  name: '',
  type: 'general',
  description: '',
  owner: '',
  scene: '',
  items: [],
  openTime: '06:00',
  closeTime: '22:00',
  discount: {
    threshold: 3,
    rate: 0.9
  }
});

onMounted(() => {
  if (route.params.id) {
    const shop = shopStore.getShopById(route.params.id);
    if (shop) {
      form.value = { ...shop };
    }
  } else {
    form.value.id = `shop_${Date.now()}`;
  }
});

const handleAddItem = () => {
  form.value.items.push({ itemId: '', price: 0, stock: 0 });
};

const handleRemoveItem = (index) => {
  form.value.items.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.name) {
    ElMessage.warning('请输入商店名称');
    return;
  }

  if (isEdit.value) {
    shopStore.updateShop(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    shopStore.addShop(form.value);
    ElMessage.success('创建成功');
  }
  router.push('/shop/list');
};

const handleCancel = () => {
  router.push('/shop/list');
};
</script>

<style scoped>
.shop-edit {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-left: 8px;
  color: #909399;
}
</style>
