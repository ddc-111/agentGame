<template>
  <div class="flow-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>流程编排 - {{ currentFlow?.name || 'NPC出门购物流程' }}</span>
          <div>
            <el-button @click="handleAddNode">添加节点</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <div class="flow-canvas" ref="canvasRef">
        <div class="canvas-toolbar">
          <el-button-group>
            <el-button @click="zoomIn">放大</el-button>
            <el-button @click="zoomOut">缩小</el-button>
            <el-button @click="resetZoom">重置</el-button>
          </el-button-group>
        </div>

        <div class="canvas-content" :style="{ transform: `scale(${zoom})` }">
          <svg class="flow-connections">
            <path
              v-for="edge in currentFlow?.edges"
              :key="edge.id"
              :d="getPath(edge)"
              stroke="#409eff"
              stroke-width="2"
              fill="none"
              marker-end="url(#arrowhead)"
            />
            <defs>
              <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="10" refY="3.5" orient="auto">
                <polygon points="0 0, 10 3.5, 0 7" fill="#409eff" />
              </marker>
            </defs>
          </svg>

          <div
            v-for="node in currentFlow?.nodes"
            :key="node.id"
            class="flow-node"
            :class="node.type"
            :style="{ left: node.position.x + 'px', top: node.position.y + 'px' }"
            @click="selectNode(node)"
          >
            <div class="node-header">
              <el-tag :type="getNodeTypeTag(node.type)" size="small">{{ node.type }}</el-tag>
            </div>
            <div class="node-body">
              {{ node.data.label }}
            </div>
            <div class="node-handles">
              <div class="handle input" />
              <div class="handle output" />
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-dialog v-model="nodeDialogVisible" title="节点编辑" width="600px">
      <el-form :model="selectedNode" label-width="100px">
        <el-form-item label="节点ID">
          <el-input v-model="selectedNode.id" disabled />
        </el-form-item>
        <el-form-item label="节点类型">
          <el-select v-model="selectedNode.type">
            <el-option label="开始" value="start" />
            <el-option label="结束" value="end" />
            <el-option label="动作" value="action" />
            <el-option label="条件" value="condition" />
            <el-option label="等待" value="wait" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="selectedNode.data.label" placeholder="节点标签" />
        </el-form-item>

        <template v-if="selectedNode.type === 'action'">
          <el-form-item label="动作类型">
            <el-select v-model="selectedNode.data.action">
              <el-option label="移动" value="move" />
              <el-option label="对话" value="dialogue" />
              <el-option label="购买" value="purchase" />
              <el-option label="进入" value="enter" />
              <el-option label="离开" value="leave" />
              <el-option label="等待" value="wait" />
            </el-select>
          </el-form-item>
        </template>

        <template v-if="selectedNode.type === 'condition'">
          <el-form-item label="条件表达式">
            <el-input v-model="selectedNode.data.condition" placeholder="例如: shop.isOpen()" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="nodeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveNode">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useTaskStore } from '@/stores';
import { ElMessage } from 'element-plus';

const taskStore = useTaskStore();

const zoom = ref(1);
const canvasRef = ref(null);
const nodeDialogVisible = ref(false);
const selectedNode = ref({
  id: '',
  type: 'action',
  position: { x: 0, y: 0 },
  data: { label: '', action: '', params: {} }
});

const currentFlow = computed(() => taskStore.flows[0]);

const getNodeTypeTag = (type) => {
  const map = {
    start: 'success',
    end: 'danger',
    action: 'primary',
    condition: 'warning',
    wait: 'info'
  };
  return map[type] || '';
};

const getPath = (edge) => {
  const sourceNode = currentFlow.value?.nodes.find(n => n.id === edge.source);
  const targetNode = currentFlow.value?.nodes.find(n => n.id === edge.target);

  if (!sourceNode || !targetNode) return '';

  const startX = sourceNode.position.x + 100;
  const startY = sourceNode.position.y + 40;
  const endX = targetNode.position.x;
  const endY = targetNode.position.y + 40;

  const midX = (startX + endX) / 2;

  return `M ${startX} ${startY} C ${midX} ${startY}, ${midX} ${endY}, ${endX} ${endY}`;
};

const selectNode = (node) => {
  selectedNode.value = { ...node };
  nodeDialogVisible.value = true;
};

const zoomIn = () => {
  zoom.value = Math.min(zoom.value + 0.1, 2);
};

const zoomOut = () => {
  zoom.value = Math.max(zoom.value - 0.1, 0.5);
};

const resetZoom = () => {
  zoom.value = 1;
};

const handleAddNode = () => {
  const newNode = {
    id: `node_${Date.now()}`,
    type: 'action',
    position: { x: 200, y: 200 },
    data: { label: '新节点', action: '', params: {} }
  };
  currentFlow.value.nodes.push(newNode);
};

const handleSaveNode = () => {
  const index = currentFlow.value.nodes.findIndex(n => n.id === selectedNode.value.id);
  if (index !== -1) {
    currentFlow.value.nodes[index] = { ...selectedNode.value };
  }
  nodeDialogVisible.value = false;
  ElMessage.success('节点已更新');
};

const handleSave = () => {
  taskStore.updateFlow(currentFlow.value.id, currentFlow.value);
  ElMessage.success('流程已保存');
};
</script>

<style scoped>
.flow-editor {
  max-width: 1600px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.flow-canvas {
  position: relative;
  height: calc(100vh - 250px);
  background-color: #f5f7fa;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: auto;
}

.canvas-toolbar {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
}

.canvas-content {
  position: relative;
  width: 2000px;
  height: 1000px;
  transform-origin: 0 0;
}

.flow-connections {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.flow-node {
  position: absolute;
  width: 200px;
  background-color: #fff;
  border: 2px solid #ebeef5;
  border-radius: 8px;
  cursor: move;
  user-select: none;
}

.flow-node:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.flow-node.start {
  border-color: #67c23a;
}

.flow-node.end {
  border-color: #f56c6c;
}

.flow-node.condition {
  border-color: #e6a23c;
}

.node-header {
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  border-radius: 6px 6px 0 0;
}

.node-body {
  padding: 12px;
  font-size: 14px;
  color: #303133;
}

.node-handles {
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  transform: translateY(-50%);
}

.handle {
  position: absolute;
  width: 12px;
  height: 12px;
  background-color: #409eff;
  border-radius: 50%;
}

.handle.input {
  left: -6px;
}

.handle.output {
  right: -6px;
}
</style>
