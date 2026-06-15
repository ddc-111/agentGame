import { createRouter, createWebHistory } from 'vue-router';
import Home from './views/Home.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/npc', component: () => import('./views/NPC.vue') },
  { path: '/player', component: () => import('./views/Player.vue') }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
