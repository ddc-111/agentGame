<template>
  <div class="provider-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>模型提供商配置</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加提供商
          </el-button>
        </div>
      </template>

      <el-table :data="llmStore.providers" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="baseUrl" label="API地址" />
        <el-table-column label="API Key" width="200">
          <template #default="{ row }">
            <span v-if="row.apiKey">{{ '••••••••' + row.apiKey.slice(-4) }}</span>
            <span v-else style="color: #f56c6c">未配置</span>
          </template>
        </el-table-column>
        <el-table-column label="模型数量" width="100">
          <template #default="{ row }">
            {{ row.models?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" text @click="handleTest(row)">测试</el-button>
            <el-button type="danger" text @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑提供商' : '添加提供商'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="提供商ID" required>
          <el-input v-model="form.id" placeholder="例如: openai" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="例如: OpenAI" />
        </el-form-item>
        <el-form-item label="API地址" required>
          <el-input v-model="form.baseUrl" placeholder="例如: https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.apiKey" type="password" show-password placeholder="输入API Key" />
        </el-form-item>

        <el-divider content-position="left">模型列表</el-divider>

        <el-form-item>
          <el-button @click="handleAddModel">添加模型</el-button>
          <div v-for="(model, index) in form.models" :key="index" class="model-item">
            <el-input v-model="model.id" placeholder="模型ID" style="width: 200px" />
            <el-input v-model="model.name" placeholder="模型名称" style="width: 200px; margin-left: 10px" />
            <el-input-number v-model="model.maxTokens" :min="100" :max="200000" placeholder="最大Token" style="width: 150px; margin-left: 10px" />
            <el-button type="danger" text @click="handleRemoveModel(index)" style="margin-left: 10px">删除</el-button>
          </div>
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
import { useLLMStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const llmStore = useLLMStore();

const dialogVisible = ref(false);
const isEdit = ref(false);

const form = ref({
  id: '',
  name: '',
  baseUrl: '',
  apiKey: '',
  models: []
});

const handleAdd = () => {
  isEdit.value = false;
  form.value = {
    id: '',
    name: '',
    baseUrl: '',
    apiKey: '',
    models: []
  };
  dialogVisible.value = true;
};

const handleEdit = (provider) => {
  isEdit.value = true;
  form.value = { ...provider };
  dialogVisible.value = true;
};

const handleTest = async (provider) => {
  ElMessage.info(`正在测试连接: ${provider.name}`);
  // TODO: 实际测试连接
  setTimeout(() => {
    ElMessage.success('连接成功');
  }, 1000);
};

const handleDelete = async (provider) => {
  try {
    await ElMessageBox.confirm('确定要删除该提供商吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    llmStore.deleteProvider(provider.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};

const handleAddModel = () => {
  form.value.models.push({ id: '', name: '', maxTokens: 4096 });
};

const handleRemoveModel = (index) => {
  form.value.models.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.id || !form.value.name) {
    ElMessage.warning('请填写必要信息');
    return;
  }

  if (isEdit.value) {
    llmStore.updateProvider(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    llmStore.addProvider(form.value);
    ElMessage.success('添加成功');
  }
  dialogVisible.value = false;
};
</script>

<style scoped>
.provider-config {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.model-item {
  display: flex;
  align-items: center;
  margin-top: 10px;
}
</style>
