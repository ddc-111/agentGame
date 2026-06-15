<template>
  <div class="npc-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑NPC' : '新建NPC' }}</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="NPC ID">
              <el-input v-model="form.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="名称" required>
              <el-input v-model="form.name" placeholder="请输入NPC名称" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="称号">
              <el-input v-model="form.title" placeholder="请输入NPC称号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入NPC描述" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="头像">
              <el-upload
                class="avatar-uploader"
                action="#"
                :auto-upload="false"
                :on-change="handleAvatarChange"
              >
                <img v-if="form.avatar" :src="form.avatar" class="avatar-preview" />
                <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="精灵图">
              <el-upload
                class="sprite-uploader"
                action="#"
                :auto-upload="false"
                :on-change="handleSpriteChange"
              >
                <img v-if="form.sprite" :src="form.sprite" class="sprite-preview" />
                <el-icon v-else class="sprite-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="所属场景">
              <el-select v-model="form.position.scene" placeholder="选择场景">
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
            <el-form-item label="初始位置X">
              <el-input-number v-model="form.position.x" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="初始位置Y">
              <el-input-number v-model="form.position.y" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="关联智能体">
              <el-select v-model="form.agentId" placeholder="选择智能体" clearable>
                <el-option
                  v-for="agent in agentStore.agents"
                  :key="agent.id"
                  :label="agent.name"
                  :value="agent.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="关联商店">
              <el-select v-model="form.shopId" placeholder="选择商店" clearable>
                <el-option
                  v-for="shop in shopStore.shops"
                  :key="shop.id"
                  :label="shop.name"
                  :value="shop.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label="行为模式">
              <el-select v-model="form.behaviors" multiple placeholder="选择行为模式">
                <el-option label="空闲" value="idle" />
                <el-option label="打招呼" value="greet" />
                <el-option label="售卖" value="sell" />
                <el-option label="聊天" value="chat" />
                <el-option label="巡逻" value="patrol" />
                <el-option label="锻造" value="forge" />
                <el-option label="休息" value="rest" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">日程安排</el-divider>

        <el-form-item>
          <el-button @click="handleAddSchedule">添加日程</el-button>
          <el-table :data="form.schedule" style="margin-top: 10px;">
            <el-table-column label="时间" width="150">
              <template #default="{ row }">
                <el-time-picker v-model="row.time" format="HH:mm" value-format="HH:mm" />
              </template>
            </el-table-column>
            <el-table-column label="动作" width="150">
              <template #default="{ row }">
                <el-select v-model="row.action" placeholder="选择动作">
                  <el-option label="开店" value="open_shop" />
                  <el-option label="关店" value="close_shop" />
                  <el-option label="回家" value="go_home" />
                  <el-option label="巡逻" value="patrol" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="目标位置" width="300">
              <template #default="{ row }">
                <el-select v-model="row.position.scene" placeholder="场景" style="width: 120px">
                  <el-option
                    v-for="scene in sceneStore.scenes"
                    :key="scene.id"
                    :label="scene.name"
                    :value="scene.id"
                  />
                </el-select>
                <el-input-number v-model="row.position.x" :min="0" placeholder="X" style="width: 80px; margin-left: 8px" />
                <el-input-number v-model="row.position.y" :min="0" placeholder="Y" style="width: 80px; margin-left: 8px" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemoveSchedule($index)">删除</el-button>
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
import { useNPCStore, useSceneStore, useAgentStore, useShopStore } from '@/stores';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();
const npcStore = useNPCStore();
const sceneStore = useSceneStore();
const agentStore = useAgentStore();
const shopStore = useShopStore();

const isEdit = computed(() => !!route.params.id);

const form = ref({
  id: '',
  name: '',
  title: '',
  description: '',
  avatar: '',
  sprite: '',
  position: { x: 0, y: 0, scene: '' },
  agentId: '',
  shopId: '',
  dialogues: [],
  behaviors: [],
  schedule: []
});

onMounted(() => {
  if (route.params.id) {
    const npc = npcStore.getNPCById(route.params.id);
    if (npc) {
      form.value = { ...npc };
    }
  } else {
    form.value.id = `npc_${Date.now()}`;
  }
});

const handleAvatarChange = (file) => {
  const reader = new FileReader();
  reader.onload = (e) => {
    form.value.avatar = e.target.result;
  };
  reader.readAsDataURL(file.raw);
};

const handleSpriteChange = (file) => {
  const reader = new FileReader();
  reader.onload = (e) => {
    form.value.sprite = e.target.result;
  };
  reader.readAsDataURL(file.raw);
};

const handleAddSchedule = () => {
  form.value.schedule.push({
    time: '08:00',
    action: 'open_shop',
    position: { x: 0, y: 0, scene: '' }
  });
};

const handleRemoveSchedule = (index) => {
  form.value.schedule.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.name) {
    ElMessage.warning('请输入NPC名称');
    return;
  }

  if (isEdit.value) {
    npcStore.updateNPC(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    npcStore.addNPC(form.value);
    ElMessage.success('创建成功');
  }
  router.push('/npc/list');
};

const handleCancel = () => {
  router.push('/npc/list');
};
</script>

<style scoped>
.npc-edit {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.avatar-uploader,
.sprite-uploader {
  width: 120px;
  height: 120px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.avatar-uploader:hover,
.sprite-uploader:hover {
  border-color: #409eff;
}

.avatar-uploader-icon,
.sprite-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-preview,
.sprite-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
