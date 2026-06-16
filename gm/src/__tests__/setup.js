import { vi } from 'vitest';
import { config } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';

globalThis.fetch = vi.fn();

globalThis.localStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
};

globalThis.navigator = {
  clipboard: {
    writeText: vi.fn()
  }
};

const elStub = {
  template: '<div><slot /></div>',
  props: ['modelValue', 'label', 'value', 'type', 'size', 'title', 'placeholder', 'rows', 'min', 'max', 'step', 'disabled', 'loading', 'closable', 'show-icon', 'content-position', 'format', 'multiple', 'clearable', 'auto-upload', 'on-change', 'action', 'data']
};

const elStubWithSlots = {
  template: '<div><slot /><slot name="header" /><slot name="default" /></div>',
  props: ['modelValue', 'label', 'value', 'type', 'size', 'title', 'placeholder', 'rows', 'min', 'max', 'step', 'disabled', 'loading', 'closable', 'show-icon', 'content-position', 'format', 'multiple', 'clearable', 'auto-upload', 'on-change', 'action', 'data', 'gutter', 'span']
};

config.global.stubs = {
  'el-card': elStubWithSlots,
  'el-form': elStubWithSlots,
  'el-form-item': elStubWithSlots,
  'el-input': elStub,
  'el-input-number': elStub,
  'el-select': elStub,
  'el-option': elStub,
  'el-button': elStub,
  'el-row': elStubWithSlots,
  'el-col': elStubWithSlots,
  'el-table': {
    template: '<div><slot /></div>',
    props: ['data', 'style']
  },
  'el-table-column': {
    template: '<div><slot /></div>',
    props: ['prop', 'label', 'width']
  },
  'el-divider': elStub,
  'el-upload': elStubWithSlots,
  'el-icon': elStubWithSlots,
  'el-tag': elStub,
  'el-alert': elStub,
  'el-empty': elStub,
  'el-time-picker': elStub,
  'Plus': { template: '<span />' },
  'MagicStick': { template: '<span />' },
  'ArrowDown': { template: '<span />' },
  'ArrowUp': { template: '<span />' },
  'Check': { template: '<span />' },
  'CopyDocument': { template: '<span />' }
};
