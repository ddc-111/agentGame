<template>
  <div class="item-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>道具管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建道具
          </el-button>
        </div>
      </template>

      <el-table :data="shopStore.items" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column label="图标" width="80">
          <template #default="{ row }">
            <el-avatar :size="40" :src="row.icon">
              {{ row.name?.charAt(0) }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="120" />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag :type="getCategoryTag(row.category)">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column label="效果" width="200">
          <template #default="{ row }">
            <span v-if="row.effect">
              <span v-for="(value, key) in row.effect" :key="key">
                {{ key }}: {{ value }}
              </span>
            </span>
            <span v-else>-</span>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑道具' : '新建道具'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="道具ID">
          <el-input v-model="form.id" disabled />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="请输入道具名称" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category">
            <el-option label="药品" value="medicine" />
            <el-option label="食物" value="food" />
            <el-option label="武器" value="weapon" />
            <el-option label="防具" value="armor" />
            <el-option label="工具" value="tool" />
            <el-option label="材料" value="material" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入道具描述" />
        </el-form-item>
        <el-form-item label="图标">
          <el-upload
            class="icon-uploader"
            action="#"
            :auto-upload="false"
            :on-change="handleIconChange"
          >
            <img v-if="form.icon" :src="form.icon" class="icon-preview" />
            <el-icon v-else class="icon-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>

        <el-divider content-position="left">效果配置</el-divider>

        <el-form-item label="效果类型">
          <el-select v-model="effectType" placeholder="选择效果类型">
            <el-option label="恢复生命" value="hp" />
            <el-option label="恢复法力" value="mp" />
            <el-option label="恢复体力" value="stamina" />
            <el-option label="增加攻击" value="attack" />
            <el-option label="增加防御" value="defense" />
            <el-option label="增加抗性" value="resist" />
          </el-select>
        </el-form-item>
        <el-form-item label="效果值">
          <el-input-number v-model="effectValue" :min="0" />
          <el-button type="primary" @click="addEffect" style="margin-left: 10px">添加</el-button>
        </el-form-item>

        <el-table :data="Object.entries(form.effect || {})" style="margin-bottom: 20px;">
          <el-table-column prop="0" label="效果类型" width="150" />
          <el-table-column prop="1" label="数值" width="100" />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button type="danger" text @click="removeEffect(row[0])">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useShopStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const shopStore = useShopStore();

const dialogVisible = ref(false);
const isEdit = ref(false);
const effectType = ref('');
const effectValue = ref(0);

const form = ref({
  id: '',
  name: '',
  category: 'medicine',
  description: '',
  icon: '',
  effect: {}
});

const getCategoryTag = (category) => {
  const map = {
    medicine: 'success',
    food: 'warning',
    weapon: 'danger',
    armor: 'primary',
    tool: 'info',
    material: '',
    other: ''
  };
  return map[category] || '';
};

const handleAdd = () => {
  isEdit.value = false;
  form.value = {
    id: `item_${Date.now()}`,
    name: '',
    category: 'medicine',
    description: '',
    icon: '',
    effect: {}
  };
  dialogVisible.value = true;
};

const handleEdit = (item) => {
  isEdit.value = true;
  form.value = { ...item };
  dialogVisible.value = true;
};

const handleDelete = async (item) => {
  try {
    await ElMessageBox.confirm('确定要删除该道具吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    shopStore.deleteItem(item.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};

const handleIconChange = (file) => {
  const reader = new FileReader();
  reader.onload = (e) => {
    form.value.icon = e.target.result;
  };
  reader.readAsDataURL(file.raw);
};

const addEffect = () => {
  if (effectType.value && effectValue.value > 0) {
    if (!form.value.effect) {
      form.value.effect = {};
    }
    form.value.effect[effectType.value] = effectValue.value;
    effectType.value = '';
    effectValue.value = 0;
  }
};

const removeEffect = (key) => {
  delete form.value.effect[key];
};

const handleSave = () => {
  if (!form.value.name) {
    ElMessage.warning('请输入道具名称');
    return;
  }

  if (isEdit.value) {
    shopStore.updateItem(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    shopStore.addItem(form.value);
    ElMessage.success('创建成功');
  }
  dialogVisible.value = false;
};
</script>

<style scoped>
.item-manager {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.icon-uploader {
  width: 100px;
  height: 100px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.icon-uploader:hover {
  border-color: #409eff;
}

.icon-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
