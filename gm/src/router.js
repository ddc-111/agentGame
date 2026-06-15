import { createRouter, createWebHistory } from 'vue-router';
import Home from './views/Home.vue';

const routes = [
  { path: '/', component: Home },
  
  // 场景编辑
  { path: '/scene/list', component: () => import('./views/scene/SceneList.vue') },
  { path: '/scene/edit/:id?', component: () => import('./views/scene/SceneEdit.vue') },
  { path: '/scene/tileset', component: () => import('./views/scene/TilesetManager.vue') },
  
  // NPC编辑
  { path: '/npc/list', component: () => import('./views/npc/NPCList.vue') },
  { path: '/npc/edit/:id?', component: () => import('./views/npc/NPCEdit.vue') },
  { path: '/npc/dialogue/:id?', component: () => import('./views/npc/DialogueTree.vue') },
  
  // 智能体配置
  { path: '/agent/list', component: () => import('./views/agent/AgentList.vue') },
  { path: '/agent/edit/:id?', component: () => import('./views/agent/AgentEdit.vue') },
  { path: '/agent/memory/:id?', component: () => import('./views/agent/MemoryConfig.vue') },
  
  // 大模型配置
  { path: '/llm/provider', component: () => import('./views/llm/ProviderConfig.vue') },
  { path: '/llm/model', component: () => import('./views/llm/ModelConfig.vue') },
  { path: '/llm/test', component: () => import('./views/llm/LLMTest.vue') },
  
  // 提示词配置
  { path: '/prompt/template', component: () => import('./views/prompt/PromptTemplate.vue') },
  { path: '/prompt/variable', component: () => import('./views/prompt/VariableManager.vue') },
  { path: '/prompt/test', component: () => import('./views/prompt/PromptTest.vue') },
  
  // 商店配置
  { path: '/shop/list', component: () => import('./views/shop/ShopList.vue') },
  { path: '/shop/items', component: () => import('./views/shop/ItemManager.vue') },
  { path: '/shop/edit/:id?', component: () => import('./views/shop/ShopEdit.vue') },
  
  // 任务系统
  { path: '/task/list', component: () => import('./views/task/TaskList.vue') },
  { path: '/task/edit/:id?', component: () => import('./views/task/TaskEdit.vue') },
  { path: '/task/flow/:id?', component: () => import('./views/task/FlowEditor.vue') },
  
  // 系统配置
  { path: '/config/game', component: () => import('./views/config/GameConfig.vue') },
  { path: '/config/export', component: () => import('./views/config/DataExport.vue') },
  { path: '/config/import', component: () => import('./views/config/DataImport.vue') }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
