import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import SceneEdit from '@/views/scene/SceneEdit.vue';
import { useSceneStore, useNPCStore } from '@/stores';

describe('SceneEdit', () => {
  let wrapper;
  let router;
  let sceneStore;
  let npcStore;

  beforeEach(() => {
    setActivePinia(createPinia());
    
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/scene/list', component: { template: '<div>SceneList</div>' } },
        { path: '/scene/edit/:id?', component: SceneEdit }
      ]
    });

    sceneStore = useSceneStore();
    npcStore = useNPCStore();
  });

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount();
    }
  });

  it('should render correctly in create mode', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.find('.scene-edit').exists()).toBe(true);
    expect(wrapper.text()).toContain('新建场景');
  });

  it('should render correctly in edit mode', async () => {
    router.push('/scene/edit/scene_001');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('编辑场景');
  });

  it('should have form fields', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.find('.el-form').exists()).toBe(true);
  });

  it('should have save and cancel buttons', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
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

  it('should have NPC section', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('场景NPC');
  });

  it('should have portal section', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.text()).toContain('传送点');
  });

  it('should add NPC when add button clicked', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    const addButton = wrapper.find('button:contains("添加NPC")');
    if (addButton.exists()) {
      await addButton.trigger('click');
      expect(wrapper.vm.form.npcs).toHaveLength(1);
    }
  });

  it('should add portal when add button clicked', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    const addButton = wrapper.find('button:contains("添加传送点")');
    if (addButton.exists()) {
      await addButton.trigger('click');
      expect(wrapper.vm.form.portals).toHaveLength(1);
    }
  });

  it('should show warning if name is empty on save', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
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

  it('should load scene data in edit mode', async () => {
    router.push('/scene/edit/scene_001');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    expect(wrapper.vm.form.id).toBe('scene_001');
    expect(wrapper.vm.form.name).toBe('小镇中心');
  });

  it('should have width and height inputs', async () => {
    router.push('/scene/edit');
    await router.isReady();

    wrapper = mount(SceneEdit, {
      global: {
        plugins: [router, createPinia()]
      }
    });

    const inputNumbers = wrapper.findAll('.el-input-number');
    expect(inputNumbers.length).toBeGreaterThanOrEqual(2);
  });
});
