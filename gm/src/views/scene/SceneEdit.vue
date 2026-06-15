<template>
  <div class="scene-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑场景' : '新建场景' }}</span>
          <div>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="场景ID">
              <el-input v-model="form.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="场景名称" required>
              <el-input v-model="form.name" placeholder="请输入场景名称" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="场景描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入场景描述" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="背景图片">
              <el-upload
                class="background-uploader"
                action="#"
                :auto-upload="false"
                :on-change="handleBackgroundChange"
              >
                <img v-if="form.background" :src="form.background" class="background-preview" />
                <el-icon v-else class="background-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="宽度">
              <el-input-number v-model="form.width" :min="800" :max="3840" :step="100" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="高度">
              <el-input-number v-model="form.height" :min="600" :max="2160" :step="100" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">场景NPC</el-form-item>

        <el-form-item>
          <el-button @click="handleAddNPC">添加NPC</el-button>
          <el-table :data="form.npcs" style="margin-top: 10px;">
            <el-table-column prop="npcId" label="NPC" width="200">
              <template #default="{ row }">
                <el-select v-model="row.npcId" placeholder="选择NPC">
                  <el-option
                    v-for="npc in npcStore.npcs"
                    :key="npc.id"
                    :label="npc.name"
                    :value="npc.id"
                  />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="位置" width="300">
              <template #default="{ row }">
                <el-input-number v-model="row.x" :min="0" :max="form.width" placeholder="X" style="width: 100px" />
                <span style="margin: 0 8px">,</span>
                <el-input-number v-model="row.y" :min="0" :max="form.height" placeholder="Y" style="width: 100px" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemoveNPC($index)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>

        <el-divider content-position="left">传送点</el-divider>

        <el-form-item>
          <el-button @click="handleAddPortal">添加传送点</el-button>
          <el-table :data="form.portals" style="margin-top: 10px;">
            <el-table-column label="位置" width="200">
              <template #default="{ row }">
                <el-input-number v-model="row.x" :min="0" :max="form.width" placeholder="X" style="width: 80px" />
                <span style="margin: 0 4px">,</span>
                <el-input-number v-model="row.y" :min="0" :max="form.height" placeholder="Y" style="width: 80px" />
              </template>
            </el-table-column>
            <el-table-column label="目标场景" width="200">
              <template #default="{ row }">
                <el-select v-model="row.targetScene" placeholder="选择场景">
                  <el-option
                    v-for="scene in sceneStore.scenes"
                    :key="scene.id"
                    :label="scene.name"
                    :value="scene.id"
                  />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="目标位置" width="200">
              <template #default="{ row }">
                <el-input-number v-model="row.targetX" :min="0" placeholder="X" style="width: 80px" />
                <span style="margin: 0 4px">,</span>
                <el-input-number v-model="row.targetY" :min="0" placeholder="Y" style="width: 80px" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ $index }">
                <el-button type="danger" text @click="handleRemovePortal($index)">删除</el-button>
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
import { useSceneStore, useNPCStore } from '@/stores';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();
const sceneStore = useSceneStore();
const npcStore = useNPCStore();

const isEdit = computed(() => !!route.params.id);

const form = ref({
  id: '',
  name: '',
  description: '',
  background: '',
  width: 1920,
  height: 1080,
  npcs: [],
  portals: []
});

onMounted(() => {
  if (route.params.id) {
    const scene = sceneStore.getSceneById(route.params.id);
    if (scene) {
      form.value = { ...scene };
    }
  } else {
    form.value.id = `scene_${Date.now()}`;
  }
});

const handleBackgroundChange = (file) => {
  const reader = new FileReader();
  reader.onload = (e) => {
    form.value.background = e.target.result;
  };
  reader.readAsDataURL(file.raw);
};

const handleAddNPC = () => {
  form.value.npcs.push({ npcId: '', x: 0, y: 0 });
};

const handleRemoveNPC = (index) => {
  form.value.npcs.splice(index, 1);
};

const handleAddPortal = () => {
  form.value.portals.push({ x: 0, y: 0, targetScene: '', targetX: 0, targetY: 0 });
};

const handleRemovePortal = (index) => {
  form.value.portals.splice(index, 1);
};

const handleSave = () => {
  if (!form.value.name) {
    ElMessage.warning('请输入场景名称');
    return;
  }

  if (isEdit.value) {
    sceneStore.updateScene(form.value.id, form.value);
    ElMessage.success('更新成功');
  } else {
    sceneStore.addScene(form.value);
    ElMessage.success('创建成功');
  }
  router.push('/scene/list');
};

const handleCancel = () => {
  router.push('/scene/list');
};
</script>

<style scoped>
.scene-edit {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.background-uploader {
  width: 200px;
  height: 120px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.background-uploader:hover {
  border-color: #409eff;
}

.background-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 200px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.background-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
