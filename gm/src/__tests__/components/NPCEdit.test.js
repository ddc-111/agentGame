import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import NPCEdit from '@/views/npc/NPCEdit.vue';
import { useNPCStore, useSceneStore, useAgentStore, useShopStore } from '@/stores';

const createTestRouter = (initialRoute = '/npc/edit') => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div>Home</div>' } },
      { path: '/npc/list', component: { template: '<div>NPCList</div>' } },
      { path: '/npc/edit/:id?', component: NPCEdit }
    ]
  });
  router.push(initialRoute);
  return router;
};

describe('NPCEdit', () => {
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

  const mountComponent = async (route = '/npc/edit') => {
    const router = createTestRouter(route);
    await router.isReady();
    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, pinia]
      }
    });
    return wrapper;
  };

  it('should render correctly in create mode', async () => {
    await mountComponent('/npc/edit');
    expect(wrapper.find('.npc-edit').exists()).toBe(true);
    expect(wrapper.text()).toContain('新建NPC');
  });

  it('should render correctly in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.text()).toContain('编辑NPC');
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

  it('should show schedule section text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('日程安排');
  });

  it('should add schedule when handleAddSchedule called', async () => {
    await mountComponent();
    expect(wrapper.vm.form.schedule).toHaveLength(0);
    wrapper.vm.handleAddSchedule();
    expect(wrapper.vm.form.schedule).toHaveLength(1);
  });

  it('should remove schedule when handleRemoveSchedule called', async () => {
    await mountComponent();
    wrapper.vm.form.schedule = [{ time: '08:00', action: 'open_shop', position: { x: 0, y: 0, scene: '' } }];
    wrapper.vm.handleRemoveSchedule(0);
    expect(wrapper.vm.form.schedule).toHaveLength(0);
  });

  it('should not save if name is empty', async () => {
    await mountComponent();
    wrapper.vm.form.name = '';
    wrapper.vm.handleSave();
  });

  it('should load NPC data in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.form.id).toBe('npc_001');
    expect(wrapper.vm.form.name).toBe('李掌柜');
  });

  it('should have behavior select text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('行为模式');
  });

  it('should have position inputs text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('初始位置X');
    expect(wrapper.text()).toContain('初始位置Y');
  });

  it('should have agent select text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('关联智能体');
  });

  it('should have shop select text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('关联商店');
  });

  it('should have avatar and sprite uploaders text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('头像');
    expect(wrapper.text()).toContain('精灵图');
  });

  it('should have scene select text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('所属场景');
  });

  it('should generate id in create mode', async () => {
    await mountComponent();
    expect(wrapper.vm.form.id).toMatch(/^npc_\d+$/);
  });

  it('should have isEdit computed', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.isEdit).toBe(true);
  });

  it('should have isEdit false in create mode', async () => {
    await mountComponent('/npc/edit');
    expect(wrapper.vm.isEdit).toBe(false);
  });

  it('should have default form values', async () => {
    await mountComponent();
    expect(wrapper.vm.form.name).toBe('');
    expect(wrapper.vm.form.title).toBe('');
    expect(wrapper.vm.form.description).toBe('');
    expect(wrapper.vm.form.behaviors).toEqual([]);
  });
});
