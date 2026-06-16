<template>
  <div class="skill-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>技能管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新建技能
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="技能列表" name="list">
          <el-table :data="skillStore.skills" style="width: 100%">
            <el-table-column prop="name" label="名称" width="120" />
            <el-table-column prop="code" label="编码" width="160" />
            <el-table-column label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="getTypeTag(row.type)">{{ getTypeLabel(row.type) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="mpCost" label="MP消耗" width="80" />
            <el-table-column prop="damage" label="伤害" width="80" />
            <el-table-column prop="heal" label="治疗" width="80" />
            <el-table-column prop="cooldown" label="冷却" width="80" />
            <el-table-column prop="level" label="等级" width="80" />
            <el-table-column prop="description" label="描述" />
            <el-table-column label="效果" width="150">
              <template #default="{ row }">
                <span v-if="row.effect && row.effect.type">{{ row.effect.type }} {{ row.effect.duration ? row.effect.duration + '回合' : '' }}</span>
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
        </el-tab-pane>

        <el-tab-pane label="技能树" name="tree">
          <div class="skill-tree-container">
            <div class="tree-toolbar">
              <el-button type="primary" size="small" @click="resetTree">重置布局</el-button>
              <el-button type="success" size="small" @click="saveTree">保存技能树</el-button>
            </div>
            <div class="skill-tree-canvas" ref="treeCanvas">
              <svg class="tree-edges">
                <line
                  v-for="(edge, i) in skillStore.skillTree.edges"
                  :key="i"
                  :x1="getNodePos(edge.source).x"
                  :y1="getNodePos(edge.source).y"
                  :x2="getNodePos(edge.target).x"
                  :y2="getNodePos(edge.target).y"
                  stroke="#409eff"
                  stroke-width="2"
                  marker-end="url(#arrowhead)"
                />
                <defs>
                  <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="10" refY="3.5" orient="auto">
                    <polygon points="0 0, 10 3.5, 0 7" fill="#409eff" />
                  </marker>
                </defs>
              </svg>
              <div
                v-for="node in skillStore.skillTree.nodes"
                :key="node.id"
                class="tree-node"
                :class="node.type"
                :style="{ left: node.x + 'px', top: node.y + 'px' }"
                @click="selectNode(node)"
              >
                <div class="node-icon">{{ node.type === 'root' ? '🌳' : '⚔️' }}</div>
                <div class="node-label">{{ node.label }}</div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="效果预览" name="preview">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-card>
                <template #header>选择技能</template>
                <el-select v-model="previewSkill" placeholder="选择技能" style="width: 100%">
                  <el-option
                    v-for="skill in skillStore.skills"
                    :key="skill.id"
                    :label="skill.name"
                    :value="skill.id"
                  />
                </el-select>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card v-if="previewSkill">
                <template #header>技能详情</template>
                <div class="skill-detail" v-if="getPreviewSkill">
                  <h3>{{ getPreviewSkill.name }}</h3>
                  <p class="skill-desc">{{ getPreviewSkill.description }}</p>
                  <div class="skill-stats">
                    <div class="stat-item">
                      <span class="stat-label">类型:</span>
                      <el-tag :type="getTypeTag(getPreviewSkill.type)">{{ getTypeLabel(getPreviewSkill.type) }}</el-tag>
                    </div>
                    <div class="stat-item">
                      <span class="stat-label">MP消耗:</span>
                      <span class="stat-value">{{ getPreviewSkill.mpCost }}</span>
                    </div>
                    <div class="stat-item">
                      <span class="stat-label">伤害:</span>
                      <span class="stat-value damage">{{ getPreviewSkill.damage }}</span>
                    </div>
                    <div class="stat-item">
                      <span class="stat-label">治疗:</span>
                      <span class="stat-value heal">{{ getPreviewSkill.heal }}</span>
                    </div>
                    <div class="stat-item">
                      <span class="stat-label">冷却:</span>
                      <span class="stat-value">{{ getPreviewSkill.cooldown }}回合</span>
                    </div>
                    <div class="stat-item">
                      <span class="stat-label">等级要求:</span>
                      <span class="stat-value">Lv.{{ getPreviewSkill.level }}</span>
                    </div>
                  </div>
                  <div v-if="getPreviewSkill.effect && getPreviewSkill.effect.type" class="skill-effect">
                    <h4>附加效果</h4>
                    <p>{{ getEffectDescription(getPreviewSkill.effect) }}</p>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑技能' : '新建技能'" width="600px">
      <el-form :model="skillForm" label-width="100px">
        <el-form-item label="技能名称">
          <el-input v-model="skillForm.name" />
        </el-form-item>
        <el-form-item label="技能编码">
          <el-input v-model="skillForm.code" />
        </el-form-item>
        <el-form-item label="技能类型">
          <el-select v-model="skillForm.type">
            <el-option label="攻击" value="attack" />
            <el-option label="治疗" value="heal" />
            <el-option label="增益" value="buff" />
            <el-option label="减益" value="debuff" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="skillForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="MP消耗">
              <el-input-number v-model="skillForm.mpCost" :min="0" :max="999" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="冷却回合">
              <el-input-number v-model="skillForm.cooldown" :min="0" :max="99" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="伤害值">
              <el-input-number v-model="skillForm.damage" :min="0" :max="9999" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="治疗值">
              <el-input-number v-model="skillForm.heal" :min="0" :max="9999" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="等级要求">
          <el-input-number v-model="skillForm.level" :min="1" :max="100" />
        </el-form-item>
        <el-divider content-position="left">附加效果</el-divider>
        <el-form-item label="效果类型">
          <el-select v-model="skillForm.effect.type" clearable>
            <el-option label="无" value="" />
            <el-option label="灼烧" value="burn" />
            <el-option label="中毒" value="poison" />
            <el-option label="冰冻" value="freeze" />
            <el-option label="攻击提升" value="attack_up" />
            <el-option label="防御提升" value="defense_up" />
            <el-option label="速度提升" value="speed_up" />
            <el-option label="攻击降低" value="attack_down" />
            <el-option label="防御降低" value="defense_down" />
          </el-select>
        </el-form-item>
        <el-row :gutter="20" v-if="skillForm.effect.type">
          <el-col :span="12">
            <el-form-item label="持续回合">
              <el-input-number v-model="skillForm.effect.duration" :min="1" :max="99" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="效果数值">
              <el-input-number v-model="skillForm.effect.value" :min="0" :max="999" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveSkill">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue';
import { useSkillStore } from '@/stores';
import { ElMessage, ElMessageBox } from 'element-plus';

const skillStore = useSkillStore();
const activeTab = ref('list');
const dialogVisible = ref(false);
const isEditing = ref(false);
const editingId = ref(null);
const previewSkill = ref(null);

const skillForm = reactive({
  name: '',
  code: '',
  type: 'attack',
  description: '',
  mpCost: 5,
  damage: 0,
  heal: 0,
  cooldown: 0,
  level: 1,
  effect: { type: '', duration: 0, value: 0 }
});

const getPreviewSkill = computed(() => {
  if (!previewSkill.value) return null;
  return skillStore.getSkillById(previewSkill.value);
});

const getTypeTag = (type) => {
  const map = { attack: 'danger', heal: 'success', buff: 'warning', debuff: 'info' };
  return map[type] || '';
};

const getTypeLabel = (type) => {
  const map = { attack: '攻击', heal: '治疗', buff: '增益', debuff: '减益' };
  return map[type] || type;
};

const getEffectDescription = (effect) => {
  if (!effect || !effect.type) return '';
  const descs = {
    burn: `灼烧：每回合造成${effect.value || 5}点伤害，持续${effect.duration}回合`,
    poison: `中毒：每回合造成${effect.value || 5}点伤害，持续${effect.duration}回合`,
    freeze: `冰冻：无法行动${effect.duration}回合`,
    attack_up: `攻击力提升${effect.value || 30}%，持续${effect.duration}回合`,
    defense_up: `防御力提升${effect.value || 50}%，持续${effect.duration}回合`,
    speed_up: `速度提升${effect.value || 50}%，持续${effect.duration}回合`,
    attack_down: `攻击力降低${effect.value || 20}%，持续${effect.duration}回合`,
    defense_down: `防御力降低${effect.value || 20}%，持续${effect.duration}回合`
  };
  return descs[effect.type] || effect.type;
};

const getNodePos = (nodeId) => {
  const node = skillStore.skillTree.nodes.find(n => n.id === nodeId);
  return node ? { x: node.x + 60, y: node.y + 30 } : { x: 0, y: 0 };
};

const selectNode = (node) => {
  if (node.skillCode) {
    const skill = skillStore.getSkillByCode(node.skillCode);
    if (skill) {
      previewSkill.value = skill.id;
      activeTab.value = 'preview';
    }
  }
};

const handleAdd = () => {
  isEditing.value = false;
  editingId.value = null;
  Object.assign(skillForm, {
    name: '', code: '', type: 'attack', description: '',
    mpCost: 5, damage: 0, heal: 0, cooldown: 0, level: 1,
    effect: { type: '', duration: 0, value: 0 }
  });
  dialogVisible.value = true;
};

const handleEdit = (skill) => {
  isEditing.value = true;
  editingId.value = skill.id;
  Object.assign(skillForm, {
    name: skill.name,
    code: skill.code,
    type: skill.type,
    description: skill.description,
    mpCost: skill.mpCost,
    damage: skill.damage,
    heal: skill.heal,
    cooldown: skill.cooldown,
    level: skill.level,
    effect: { ...(skill.effect || { type: '', duration: 0, value: 0 }) }
  });
  dialogVisible.value = true;
};

const handleDelete = async (skill) => {
  try {
    await ElMessageBox.confirm('确定要删除该技能吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    skillStore.deleteSkill(skill.id);
    ElMessage.success('删除成功');
  } catch {
    // 取消
  }
};

const saveSkill = () => {
  if (isEditing.value) {
    skillStore.updateSkill(editingId.value, { ...skillForm });
  } else {
    skillStore.addSkill({ ...skillForm });
  }
  dialogVisible.value = false;
  ElMessage.success('保存成功');
};

const resetTree = () => {
  ElMessage.info('技能树布局已重置');
};

const saveTree = () => {
  ElMessage.success('技能树已保存');
};
</script>

<style scoped>
.skill-editor {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.skill-tree-container {
  position: relative;
}

.tree-toolbar {
  margin-bottom: 16px;
}

.skill-tree-canvas {
  position: relative;
  width: 100%;
  height: 500px;
  background: #f5f7fa;
  border-radius: 8px;
  overflow: auto;
}

.tree-edges {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.tree-node {
  position: absolute;
  width: 120px;
  text-align: center;
  cursor: pointer;
  transition: transform 0.2s;
}

.tree-node:hover {
  transform: scale(1.05);
}

.tree-node .node-icon {
  width: 60px;
  height: 60px;
  margin: 0 auto 8px;
  background: #fff;
  border: 3px solid #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.tree-node.root .node-icon {
  border-color: #67c23a;
  background: #f0f9eb;
}

.tree-node .node-label {
  font-size: 12px;
  color: #303133;
  font-weight: 500;
}

.skill-detail h3 {
  margin: 0 0 8px;
  font-size: 18px;
}

.skill-desc {
  color: #909399;
  margin-bottom: 16px;
}

.skill-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-label {
  color: #909399;
  font-size: 13px;
}

.stat-value {
  font-weight: 500;
}

.stat-value.damage {
  color: #f56c6c;
}

.stat-value.heal {
  color: #67c23a;
}

.skill-effect {
  margin-top: 16px;
  padding: 12px;
  background: #fdf6ec;
  border-radius: 4px;
}

.skill-effect h4 {
  margin: 0 0 8px;
  color: #e6a23c;
}
</style>
