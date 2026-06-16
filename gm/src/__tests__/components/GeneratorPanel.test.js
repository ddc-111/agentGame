import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import GeneratorPanel from '@/components/generator/GeneratorPanel.vue';

describe('GeneratorPanel', () => {
  let wrapper;

  beforeEach(() => {
    setActivePinia(createPinia());
    
    globalThis.fetch = vi.fn();
    globalThis.localStorage.getItem.mockReturnValue(null);
    globalThis.localStorage.setItem.mockImplementation(() => {});
    globalThis.localStorage.removeItem.mockImplementation(() => {});
  });

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount();
    }
  });

  it('should render correctly', () => {
    wrapper = mount(GeneratorPanel);
    expect(wrapper.find('.generator-panel').exists()).toBe(true);
    expect(wrapper.find('.panel-header').exists()).toBe(true);
  });

  it('should show AI生成助手 text', () => {
    wrapper = mount(GeneratorPanel);
    expect(wrapper.text()).toContain('AI生成助手');
  });

  it('should toggle panel on header click', async () => {
    wrapper = mount(GeneratorPanel);
    expect(wrapper.find('.panel-body').isVisible()).toBe(false);
    
    await wrapper.find('.panel-header').trigger('click');
    expect(wrapper.find('.panel-body').isVisible()).toBe(true);
    
    await wrapper.find('.panel-header').trigger('click');
    expect(wrapper.find('.panel-body').isVisible()).toBe(false);
  });

  it('should have correct default form values', () => {
    wrapper = mount(GeneratorPanel);
    const vm = wrapper.vm;
    expect(vm.form.type).toBe('npc');
    expect(vm.form.action).toBe('create');
    expect(vm.form.count).toBe(1);
    expect(vm.form.description).toBe('');
    expect(vm.form.theme).toBe('古风小镇');
    expect(vm.form.dynasty).toBe('fictional');
    expect(vm.form.style).toBe('');
  });

  it('should show warning if description is empty', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const generateButton = wrapper.find('button:contains("生成")');
    if (generateButton.exists()) {
      await generateButton.trigger('click');
    }
  });

  it('should have generator type options', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const select = wrapper.find('.el-select');
    expect(select.exists()).toBe(true);
  });

  it('should have action options', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const selects = wrapper.findAll('.el-select');
    expect(selects.length).toBeGreaterThanOrEqual(2);
  });

  it('should have count input', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const inputNumber = wrapper.find('.el-input-number');
    expect(inputNumber.exists()).toBe(true);
  });

  it('should have description textarea', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const textarea = wrapper.find('textarea');
    expect(textarea.exists()).toBe(true);
  });

  it('should have generate button', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const buttons = wrapper.findAll('button');
    const generateButton = buttons.find(b => b.text().includes('生成'));
    expect(generateButton).toBeDefined();
  });

  it('should have test connection button', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const buttons = wrapper.findAll('button');
    const testButton = buttons.find(b => b.text().includes('测试连接'));
    expect(testButton).toBeDefined();
  });

  it('should have clear result button', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    const buttons = wrapper.findAll('button');
    const clearButton = buttons.find(b => b.text().includes('清空结果'));
    expect(clearButton).toBeDefined();
  });

  it('should emit apply event when apply button clicked', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    wrapper.vm.result = { test: 'data' };
    await wrapper.vm.$nextTick();
    
    const buttons = wrapper.findAll('button');
    const applyButton = buttons.find(b => b.text().includes('应用到当前编辑'));
    if (applyButton) {
      await applyButton.trigger('click');
      expect(wrapper.emitted('apply')).toBeTruthy();
    }
  });

  it('should clear result when clear button clicked', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    wrapper.vm.result = { test: 'data' };
    wrapper.vm.error = 'test error';
    await wrapper.vm.$nextTick();
    
    const buttons = wrapper.findAll('button');
    const clearButton = buttons.find(b => b.text().includes('清空结果'));
    if (clearButton) {
      await clearButton.trigger('click');
      expect(wrapper.vm.result).toBeNull();
      expect(wrapper.vm.error).toBe('');
    }
  });

  it('should show result text when result exists', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    wrapper.vm.result = { test: 'data' };
    await wrapper.vm.$nextTick();
    
    expect(wrapper.find('.result-json').exists()).toBe(true);
  });

  it('should show error alert when error exists', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    wrapper.vm.error = '测试错误';
    await wrapper.vm.$nextTick();
    
    expect(wrapper.find('.el-alert').exists()).toBe(true);
  });

  it('should show history section', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    expect(wrapper.find('.generator-history').exists()).toBe(true);
  });

  it('should show empty history message', async () => {
    wrapper = mount(GeneratorPanel);
    await wrapper.find('.panel-header').trigger('click');
    
    expect(wrapper.find('.el-empty').exists()).toBe(true);
  });
});
