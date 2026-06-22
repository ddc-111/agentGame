import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import { usePromptStore } from './prompt.js';

describe('usePromptStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = usePromptStore();
    vi.restoreAllMocks();
  });

  it('should have correct initial templates', () => {
    expect(store.templates).toHaveLength(3);
    expect(store.templates[0].id).toBe('template_001');
    expect(store.templates[0].name).toBe('NPC基础人设');
    expect(store.templates[0].category).toBe('system');
    expect(store.templates[0].content).toContain('你是{{npc_name}}');
    expect(store.templates[0].variables).toHaveLength(10);
  });

  it('should have correct initial variables', () => {
    expect(store.variables).toHaveLength(10);
    expect(store.variables[0]).toEqual({
      name: 'npc_name',
      type: 'string',
      description: 'NPC名称',
      source: 'npc'
    });
  });

  it('should add a template', () => {
    vi.spyOn(Date, 'now').mockReturnValue(1234567890);
    const newTemplate = {
      name: 'Test Template',
      category: 'custom',
      content: 'Test content {{var1}}',
      variables: [{ name: 'var1', type: 'string', description: 'Variable 1' }]
    };
    store.addTemplate(newTemplate);
    expect(store.templates).toHaveLength(4);
    const added = store.templates[3];
    expect(added.id).toBe('template_1234567890');
    expect(added.name).toBe('Test Template');
    expect(added.category).toBe('custom');
    expect(added.content).toBe('Test content {{var1}}');
    expect(added.variables).toEqual([{ name: 'var1', type: 'string', description: 'Variable 1' }]);
  });

  it('should update a template', () => {
    const id = 'template_001';
    const updateData = { name: 'Updated Template', category: 'updated' };
    store.updateTemplate(id, updateData);
    const updated = store.templates.find(t => t.id === id);
    expect(updated.name).toBe('Updated Template');
    expect(updated.category).toBe('updated');
    expect(updated.content).toContain('你是{{npc_name}}');
    expect(updated.variables).toHaveLength(10);
  });

  it('should not update if template id not found', () => {
    const initialLength = store.templates.length;
    store.updateTemplate('non_existent_id', { name: 'Test' });
    expect(store.templates).toHaveLength(initialLength);
  });

  it('should delete a template', () => {
    const idToDelete = 'template_002';
    store.deleteTemplate(idToDelete);
    expect(store.templates).toHaveLength(2);
    expect(store.templates.find(t => t.id === idToDelete)).toBeUndefined();
  });

  it('should add a variable', () => {
    const newVariable = {
      name: 'new_var',
      type: 'number',
      description: 'New variable',
      source: 'custom'
    };
    store.addVariable(newVariable);
    expect(store.variables).toHaveLength(11);
    const added = store.variables[10];
    expect(added).toEqual(newVariable);
  });

  it('should update a variable', () => {
    const name = 'npc_name';
    const updateData = { description: 'Updated description', source: 'updated' };
    store.updateVariable(name, updateData);
    const updated = store.variables.find(v => v.name === name);
    expect(updated.description).toBe('Updated description');
    expect(updated.source).toBe('updated');
    expect(updated.type).toBe('string');
  });

  it('should not update if variable name not found', () => {
    const initialLength = store.variables.length;
    store.updateVariable('non_existent_name', { description: 'Test' });
    expect(store.variables).toHaveLength(initialLength);
  });

  it('should delete a variable', () => {
    const nameToDelete = 'mood';
    store.deleteVariable(nameToDelete);
    expect(store.variables).toHaveLength(9);
    expect(store.variables.find(v => v.name === nameToDelete)).toBeUndefined();
  });
});