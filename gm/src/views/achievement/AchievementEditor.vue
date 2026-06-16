<template>
  <div class="achievement-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>成就管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新建成就
          </el-button>
        </div>
      </template>

      <el-table :data="achievementStore.achievements" style="width: 100%">
        <el-table-column prop="icon" label="图标" width="60">
          <template #default="{ row }">
            <span class="achievement-icon">{{ row.icon }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="120" />
        <el-table-column prop="code" label="编码" width="160" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="条件" width="200">
          <template #default="{ row }">
            <el-tag>{{ getConditionLabel(row.condition?.type) }}</el-tag>
            <span class="condition-value">{{ row.condition?.value }}</span>
          </template>
        </el-table-column>
        <el-table-column label="奖励" width="180">
          <template #default="{ row }">
            <span v-if="row.reward">经验:{{ row.reward.exp }} 金币:{{ row.reward.gold }}</span>
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

    <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑成就' : '新建成就'" width="600px">
      <el-form :model="achievementForm" label-width="100px">
        <el-form-item label="成就名称">
          <el-input v-model="achievementForm.name" />
        </el-form-item>
        <el-form-item label="成就编码">
          <el-input v-model="achievementForm.code" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="achievementForm.icon" style="width: 80px" />
          <span class="icon-preview">{{ achievementForm.icon }}</span>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="achievementForm.description" type="textarea" :rows="2" />
        </el-form-item>

        <el-divider content-position="left">达成条件</el-divider>
        <el-form-item label="条件类型">
          <el-select v-model="achievementForm.condition.type" style="width: 100%">
            <el-option
              v-for="ct in achievementStore.conditionTypes"
              :key="ct.value"
              :label="ct.label"
              :value="ct.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="目标值">
          <el-input-number v-model="achievementForm.condition.value" :min="1" :max="999999" />
        </el-form-item>

        <el-divider content-position="left">奖励</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="经验奖励">
              <el-input-number v-model="achievementForm.reward.exp" :min="0" :max="99999" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="金币奖励">
              <el-input-number v-model="achievementForm.reward.gold" :min="0" :max="99999" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAchievement">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useAchievementStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const achievementStore = useAchievementStore();
const dialogVisible = ref(false);
const isEditing = ref(false);
const editingId = ref(null);

const achievementForm = reactive({
  name: '',
  code: '',
  icon: '⭐',
  description: '',
  condition: { type: 'quest_complete', value: 1 },
  reward: { exp: 50, gold: 100 }
});

const getConditionLabel = (type) => {
  const ct = achievementStore.conditionTypes.find(c => c.value === type);
  return ct ? ct.label : type || '-';
};

const handleAdd = () => {
  isEditing.value = false;
  editingId.value = null;
  Object.assign(achievementForm, {
    name: '', code: '', icon: '⭐', description: '',
    condition: { type: 'quest_complete', value: 1 },
    reward: { exp: 50, gold: 100 }
  });
  dialogVisible.value = true;
};

const handleEdit = (ach) => {
  isEditing.value = true;
  editingId.value = ach.id;
  Object.assign(achievementForm, {
    name: ach.name,
    code: ach.code,
    icon: ach.icon,
    description: ach.description,
    condition: { ...(ach.condition || { type: 'quest_complete', value: 1 }) },
    reward: { ...(ach.reward || { exp: 50, gold: 100 }) }
  });
  dialogVisible.value = true;
};

const handleDelete = async (ach) => {
  try {
    await ElMessageBox.confirm('确定要删除该成就吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    achievementStore.deleteAchievement(ach.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消
  }
};

const saveAchievement = () => {
  if (isEditing.value) {
    achievementStore.updateAchievement(editingId.value, { ...achievementForm });
  } else {
    achievementStore.addAchievement({ ...achievementForm });
  }
  dialogVisible.value = false;
  ElMessage.success('保存成功');
};
</script>

<style scoped>
.achievement-editor {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.achievement-icon {
  font-size: 24px;
}

.condition-value {
  margin-left: 8px;
  color: #409eff;
  font-weight: 500;
}

.icon-preview {
  margin-left: 16px;
  font-size: 24px;
}
</style>
