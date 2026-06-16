<template>
  <div class="behavior-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>NPC行为编辑</span>
          <el-select v-model="selectedNPC" placeholder="选择NPC" style="width: 200px">
            <el-option
              v-for="npc in npcList"
              :key="npc.code"
              :label="npc.name"
              :value="npc.code"
            />
          </el-select>
        </div>
      </template>

      <el-tabs v-model="activeTab" v-if="selectedNPC">
        <el-tab-pane label="日程表" name="schedule">
          <div class="timeline-editor">
            <div class="timeline-header">
              <span>24小时日程表</span>
              <el-button type="primary" size="small" @click="addScheduleItem">添加时段</el-button>
            </div>
            <div class="timeline-bar">
              <div
                v-for="(item, index) in schedule"
                :key="index"
                class="timeline-block"
                :style="getBlockStyle(item)"
                @click="editScheduleItem(index)"
              >
                <div class="block-time">{{ item.time }}</div>
                <div class="block-action">{{ item.action }}</div>
              </div>
            </div>
            <div class="timeline-labels">
              <span v-for="h in 24" :key="h" class="hour-label">{{ (h-1).toString().padStart(2, '0') }}</span>
            </div>

            <el-table :data="schedule" style="width: 100%; margin-top: 20px">
              <el-table-column prop="time" label="时间" width="100" />
              <el-table-column prop="action" label="行为" width="150" />
              <el-table-column prop="scene" label="场景" />
              <el-table-column label="操作" width="120">
                <template #default="{ $index }">
                  <el-button type="primary" text @click="editScheduleItem($index)">编辑</el-button>
                  <el-button type="danger" text @click="removeScheduleItem($index)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="情绪配置" name="mood">
          <el-form label-width="120px">
            <el-divider content-position="left">基础情绪</el-divider>
            <el-form-item label="默认情绪">
              <el-select v-model="moodConfig.default">
                <el-option label="平静" value="calm" />
                <el-option label="开心" value="happy" />
                <el-option label="悲伤" value="sad" />
                <el-option label="愤怒" value="angry" />
                <el-option label="兴奋" value="excited" />
                <el-option label="紧张" value="nervous" />
              </el-select>
            </el-form-item>

            <el-divider content-position="left">情绪触发器</el-divider>
            <div v-for="(trigger, index) in moodConfig.triggers" :key="index" class="mood-trigger">
              <el-row :gutter="16">
                <el-col :span="8">
                  <el-input v-model="trigger.event" placeholder="触发事件" />
                </el-col>
                <el-col :span="6">
                  <el-select v-model="trigger.mood" placeholder="情绪">
                    <el-option label="开心" value="happy" />
                    <el-option label="悲伤" value="sad" />
                    <el-option label="愤怒" value="angry" />
                    <el-option label="兴奋" value="excited" />
                    <el-option label="紧张" value="nervous" />
                  </el-select>
                </el-col>
                <el-col :span="6">
                  <el-input v-model="trigger.dialogue" placeholder="触发台词" />
                </el-col>
                <el-col :span="4">
                  <el-button type="danger" text @click="removeMoodTrigger(index)">删除</el-button>
                </el-col>
              </el-row>
            </div>
            <el-button type="primary" size="small" @click="addMoodTrigger">添加触发器</el-button>

            <el-divider content-position="left">情绪变化规则</el-divider>
            <el-form-item label="情绪衰减">
              <el-switch v-model="moodConfig.decay" />
              <span class="form-tip">情绪会随时间恢复到默认值</span>
            </el-form-item>
            <el-form-item v-if="moodConfig.decay" label="衰减速率">
              <el-slider v-model="moodConfig.decayRate" :min="1" :max="100" show-input />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="行为树" name="behavior">
          <div class="behavior-tree">
            <div class="tree-toolbar">
              <el-button type="primary" size="small" @click="addBehaviorNode">添加节点</el-button>
              <el-button type="success" size="small" @click="saveBehaviorTree">保存</el-button>
            </div>
            <div class="tree-canvas">
              <div
                v-for="node in behaviorTree"
                :key="node.id"
                class="behavior-node"
                :class="node.type"
                :style="{ left: node.x + 'px', top: node.y + 'px' }"
              >
                <div class="node-header">
                  <el-tag :type="getNodeTypeTag(node.type)" size="small">{{ getNodeTypeLabel(node.type) }}</el-tag>
                  <el-button type="danger" text size="small" @click="removeBehaviorNode(node.id)">x</el-button>
                </div>
                <div class="node-content">{{ node.label }}</div>
                <div class="node-params" v-if="node.params">
                  <div v-for="(v, k) in node.params" :key="k" class="param-item">
                    <span class="param-key">{{ k }}:</span>
                    <span class="param-value">{{ v }}</span>
                  </div>
                </div>
              </div>
              <svg class="tree-connections">
                <line
                  v-for="(conn, i) in behaviorConnections"
                  :key="i"
                  :x1="getNodeCenter(conn.from).x"
                  :y1="getNodeCenter(conn.from).y"
                  :x2="getNodeCenter(conn.to).x"
                  :y2="getNodeCenter(conn.to).y"
                  stroke="#409eff"
                  stroke-width="2"
                  marker-end="url(#arrow)"
                />
                <defs>
                  <marker id="arrow" markerWidth="10" markerHeight="7" refX="10" refY="3.5" orient="auto">
                    <polygon points="0 0, 10 3.5, 0 7" fill="#409eff" />
                  </marker>
                </defs>
              </svg>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>

      <el-empty v-else description="请先选择一个NPC" />
    </el-card>

    <el-dialog v-model="scheduleDialogVisible" :title="editingScheduleIndex === -1 ? '添加时段' : '编辑时段'" width="400px">
      <el-form :model="scheduleForm" label-width="80px">
        <el-form-item label="时间">
          <el-time-picker v-model="scheduleForm.time" format="HH:mm" value-format="HH:mm" />
        </el-form-item>
        <el-form-item label="行为">
          <el-select v-model="scheduleForm.action" filterable allow-create>
            <el-option label="站岗" value="stand" />
            <el-option label="巡逻" value="patrol" />
            <el-option label="开店" value="open_shop" />
            <el-option label="关店" value="close_shop" />
            <el-option label="回家" value="go_home" />
            <el-option label="打猎" value="hunt" />
            <el-option label="休息" value="rest" />
            <el-option label="聊天" value="chat" />
          </el-select>
        </el-form-item>
        <el-form-item label="场景">
          <el-input v-model="scheduleForm.scene" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scheduleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveScheduleItem">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { ElMessage } from 'element-plus';

