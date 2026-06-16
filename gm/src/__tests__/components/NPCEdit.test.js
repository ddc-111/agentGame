import { describe, it, expect, beforeEach, vi } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import NPCEdit from '@/views/npc/NPCEdit.vue';

const createTestRouter = async (route = '/npc/edit') => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/npc/list', component: { template: '<div />' } },
      { path: '/npc/edit/:id?', component: NPCEdit }
    ]
  });
  router.push(route);
  await router.isReady();
  return router;
};

describe('NPCEdit', () => {
  let wrapper;
  let pinia;

  beforeEach(() => {
    pinia = createPinia();
    setActivePinia(pinia);
  });

  afterEach(() => {
    if (wrapper) wrapper.unmount();
  });

  const mountComponent = async (route = '/npc/edit') => {
    const router = await createTestRouter(route);
    wrapper = shallowMount(NPCEdit, {
      global: { plugins: [router, pinia] }
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

  it('should show schedule section text', async () => {
    await mountComponent();
    expect(wrapper.text()).toContain('日程安排');
  });

  it('should add schedule when handleAddSchedule called', async () => {
    await mountComponent();
    expect(wrapper.vm.form.schedule).toHaveLength(0);
    wrapper.vm.handleAddSchedule();
    expect(wrapper.vm.form.schedule).toHaveLength(1);
    expect(wrapper.vm.form.schedule[0].time).toBe('08:00');
    expect(wrapper.vm.form.schedule[0].action).toBe('open_shop');
  });

  it('should remove schedule when handleRemoveSchedule called', async () => {
    await mountComponent();
    wrapper.vm.form.schedule = [{ time: '08:00', action: 'open_shop', position: { x: 0, y: 0, scene: '' } }];
    wrapper.vm.handleRemoveSchedule(0);
    expect(wrapper.vm.form.schedule).toHaveLength(0);
  });

  it('should load NPC data in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.form.id).toBe('npc_001');
    expect(wrapper.vm.form.name).toBe('李掌柜');
    expect(wrapper.vm.form.title).toBe('杂货铺老板');
    expect(wrapper.vm.form.description).toBe('一位精明的中年商人，经营着镇上最大的杂货铺');
  });

  it('should have default form values in create mode', async () => {
    await mountComponent();
    expect(wrapper.vm.form.name).toBe('');
    expect(wrapper.vm.form.title).toBe('');
    expect(wrapper.vm.form.description).toBe('');
    expect(wrapper.vm.form.behaviors).toEqual([]);
    expect(wrapper.vm.form.schedule).toEqual([]);
    expect(wrapper.vm.form.dialogues).toEqual([]);
    expect(wrapper.vm.form.agentId).toBe('');
    expect(wrapper.vm.form.shopId).toBe('');
  });

  it('should generate id in create mode', async () => {
    await mountComponent();
    expect(wrapper.vm.form.id).toMatch(/^npc_\d+$/);
  });

  it('should have isEdit true in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.isEdit).toBe(true);
  });

  it('should have isEdit false in create mode', async () => {
    await mountComponent('/npc/edit');
    expect(wrapper.vm.isEdit).toBe(false);
  });

  it('should have default position values', async () => {
    await mountComponent();
    expect(wrapper.vm.form.position).toEqual({ x: 0, y: 0, scene: '' });
  });

  it('should load schedule data in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.form.schedule).toEqual([
      { time: '06:00', action: 'open_shop', position: { x: 400, y: 300, scene: 'scene_002' } },
      { time: '22:00', action: 'close_shop', position: { x: 200, y: 200, scene: 'scene_001' } }
    ]);
  });

  it('should load behaviors in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.form.behaviors).toEqual(['idle', 'greet', 'sell']);
  });

  it('should load agentId in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.form.agentId).toBe('agent_001');
  });

  it('should load shopId in edit mode', async () => {
    await mountComponent('/npc/edit/npc_001');
    expect(wrapper.vm.form.shopId).toBe('shop_001');
  });
});
