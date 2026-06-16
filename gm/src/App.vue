<template>
  <el-container class="app-container">
    <el-aside width="220px" class="app-aside">
      <div class="logo">
        <h2>古风RPG编辑器</h2>
      </div>
      <el-menu
        :default-active="currentRoute"
        :router="true"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>

        <el-sub-menu index="scene">
          <template #title>
            <el-icon><Picture /></el-icon>
            <span>场景编辑</span>
          </template>
          <el-menu-item index="/scene/list">场景列表</el-menu-item>
          <el-menu-item index="/scene/edit">场景编辑器</el-menu-item>
          <el-menu-item index="/scene/tileset">图块管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="npc">
          <template #title>
            <el-icon><User /></el-icon>
            <span>NPC编辑</span>
          </template>
          <el-menu-item index="/npc/list">NPC列表</el-menu-item>
          <el-menu-item index="/npc/edit">NPC编辑器</el-menu-item>
          <el-menu-item index="/npc/dialogue">对话树</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="agent">
          <template #title>
            <el-icon><ChatDotRound /></el-icon>
            <span>智能体配置</span>
          </template>
          <el-menu-item index="/agent/list">智能体列表</el-menu-item>
          <el-menu-item index="/agent/edit">智能体编辑</el-menu-item>
          <el-menu-item index="/agent/memory">记忆配置</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="llm">
          <template #title>
            <el-icon><Cpu /></el-icon>
            <span>大模型配置</span>
          </template>
          <el-menu-item index="/llm/provider">模型提供商</el-menu-item>
          <el-menu-item index="/llm/model">模型配置</el-menu-item>
          <el-menu-item index="/llm/test">连接测试</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="prompt">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>提示词配置</span>
          </template>
          <el-menu-item index="/prompt/template">提示词模板</el-menu-item>
          <el-menu-item index="/prompt/variable">变量管理</el-menu-item>
          <el-menu-item index="/prompt/test">提示词测试</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="shop">
          <template #title>
            <el-icon><ShoppingCart /></el-icon>
            <span>商店配置</span>
          </template>
          <el-menu-item index="/shop/list">商店列表</el-menu-item>
          <el-menu-item index="/shop/items">道具管理</el-menu-item>
          <el-menu-item index="/shop/edit">商店编辑</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="task">
          <template #title>
            <el-icon><List /></el-icon>
            <span>任务系统</span>
          </template>
          <el-menu-item index="/task/list">任务列表</el-menu-item>
          <el-menu-item index="/task/edit">任务编辑</el-menu-item>
          <el-menu-item index="/task/flow">流程编排</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="config">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统配置</span>
          </template>
          <el-menu-item index="/config/game">游戏配置</el-menu-item>
          <el-menu-item index="/config/export">数据导出</el-menu-item>
          <el-menu-item index="/config/import">数据导入</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="app-header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute !== '/'">{{ currentPageTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-button type="warning" @click="toggleGenerator">
            <el-icon><MagicStick /></el-icon>
            AI助手
          </el-button>
          <el-button type="primary" @click="handleSave">
            <el-icon><Check /></el-icon>
            保存
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </el-header>

      <el-main class="app-main">
        <router-view @apply-generator="handleApplyGenerator" />
      </el-main>
    </el-container>

    <!-- 生成智能体面板 -->
    <GeneratorPanel ref="generatorPanel" @apply="handleApplyGenerator" />
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import GeneratorPanel from '@/components/generator/GeneratorPanel.vue';

const route = useRoute();
const generatorPanel = ref(null);
const currentRoute = computed(() => route.path);

const pageTitles = {
  '/': '首页',
  '/scene/list': '场景列表',
  '/scene/edit': '场景编辑器',
  '/scene/tileset': '图块管理',
  '/npc/list': 'NPC列表',
  '/npc/edit': 'NPC编辑器',
  '/npc/dialogue': '对话树',
  '/agent/list': '智能体列表',
  '/agent/edit': '智能体编辑',
  '/agent/memory': '记忆配置',
  '/llm/provider': '模型提供商',
  '/llm/model': '模型配置',
  '/llm/test': '连接测试',
  '/prompt/template': '提示词模板',
  '/prompt/variable': '变量管理',
  '/prompt/test': '提示词测试',
  '/shop/list': '商店列表',
  '/shop/items': '道具管理',
  '/shop/edit': '商店编辑',
  '/task/list': '任务列表',
  '/task/edit': '任务编辑',
  '/task/flow': '流程编排',
  '/config/game': '游戏配置',
  '/config/export': '数据导出',
  '/config/import': '数据导入'
};

const currentPageTitle = computed(() => pageTitles[route.path] || '');

const toggleGenerator = () => {
  if (generatorPanel.value) {
    generatorPanel.value.togglePanel();
  }
};

const handleApplyGenerator = (data) => {
  console.log('Apply generator data:', data);
  ElMessage.success('已应用生成内容');
};

const handleSave = () => {
  ElMessage.success('保存成功');
};

const handleExport = () => {
  ElMessage.success('导出成功');
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
}

.app-container {
  height: 100vh;
}

.app-aside {
  background-color: #304156;
  overflow-y: auto;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background-color: #263445;
}

.logo h2 {
  font-size: 16px;
  font-weight: 500;
}

.el-menu {
  border-right: none;
}

.app-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.app-main {
  background-color: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
  padding-bottom: 60px;
}

.header-right {
  display: flex;
  gap: 10px;
}
</style>
