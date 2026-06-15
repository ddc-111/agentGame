<template>
  <div class="task-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑任务' : '新建任务' }}</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="任务ID">
              <el-input v-model="form.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="任务名称" required>
              <el-input v-model="form.name" placeholder="请输入任务名称" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="任务类型">
              <el-select v-model="form.type">
                <el-option label="主线任务" value="main" />
                <el-option label="支线任务" value="side" />
                <el-option label="日常任务" value="daily" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="任务描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入任务描述" />
        </el-form-item>

        <el-divider content-position="left">触发条件</el-divider>

        <el-form-item label="触发类型">
          <el-select v-model="form.trigger.type">
            <el-option label="自动触发" value="auto" />
            <el-option label="NPC对话触发" value="dialogue" />
            <el-option label="任务完成触发" value="task_complete" />
            <el-option label="物品获得触发" value="item_collect" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发条件">
          <div v-for="(condition, index) in form.trigger.conditions" :key="index" class="condition-item">
            <el-select v-model="condition.type" placeholder="条件类型" style="width: 150px">
              <el-option label="玩家等级" value="player_level" />
              <el-option label="任务完成" value="task_id" />
              <el-option label="拥有物品" value="has_item" />
            </el-select>
            <el-select v-model="condition.operator" placeholder="运算符" style="width: 100px; margin-left: 10px">
              <el-option label="==" value="==" />
              <el-option label=">=" value=">=" />
              <el-option label="<=" value="<=" />
              <el-option label=">" value=">" />
              <el-option label="<" value="<" />
            </el-select>
            <el-input v-model="condition.value" placeholder="值" style="width: 150px; margin-left: 10px" />
            <el-button type="danger" text @click="removeCondition(index)" style="margin-left: 10px">删除</el-button>
          </div>
          <el-button @click="addCondition" style="margin-top: 10px">添加条件</el-button>
        </el-form-item>

        <el-divider content-position="left">任务目标</el-divider>

        <el-form-item>
          <el-button @click="handleAddObjective">添加目标</el-button>
          <el-table :data="form.objectives" style="margin-top: 10px;">
            <el-table-column prop="id" label="ID" width="120" />
            <el-table-column label="类型" width="120">
              <template #default="{ row }">
                <el-select v-model="row.type">
                  <el-option label="对话" value="dialogue" />
                  <el-option label="收集" value="collect" />
                  <el-option label="击杀" value="kill" />
                  <el-option label="到达" value="reach" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="目标" width="200">
              <template #default="{ row }">
                <el-input v-model="row.target" placeholder="目标ID" />
              </template>
            </el-table-column>
            <el-table-column label="数量" width="100">
              <template #default="{ row }">
                <el-input-number v-model="row.count" :min="1" />
              </template>
            </el-table-column>
            <el-table-column label="描述">
              <template #default="{ row }">
                <el-input v-model="row.description" placeholder="目标描述" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemoveObjective($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>

        <el-divider content-position="left">任务奖励</el-divider>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="经验奖励">
              <el-input-number v-model="form.rewards.exp" :min="0" :step="100" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="金币奖励">
              <el-input-number v-model="form.rewards.gold" :min="0" :step="100" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="下一任务">
              <el-select v-model="form.nextTask" placeholder="选择下一任务" clearable>
                <el-option
                  v-for="task in taskStore.tasks"
                  :key="task.id"
                  :label="task.name"
                  :value="task.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useTaskStore } from '@/stores';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();
const taskStore = useTaskStore();

const isEdit = computed(() => !!route.params.id);

const form = ref({
  id: '',
  name: '',
  type: 'main',
  description: '',
  status: 'active',
  trigger: {
    type: 'auto',
    conditions: []
  },
  objectives: [],
  rewards: {
    exp: 0,
    gold: 0,
    items: []
  },
  nextTask: '',
  dialogue: ''
});

onMounted(() => {
  if (route.params.id) {
    const task = taskStore.getTaskById(route.params.id);
    if (task) {
      form.value = { ...task };
    }
  } else {
    form.value.id = `task_${Date.now()}`;
  }
});

const addCondition = () => {
  form.value.trigger.conditions.push({ type: '', operator: '==', value: '' });
};

const removeCondition = (index) => {
  form.value.trigger.conditions.splice(index, 1);
};

const handleAddObjective = () => {
  form.value.objectives.push({
    id: `obj_${Date.now()}`,
    type: 'dialogue',
    target: '',
    count: 1,
    description: '',
    completed: false
  });
};

const handleRemoveObjective = (index) => {
  form.value.objectives.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.name) {
    ElMessage.warning('请输入任务名称');
    return;
  }

  if (isEdit.value) {
    taskStore.updateTask(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    taskStore.addTask(form.value);
    ElMessage.success('创建成功');
  }
  router.push('/task/list');
};

const handleCancel = () => {
  router.push('/task/list');
};
</script>

<style scoped>
.task-edit {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.condition-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
</style>
