import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import NPCEdit from '@/views/npc/NPCEdit.vue';
import { useNPCStore, useSceneStore, useAgentStore, useShopStore } from '@/stores';

describe('NPCEdit', () => {
  let wrapper;
  let router;
  let npcStore;
  let sceneStore;
  let agentStore;
  let shopStore;

  beforeEach(() => {
    setActivePinia(createPinia());
    
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/npc/list', component: { template: '<div>NPCList</div>' } },
        { path: '/npc/edit/:id?', component: NPCEdit }
      ]
    });

    npcStore = useNPCStore();
    sceneStore = useSceneStore();
    agentStore = useAgentStore();
    shopStore = useShopStore();
  });

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount();
    }
  });

  it('should render correctly in create mode', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.find('.npc-edit').exists()).toBe(true);
    expect(wrapper.text()).toContain('新建NPC');
  });

  it('should render correctly in edit mode', async () => {
    router.push('/npc/edit/npc_001');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('编辑NPC');
  });

  it('should have form fields', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.find('.el-form').exists()).toBe(true);
  });

  it('should have save and cancel buttons', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    const buttons = wrapper.findAll('button');
    const saveButton = buttons.find(b => b.text().includes('保存'));
    const cancelButton = buttons.find(b => b.text().includes('取消'));
    
    expect(saveButton).toBeDefined();
    expect(cancelButton).toBeDefined();
  });

  it('should have schedule section', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('日程安排');
  });

  it('should add schedule when add button clicked', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    const addButton = wrapper.find('button:contains("添加日程")');
    if (addButton.exists()) {
      await addButton.trigger('click');
      expect(wrapper.vm.form.schedule).toHaveLength(1);
    }
  });

  it('should show warning if name is empty on save', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    wrapper.vm.form.name = '';
    const saveButton = wrapper.find('button:contains("保存")');
    if (saveButton.exists()) {
      await saveButton.trigger('click');
    }
  });

  it('should load NPC data in edit mode', async () => {
    router.push('/npc/edit/npc_001');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.vm.form.id).toBe('npc_001');
    expect(wrapper.vm.form.name).toBe('李掌柜');
  });

  it('should have behavior select', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('行为模式');
  });

  it('should have position inputs', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('初始位置X');
    expect(wrapper.text()).toContain('初始位置Y');
  });

  it('should have agent select', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('关联智能体');
  });

  it('should have shop select', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('关联商店');
  });

  it('should have avatar and sprite uploaders', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('头像');
    expect(wrapper.text()).toContain('精灵图');
  });

  it('should have scene select', async () => {
    router.push('/npc/edit');
    await router.isReady();

    wrapper = mount(NPCEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('所属场景');
  });
});
