<template>
  <div class="prompt-template">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>提示词模板</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新建模板
          </el-button>
        </div>
      </template>

      <el-table :data="promptStore.templates" style="width: 100%">
        <el-table-column prop="id" label="ID" width="120" />
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryTag(row.category)">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容">
          <template #default="{ row }">
            {{ row.content?.substring(0, 50) }}...
          </template>
        </el-table-column>
        <el-table-column label="变量" width="100">
          <template #default="{ row }">
            {{ row.variables?.length || 0 }} 个
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑模板' : '新建模板'" width="800px">
      <el-form :model="form" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模板ID">
              <el-input v-model="form.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类">
              <el-select v-model="form.category">
                <el-option label="系统提示" value="system" />
                <el-option label="对话模板" value="dialogue" />
                <el-option label="摘要模板" value="summary" />
                <el-option label="任务模板" value="task" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="模板名称" required>
          <el-input v-model="form.name" placeholder="请输入模板名称" />
        </el-form-item>

        <el-form-item label="模板内容" required>
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="12"
            placeholder="输入模板内容，使用 {{变量名}} 插入变量"
          />
        </el-form-item>

        <el-divider content-position="left">变量定义</el-divider>

        <el-form-item>
          <el-button @click="handleAddVariable">添加变量</el-button>
          <div v-for="(variable, index) in form.variables" :key="index" class="variable-item">
            <el-input v-model="variable.name" placeholder="变量名" style="width: 150px" />
            <el-select v-model="variable.type" placeholder="类型" style="width: 120px; margin-left: 10px">
              <el-option label="字符串" value="string" />
              <el-option label="数字" value="number" />
              <el-option label="布尔" value="boolean" />
              <el-option label="数组" value="array" />
            </el-select>
            <el-input v-model="variable.description" placeholder="描述" style="width: 200px; margin-left: 10px" />
            <el-button type="danger" text @click="handleRemoveVariable(index)" style="margin-left: 10px">删除</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="previewVisible" title="模板预览" width="600px">
      <div class="preview-content">
        <pre>{{ previewContent }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { usePromptStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const promptStore = usePromptStore();

const dialogVisible = ref(false);
const previewVisible = ref(false);
const isEdit = ref(false);
const previewContent = ref('');

const form = ref({
  id: '',
  name: '',
  category: 'system',
  content: '',
  variables: []
});

const getCategoryTag = (category) => {
  const map = {
    system: 'primary',
    dialogue: 'success',
    summary: 'warning',
    task: 'danger'
  };
  return map[category] || '';
};

const handleAdd = () => {
  isEdit.value = false;
  form.value = {
    id: `template_${Date.now()}`,
    name: '',
    category: 'system',
    content: '',
    variables: []
  };
  dialogVisible.value = true;
};

const handleEdit = (template) => {
  isEdit.value = true;
  form.value = { ...template };
  dialogVisible.value = true;
};

const handlePreview = (template) => {
  previewContent.value = template.content;
  previewVisible.value = true;
};

const handleDelete = async (template) => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    promptStore.deleteTemplate(template.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消操作
  }
};

const handleAddVariable = () => {
  form.value.variables.push({ name: '', type: 'string', description: '' });
};

const handleRemoveVariable = (index) => {
  form.value.variables.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.name || !form.value.content) {
    ElMessage.warning('请填写必要信息');
    return;
  }

  if (isEdit.value) {
    promptStore.updateTemplate(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    promptStore.addTemplate(form.value);
    ElMessage.success('创建成功');
  }
  dialogVisible.value = false;
};
</script>

<style scoped>
.prompt-template {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.variable-item {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.preview-content {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
}

.preview-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
