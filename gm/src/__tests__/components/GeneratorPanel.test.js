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
    expect(wrapper.vm.isExpanded).toBe(false);

    await wrapper.find('.panel-header').trigger('click');
    expect(wrapper.vm.isExpanded).toBe(true);

    await wrapper.find('.panel-header').trigger('click');
    expect(wrapper.vm.isExpanded).toBe(false);
  });

  it('should have correct default form values', () => {
    wrapper = mount(GeneratorPanel);
    expect(wrapper.vm.form.type).toBe('npc');
    expect(wrapper.vm.form.action).toBe('create');
    expect(wrapper.vm.form.count).toBe(1);
    expect(wrapper.vm.form.description).toBe('');
    expect(wrapper.vm.form.theme).toBe('古风小镇');
    expect(wrapper.vm.form.dynasty).toBe('fictional');
    expect(wrapper.vm.form.style).toBe('');
  });

  it('should have panel-body element', () => {
    wrapper = mount(GeneratorPanel);
    expect(wrapper.find('.panel-body').exists()).toBe(true);
  });

  it('should have generator-form element', async () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.isExpanded = true;
    await wrapper.vm.$nextTick();
    expect(wrapper.find('.generator-form').exists()).toBe(true);
  });

  it('should have resultText computed property', () => {
    wrapper = mount(GeneratorPanel);
    expect(wrapper.vm.resultText).toBe('');

    wrapper.vm.result = { test: 'data' };
    expect(wrapper.vm.resultText).toBe(JSON.stringify({ test: 'data' }, null, 2));
  });

  it('should emit apply event when handleApply called', async () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.result = { test: 'data' };
    wrapper.vm.form.type = 'npc';
    wrapper.vm.handleApply();
    expect(wrapper.emitted('apply')).toBeTruthy();
    expect(wrapper.emitted('apply')[0][0]).toEqual({
      type: 'npc',
      data: { test: 'data' }
    });
  });

  it('should clear result and error when handleClear called', () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.result = { test: 'data' };
    wrapper.vm.error = 'test error';

    wrapper.vm.handleClear();
    expect(wrapper.vm.result).toBeNull();
    expect(wrapper.vm.error).toBe('');
  });

  it('should have result-json class when result exists', async () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.result = { test: 'data' };
    await wrapper.vm.$nextTick();
    expect(wrapper.find('.result-json').exists()).toBe(true);
  });

  it('should have generator-history element', async () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.isExpanded = true;
    await wrapper.vm.$nextTick();
    expect(wrapper.find('.generator-history').exists()).toBe(true);
  });

  it('should have history-list element', async () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.isExpanded = true;
    await wrapper.vm.$nextTick();
    expect(wrapper.find('.history-list').exists()).toBe(true);
  });

  it('should load history from storage on mount', () => {
    const history = [{ type: 'npc', action: 'create', description: 'test', result: {}, time: '12:00' }];
    globalThis.localStorage.getItem.mockReturnValue(JSON.stringify(history));

    wrapper = mount(GeneratorPanel);
    expect(wrapper.vm.history).toHaveLength(1);
  });

  it('should handle invalid history in storage', () => {
    globalThis.localStorage.getItem.mockReturnValue('invalid json');

    wrapper = mount(GeneratorPanel);
    expect(wrapper.vm.history).toHaveLength(0);
  });

  it('should clear history', () => {
    wrapper = mount(GeneratorPanel);
    wrapper.vm.history = [{ type: 'npc' }];
    wrapper.vm.clearHistory();
    expect(wrapper.vm.history).toHaveLength(0);
  });

  it('should load history item into form', () => {
    wrapper = mount(GeneratorPanel);
    const item = {
      type: 'scene',
      action: 'expand',
      description: 'test description',
      result: { data: 'test' }
    };
    wrapper.vm.loadHistory(item);
    expect(wrapper.vm.form.type).toBe('scene');
    expect(wrapper.vm.form.action).toBe('expand');
    expect(wrapper.vm.form.description).toBe('test description');
    expect(wrapper.vm.result).toEqual({ data: 'test' });
  });
});
