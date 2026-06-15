<template>
  <div class="variable-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>变量管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加变量
          </el-button>
        </div>
      </template>

      <el-table :data="promptStore.variables" style="width: 100%">
        <el-table-column prop="name" label="变量名" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="source" label="来源" width="120">
          <template #default="{ row }">
            <el-tag :type="getSourceTag(row.source)">{{ row.source }}</el-tag>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑变量' : '添加变量'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="变量名" required>
          <el-input v-model="form.name" placeholder="例如: player_name" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="form.type">
            <el-option label="字符串" value="string" />
            <el-option label="数字" value="number" />
            <el-option label="布尔" value="boolean" />
            <el-option label="数组" value="array" />
            <el-option label="对象" value="object" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" required>
          <el-input v-model="form.description" placeholder="输入变量描述" />
        </el-form-item>
        <el-form-item label="来源">
          <el-select v-model="form.source" placeholder="选择变量来源">
            <el-option label="游戏系统" value="game" />
            <el-option label="玩家信息" value="player" />
            <el-option label="NPC信息" value="npc" />
            <el-option label="智能体" value="agent" />
            <el-option label="商店" value="shop" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
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
import { usePromptStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const promptStore = usePromptStore();

const dialogVisible = ref(false);
const isEdit = ref(false);

const form = ref({
  name: '',
  type: 'string',
  description: '',
  source: 'custom'
});

const getSourceTag = (source) => {
  const map = {
    game: 'primary',
    player: 'success',
    npc: 'warning',
    agent: 'danger',
    shop: 'info',
    custom: ''
  };
  return map[source] || '';
};

const handleAdd = () => {
  isEdit.value = false;
  form.value = {
    name: '',
    type: 'string',
    description: '',
    source: 'custom'
  };
  dialogVisible.value = true;
};

const handleEdit = (variable) => {
  isEdit.value = true;
  form.value = { ...variable };
  dialogVisible.value = true;
};

const handleDelete = async (variable) => {
  try {
    await ElMessageBox.confirm('确定要删除该变量吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    promptStore.deleteVariable(variable.name);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};

const handleSave = () => {
  if (!form.value.name || !form.value.description) {
    ElMessage.warning('请填写必要信息');
    return;
  }

  if (isEdit.value) {
    promptStore.updateVariable(form.value.name, form.value);
    ElMessage.success('更新成功');
  } else {
    promptStore.addVariable(form.value);
    ElMessage.success('添加成功');
  }
  dialogVisible.value = false;
};
</script>

<style scoped>
.variable-manager {
  max-width: 1000px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