const activeTab = ref('schedule');
const selectedNPC = ref('');

const npcList = ref([
  { name: '老村长', code: 'npc_chief_chen' },
  { name: '李掌柜', code: 'npc_merchant_li' },
  { name: '王大娘', code: 'npc_tea_wang' },
  { name: '张铁匠', code: 'npc_blacksmith_zhang' },
  { name: '猎户老周', code: 'npc_hunter_zhou' },
  { name: '小石头', code: 'npc_kid_stone' }
]);

const schedule = ref([
  { time: '06:00', action: 'stand', scene: 'scene_village_center' },
  { time: '08:00', action: 'open_shop', scene: 'scene_general_store' },
  { time: '12:00', action: 'rest', scene: 'scene_general_store' },
  { time: '13:00', action: 'open_shop', scene: 'scene_general_store' },
  { time: '22:00', action: 'close_shop', scene: 'scene_village_center' },
  { time: '23:00', action: 'go_home', scene: '' }
]);

const moodConfig = reactive({
  default: 'calm',
  decay: true,
  decayRate: 30,
  triggers: [
    { event: '玩家购买商品', mood: 'happy', dialogue: '谢谢客官！' },
    { event: '玩家离开', mood: 'calm', dialogue: '客官慢走！' },
    { event: '玩家询问狼群', mood: 'nervous', dialogue: '别提了，最近...' }
  ]
});

const behaviorTree = ref([
  { id: 'root', type: 'selector', label: 'NPC行为根节点', x: 400, y: 20, params: {} },
  { id: 'combat', type: 'condition', label: '是否在战斗中？', x: 200, y: 120, params: { check: 'inCombat' } },
  { id: 'flee', type: 'action', label: '逃跑', x: 50, y: 220, params: { action: 'flee', priority: 'high' } },
  { id: 'shop', type: 'condition', label: '是否在商店？', x: 600, y: 120, params: { check: 'inShop' } },
  { id: 'sell', type: 'action', label: '与玩家交易', x: 450, y: 220, params: { action: 'trade' } },
  { id: 'idle', type: 'action', label: '随机闲逛', x: 750, y: 220, params: { action: 'wander', radius: 100 } }
]);

const behaviorConnections = ref([
  { from: 'root', to: 'combat' },
  { from: 'root', to: 'shop' },
  { from: 'combat', to: 'flee' },
  { from: 'shop', to: 'sell' },
  { from: 'shop', to: 'idle' }
]);

