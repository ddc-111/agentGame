import { describe, it, expect, beforeEach, vi } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import SceneEdit from '@/views/scene/SceneEdit.vue';
import { useSceneStore } from '@/stores/scene';

const createTestRouter = async (route = '/scene/edit') => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/scene/list', component: { template: '<div />' } },
      { path: '/scene/edit/:id?', component: SceneEdit }
    ]
  });
  router.push(route);
  await router.isReady();
  return router;
};

describe('SceneEdit', () => {
  let wrapper;
  let pinia;

  beforeEach(() => {
    pinia = createPinia();
    setActivePinia(pinia);
  });

  afterEach(() => {
    if (wrapper) wrapper.unmount();
  });

  const mountComponent = async (route = '/scene/edit') => {
    const router = await createTestRouter(route);
    wrapper = shallowMount(SceneEdit, {
      global: { plugins: [router, pinia] }
    });
    return wrapper;
  };

  it('should render correctly in create mode', async () => {
    await mountComponent('/scene/edit');
    expect(wrapper.find('.scene-edit').exists()).toBe(true);
    expect(wrapper.text()).toContain('新建场景');
  });

  it('should render correctly in edit mode', async () => {
    await mountComponent('/scene/edit/scene_001');
    expect(wrapper.text()).toContain('编辑场景');
  });

  it('should show NPC section text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('场景NPC');
  });

  it('should show portal section text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('传送点');
  });

  it('should add NPC when handleAddNPC called', async () => {
    await mountComponent();
    expect(wrapper.vm.form.npcs).toHaveLength(0);
    wrapper.vm.handleAddNPC();
    expect(wrapper.vm.form.npcs).toHaveLength(1);
    expect(wrapper.vm.form.npcs[0]).toEqual({ npcId: '', x: 0, y: 0 });
  });

  it('should remove NPC when handleRemoveNPC called', async () => {
    await mountComponent();
    wrapper.vm.form.npcs = [{ npcId: 'npc_001', x: 0, y: 0 }];
    wrapper.vm.handleRemoveNPC(0);
    expect(wrapper.vm.form.npcs).toHaveLength(0);
  });

  it('should add portal when handleAddPortal called', async () => {
    await mountComponent();
    expect(wrapper.vm.form.portals).toHaveLength(0);
    wrapper.vm.handleAddPortal();
    expect(wrapper.vm.form.portals).toHaveLength(1);
    expect(wrapper.vm.form.portals[0]).toEqual({ x: 0, y: 0, targetScene: '', targetX: 0, targetY: 0 });
  });

  it('should remove portal when handleRemovePortal called', async () => {
    await mountComponent();
    wrapper.vm.form.portals = [{ x: 0, y: 0, targetScene: '', targetX: 0, targetY: 0 }];
    wrapper.vm.handleRemovePortal(0);
    expect(wrapper.vm.form.portals).toHaveLength(0);
  });

  it('should load scene data in edit mode', async () => {
    await mountComponent('/scene/edit/scene_001');
    expect(wrapper.vm.form.id).toBe('scene_001');
    expect(wrapper.vm.form.name).toBe('小镇中心');
    expect(wrapper.vm.form.description).toBe('古风小镇的中心广场，人来人往');
  });

  it('should have default form values in create mode', async () => {
    await mountComponent();
    expect(wrapper.vm.form.width).toBe(1920);
    expect(wrapper.vm.form.height).toBe(1080);
    expect(wrapper.vm.form.npcs).toEqual([]);
    expect(wrapper.vm.form.portals).toEqual([]);
  });

  it('should generate id in create mode', async () => {
    await mountComponent();
    expect(wrapper.vm.form.id).toMatch(/^scene_\d+$/);
  });

  it('should have isEdit true in edit mode', async () => {
    await mountComponent('/scene/edit/scene_001');
    expect(wrapper.vm.isEdit).toBe(true);
  });

  it('should have isEdit false in create mode', async () => {
    await mountComponent('/scene/edit');
    expect(wrapper.vm.isEdit).toBe(false);
  });
});
