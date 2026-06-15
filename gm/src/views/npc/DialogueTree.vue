<template>
  <div class="dialogue-tree">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>对话树编辑 - {{ npc?.name || '未选择NPC' }}</span>
          <div>
            <el-button @click="handleAddNode">添加节点</el-button>
            <el-button type="primary" @click="handleSave">保存</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="node-list-card">
            <template #header>对话节点</template>
            <div class="node-list">
              <div
                v-for="node in dialogueNodes"
                :key="node.id"
                class="node-item"
                :class="{ active: currentNode?.id === node.id }"
                @click="selectNode(node)"
              >
                <div class="node-type">
                  <el-tag :type="getNodeTypeTag(node.type)">{{ node.type }}</el-tag>
                </div>
                <div class="node-content">{{ node.content?.substring(0, 30) }}...</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="16">
          <el-card v-if="currentNode">
            <template #header>节点编辑</template>
            <el-form :model="currentNode" label-width="100px">
              <el-form-item label="节点ID">
                <el-input v-model="currentNode.id" disabled />
              </el-form-item>

              <el-form-item label="节点类型">
                <el-select v-model="currentNode.type">
                  <el-option label="NPC对话" value="npc_say" />
                  <el-option label="玩家选择" value="player_choice" />
                  <el-option label="条件判断" value="condition" />
                  <el-option label="执行动作" value="action" />
                  <el-option label="结束对话" value="end" />
                </el-select>
              </el-form-item>

              <el-form-item v-if="currentNode.type === 'npc_say'" label="NPC台词">
                <el-input v-model="currentNode.content" type="textarea" :rows="4" placeholder="输入NPC要说的话" />
              </el-form-item>

              <el-form-item v-if="currentNode.type === 'player_choice'" label="玩家选项">
                <div v-for="(choice, index) in currentNode.choices" :key="index" class="choice-item">
                  <el-input v-model="choice.text" placeholder="选项文本" style="width: 200px" />
                  <el-select v-model="choice.nextNode" placeholder="跳转节点" style="width: 150px; margin-left: 10px">
                    <el-option
                      v-for="n in dialogueNodes"
                      :key="n.id"
                      :label="n.id"
                      :value="n.id"
                    />
                  </el-select>
                  <el-button type="danger" text @click="removeChoice(index)">删除</el-button>
                </div>
                <el-button @click="addChoice">添加选项</el-button>
              </el-form-item>

              <el-form-item v-if="currentNode.type === 'condition'" label="条件表达式">
                <el-input v-model="currentNode.condition" placeholder="例如: player.gold >= 100" />
              </el-form-item>

              <el-form-item v-if="currentNode.type === 'action'" label="动作类型">
                <el-select v-model="currentNode.action">
                  <el-option label="给予物品" value="give_item" />
                  <el-option label="扣除金币" value="deduct_gold" />
                  <el-option label="触发任务" value="trigger_task" />
                  <el-option label="打开商店" value="open_shop" />
                </el-select>
              </el-form-item>

              <el-form-item label="下一节点">
                <el-select v-model="currentNode.nextNode" placeholder="选择下一节点" clearable>
                  <el-option
                    v-for="n in dialogueNodes"
                    :key="n.id"
                    :label="n.id"
                    :value="n.id"
                  />
                </el-select>
              </el-form-item>
            </el-form>
          </el-card>
          <el-empty v-else description="请选择一个节点" />
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useNPCStore } from '@/stores';

const route = useRoute();
const npcStore = useNPCStore();

const npc = computed(() => npcStore.getNPCById(route.params.id));

const dialogueNodes = ref([
  {
    id: 'node_001',
    type: 'npc_say',
    content: '客官好！欢迎光临小店，看看需要些什么？',
    nextNode: 'node_002'
  },
  {
    id: 'node_002',
    type: 'player_choice',
    choices: [
      { text: '我想买些草药', nextNode: 'node_003' },
      { text: '随便看看', nextNode: 'node_004' },
      { text: '再见', nextNode: 'node_005' }
    ]
  },
  {
    id: 'node_003',
    type: 'npc_say',
    content: '草药100文一份，要几份？',
    nextNode: 'node_006'
  },
  {
    id: 'node_004',
    type: 'npc_say',
    content: '好的，客官慢慢看，有需要随时叫我。',
    nextNode: null
  },
  {
    id: 'node_005',
    type: 'end',
    content: '客官慢走，欢迎下次再来！',
    nextNode: null
  },
  {
    id: 'node_006',
    type: 'player_choice',
    choices: [
      { text: '来3份', nextNode: 'node_007' },
      { text: '来5份', nextNode: 'node_008' },
      { text: '算了', nextNode: 'node_004' }
    ]
  },
  {
    id: 'node_007',
    type: 'action',
    action: 'give_item',
    params: { itemId: 'item_001', count: 3 },
    nextNode: 'node_009'
  },
  {
    id: 'node_008',
    type: 'action',
    action: 'give_item',
    params: { itemId: 'item_001', count: 5 },
    nextNode: 'node_009'
  },
  {
    id: 'node_009',
    type: 'npc_say',
    content: '好的，已经为客官准备好了。还需要别的吗？',
    nextNode: 'node_002'
  }
]);

const currentNode = ref(null);

const selectNode = (node) => {
  currentNode.value = node;
};

const getNodeTypeTag = (type) => {
  const map = {
    npc_say: 'primary',
    player_choice: 'success',
    condition: 'warning',
    action: 'danger',
    end: 'info'
  };
  return map[type] || '';
};

const addChoice = () => {
  if (!currentNode.value.choices) {
    currentNode.value.choices = [];
  }
  currentNode.value.choices.push({ text: '', nextNode: '' });
};

const removeChoice = (index) => {
  currentNode.value.choices.splice(index, 1);
};

const handleAddNode = () => {
  const newNode = {
    id: `node_${Date.now()}`,
    type: 'npc_say',
    content: '',
    nextNode: null
  };
  dialogueNodes.value.push(newNode);
  currentNode.value = newNode;
};

const handleSave = () => {
  console.log('Save dialogue:', dialogueNodes.value);
};
</script>

<style scoped>
.dialogue-tree {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.node-list-card {
  height: calc(100vh - 250px);
  overflow-y: auto;
}

.node-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.node-item {
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.node-item:hover {
  background-color: #f5f7fa;
}

.node-item.active {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.node-type {
  margin-bottom: 8px;
}

.node-content {
  font-size: 12px;
  color: #606266;
}

.choice-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}
</style>
