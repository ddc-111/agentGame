import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import SceneEdit from '@/views/scene/SceneEdit.vue';
import { useSceneStore, useNPCStore } from '@/stores';

const createTestRouter = (initialRoute = '/scene/edit') => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div>Home</div>' } },
      { path: '/scene/list', component: { template: '<div>SceneList</div>' } },
      { path: '/scene/edit/:id?', component: SceneEdit }
    ]
  });
  router.push(initialRoute);
  return router;
};

describe('SceneEdit', () => {
  let wrapper;
  let pinia;

  beforeEach(async () => {
    pinia = createPinia();
    setActivePinia(pinia);
  });

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount();
    }
  });

  const mountComponent = async (route = '/scene/edit') => {
    const router = createTestRouter(route);
    await router.isReady();
    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, pinia]
      }
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

  it('should have form element', async () => {
    await mountComponent();
    expect(wrapper.findComponent({ name: 'el-form' }).exists()).toBe(true);
  });

  it('should have save and cancel buttons', async () => {
    await mountComponent();
    const buttons = wrapper.findAllComponents({ name: 'el-button' });
    expect(buttons.length).toBeGreaterThanOrEqual(2);
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
  });

  it('should remove portal when handleRemovePortal called', async () => {
    await mountComponent();
    wrapper.vm.form.portals = [{ x: 0, y: 0, targetScene: '', targetX: 0, targetY: 0 }];
    wrapper.vm.handleRemovePortal(0);
    expect(wrapper.vm.form.portals).toHaveLength(0);
  });

  it('should not save if name is empty', async () => {
    await mountComponent();
    wrapper.vm.form.name = '';
    wrapper.vm.handleSave();
  });

  it('should load scene data in edit mode', async () => {
    await mountComponent('/scene/edit/scene_001');
    expect(wrapper.vm.form.id).toBe('scene_001');
    expect(wrapper.vm.form.name).toBe('小镇中心');
  });

  it('should have width and height in form defaults', async () => {
    await mountComponent();
    expect(wrapper.vm.form.width).toBe(1920);
    expect(wrapper.vm.form.height).toBe(1080);
  });

  it('should generate id in create mode', async () => {
    await mountComponent();
    expect(wrapper.vm.form.id).toMatch(/^scene_\d+$/);
  });

  it('should have isEdit computed', async () => {
    await mountComponent('/scene/edit/scene_001');
    expect(wrapper.vm.isEdit).toBe(true);
  });

  it('should have isEdit false in create mode', async () => {
    await mountComponent('/scene/edit');
    expect(wrapper.vm.isEdit).toBe(false);
  });
});
