<template>
  <div class="generator-panel" :class="{ expanded: isExpanded }">
    <div class="panel-header" @click="togglePanel">
      <div class="header-left">
        <el-icon><MagicStick /></el-icon>
        <span>AI生成助手</span>
        <el-tag v-if="status.enabled" type="success" size="small">已连接</el-tag>
        <el-tag v-else type="info" size="small">未连接</el-tag>
      </div>
      <div class="header-right">
        <el-icon v-if="isExpanded"><ArrowDown /></el-icon>
        <el-icon v-else><ArrowUp /></el-icon>
      </div>
    </div>

    <div v-show="isExpanded" class="panel-body">
      <div class="generator-form">
        <el-form :model="form" label-position="top" size="small">
          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item label="生成类型">
                <el-select v-model="form.type" placeholder="选择类型">
                  <el-option label="NPC" value="npc" />
                  <el-option label="场景" value="scene" />
                  <el-option label="任务" value="task" />
                  <el-option label="商店" value="shop" />
                  <el-option label="道具" value="item" />
                  <el-option label="智能体" value="agent" />
                  <el-option label="对话" value="dialogue" />
                  <el-option label="流程" value="flow" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="操作">
                <el-select v-model="form.action" placeholder="选择操作">
                  <el-option label="新建" value="create" />
                  <el-option label="补全" value="complete" />
                  <el-option label="扩展" value="expand" />
                  <el-option label="古风化" value="translate" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="数量">
                <el-input-number v-model="form.count" :min="1" :max="10" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="描述/需求">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="3"
              placeholder="描述你想要生成的内容，例如：一个卖包子的老大爷，性格豪爽..."
            />
          </el-form-item>

          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item label="主题">
                <el-input v-model="form.theme" placeholder="如：古风小镇" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="朝代背景">
                <el-select v-model="form.dynasty" placeholder="选择朝代" clearable>
                  <el-option label="唐朝" value="tang" />
                  <el-option label="宋朝" value="song" />
                  <el-option label="明朝" value="ming" />
                  <el-option label="架空" value="fictional" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="风格">
                <el-select v-model="form.style" placeholder="选择风格" clearable>
                  <el-option label="正经严肃" value="serious" />
                  <el-option label="轻松幽默" value="humorous" />
                  <el-option label="神秘悬疑" value="mysterious" />
                  <el-option label="浪漫唯美" value="romantic" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item>
            <el-button type="primary" @click="handleGenerate" :loading="loading">
              <el-icon><MagicStick /></el-icon>
              生成
            </el-button>
            <el-button @click="handleTest">测试连接</el-button>
            <el-button @click="handleClear">清空结果</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="generator-result" v-if="result || error">
        <el-divider content-position="left">生成结果</el-divider>

        <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" />

        <template v-if="result">
          <div class="result-actions">
            <el-button type="primary" size="small" @click="handleApply">
              <el-icon><Check /></el-icon>
              应用到当前编辑
            </el-button>
            <el-button size="small" @click="handleCopy">
              <el-icon><CopyDocument /></el-icon>
              复制JSON
            </el-button>
            <el-button size="small" @click="handleSaveAsItem">
              <el-icon><Plus /></el-icon>
              保存为新记录
            </el-button>
          </div>

          <el-input
            v-model="resultText"
            type="textarea"
            :rows="10"
            readonly
            class="result-json"
          />
        </template>
      </div>

      <div class="generator-history">
        <el-divider content-position="left">
          历史记录
          <el-button type="primary" text size="small" @click="clearHistory">清空</el-button>
        </el-divider>
        <div class="history-list">
          <div
            v-for="(item, index) in history"
            :key="index"
            class="history-item"
            @click="loadHistory(item)"
          >
            <div class="history-header">
              <el-tag size="small">{{ item.type }}</el-tag>
              <span class="history-action">{{ item.action }}</span>
              <span class="history-time">{{ item.time }}</span>
            </div>
            <div class="history-desc">{{ item.description?.substring(0, 50) }}...</div>
          </div>
          <el-empty v-if="history.length === 0" description="暂无历史记录" :image-size="60" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { ElMessage } from 'element-plus';

const emit = defineEmits(['apply']);

const isExpanded = ref(false);
const loading = ref(false);
const result = ref(null);
const error = ref('');
const history = ref([]);

const form = ref({
  type: 'npc',
  action: 'create',
  count: 1,
  description: '',
  theme: '古风小镇',
  dynasty: 'fictional',
  style: ''
});

