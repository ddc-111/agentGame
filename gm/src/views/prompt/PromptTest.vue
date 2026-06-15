<template>
  <div class="prompt-test">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>提示词测试</span>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>测试配置</template>
            <el-form :model="testConfig" label-width="100px">
              <el-form-item label="选择模板">
                <el-select v-model="testConfig.templateId" placeholder="选择模板" @change="handleTemplateChange">
                  <el-option
                    v-for="template in promptStore.templates"
                    :key="template.id"
                    :label="template.name"
                    :value="template.id"
                  />
                </el-select>
              </el-form-item>

              <el-divider content-position="left">变量值</el-divider>

              <el-form-item
                v-for="variable in currentVariables"
                :key="variable.name"
                :label="variable.name"
              >
                <el-input
                  v-model="testConfig.variables[variable.name]"
                  :placeholder="variable.description"
                />
              </el-form-item>

              <el-form-item>
                <el-button type="primary" @click="handleRender">渲染模板</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>

        <el-col :span="12">
          <el-card>
            <template #header>渲染结果</template>
            <div v-if="renderedContent" class="render-result">
              <pre>{{ renderedContent }}</pre>
            </div>
            <el-empty v-else description="选择模板并填入变量值后点击渲染" />
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { usePromptStore } from '@/stores';

const promptStore = usePromptStore();

const testConfig = ref({
  templateId: '',
  variables: {}
});

const renderedContent = ref('');

const currentTemplate = computed(() => {
  return promptStore.templates.find(t => t.id === testConfig.value.templateId);
});

const currentVariables = computed(() => {
  return currentTemplate.value?.variables || [];
});

const handleTemplateChange = () => {
  testConfig.value.variables = {};
  renderedContent.value = '';
};

const handleRender = () => {
  if (!currentTemplate.value) return;

  let content = currentTemplate.value.content;

  // 替换变量
  for (const [key, value] of Object.entries(testConfig.value.variables)) {
    const regex = new RegExp(`\\{\\{${key}\\}\\}`, 'g');
    content = content.replace(regex, value);
  }

  renderedContent.value = content;
};
</script>

<style scoped>
.prompt-test {
  max-width: 1400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.render-result {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  max-height: 600px;
  overflow-y: auto;
}

.render-result pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