const scheduleDialogVisible = ref(false);
const editingScheduleIndex = ref(-1);
const scheduleForm = reactive({ time: '08:00', action: 'stand', scene: '' });

const getBlockStyle = (item) => {
  const hour = parseInt(item.time.split(':')[0]);
  const left = (hour / 24) * 100;
  return { left: left + '%', width: '4.16%' };
};

const addScheduleItem = () => {
  editingScheduleIndex.value = -1;
  Object.assign(scheduleForm, { time: '08:00', action: 'stand', scene: '' });
  scheduleDialogVisible.value = true;
};

const editScheduleItem = (index) => {
  editingScheduleIndex.value = index;
  Object.assign(scheduleForm, schedule.value[index]);
  scheduleDialogVisible.value = true;
};

const removeScheduleItem = (index) => {
  schedule.value.splice(index, 1);
};

const saveScheduleItem = () => {
  if (editingScheduleIndex.value === -1) {
    schedule.value.push({ ...scheduleForm });
    schedule.value.sort((a, b) => a.time.localeCompare(b.time));
  } else {
    schedule.value[editingScheduleIndex.value] = { ...scheduleForm };
  }
  scheduleDialogVisible.value = false;
  ElMessage.success('保存成功');
};

const addMoodTrigger = () => {
  moodConfig.triggers.push({ event: '', mood: 'calm', dialogue: '' });
};

const removeMoodTrigger = (index) => {
  moodConfig.triggers.splice(index, 1);
};

const getNodeTypeTag = (type) => {
  const map = { selector: 'primary', sequence: 'success', condition: 'warning', action: 'danger' };
  return map[type] || '';
};

const getNodeTypeLabel = (type) => {
  const map = { selector: '选择', sequence: '序列', condition: '条件', action: '动作' };
  return map[type] || type;
};

const getNodeCenter = (nodeId) => {
  const node = behaviorTree.value.find(n => n.id === nodeId);
  return node ? { x: node.x + 75, y: node.y + 40 } : { x: 0, y: 0 };
};

const addBehaviorNode = () => {
  const id = 'node_' + Date.now();
  behaviorTree.value.push({ id, type: 'action', label: '新行为', x: 400, y: 300, params: {} });
};

const removeBehaviorNode = (id) => {
  behaviorTree.value = behaviorTree.value.filter(n => n.id !== id);
  behaviorConnections.value = behaviorConnections.value.filter(c => c.from !== id && c.to !== id);
};

const saveBehaviorTree = () => {
  ElMessage.success('行为树已保存');
};
</script>

<style scoped>
.behavior-editor {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.timeline-editor {
  margin-top: 16px;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.timeline-bar {
  position: relative;
  height: 60px;
  background: #f5f7fa;
  border-radius: 4px;
  overflow: hidden;
}

.timeline-block {
  position: absolute;
  top: 0;
  height: 100%;
  background: #409eff;
  color: white;
  padding: 4px;
  cursor: pointer;
  overflow: hidden;
  border-right: 1px solid #fff;
}

.timeline-block:hover {
  opacity: 0.85;
}

.block-time {
  font-size: 10px;
  opacity: 0.8;
}

.block-action {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.timeline-labels {
  display: flex;
  margin-top: 4px;
}

.hour-label {
  flex: 1;
  text-align: center;
  font-size: 10px;
  color: #909399;
}

.mood-trigger {
  margin-bottom: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}

.behavior-tree {
  margin-top: 16px;
}

.tree-toolbar {
  margin-bottom: 16px;
}

.tree-canvas {
  position: relative;
  width: 100%;
  height: 400px;
  background: #f5f7fa;
  border-radius: 8px;
  overflow: auto;
}

.tree-connections {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.behavior-node {
  position: absolute;
  width: 150px;
  background: #fff;
  border: 2px solid #dcdfe6;
  border-radius: 8px;
  padding: 8px;
  cursor: pointer;
}

.behavior-node.selector { border-color: #409eff; }
.behavior-node.sequence { border-color: #67c23a; }
.behavior-node.condition { border-color: #e6a23c; }
.behavior-node.action { border-color: #f56c6c; }

.node-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.node-content {
  font-size: 12px;
  font-weight: 500;
  color: #303133;
}

.node-params {
  margin-top: 4px;
  font-size: 11px;
  color: #909399;
}

.param-item {
  display: flex;
  gap: 4px;
}

.param-key {
  font-weight: 500;
}
</style>