const status = ref({
  enabled: false,
  provider: '',
  model: '',
  base_url: ''
});

const resultText = computed(() => {
  if (!result.value) return '';
  return JSON.stringify(result.value, null, 2);
});

onMounted(() => {
  checkStatus();
  loadHistoryFromStorage();
});

const togglePanel = () => {
  isExpanded.value = !isExpanded.value;
};

const checkStatus = async () => {
  try {
    const resp = await fetch('/api/generator/status');
    const data = await resp.json();
    status.value = data;
  } catch (e) {
    status.value = { enabled: false };
  }
};

const handleGenerate = async () => {
  if (!form.value.description) {
    ElMessage.warning('请输入描述');
    return;
  }

  loading.value = true;
  error.value = '';
  result.value = null;

  try {
    const resp = await fetch('/api/generator/generate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        type: form.value.type,
        action: form.value.action,
        params: {
          description: form.value.description,
          theme: form.value.theme,
          dynasty: form.value.dynasty,
          style: form.value.style,
          count: form.value.count
        }
      })
    });

    const data = await resp.json();

    if (data.success) {
      result.value = data.data;
      addToHistory();
      ElMessage.success('生成成功');
    } else {
      error.value = data.error || '生成失败';
    }
  } catch (e) {
    error.value = `请求失败: ${e.message}`;
  } finally {
    loading.value = false;
  }
};

const handleTest = async () => {
  loading.value = true;
  error.value = '';

  try {
    const resp = await fetch('/api/generator/test', { method: 'POST' });
    const data = await resp.json();

    if (data.success) {
      result.value = data.data;
      ElMessage.success('连接测试成功');
    } else {
      error.value = data.error || '测试失败';
    }
  } catch (e) {
    error.value = `测试失败: ${e.message}`;
  } finally {
    loading.value = false;
  }
};

const handleClear = () => {
  result.value = null;
  error.value = '';
};

const handleApply = () => {
  emit('apply', {
    type: form.value.type,
    data: result.value
  });
  ElMessage.success('已应用到当前编辑');
};

const handleCopy = () => {
  navigator.clipboard.writeText(resultText.value);
  ElMessage.success('已复制到剪贴板');
};

const handleSaveAsItem = () => {
  // TODO: 保存到数据库
  ElMessage.success('已保存为新记录');
};

const addToHistory = () => {
  const item = {
    type: form.value.type,
    action: form.value.action,
    description: form.value.description,
    result: result.value,
    time: new Date().toLocaleTimeString()
  };

  history.value.unshift(item);
  if (history.value.length > 20) {
    history.value.pop();
  }

  saveHistoryToStorage();
};

const loadHistory = (item) => {
  form.value.type = item.type;
  form.value.action = item.action;
  form.value.description = item.description;
  result.value = item.result;
};

const clearHistory = () => {
  history.value = [];
  localStorage.removeItem('generator_history');
};

const saveHistoryToStorage = () => {
  localStorage.setItem('generator_history', JSON.stringify(history.value));
};

const loadHistoryFromStorage = () => {
  try {
    const saved = localStorage.getItem('generator_history');
    if (saved) {
      history.value = JSON.parse(saved);
    }
  } catch (e) {
    // ignore
  }
};
</script>

<style scoped>
.generator-panel {
  position: fixed;
  bottom: 0;
  right: 20px;
  width: 480px;
  background: #fff;
  border-radius: 8px 8px 0 0;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  transition: all 0.3s;
}

.generator-panel.expanded {
  width: 600px;
  max-height: 80vh;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 8px 8px 0 0;
  cursor: pointer;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-left .el-tag {
  margin-left: 8px;
}

.panel-body {
  padding: 16px;
  max-height: 70vh;
  overflow-y: auto;
}

.generator-form {
  margin-bottom: 16px;
}

.result-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.result-json {
  font-family: monospace;
}

.generator-history {
  margin-top: 16px;
}

.history-list {
  max-height: 200px;
  overflow-y: auto;
}

.history-item {
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.history-item:hover {
  background: #f5f7fa;
  border-color: #409eff;
}

.history-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.history-action {
  color: #909399;
  font-size: 12px;
}

.history-time {
  margin-left: auto;
  color: #c0c4cc;
  font-size: 12px;
}

.history-desc {
  color: #606266;
  font-size: 13px;
}
</style>
